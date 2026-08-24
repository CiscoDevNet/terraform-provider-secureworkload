package secureworkload

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"terraform-provider-secureworkload/secureworkload/signer"
)

// newTestClient builds a Client pointed at the given test server URL,
// bypassing New() so we don't need a real 40-byte API secret validated path.
func newTestClient(t *testing.T, maxRetries int) Client {
	t.Helper()
	s, err := signer.New("key", "0123456789012345678901234567890123456789")
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}
	return Client{
		Config: Config{MaxRetries: maxRetries},
		client: http.DefaultClient,
		signer: s,
	}
}

func TestDoRetriesOn429ThenSucceeds(t *testing.T) {
	var attempts int32
	var lastBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		body, _ := io.ReadAll(r.Body)
		lastBody = string(body)
		if n < 3 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	c := newTestClient(t, 5)

	req, err := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader([]byte(`{"hello":"world"}`)))
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	var result map[string]interface{}
	if err := c.Do(req, &result); err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("expected 3 attempts, got %d", got)
	}
	if lastBody != `{"hello":"world"}` {
		t.Fatalf("expected replayed request body, got %q", lastBody)
	}
	if result["ok"] != true {
		t.Fatalf("expected decoded result ok=true, got %+v", result)
	}
}

func TestDoGivesUpAfterMaxRetries(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	c := newTestClient(t, 3)

	req, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}

	err = c.Do(req, nil)
	if err == nil {
		t.Fatalf("expected error after exhausting retries, got nil")
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("expected 3 attempts (maxRetries), got %d", got)
	}
}

func TestShouldRetry(t *testing.T) {
	cases := map[int]bool{
		http.StatusOK:                  false,
		http.StatusBadRequest:          false,
		http.StatusUnauthorized:        false,
		http.StatusTooManyRequests:     true,
		http.StatusInternalServerError: true,
		http.StatusBadGateway:          true,
		http.StatusServiceUnavailable:  true,
		http.StatusGatewayTimeout:      true,
	}
	for code, want := range cases {
		if got := shouldRetry(code); got != want {
			t.Errorf("shouldRetry(%d) = %v, want %v", code, got, want)
		}
	}
}
