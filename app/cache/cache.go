package cache

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"

	"github.com/superboomer/map-tile-provider/app/tile"
)

//go:generate moq -out cache_mock.go . Cache

// Cache describe basic cache for tiles
type Cache interface {
	SaveTile(vendor string, t *tile.Tile) error
	LoadTile(vendor string, t *tile.Tile) ([]byte, error)
}

// MapCache manages cached tiles with both in-memory and persistent storage
type MapCache struct {
	db    *bbolt.DB
	path  string
	alive time.Duration
	mutex sync.RWMutex
}

// NewCache initializes a new Cache instance
func NewCache(path string, alive time.Duration) (*MapCache, error) {
	err := os.MkdirAll(path, 0o700)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	dbPath := filepath.Join(path, "index.db")
	db, err := bbolt.Open(dbPath, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	return &MapCache{db: db, path: path, alive: alive, mutex: sync.RWMutex{}}, nil
}

// SaveTile saves a tile to both BoltDB and disk storage
func (c *MapCache) SaveTile(vendor string, t *tile.Tile) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(vendor))
		if err != nil {
			return err
		}

		key := []byte(fmt.Sprintf("%d_%d_%d", t.X, t.Y, t.Z))
		value := unixTimeEncode(time.Now())
		return bucket.Put(key, value)
	})

	if err != nil {
		return fmt.Errorf("failed to update db: %w", err)
	}

	return c.saveImage(vendor, t)
}

// LoadTile attempts to load a tile from cache, checking both BoltDB and disk storage
func (c *MapCache) LoadTile(vendor string, t *tile.Tile) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var img []byte
	var err error

	err = c.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(vendor))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		key := []byte(fmt.Sprintf("%d_%d_%d", t.X, t.Y, t.Z))
		value := bucket.Get(key)
		if value == nil || !time.Now().Before(unixTimeDecode(value).Add(c.alive)) {
			return fmt.Errorf("tile not found or expired")
		}

		img, err = os.ReadFile(filepath.Clean(filepath.Join(c.path, vendor, fmt.Sprintf("%d", t.Z), fmt.Sprintf("%d_%d.jpeg", t.X, t.Y))))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load tile: %w", err)
	}

	return img, nil
}

// saveImage saves an image file to disk
func (c *MapCache) saveImage(vendor string, t *tile.Tile) error {
	dirPath := filepath.Join(c.path, vendor, fmt.Sprintf("%d", t.Z))
	filePath := filepath.Join(dirPath, fmt.Sprintf("%d_%d.jpeg", t.X, t.Y))

	err := os.MkdirAll(dirPath, 0o700)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.Write(t.Image)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}

	return nil
}

// unixTimeEncode encodes time.Time to []byte
func unixTimeEncode(t time.Time) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(t.Unix()))
	return buf
}

// unixTimeDecode decodes []byte to time.Time
func unixTimeDecode(b []byte) time.Time {
	return time.Unix(int64(binary.BigEndian.Uint64(b)), 0)
}
