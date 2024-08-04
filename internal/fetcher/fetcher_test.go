package fetcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"id": 123, "title": "Test Job", "text": "This is a test job posting"}`))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Save the original FetchURL function and replace it with our test version
	originalFetchURL := FetchURL
	defer func() { FetchURL = originalFetchURL }()

	FetchURL = func(itemId int64) (string, error) {
		return server.URL, nil
	}

	// Test the Fetch function
	posting, err := Fetch(123)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if posting.ID != 123 {
		t.Errorf("Expected ID 123, got %d", posting.ID)
	}
	if posting.Title != "Test Job" {
		t.Errorf("Expected title 'Test Job', got '%s'", posting.Title)
	}
	if posting.Text != "This is a test job posting" {
		t.Errorf("Expected text 'This is a test job posting', got '%s'", posting.Text)
	}
}
