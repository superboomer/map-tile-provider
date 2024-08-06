package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJSON_FromFile(t *testing.T) {
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
	result, err := loadJSON(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, sampleData, result)
}

func TestLoadJSONFromFile_FailedJSONUnmarshal(t *testing.T) {

	// Save the JSON data to a temporary file
	tmpFile, err := os.CreateTemp("", "sample-schema-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	_, err = tmpFile.WriteString("invalid json")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if tmpErr := tmpFile.Close(); tmpErr != nil {
		t.Fatalf("Failed to close temp file: %v", tmpErr)
	}

	// Call the function under test
	result, err := loadJSONFromFile(tmpFile.Name())
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal JSON")
}

func TestLoadJSON_FromHTTP(t *testing.T) {
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

	// Call the function under test with the mock client
	result, err := loadJSON(testServer.URL)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "MockProvider", result[0].Name)
	assert.Equal(t, "mp", result[0].ID)
}

func TestLoadJSONFromHTTP_FailURL(t *testing.T) {
	// Call the function under test with the mock client
	result, err := loadJSONFromHTTP(http.DefaultClient, "not_valid_uri")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to parse URL")
}

func TestLoadJSONFromHTTP_FailGet(t *testing.T) {
	// Call the function under test with the mock client
	result, err := loadJSONFromHTTP(http.DefaultClient, "ftp://example.com")

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to send http query")
}

func TestLoadJSONFromHTTP_FailStatusCode(t *testing.T) {
	// Setup a test HTTP server that returns a predefined Status Code
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
	}))
	defer testServer.Close()

	// Call the function under test with the mock client
	result, err := loadJSON(testServer.URL)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to fetch JSON: received status code")
}

func TestLoadJSONFromHTTP_FailJSONUnmarshal(t *testing.T) {
	// Setup a test HTTP server that returns a predefined Status Code
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Call the function under test with the mock client
	result, err := loadJSON(testServer.URL)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal JSON")
}
