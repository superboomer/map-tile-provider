package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/tile"
	"go.etcd.io/bbolt"
)

func TestNewCache_Success(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test")
	defer os.RemoveAll(tmpDir)

	cache, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)
	assert.NotNil(t, cache)
}

func TestNewCache_FailedDir(t *testing.T) {
	cache, err := NewCache("", time.Hour, nil)
	assert.Error(t, err)
	assert.Nil(t, cache)
}

func TestNewCache_FailedIndex(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test")
	defer os.RemoveAll(tmpDir)

	cache, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)
	assert.NotNil(t, cache)

	cache2, err := NewCache(tmpDir, time.Hour, &bbolt.Options{Timeout: time.Second})
	assert.Error(t, err)
	assert.Nil(t, cache2)
}

func TestSaveTile_Success(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test-save")
	defer os.RemoveAll(tmpDir)

	cache, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)

	testTile := &tile.Tile{X: 1, Y: 2, Z: 3, Image: []byte("test-image")}
	err = cache.SaveTile("vendor", testTile)
	assert.NoError(t, err)
}

func TestSaveTile_FailedReadOnly(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test-save")
	defer os.RemoveAll(tmpDir)

	createCacheFile, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)

	assert.NoError(t, createCacheFile.Close())

	cache, err := NewCache(tmpDir, time.Hour, &bbolt.Options{ReadOnly: true})
	assert.NoError(t, err)

	testTile := &tile.Tile{X: 1, Y: 2, Z: 3, Image: []byte("test-image")}
	err = cache.SaveTile("vendor", testTile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update db:")
}

func TestLoadTile_Success(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test-load")
	defer os.RemoveAll(tmpDir)

	cache, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)

	testTile := &tile.Tile{X: 1, Y: 2, Z: 3, Image: []byte("test-image")}
	err = cache.SaveTile("vendor", testTile)
	assert.NoError(t, err)

	loadedTile, err := cache.LoadTile("vendor", &tile.Tile{X: 1, Y: 2, Z: 3})
	assert.NoError(t, err)
	assert.NotNil(t, loadedTile)
}

func TestLoadTile_FailedBucketError(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test-load")
	defer os.RemoveAll(tmpDir)

	cache, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)

	loadedTile, err := cache.LoadTile("vendor", &tile.Tile{X: 1, Y: 2, Z: 3})
	assert.Error(t, err)
	assert.Nil(t, loadedTile)
}

func TestLoadTile_FailedTileError(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test-load")
	defer os.RemoveAll(tmpDir)

	cache, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)

	testTile := &tile.Tile{X: 1, Y: 1, Z: 1, Image: []byte("test-image")}
	err = cache.SaveTile("vendor", testTile)
	assert.NoError(t, err)

	loadedTile, err := cache.LoadTile("vendor", &tile.Tile{X: 2, Y: 2, Z: 2})
	assert.Error(t, err)
	assert.Nil(t, loadedTile)
}

func TestSaveImage_Success(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "map-tile-provider-test-save-image")
	defer os.RemoveAll(tmpDir)

	cache, err := NewCache(tmpDir, time.Hour, nil)
	assert.NoError(t, err)

	testTile := &tile.Tile{X: 1, Y: 2, Z: 3, Image: []byte("test-image")}
	err = cache.saveImage("vendor", testTile)
	assert.NoError(t, err)

	imagePath := filepath.Join(tmpDir, "vendor", "3", "1_2.jpeg")
	_, err = os.Stat(imagePath)
	assert.NoError(t, err)
}

func TestUnixTimeEncodeDecode(t *testing.T) {
	now := time.Now()
	encoded := unixTimeEncode(now)
	decoded := unixTimeDecode(encoded)

	assert.Equal(t, now.Unix(), decoded.Unix())
}
