package secureworkload

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

// APIError is returned by Client.Do for terminal (non-retryable, or
// retries-exhausted) non-2xx responses. It preserves the HTTP status code
// so callers can make behavioral decisions (e.g. treating 404 as a signal
// that a remote resource no longer exists) without parsing error strings.
type APIError struct {
	StatusCode int
	Method     string
	URL        string
	Body       interface{}
	Attempts   int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Request %s %s\n failed with status code %d\n response %+v\n (after %d attempts)",
		e.Method, e.URL, e.StatusCode, e.Body, e.Attempts)
}

// IsNotFound reports whether err is (or wraps) an *APIError with a 404
// Not Found status code.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// Configuration for creating a SecureWorkload API client
type Config struct {
	APIKey                 string
	APISecret              string
	APIURL                 string
	DisableTLSVerification bool
	// MaxRetries is the maximum number of total attempts (1 initial + retries)
	// made for a request that receives a 429 or transient server error
	// response. If zero, DefaultMaxRetries is used.
	MaxRetries int
	// RetryMaxWait, if non-zero, overrides the maximum per-attempt wait
	// duration used by the retry backoff.
	RetryMaxWait time.Duration
}

// A client for making signed HTTP requests to a SecureWorkload API
type Client struct {
	Config Config
	client *http.Client
	signer signer.Signer
}

// New creates a new SecureWorkload client based off the provided
// config, returning the client and error (if any).
func New(config Config) (Client, error) {
	signer, err := signer.New(config.APIKey, config.APISecret)
	if err != nil {
		return Client{}, err
	}
	// Remove any trailing slash to be more forgiving of user input
	config.APIURL = strings.TrimSuffix(config.APIURL, "/")
	client := Client{
		Config: config,
		signer: signer,
	}
	if config.DisableTLSVerification {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.client = &http.Client{Transport: transport}
	} else {
		client.client = http.DefaultClient
	}
	return client, nil
}

// Do signs and sends a request, retrying on HTTP 429 (Too Many Requests)
// and transient 5xx server errors with backoff. If the provided result
// interface is not nil, the response will be json decoded to the provided interface.
func (c *Client) Do(request *http.Request, result interface{}) error {
	maxRetries := c.Config.MaxRetries
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}

	// Buffer the request body (if any) so it can be replayed on retries,
	// since the body stream is consumed by both signing and sending.
	var bodyBytes []byte
	if request.Body != nil && request.Body != http.NoBody {
		var err error
		bodyBytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
		request.Body.Close()
	}

	var lastErr error
	var lastResponse *http.Response
	attempts := 0

	for attempt := 0; attempt < maxRetries; attempt++ {
		attempts++

		// Abort immediately if the caller's context has been cancelled.
		if err := request.Context().Err(); err != nil {
			return err
		}

		// Refill the (possibly nil) request body for this attempt.
		if bodyBytes != nil {
			request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			request.ContentLength = int64(len(bodyBytes))
		}

		// The signature covers a per-request Timestamp header and a body
		// checksum, so the request must be re-signed on every attempt,
		// not just the first.
		if err := c.signer.Sign(request); err != nil {
			return err
		}

		response, err := c.client.Do(request)
		if err != nil {
			lastErr = err
			// Network-level errors: back off and retry like a transient failure,
			// unless the context was cancelled/expired.
			if ctxErr := request.Context().Err(); ctxErr != nil {
				return ctxErr
			}
			if attempt < maxRetries-1 {
				time.Sleep(c.retryDelay(nil, attempt))
				continue
			}
			return lastErr
		}

		if shouldRetry(response.StatusCode) && attempt < maxRetries-1 {
			delay := c.retryDelay(response, attempt)
			// Drain and close the body so the connection can be reused.
			io.Copy(io.Discard, response.Body)
			response.Body.Close()
			lastResponse = response
			time.Sleep(delay)
			continue
		}

		if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
			defer response.Body.Close()
			var rawBodyBuffer bytes.Buffer
			// Decode raw response, usually contains
			// additional error details
			body := io.TeeReader(response.Body, &rawBodyBuffer)
			var responseBody interface{}
			json.NewDecoder(body).Decode(&responseBody)
			return &APIError{
				StatusCode: response.StatusCode,
				Method:     request.Method,
				URL:        request.URL.String(),
				Body:       responseBody,
				Attempts:   attempts,
			}
		}

		defer response.Body.Close()
		// If no result is expected, don't attempt to decode a potentially
		// empty response stream and avoid incurring EOF errors
		if result == nil {
			return nil
		}
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return err
		}
		return nil
	}

	if lastResponse != nil {
		return fmt.Errorf("Request %+v\n failed with status code %d after %d attempts\n%+v", request,
			lastResponse.StatusCode, attempts, lastResponse)
	}
	if lastErr != nil {
		return fmt.Errorf("request failed after %d attempts: %w", attempts, lastErr)
	}
	return fmt.Errorf("request failed after %d attempts", attempts)
}

// retryDelay computes the wait duration before the next attempt, honoring
// any Config.RetryMaxWait override for the caller.
func (c *Client) retryDelay(response *http.Response, attempt int) time.Duration {
	delay := retryAfterDelay(response, attempt)
	if c.Config.RetryMaxWait > 0 && delay > c.Config.RetryMaxWait {
		delay = c.Config.RetryMaxWait
	}
	return delay
}
