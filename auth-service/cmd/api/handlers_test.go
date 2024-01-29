package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// RoundTripFunc is a function type representing the RoundTrip method of the http.RoundTripper interface.
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip executes a single HTTP transaction.
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewTestClient creates a new HTTP client with a custom transport.
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func Test_authenticate(t *testing.T) {
	jsonToReturn := `
{
	"error": false,
	"message": "some message"
}
`
	client := NewTestClient(func(req *http.Request) (*http.Response, error) {
		// Create a response with a status code 200 and the specified JSON data
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
			Header:     make(http.Header),
		}, nil // No error occurred, so return nil
	})

	testApp.Client = client

	postBody := map[string]interface{}{
		"email":    "me@here.com",
		"password": "verysecret",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("Expected %v but got %d", http.StatusAccepted, rr.Code)
	}
}
