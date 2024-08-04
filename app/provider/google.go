package provider

import (
	"fmt"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/tile"
)

// GoogleLoader implement Provider interface for Google provider
type GoogleLoader struct {
	name    string
	key     string
	maxJobs int
	maxZoom int
	url     string
}

// Google return Google provider
func Google() Provider {
	return GoogleLoader{
		name:    "Google Maps (Satellite)",
		key:     "google",
		maxJobs: 5,
		maxZoom: 21,
		url:     "https://mts1.google.com/vt/lyrs=s",
	}
}

// GetTile calculate tile XYZ
func (l GoogleLoader) GetTile(lat, long, scale float64) tile.Tile {

	tileX, tileY := tile.ConvertToTile(lat, long, scale, &tile.ElipsSpherical)

	return tile.Tile{
		X: tileX,
		Y: tileY,
		Z: int(scale),
	}
}

// MaxJobs return count of max tile downloading per request
func (l GoogleLoader) MaxJobs() int {
	return l.maxJobs
}

// MaxZoom return max zoom for specified provider
func (l GoogleLoader) MaxZoom() int {
	return l.maxZoom
}

// Name return provider name
func (l GoogleLoader) Name() string {
	return l.name
}

// Key return provider key
func (l GoogleLoader) Key() string {
	return l.key
}

// GetRequest build http request for specified Tile
func (l GoogleLoader) GetRequest(t *tile.Tile) *http.Request {

	req, _ := http.NewRequest(http.MethodGet, l.url, http.NoBody)

	q := req.URL.Query()
	q.Add("x", fmt.Sprintf("%d", t.X))
	q.Add("y", fmt.Sprintf("%d", t.Y))
	q.Add("z", fmt.Sprintf("%d", t.Z))
	req.URL.RawQuery = q.Encode()

	return req
}
