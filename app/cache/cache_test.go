package cache_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/superboomer/map-tile-provider/app/cache"
	"github.com/superboomer/map-tile-provider/app/tile"
)

func TestLoadCache(t *testing.T) {
	// Setup test environment
	tmpDir := "./tmp"
	cachePath := filepath.Join(tmpDir, "cache_test")
	os.MkdirAll(cachePath, 0o600)
	defer os.RemoveAll(cachePath) // Clean up after test

	c, err := cache.LoadCache(cachePath, time.Hour*24)
	if err != nil {
		t.Errorf("LoadCache() error = %v", err)
		return
	}
	if c == nil {
		t.Error("Expected cache to be initialized, got nil")
	}
}

func TestSaveTile(t *testing.T) {
	// Setup test environment
	tmpDir := "./tmp"
	cachePath := filepath.Join(tmpDir, "save_tile_test")
	os.MkdirAll(cachePath, 0o600)
	defer os.RemoveAll(cachePath) // Clean up after test

	c, _ := cache.LoadCache(cachePath, time.Hour*24)

	t.Run("valid tile save", func(t *testing.T) {
		testTile := &tile.Tile{X: 1, Y: 2, Z: 3}
		img := []byte("test image")
		testTile.Image = img

		err := c.SaveTile("vendor", testTile)
		if err != nil {
			t.Errorf("SaveTile() error = %v", err)
			return
		}

		expectedImagePath := filepath.Join(cachePath, "vendor", fmt.Sprintf("%d", testTile.Z), fmt.Sprintf("%d_%d.jpeg", testTile.X, testTile.Y))
		if _, errStat := os.Stat(expectedImagePath); os.IsNotExist(errStat) {
			t.Errorf("Expected image file does not exist at path: %s", expectedImagePath)
			return
		}

		actualImg, err := os.ReadFile(expectedImagePath)
		if err != nil {
			t.Errorf("Failed to read saved image file: %v", err)
			return
		}
		if !bytes.Equal(actualImg, img) {
			t.Errorf("Saved image does not match expected content")
		}
	})

	t.Run("image not provided", func(t *testing.T) {
		testTile := &tile.Tile{X: 1, Y: 2, Z: 3}
		err := c.SaveTile("vendor", testTile)
		if err == nil {
			t.Error("Expected SaveTile() to return an error for missing image")
		}
	})
}

func TestLoadFile(t *testing.T) {
	// Setup test environment
	tmpDir := os.TempDir()
	cachePath := filepath.Join(tmpDir, "load_file_test")
	os.MkdirAll(cachePath, 0o600)
	defer os.RemoveAll(cachePath) // Clean up after test

	c, _ := cache.LoadCache(cachePath, time.Hour*24)

	t.Run("file exists", func(t *testing.T) {
		testTile := &tile.Tile{X: 1, Y: 2, Z: 3}
		img := []byte("test image")
		testTile.Image = img

		// Save tile first
		err := c.SaveTile("vendor", testTile)
		if err != nil {
			t.Errorf("SaveTile() error = %v", err)
			return
		}

		loadedImg, err := c.LoadFile("vendor", testTile)
		if err != nil {
			t.Errorf("LoadFile() error = %v", err)
			return
		}
		if !bytes.Equal(loadedImg, img) {
			t.Errorf("Loaded image does not match saved image")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		testTile := &tile.Tile{X: 4, Y: 5, Z: 6}
		_, err := c.LoadFile("vendor", testTile)
		if err == nil {
			t.Error("Expected LoadFile() to return an error for missing file")
		}
	})
}
