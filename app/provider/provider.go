package provider

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/superboomer/maptile/app/tile"
)

//go:generate moq -out provider_mock.go -fmt goimports . Provider

// Provider is an interface which implement all necessary stuff for map provider
type Provider interface {
	GetTile(lat, long, scale float64) tile.Tile

	ID() string
	Name() string
	MaxJobs() int
	MaxZoom() int

	GetRequest(t *tile.Tile) *http.Request
}

// MapProvider contains all data about provider
type MapProvider struct {
	name       string
	id         string
	url        string
	headers    *http.Header
	maxJobs    int
	maxZoom    int
	projection *tile.Elips
}

// createProvider create new provider by specified Schema
func createProvider(schema *schema) (Provider, error) {

	p := &MapProvider{
		name:    schema.Name,
		id:      schema.ID,
		url:     schema.Request.URL,
		maxJobs: schema.MaxJobs,
		maxZoom: schema.MaxZoom,
	}
	switch schema.Projection {
	case "wgs84":
		p.projection = &tile.ElipsWGS84
	case "spherical":
		p.projection = &tile.ElipsSpherical
	default:
		return nil, fmt.Errorf("projection %v not found for provider %v", schema.Projection, schema.Name)
	}

	buildHeaders := &http.Header{}

	for _, h := range schema.Request.Headers {
		buildHeaders.Set(h.Key, h.Value)
	}

	p.headers = buildHeaders

	return p, nil
}

// GetTile calculate tile XYZ
func (p *MapProvider) GetTile(lat, long, scale float64) tile.Tile {

	tileX, tileY := tile.ConvertToTile(lat, long, scale, p.projection)

	return tile.Tile{
		X: tileX,
		Y: tileY,
		Z: int(scale),
	}
}

// MaxJobs return count of max tile downloading per request
func (p *MapProvider) MaxJobs() int {
	return p.maxJobs
}

// MaxZoom return max zoom for specified provider
func (p *MapProvider) MaxZoom() int {
	return p.maxZoom
}

// Name return provider name
func (p *MapProvider) Name() string {
	return p.name
}

// ID return provider ID
func (p *MapProvider) ID() string {
	return p.id
}

// GetRequest build http request for specified Tile
func (p *MapProvider) GetRequest(t *tile.Tile) *http.Request {

	replacer := strings.NewReplacer("{x}", fmt.Sprint(t.X), "{y}", fmt.Sprint(t.Y), "{z}", fmt.Sprint(t.Z))
	req, _ := http.NewRequest(http.MethodGet, replacer.Replace(p.url), http.NoBody)

	if p.headers != nil {
		req.Header = *p.headers
	}

	return req
}
