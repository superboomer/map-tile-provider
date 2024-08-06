package api_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/options"
	"github.com/superboomer/map-tile-provider/app/server/api"
	"go.uber.org/zap"
)

func TestCreateAPI_Success(t *testing.T) {
	// Execute
	res, err := api.CreateAPI(zap.NewNop(), &options.Cache{Enable: false}, "./../../../example/providers.json", 512)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 512, res.MaxSide)
	assert.Nil(t, res.Cache)
	assert.Equal(t, res.Providers.GetAllID(), []string{"google", "osm", "arcgis"})
}

func TestCreateAPI_EnableCacheSuccess(t *testing.T) {
	tmpDir := "./tmp/cache"
	// Execute
	res, err := api.CreateAPI(zap.NewNop(), &options.Cache{Enable: true, Path: "./tmp/cache", Alive: 60}, "./../../../example/providers.json", 512)
	defer os.RemoveAll(tmpDir)
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 512, res.MaxSide)
	assert.NotNil(t, res.Cache)
}

func TestCreateAPI_EnableCacheFailure(t *testing.T) {
	// Execute
	res, err := api.CreateAPI(zap.NewNop(), &options.Cache{Enable: true, Path: "", Alive: 60}, "./../../../example/providers.json", 512)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCreateAPI_LoadProviderListFailure(t *testing.T) {
	// Execute
	res, err := api.CreateAPI(zap.NewNop(), &options.Cache{Enable: false}, "provider/source/invalid", 512)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can't load provider list")
	assert.Nil(t, res)
}
