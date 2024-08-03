package cache

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/superboomer/map-tile-provider/app/tile"
	bolt "go.etcd.io/bbolt"
)

// Cache contains all necessary stuff for cache
type Cache struct {
	db    *bolt.DB
	path  string
	alive time.Duration
}

func unixTimeEncode(t time.Time) []byte {
	buf := make([]byte, 8)
	u := uint64(t.Unix())
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func unixTimeDecode(b []byte) time.Time {
	i := int64(binary.BigEndian.Uint64(b))
	return time.Unix(i, 0)
}

// LoadCache create Cache struct
func LoadCache(path string, alive time.Duration) (*Cache, error) {
	err := os.MkdirAll(path, 0o600)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(filepath.Join(path, "index.db"), 0o600, nil)
	if err != nil {
		return nil, err
	}

	return &Cache{db: db, path: path, alive: alive}, nil
}

// LoadFile load tile image from cache
func (c *Cache) LoadFile(vendor string, t *tile.Tile) ([]byte, error) {
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(vendor))

		if b == nil {
			return fmt.Errorf("cache is invalid")
		}

		v := b.Get([]byte(fmt.Sprintf("%d_%d_%d", t.X, t.Y, t.Z)))
		if v == nil {
			return fmt.Errorf("cache is invalid")
		}

		if !time.Now().Before(unixTimeDecode(v).Add(c.alive)) {
			return fmt.Errorf("cache is invalid")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	img, err := os.ReadFile(filepath.Clean(filepath.Join(c.path, vendor, fmt.Sprintf("%d", t.Z), fmt.Sprintf("%d_%d.jpeg", t.X, t.Y))))
	if err != nil {
		return nil, fmt.Errorf("can't read image from cache: %w", err)
	}

	return img, nil
}

// SaveTile save tile to boldDB and write image on disk
func (c *Cache) SaveTile(vendor string, t *tile.Tile) error {
	if t.Image == nil {
		return fmt.Errorf("image not provided")
	}

	return c.db.Update(func(tx *bolt.Tx) error {
		err := c.saveImage(vendor, t)
		if err != nil {
			return err
		}

		b, err := tx.CreateBucketIfNotExists([]byte(vendor))
		if err != nil {
			return err
		}
		return b.Put([]byte(fmt.Sprintf("%d_%d_%d", t.X, t.Y, t.Z)), unixTimeEncode(time.Now()))
	})
}

// saveImage create cache folder if need and save image on disk
func (c *Cache) saveImage(vendor string, t *tile.Tile) error {
	err := os.MkdirAll(filepath.Join(c.path, vendor, fmt.Sprintf("%d", t.Z)), 0o600)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Clean(filepath.Join(c.path, vendor, fmt.Sprintf("%d", t.Z), fmt.Sprintf("%d_%d.jpeg", t.X, t.Y))))
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(t.Image)
	return err
}
