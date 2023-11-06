// +build all unittests

package signer

import (
	"net/http"
	"testing"
)

var (
	APIKey    = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	APISecret = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
)

func TestCalculateBodyCheckSumReturnsCorrectCheckSumForBody(t *testing.T) {
	var calculatedCheckSumTests = []struct {
		payload  string
		checksum string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},                                                           // empty payload
		{`{"Name":"Alice","Body":"Hello","Time":1294706395881547000}`, "feb73c8ccdb1bce2e64e54888acfae104186458376662336972f2105d7033899"}, // non-empty payload
	}
	for _, test := range calculatedCheckSumTests {
		calculatedCheckSum := CalculateBodyChecksum([]byte(test.payload))
		if calculatedCheckSum != test.checksum {
			t.Errorf("Expected digest to equal %+v for %s, got %s", test.checksum, test.payload, calculatedCheckSum)
		}
	}
}

func TestCalculateRequestSignatureReturnsCorrectSignatureForAGivenRequest(t *testing.T) {
	// Simulate standard GET request with no body
	request, err := http.NewRequest("GET", "https://acme.secureworkloadpreview.com/openapi/v1/app_scopes", nil)
	request.Header.Set(ContentTypeHeaderKey, JSONContentType)
	request.Header.Set(TimestampHeaderKey, "2020-04-21T18:23:37+0000")
	signer, err := New(APIKey, APISecret)
	if err != nil {
		t.Error(err)
	}
	requestSignature, err := CalculateRequestSignature(request, *signer.rawAPISecret)
	if err != nil {
		t.Error(err)
	}
	expectedSignature := "uTixNCblmmldnLBADYcuzfxU/LHy3a2RM3jLrQwxFtc="
	if requestSignature != expectedSignature {
		t.Errorf("Expected signature for request %+v\n to be %s, got %s", request, expectedSignature, requestSignature)
	}
}

type TestPayload struct {
	Name string `json:"Name"`
	Body string `json:"Body"`
	Time int    `json:"Time"`
}

func TestSignAddsCorrectBodyChecksumAsHeader(t *testing.T) {
	var signAddsCorrectBodyChecksumAsHeaderTests = []struct {
		payload  interface{}
		checksum string
		method   string
		url      string
	}{
		{payload: nil,
			checksum: "38e0b9de817f645c4bec37c0d4a3e58baecccb040f5718dc069a72c7385a0bed",
			method:   "GET",
			url:      "https://acme.secureworkloadpreview.com/openapi/v1/app_scopes"}, // empty payload
		{
			payload: TestPayload{
				Name: "Alice",
				Body: "Hello",
				Time: 1294706395881547000,
			},
			checksum: "aa38a6462c07c34671ece67fe7933d7e1d77a9540027d92512daf0a59f6d430b",
			method:   "POST",
			url:      "https://acme.secureworkloadpreview.com/openapi/v1/app_scopes",
		}, // non-empty payload
	}
	for _, test := range signAddsCorrectBodyChecksumAsHeaderTests {
		request, err := CreateJSONRequest(test.method, test.url, test.payload)
		if err != nil {
			t.Error(err)
		}
		signer, err := New(APIKey, APISecret)
		if err != nil {
			t.Error(err)
		}
		err = signer.Sign(request)
		if err != nil {
			t.Error(err)
		}
		if test.payload != nil && request.Header[ChecksumHeaderKey][0] != test.checksum {
			t.Errorf("Expected request checksum to equal %s, got %+v", test.checksum, request.Header[ChecksumHeaderKey][0])
		}
	}
}

func TestCreateJSONRequest(t *testing.T) {
	requestMethod := http.MethodGet
	requestURL := "http://example.com/path"
	requestBody := TestPayload{
		Name: "Alice",
		Body: "Hello",
		Time: 1294706395881547000,
	}
	request, err := CreateJSONRequest(requestMethod, requestURL, requestBody)
	if err != nil {
		t.Error(err)
	}
	if request.Header[ContentTypeHeaderKey][0] != JSONContentType {
		t.Errorf("Expected created request content type to be %s, got %s", JSONContentType, request.Header[ContentTypeHeaderKey][0])
	}
	if request.Method != requestMethod {
		t.Errorf("Expected created request method to be %s, got %s", requestMethod, request.Method)
	}
}
