package secureworkload

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

// Configuration for creating a SecureWorkload API client
type Config struct {
	APIKey                 string
	APISecret              string
	APIURL                 string
	DisableTLSVerification bool
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

// Do signs and sends a request, if the provided result
// interface is not nil, the response will be json decoded to the provided interface
func (c *Client) Do(request *http.Request, result interface{}) error {
	err := c.signer.Sign(request)
	if err != nil {
		return err
	}
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		var rawBodyBuffer bytes.Buffer
		// Decode raw response, usually contains
		// additional error details
		body := io.TeeReader(response.Body, &rawBodyBuffer)
		var responseBody interface{}
		json.NewDecoder(body).Decode(&responseBody)
		return fmt.Errorf("Request %+v\n failed with status code %d\n response %+v\n%+v", request,
			response.StatusCode, responseBody,
			response)
	}
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
