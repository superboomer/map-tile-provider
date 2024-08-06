package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// schema contains all data about provider
type schema struct {
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	MaxJobs    int       `json:"max_jobs"`
	MaxZoom    int       `json:"max_zoom"`
	Projection string    `json:"proj"`
	Request    reqSchema `json:"request"`
}

type reqSchema struct {
	URL     string          `json:"url"`
	Headers []headersSchema `json:"headers"`
}

type headersSchema struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// loadJSON loads a JSON file from an HTTP URL or from the local filesystem based on the protocol.
func loadJSON(source string) ([]schema, error) {
	// Check if the source starts with "http://" or "https://"
	if len(source) > 7 && (source[:7] == "http://" || source[:8] == "https://") {
		return loadJSONFromHTTP(http.DefaultClient, source)
	}

	// Otherwise, treat it as a local file path
	return loadJSONFromFile(source)
}

// loadJSONFromHTTP loads JSON data from an HTTP URL.
func loadJSONFromHTTP(client *http.Client, urlStr string) ([]schema, error) {

	// Parse the URL
	parsedUIrl, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	resp, err := client.Get(parsedUIrl.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JSON from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JSON: received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result []schema

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

// loadJSONFromFile loads JSON data from a local file.
func loadJSONFromFile(path string) ([]schema, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var result []schema

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}
