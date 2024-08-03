package provider

import (
	"fmt"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/tile"
)

type Provider interface {
	GetTile(lat, long, scale float64) tile.Tile

	Name() string
	MaxJobs() int
	MaxZoom() int

	GetRequest(t *tile.Tile) *http.Request
}

type ProviderList map[string]Provider

func CreateProviderList() *ProviderList {
	pl := make(ProviderList)
	return &pl
}

func (pl ProviderList) Register(p Provider) error {
	_, err := pl.Get(p.Name())
	if err == nil {
		return fmt.Errorf("provider %s already exist", p.Name())
	}

	pl[p.Name()] = p

	return nil
}

func (pl ProviderList) Get(name string) (Provider, error) {
	provider, exists := pl[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

func (pl ProviderList) GetAllNames() []string {
	names := make([]string, 0, len(pl))
	for name := range pl {
		names = append(names, name)
	}
	return names
}
