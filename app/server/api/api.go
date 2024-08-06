package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/superboomer/map-tile-provider/app/cache"
	"github.com/superboomer/map-tile-provider/app/downloader"
	"github.com/superboomer/map-tile-provider/app/options"
	"github.com/superboomer/map-tile-provider/app/provider"
	"go.uber.org/zap"
)

// API represent struct for business logic
type API struct {
	Cache      cache.Cache
	Providers  provider.List
	Downloader downloader.Downloader

	Logger *zap.Logger

	MaxSide int // max side value
}

// CreateAPI create API struct
func CreateAPI(logger *zap.Logger, cacheOpts *options.Cache, providerSource string, maxSide int) (*API, error) {

	pl, err := provider.LoadProviderList(providerSource)
	if err != nil {
		return nil, fmt.Errorf("can't load provider list: %w", err)
	}

	api := &API{
		Cache:      nil,
		Logger:     logger,
		Providers:  pl,
		MaxSide:    maxSide,
		Downloader: downloader.NewMapDownloader(http.DefaultClient),
	}

	if cacheOpts.Enable {
		logger.Info("cache enabled", zap.String("path", cacheOpts.Path), zap.Duration("alive", time.Minute*time.Duration(cacheOpts.Alive)))
		с, err := cache.NewCache(cacheOpts.Path, time.Minute*time.Duration(cacheOpts.Alive))
		if err != nil {
			return nil, fmt.Errorf("can't load cache: %w", err)
		}
		api.Cache = с
	}

	return api, nil
}
