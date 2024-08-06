package downloader

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/cache"
	"github.com/superboomer/map-tile-provider/app/provider"
	"github.com/superboomer/map-tile-provider/app/tile"
)

func TestDownload_SuccessfulWithoutCacheHit(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("image data"))
	}))
	defer ts.Close()

	mockProvider := &provider.ProviderMock{
		GetRequestFunc: func(*tile.Tile) *http.Request {
			req, _ := http.NewRequest(http.MethodGet, ts.URL, http.NoBody)
			return req
		},
		MaxJobsFunc: func() int { return 2 },
		IDFunc:      func() string { return "name" },
	}

	mockCache := &cache.CacheMock{
		LoadTileFunc: func(string, *tile.Tile) ([]byte, error) {
			return nil, errors.New("not found")
		},
		SaveTileFunc: func(string, *tile.Tile) error { return nil },
	}

	downloader := NewMapDownloader(http.DefaultClient)

	tiles := []tile.Tile{{X: 1, Y: 2, Z: 3}}
	downloadedTiles, err := downloader.Download(mockCache, mockProvider, tiles...)

	assert.NoError(t, err)
	assert.Equal(t, len(tiles), len(downloadedTiles))
	// Assert interactions
	assert.Len(t, mockProvider.GetRequestCalls(), 1) // Assuming 1 calls expected
}

func TestDownload_FailedDownload(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("image data"))
	}))
	defer ts.Close()

	mockProvider := &provider.ProviderMock{
		GetRequestFunc: func(testTile *tile.Tile) *http.Request {
			if testTile.X == 1 {
				req, _ := http.NewRequest(http.MethodGet, ts.URL, http.NoBody)
				return req
			}
			return nil
		},
		MaxJobsFunc: func() int { return 1 },
		IDFunc:      func() string { return "name" },
	}

	mockCache := &cache.CacheMock{
		LoadTileFunc: func(string, *tile.Tile) ([]byte, error) { return nil, fmt.Errorf("not found") }, // Cache miss
		SaveTileFunc: func(string, *tile.Tile) error { return nil },                                    // Assume save succeeds
	}

	downloader := NewMapDownloader(http.DefaultClient)

	_, _ = downloader.Download(mockCache, mockProvider, []tile.Tile{{X: 1, Y: 2, Z: 3}}...)

	_, err := downloader.Download(mockCache, mockProvider, []tile.Tile{{X: 4, Y: 5, Z: 6}}...)
	assert.Error(t, err) // Expect an error due to download failure
	// Assert interactions
	assert.Contains(t, err.Error(), "request is empty")
	assert.Len(t, mockProvider.GetRequestCalls(), 2) // Still expect 2 calls despite failure
}

func TestDownload_SuccessfulLoadFromCache(t *testing.T) {
	mockProvider := &provider.ProviderMock{
		MaxJobsFunc: func() int { return 2 },
		IDFunc:      func() string { return "name" },
		GetRequestFunc: func(testTile *tile.Tile) *http.Request {
			return &http.Request{}
		},
	}

	mockCache := &cache.CacheMock{
		LoadTileFunc: func(string, *tile.Tile) ([]byte, error) {
			return []byte{}, nil
		},
		SaveTileFunc: func(string, *tile.Tile) error { return nil },
	}

	downloader := NewMapDownloader(http.DefaultClient)

	tiles := []tile.Tile{{X: 1, Y: 2, Z: 3}}
	downloadedTiles, err := downloader.Download(mockCache, mockProvider, tiles...)

	assert.NoError(t, err)
	assert.Equal(t, len(tiles), len(downloadedTiles))
	// Assert interactions
	assert.Len(t, mockCache.LoadTileCalls(), 1)      // Assuming 1 calls expected
	assert.Len(t, mockProvider.GetRequestCalls(), 1) // Assuming 1 calls expected
	assert.Len(t, mockCache.SaveTileCalls(), 0)      // Assuming 0 calls expected

}

func TestDownload_FailedRequest(t *testing.T) {

	mockProvider := &provider.ProviderMock{
		GetRequestFunc: func(testTile *tile.Tile) *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "", http.NoBody)
			return req
		},
		MaxJobsFunc: func() int { return 1 },
		IDFunc:      func() string { return "name" },
	}

	downloader := NewMapDownloader(http.DefaultClient)

	_, err := downloader.Download(nil, mockProvider, []tile.Tile{{X: 4, Y: 5, Z: 6}}...)
	assert.Error(t, err) // Expect an error due to download failure
	// Assert interactions
	assert.Contains(t, err.Error(), "error occurred when sending request to the server")
	assert.Len(t, mockProvider.GetRequestCalls(), 1) //  expect 1 calls despite failure
}

func TestDownload_FailedServerReturnedInvalidCode(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	mockProvider := &provider.ProviderMock{
		GetRequestFunc: func(testTile *tile.Tile) *http.Request {
			req, _ := http.NewRequest(http.MethodGet, ts.URL, http.NoBody)
			return req
		},
		MaxJobsFunc: func() int { return 1 },
		IDFunc:      func() string { return "name" },
	}

	downloader := NewMapDownloader(http.DefaultClient)

	_, err := downloader.Download(nil, mockProvider, []tile.Tile{{X: 4, Y: 5, Z: 6}}...)
	assert.Error(t, err) // Expect an error due to download failure
	// Assert interactions
	assert.Contains(t, err.Error(), "server returned invalid status code")
	assert.Len(t, mockProvider.GetRequestCalls(), 1) // Still expect 2 calls despite failure
}
