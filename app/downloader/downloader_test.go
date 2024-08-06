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
	assert.Len(t, mockCache.SaveTileCalls(), 1)      // Assuming 1 saves due to cache misses
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

	tiles := []tile.Tile{{X: 1, Y: 2, Z: 3}, {X: 4, Y: 5, Z: 6}}
	_, err := downloader.Download(mockCache, mockProvider, tiles...)

	assert.Error(t, err) // Expect an error due to download failure
	// Assert interactions
	assert.Contains(t, err.Error(), "request is empty")
	assert.Len(t, mockProvider.GetRequestCalls(), 2) // Still expect 2 calls despite failure
	assert.Len(t, mockCache.SaveTileCalls(), 1)      // One save due to cache miss and successful download
}
