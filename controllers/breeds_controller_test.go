package controllers

import (
	//"encoding/json"
	"fmt"
	"io"
	"strings"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"
)

type BreedsControllerTest struct {
	web.Controller
}

// Mocked HTTP Client
type MockClient struct{}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	// Simulate the successful response
	if req.URL.String() == "https://api.thecatapi.com/v1/breeds" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[{"id": "1", "name": "Abyssinian", "temperament": "Active", "origin": "Egypt", "description": "A playful and curious cat breed", "image": {"url": "https://link-to-image.com"}}]`)),
			Header:     make(http.Header),
		}, nil
	}

	// Simulate an error response
	return nil, fmt.Errorf("network error")
}

func TestFetchBreeds_Success(t *testing.T) {
	// Mock HTTP server to simulate the cat API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id": "1", "name": "Abyssinian", "temperament": "Active", "origin": "Egypt", "description": "A playful and curious cat breed", "image": {"url": "https://link-to-image.com"}}]`))
	}))
	defer mockServer.Close()

	// Mock API key
	apiKey := "test-api-key"

	// Create channels for the result and error
	ch := make(chan []Breed)
	errCh := make(chan error)

	// Run the fetch function in a separate goroutine
	go fetchBreeds(apiKey, ch, errCh)

	// Wait for result
	select {
	case breeds := <-ch:
		// Assert that the response contains the correct breed data
		assert.Len(t, breeds, 1)
		assert.Equal(t, "Abyssinian", breeds[0].Name)
	case err := <-errCh:
		t.Fatalf("Unexpected error: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestFetchBreeds_Failure(t *testing.T) {
	// Mock HTTP server to simulate a failure scenario
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a failure (e.g., 500 Internal Server Error)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// Mock API key
	apiKey := "test-api-key"

	// Create channels for the result and error
	ch := make(chan []Breed)
	errCh := make(chan error)

	// Run the fetch function in a separate goroutine
	go fetchBreeds(apiKey, ch, errCh)

	// Wait for result
	select {
	case err := <-errCh:
		// Assert that an error was received
		assert.Contains(t, err.Error(), "error fetching breeds")
	case <-ch:
		t.Fatal("Expected error, but received result")
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}
