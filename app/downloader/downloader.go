package downloader

import (
	"fmt"
	"io"
	"net/http"

	"github.com/superboomer/maptile/app/cache"
	"github.com/superboomer/maptile/app/provider"
	"github.com/superboomer/maptile/app/tile"
)

//go:generate moq -out downloader_mock.go  -fmt goimports . Downloader

// Downloader implements basic downloader interface
type Downloader interface {
	Download(c cache.Cache, l provider.Provider, tiles ...tile.Tile) ([]tile.Tile, error)
	Merge(side int, centerTile tile.Tile, tiles ...tile.Tile) ([]byte, error)
}

// MapDownloader implements interface Downloader
type MapDownloader struct {
	client *http.Client
}

// NewMapDownloader create new MapDownloader with specified httpClient
func NewMapDownloader(client *http.Client) *MapDownloader {
	return &MapDownloader{client: client}
}

type downloadQuery struct {
	Tile    tile.Tile
	Request *http.Request
	Error   error
}

// Download orchestrates the concurrent downloading of multiple tiles using a specified provider
func (m *MapDownloader) Download(c cache.Cache, l provider.Provider, tiles ...tile.Tile) ([]tile.Tile, error) {
	jobs := make(chan downloadQuery, l.MaxJobs())
	results := make(chan downloadQuery, len(tiles))

	for w := 1; w <= l.MaxJobs(); w++ {
		go m.worker(c, l.ID(), jobs, results)
	}

	for _, p := range tiles {
		jobs <- downloadQuery{Tile: p, Request: l.GetRequest(&p), Error: nil}
	}

	close(jobs)

	var result = make([]tile.Tile, 0)

	for a := 1; a <= len(tiles); a++ {
		r := <-results

		if r.Error != nil {
			return result, r.Error
		}

		result = append(result, r.Tile)
	}

	return result, nil
}

// worker download image
func (m *MapDownloader) worker(c cache.Cache, vendor string, jobs <-chan downloadQuery, results chan<- downloadQuery) {
	for j := range jobs {
		if c != nil {
			cacheImg, err := c.LoadTile(vendor, &j.Tile)
			if err == nil {
				j.Tile.Image = cacheImg
				results <- j
				continue
			}
		}

		if j.Request == nil {
			j.Error = fmt.Errorf("request is empty")
			results <- j
			continue
		}

		resp, err := m.client.Do(j.Request)
		if err != nil {
			j.Error = fmt.Errorf("error occurred when sending request to the server: err=%w", err)
			results <- j
			continue
		}

		img, err := io.ReadAll(resp.Body)
		if err != nil {
			j.Error = fmt.Errorf("can't readAll body from server answer: err=%w", err)
			results <- j
			_ = resp.Body.Close()
			continue
		}

		if resp.StatusCode != 200 {
			j.Error = fmt.Errorf("server returned invalid status code: code=%d", resp.StatusCode)
			results <- j
			_ = resp.Body.Close()
			continue
		}

		j.Tile.Image = img
		_ = resp.Body.Close()
		results <- j

		if c != nil {
			go c.SaveTile(vendor, &j.Tile)
		}
	}
}
