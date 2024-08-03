package downloader

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/superboomer/map-tile-provider/app/tile"
)

type MockProvider struct {
	MaxJobsCount int
	RequestFunc  func(t *tile.Tile) *http.Request
}

func (mp *MockProvider) MaxJobs() int {
	return mp.MaxJobsCount
}

func (mp *MockProvider) MaxZoom() int {
	return 20 // Return a fixed number of jobs for testing.
}

func (mp *MockProvider) Name() string {
	return "name"
}

func (mp *MockProvider) GetTile(lat, long, scale float64) tile.Tile {
	return tile.Tile{}
}

func (mp *MockProvider) GetRequest(t *tile.Tile) *http.Request {
	return mp.RequestFunc(t)
}

func TestStartMultiDownload_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("image data"))
	}))
	defer ts.Close()

	mockProvider := &MockProvider{
		MaxJobsCount: 2,
		RequestFunc: func(t *tile.Tile) *http.Request {
			req, _ := http.NewRequest("GET", ts.URL, http.NoBody)
			return req
		},
	}

	tiles := []tile.Tile{
		{X: 1},
		{X: 2},
	}

	results, err := StartMultiDownload(nil, mockProvider, tiles...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != len(tiles) {
		t.Fatalf("expected %d results, got %d", len(tiles), len(results))
	}

	for _, result := range results {
		if string(result.Image) != "image data" {
			t.Errorf("expected image data to be 'image data', got %s", string(result.Image))
		}
	}
}

func TestStartMultiDownload_ErrorHandling(t *testing.T) {
	mockProvider := &MockProvider{
		MaxJobsCount: 1,
		RequestFunc: func(t *tile.Tile) *http.Request {
			req, _ := http.NewRequest("GET", "http://invalid.url", http.NoBody)
			return req
		},
	}

	tiles := []tile.Tile{
		{X: 1},
	}

	results, err := StartMultiDownload(nil, mockProvider, tiles...)
	if err == nil {
		t.Fatal("expected an error but got none")
	}

	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestStartMultiDownload_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	mockProvider := &MockProvider{
		MaxJobsCount: 1,
		RequestFunc: func(t *tile.Tile) *http.Request {
			req, _ := http.NewRequest("GET", ts.URL, http.NoBody)
			return req
		},
	}

	tiles := []tile.Tile{
		{X: 1},
		{X: 2},
	}

	results, err := StartMultiDownload(nil, mockProvider, tiles...)
	if err == nil {
		t.Fatal("expected an error due to 404 but got none")
	}

	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestStartMultiDownload_BadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	mockProvider := &MockProvider{
		MaxJobsCount: 1,
		RequestFunc: func(t *tile.Tile) *http.Request {
			req, _ := http.NewRequest("GET", ts.URL, http.NoBody)
			return req
		},
	}

	tiles := []tile.Tile{
		{X: 1},
		{X: 2},
	}

	results, err := StartMultiDownload(nil, mockProvider, tiles...)
	if err == nil {
		t.Fatal("expected an error due to 400 Bad Request but got none")
	}

	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestStartMultiDownload_ConcurrentDownloads(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("image data"))
	}))
	defer ts.Close()

	mockProvider := &MockProvider{
		MaxJobsCount: 2,
		RequestFunc: func(t *tile.Tile) *http.Request {
			req, _ := http.NewRequest("GET", ts.URL, http.NoBody)
			return req
		},
	}

	tiles := []tile.Tile{
		{X: 1},
		{X: 2},
	}

	results, err := StartMultiDownload(nil, mockProvider, tiles...)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != len(tiles) {
		t.Fatalf("expected %d results, got %d", len(tiles), len(results))
	}

	for _, result := range results {
		if string(result.Image) != "image data" {
			t.Errorf("expected image data to be 'image data', got %s", string(result.Image))
		}
	}
}
