package downloader

import (
	"fmt"
	"io"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/cache"
	"github.com/superboomer/map-tile-provider/app/provider"
	"github.com/superboomer/map-tile-provider/app/tile"
)

type downloadQuery struct {
	Tile    tile.Tile
	Request *http.Request
	Error   error
}

func StartMultiDownload(cache *cache.Cache, l provider.Provider, tiles ...tile.Tile) ([]tile.Tile, error) {
	jobs := make(chan downloadQuery, l.MaxJobs())
	results := make(chan downloadQuery, len(tiles))

	for w := 1; w <= l.MaxJobs(); w++ {
		go worker(cache, l.Name(), jobs, results)
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

func worker(cache *cache.Cache, vendor string, jobs <-chan downloadQuery, results chan<- downloadQuery) {
	for j := range jobs {
		if cache != nil {
			cacheImg, err := cache.LoadFile(vendor, &j.Tile)
			if err == nil {
				j.Tile.Image = cacheImg
				results <- j
				continue
			}
		}

		resp, err := http.DefaultClient.Do(j.Request)
		if err != nil {
			j.Error = fmt.Errorf("error occured when sending request to the server: err=%w", err)
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
			j.Error = fmt.Errorf("server returned invalid status code: code=%d, body=%s", resp.StatusCode, string(img))
			results <- j
			_ = resp.Body.Close()
			continue
		}

		j.Tile.Image = img
		_ = resp.Body.Close()
		results <- j

		if cache != nil {
			go cache.SaveTile(vendor, &j.Tile)
		}
	}
}
