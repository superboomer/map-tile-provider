package provider

import (
	"fmt"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/tile"
)

type GoogleLoader struct {
	name    string
	maxJobs int
	maxZoom int
	url     string
}

func Google() Provider {
	return GoogleLoader{
		name:    "google",
		maxJobs: 5,
		maxZoom: 21,
		url:     "https://mts1.google.com/vt/lyrs=s",
	}
}

func (l GoogleLoader) GetTile(lat, long, scale float64) tile.Tile {

	tileX, tileY := tile.ConvertToTile(lat, long, scale, &tile.ElipsSpherical)

	return tile.Tile{
		X: tileX,
		Y: tileY,
		Z: int(scale),
	}
}

func (l GoogleLoader) MaxJobs() int {
	return l.maxJobs
}

func (l GoogleLoader) MaxZoom() int {
	return l.maxZoom
}

func (l GoogleLoader) Name() string {
	return l.name
}

func (l GoogleLoader) GetRequest(t *tile.Tile) *http.Request {

	req, _ := http.NewRequest(http.MethodGet, l.url, nil)

	q := req.URL.Query()
	q.Add("x", fmt.Sprintf("%d", t.X))
	q.Add("y", fmt.Sprintf("%d", t.Y))
	q.Add("z", fmt.Sprintf("%d", t.Z))
	req.URL.RawQuery = q.Encode()

	return req
}
