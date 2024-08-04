package provider

import (
	"fmt"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/tile"
)

// ArcGISLoader implement Provider interface for ArcGIS provider
type ArcGISLoader struct {
	maxJobs int
	maxZoom int
	url     string
	name    string
	key     string
}

// ArcGIS return ArcGIS provider
func ArcGIS() Provider {
	return ArcGISLoader{
		name:    "ArcGIS (Satellite)",
		key:     "arcgis",
		maxJobs: 5,
		maxZoom: 19,
		url:     "https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile",
	}
}

// GetTile calculate tile XYZ
func (l ArcGISLoader) GetTile(lat, long, scale float64) tile.Tile {

	tileX, tileY := tile.ConvertToTile(lat, long, scale, &tile.ElipsSpherical)

	return tile.Tile{
		X: tileX,
		Y: tileY,
		Z: int(scale),
	}
}

// MaxJobs return count of max tile downloading per request
func (l ArcGISLoader) MaxJobs() int {
	return l.maxJobs
}

// MaxZoom return max zoom for specified provider
func (l ArcGISLoader) MaxZoom() int {
	return l.maxZoom
}

// Name return provider name
func (l ArcGISLoader) Name() string {
	return l.name
}

// Key return provider key
func (l ArcGISLoader) Key() string {
	return l.key
}

// GetRequest build http request for specified Tile
func (l ArcGISLoader) GetRequest(t *tile.Tile) *http.Request {

	buildRequest := fmt.Sprintf("%s/%d/%d/%d", l.url, t.Z, t.Y, t.X)
	req, _ := http.NewRequest(http.MethodGet, buildRequest, http.NoBody)

	return req
}
