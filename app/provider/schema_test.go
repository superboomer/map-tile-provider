package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJSONFromFile(t *testing.T) {
	// Prepare sample JSON data
	sampleData := []schema{
		{Name: "OpenStreetMap", ID: "osm", MaxJobs: 100, MaxZoom: 19, Projection: "EPSG:3857", Request: reqSchema{URL: "https://example.com"}},
	}
	jsonData, _ := json.Marshal(sampleData)

	// Save the JSON data to a temporary file
	tmpFile, err := os.CreateTemp("", "sample-schema-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	_, err = tmpFile.Write(jsonData)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if tmpErr := tmpFile.Close(); tmpErr != nil {
		t.Fatalf("Failed to close temp file: %v", tmpErr)
	}

	// Call the function under test
	result, err := loadJSONFromFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, sampleData, result)
}

// TestLoadJSONFromHTTP tests the loadJSONFromHTTP function with a mocked HTTP response
func TestLoadJSONFromHTTP(t *testing.T) {
	// Setup a test HTTP server that returns a predefined JSON response
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		jsonPayload := []schema{
			{Name: "MockProvider", ID: "mp", MaxJobs: 100, MaxZoom: 19, Projection: "EPSG:3857", Request: reqSchema{URL: "https://example.com/data"}},
		}
		json.NewEncoder(rw).Encode(jsonPayload)
	}))
	defer testServer.Close()

	// Create a mock HTTP client
	mockClient := &http.Client{}

	// Call the function under test with the mock client
	result, err := loadJSONFromHTTP(mockClient, testServer.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "MockProvider", result[0].Name)
	assert.Equal(t, "mp", result[0].ID)
}
