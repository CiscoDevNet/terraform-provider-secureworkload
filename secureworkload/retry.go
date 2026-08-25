package secureworkload

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// DefaultMaxRetries is the default number of total attempts (1 initial + retries)
// made by Client.Do when it encounters a 429 or transient server error.
const DefaultMaxRetries = 5

// baseRetryWait and maxRetryWait bound the exponential backoff with jitter
// used when the server doesn't provide a Retry-After header.
const (
	baseRetryWait = 1 * time.Second
	maxRetryWait  = 30 * time.Second
)

// shouldRetry reports whether the given HTTP status code represents a
// transient failure worth retrying: 429 (Too Many Requests) or a
// 5xx server error that is typically transient (500, 502, 503, 504).
func shouldRetry(statusCode int) bool {
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// retryAfterDelay determines how long to wait before the next retry attempt.
// It honors a Retry-After response header if present (either as an integer
// number of seconds or an HTTP-date), otherwise it falls back to exponential
// backoff with jitter, capped at maxRetryWait.
func retryAfterDelay(resp *http.Response, attempt int) time.Duration {
	if resp != nil {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, err := strconv.Atoi(ra); err == nil {
				d := time.Duration(secs) * time.Second
				if d > maxRetryWait {
					d = maxRetryWait
				}
				if d > 0 {
					return d
				}
			} else if t, err := http.ParseTime(ra); err == nil {
				d := time.Until(t)
				if d > maxRetryWait {
					d = maxRetryWait
				}
				if d > 0 {
					return d
				}
			}
		}
	}
	// Exponential backoff with jitter: base * 2^attempt, +/- up to 50% jitter,
	// capped at maxRetryWait.
	backoff := baseRetryWait << uint(attempt)
	if backoff > maxRetryWait {
		backoff = maxRetryWait
	}
	jitter := time.Duration(rand.Int63n(int64(backoff) + 1))
	delay := backoff/2 + jitter/2
	if delay > maxRetryWait {
		delay = maxRetryWait
	}
	return delay
}

// Ready is a type of function that reports
// readiness of some state or action, returning
// bool for readiness and error (if any).
type Ready func() bool

// Await waits until the ready function is ready
// or errors, returning success and error (if any).
// To stop waiting, send on the stop channel.z
// It checks if the function is ready once and then retries
// the specified number of times with an exponential backoff between each attempt
func Await(ready Ready, maxRetries int) bool {
	for tries := 0; tries <= maxRetries; tries++ {
		success := ready()
		if !success {
			if tries != maxRetries {
				// exponentially back off before the next attempt
				// https://github.com/adonovan/gopl.io/blob/77e9f810f3c2502e9c641b97e09f9721424090f5/ch5/wait/wait.go#L30
				time.Sleep((1 * time.Second) << uint(tries))
			}
			continue
		}
		return true
	}
	return false
}
