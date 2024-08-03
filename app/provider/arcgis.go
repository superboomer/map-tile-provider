package provider

import (
	"fmt"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/tile"
)

type ArcGISLoader struct {
	maxJobs int
	maxZoom int
	url     string
	name    string
}

func ArcGIS() Provider {
	return ArcGISLoader{
		name:    "arcgis",
		maxJobs: 5,
		maxZoom: 19,
		url:     "https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile",
	}
}

func (l ArcGISLoader) GetTile(lat, long, scale float64) tile.Tile {

	tileX, tileY := tile.ConvertToTile(lat, long, scale, &tile.ElipsSpherical)

	return tile.Tile{
		X: tileX,
		Y: tileY,
		Z: int(scale),
	}
}

func (l ArcGISLoader) MaxJobs() int {
	return l.maxJobs
}

func (l ArcGISLoader) MaxZoom() int {
	return l.maxZoom
}

func (l ArcGISLoader) Name() string {
	return l.name
}

func (l ArcGISLoader) GetRequest(t *tile.Tile) *http.Request {

	buildRequest := fmt.Sprintf("%s/%d/%d/%d", l.url, t.Z, t.Y, t.X)
	req, _ := http.NewRequest(http.MethodGet, buildRequest, nil)

	return req
}
