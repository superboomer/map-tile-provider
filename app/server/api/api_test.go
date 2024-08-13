package api_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/maptile/app/options"
	"github.com/superboomer/maptile/app/server/api"
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
	assert.Contains(t, res.Providers.GetAllID(), "google")
	assert.Contains(t, res.Providers.GetAllID(), "arcgis")
	assert.Contains(t, res.Providers.GetAllID(), "osm")
}

func TestCreateAPI_EnableCacheSuccess(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "cache-test")
	// Execute
	res, err := api.CreateAPI(zap.NewNop(), &options.Cache{Enable: true, Path: tmpDir, Alive: 60}, "./../../../example/providers.json", 512)
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
