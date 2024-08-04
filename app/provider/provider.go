package provider

import (
	"fmt"
	"net/http"

	"github.com/superboomer/map-tile-provider/app/tile"
)

// Provider is an interface which implement all necessary stuff for map provider
type Provider interface {
	GetTile(lat, long, scale float64) tile.Tile

	Key() string
	Name() string
	MaxJobs() int
	MaxZoom() int

	GetRequest(t *tile.Tile) *http.Request
}

// List is a map with all registered providers
type List map[string]Provider

// CreateProviderList create empty ProviderList
func CreateProviderList() *List {
	pl := make(List)
	return &pl
}

// Register new Provider in ProviderList
func (pl List) Register(p Provider) error {
	_, err := pl.Get(p.Key())
	if err == nil {
		return fmt.Errorf("provider %s already exist", p.Name())
	}

	pl[p.Key()] = p

	return nil
}

// Get return specified by name provider
func (pl List) Get(key string) (Provider, error) {
	provider, exists := pl[key]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", key)
	}
	return provider, nil
}

// GetAllKey return all regisitered providers name
func (pl List) GetAllKey() []string {
	keys := make([]string, 0, len(pl))
	for key := range pl {
		keys = append(keys, key)
	}
	return keys
}
