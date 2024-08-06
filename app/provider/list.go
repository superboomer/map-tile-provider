package provider

import "fmt"

//go:generate moq -out list_mock.go  -fmt goimports . List

// List is a provider list
type List interface {
	Register(p Provider) error

	GetAllID() []string
	Get(key string) (Provider, error)
}

// MapList is a map with all registered providers
type MapList map[string]Provider

// createProviderList create empty ProviderList
func createProviderList() *MapList {
	pl := make(MapList)
	return &pl
}

// Register new Provider in ProviderList
func (pl MapList) Register(p Provider) error {
	_, err := pl.Get(p.ID())
	if err == nil {
		return fmt.Errorf("provider %s (%s) already exist", p.Name(), p.ID())
	}

	pl[p.ID()] = p

	return nil
}

// Get return specified by name provider
func (pl MapList) Get(key string) (Provider, error) {
	provider, exists := pl[key]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", key)
	}
	return provider, nil
}

// GetAllID return all regisitered providers ids
func (pl MapList) GetAllID() []string {
	ids := make([]string, 0, len(pl))
	for id := range pl {
		ids = append(ids, id)
	}
	return ids
}
