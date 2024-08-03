package provider

import (
	"fmt"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/tile"
)

// OSMLoader implement Provider interface for OpenStreetMap provider
type OSMLoader struct {
	name    string
	maxJobs int
	maxZoom int
	url     string
}

// OSM return OpenStreetMap provider
func OSM() Provider {
	return OSMLoader{
		name:    "openstreetmap",
		maxJobs: 5,
		maxZoom: 19,
		url:     "https://tile.openstreetmap.org",
	}
}

// GetTile calculate tile XYZ
func (l OSMLoader) GetTile(lat, long, scale float64) tile.Tile {

	tileX, tileY := tile.ConvertToTile(lat, long, scale, &tile.ElipsSpherical)

	return tile.Tile{
		X: tileX,
		Y: tileY,
		Z: int(scale),
	}
}

// MaxJobs return count of max tile downloading per request
func (l OSMLoader) MaxJobs() int {
	return l.maxJobs
}

// MaxZoom return max zoom for specified provider
func (l OSMLoader) MaxZoom() int {
	return l.maxZoom
}

// Name return provider name
func (l OSMLoader) Name() string {
	return l.name
}

// GetRequest build http request for specified Tile
func (l OSMLoader) GetRequest(t *tile.Tile) *http.Request {

	buildRequest := fmt.Sprintf("%s/%d/%d/%d.png", l.url, t.Z, t.X, t.Y)
	req, _ := http.NewRequest(http.MethodGet, buildRequest, http.NoBody)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	return req
}
