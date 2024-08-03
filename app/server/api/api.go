package api

import (
	"fmt"
	"time"

	"github.com/superboomer/map-tile-provider/app/cache"
	"github.com/superboomer/map-tile-provider/app/options"
	"github.com/superboomer/map-tile-provider/app/provider"
	"go.uber.org/zap"
)

// API represent struct for business logic
type API struct {
	Cache     *cache.Cache
	Providers *provider.List
	Logger    *zap.Logger
}

// Init connect to kernel and create API struct
func Init(logger *zap.Logger, cacheOpts *options.Cache) (*API, error) {

	pl := provider.CreateProviderList()

	// Define a slice of provider names and their corresponding functions.
	providers := map[string]func() provider.Provider{
		"google":        provider.Google,
		"arcgis":        provider.ArcGIS,
		"openstreetmap": provider.OSM,
	}

	// Loop through the providers and register each one.
	for name, f := range providers {
		if err := pl.Register(f()); err != nil {
			return nil, fmt.Errorf("error occurred when registering new provider %s: %w", name, err)
		}
	}

	api := &API{
		Cache:     nil,
		Logger:    logger,
		Providers: pl,
	}

	if cacheOpts.Enable {
		logger.Info("cache enabled", zap.String("path", cacheOpts.Path), zap.Duration("alive", time.Minute*time.Duration(cacheOpts.Alive)))
		с, err := cache.LoadCache(cacheOpts.Path, time.Minute*time.Duration(cacheOpts.Alive))
		if err != nil {
			return nil, fmt.Errorf("can't load cache: %w", err)
		}
		api.Cache = с
	}

	return api, nil
}
