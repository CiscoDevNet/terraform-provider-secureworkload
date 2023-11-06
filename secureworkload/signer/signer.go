// Package signer implements the ability to sign a request for the purpose of
// authenticating a request with a SecureWorkload API endpoint.
// Based on https://pypi.org/project/tetpyclient/ 
package signer
import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// HTTP Request Header for request body checksum.
	ChecksumHeaderKey    = "X-Tetration-Cksum"
	UserAgentHeaderKey   = "User-Agent"
	UserAgentHeaderValue = "Cisco Tetration Golang Client"
	TimestampHeaderKey   = "Timestamp"
	// ISO8601 UTC Datetime
	TimestampFormatString  = "2006-01-02T15:04:05-0700"
	UserIdHeaderKey        = "Id"
	AuthorizationHeaderKey = "Authorization"
	ContentTypeHeaderKey   = "Content-Type"
	JSONContentType        = "application/json"
	APISecretByteLength    = 40
)

// Signer contains information and methods for signing a request.
type Signer struct {
	apiKey       string
	apiSecret    string
	rawAPISecret *[APISecretByteLength]byte
}

// New returns a new signer capable of
// signing a request and error (if any).
func New(apiKey string, apiSecret string) (Signer, error) {
	rawApiSecret := &[APISecretByteLength]byte{}
	copied := copy(rawApiSecret[:], []byte(apiSecret))
	if copied != APISecretByteLength {
		return Signer{}, fmt.Errorf("invalid number %d of api secret key bytes, required %d", copied, APISecretByteLength)
	}
	signer := Signer{
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		rawAPISecret: rawApiSecret,
	}
	return signer, nil
}

// Sign signs a request, modifying the request as needed including
// adding the signed authorization header and returning error (if any).
func (s *Signer) Sign(request *http.Request) error {
	// Calculate body checksum if body present
	if request.Body != http.NoBody {
		var rawBodyBuffer bytes.Buffer
		// Duplicate any bytes read from the body so we can refill
		// the request body after reading it in order to calculate a checksum
		body := io.TeeReader(request.Body, &rawBodyBuffer)
		var bodyBytes []byte
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		// Calculate and set checksum based off request body
		bodyChecksum := CalculateBodyChecksum(bodyBytes)
		request.Header.Set(ChecksumHeaderKey, bodyChecksum)
		// Refill request body with the data read
		request.Body = ioutil.NopCloser(&rawBodyBuffer)
	}
	// Hello from Golan(d)g!
	request.Header.Set(UserAgentHeaderKey, UserAgentHeaderValue)
	// Identify request as coming from the user identified by this API key
	request.Header.Set(UserIdHeaderKey, s.apiKey)
	// Set request timestamp value as close to the end of this function to
	// minimize drift between calculation and request processing
	request.Header.Set(TimestampHeaderKey, time.Now().UTC().Format(TimestampFormatString))
	// Calculate signature for authorizing request
	requestSignature, err := CalculateRequestSignature(request, *s.rawAPISecret)
	if err != nil {
		return err
	}
	request.Header.Set(AuthorizationHeaderKey, requestSignature)
	return nil
}

// CalculateRequestSignature calculates
// the signature for the given request
// and private key signing over:
// the request method,
// the request url path including query params
// the request body checksum header value,
// the request content type header value
// and request timestamp header value
// returning the string signature and error (if any).
func CalculateRequestSignature(request *http.Request, privateKey [APISecretByteLength]byte) (string, error) {
	mac := hmac.New(sha256.New, privateKey[:])
	mac.Write([]byte(fmt.Sprintf("%s\n", request.Method)))
	mac.Write([]byte(fmt.Sprintf("%s\n", fmt.Sprintf("%s", request.URL.RequestURI()))))
	mac.Write([]byte(fmt.Sprintf("%s\n", request.Header.Get(ChecksumHeaderKey))))
	mac.Write([]byte(fmt.Sprintf("%s\n", request.Header.Get(ContentTypeHeaderKey))))
	mac.Write([]byte(fmt.Sprintf("%s\n", request.Header.Get(TimestampHeaderKey))))
	rawRequestSignature := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(rawRequestSignature), nil
}

// CalculateBodyChecksum returns the sha256 hash
// for the given body.
func CalculateBodyChecksum(body []byte) string {
	hash := sha256.Sum256(body)
	return fmt.Sprintf("%x", hash)
}

// CreateJSONRequest isolates duplicate code in creating
// HTTP requests for JSON APIs, returning a ready to send
// http request for the given method and url with the params
// as JSON encoded body into the body and error (if any).
func CreateJSONRequest(method string, url string, params interface{}) (*http.Request, error) {
	var request *http.Request
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&params)
	if err != nil {
		return request, err
	}
	request, err = http.NewRequest(method, url, &buf)
	if err != nil {
		return request, err
	}
	if params == nil {
		request.Body = http.NoBody
	}
	request.Header.Set(ContentTypeHeaderKey, JSONContentType)
	return request, nil
}
