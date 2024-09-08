package walrus

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWalrusClient_PutData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected method PUT, got %s", r.Method)
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		expectedBody := []byte("test data")
		if !bytes.Equal(body, expectedBody) {
			t.Errorf("Expected body %s, got %s", expectedBody, body)
		}

		response := `{"newlyCreated":{"blobObject":{"blobId":"test-blob-id"}}}`
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := NewWalrusClient(server.URL)

	putData := []byte("test data")
	blobID, err := client.PutData(putData)
	if err != nil {
		t.Fatalf("Failed to put data: %v", err)
	}

	expectedBlobID := "test-blob-id"
	if blobID != expectedBlobID {
		t.Errorf("Expected blob ID %s, but got %s", expectedBlobID, blobID)
	}
}

func TestWalrusClient_GetData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/v1/") {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write([]byte("some byte data"))
			return
		}

		http.Error(w, "Not Found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewWalrusClient(server.URL)

	blobID := "test-blob-id"
	getResponse, err := client.GetData(blobID)
	if err != nil {
		t.Fatalf("Failed to get data: %v", err)
	}

	expectedGetResponse := "some byte data"
	if string(getResponse) != expectedGetResponse {
		t.Errorf("Expected GET response %s, but got %s", expectedGetResponse, string(getResponse))
	}
}
