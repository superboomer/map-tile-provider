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

func LoadCache(path string, alive time.Duration) (*Cache, error) {
	err := os.MkdirAll(path, 0600)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(filepath.Join(path, "index.db"), 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Cache{db: db, path: path, alive: alive}, nil
}

func (c *Cache) LoadFile(vendor string, tile *tile.Tile) ([]byte, error) {
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(vendor))

		if b == nil {
			return fmt.Errorf("cache is invalid")
		}

		v := b.Get([]byte(fmt.Sprintf("%d_%d_%d", tile.X, tile.Y, tile.Z)))
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

	img, err := os.ReadFile(filepath.Join(c.path, vendor, fmt.Sprintf("%d", tile.Z), fmt.Sprintf("%d_%d.jpeg", tile.X, tile.Y)))
	if err != nil {
		return nil, fmt.Errorf("can't read image from cache: %w", err)
	}

	return img, nil
}

func (c *Cache) SaveTile(vendor string, tile *tile.Tile) error {
	if tile.Image == nil {
		return fmt.Errorf("image not provided")
	}

	return c.db.Update(func(tx *bolt.Tx) error {
		err := c.saveImage(vendor, tile)
		if err != nil {
			return err
		}

		b, err := tx.CreateBucketIfNotExists([]byte(vendor))
		if err != nil {
			return err
		}
		return b.Put([]byte(fmt.Sprintf("%d_%d_%d", tile.X, tile.Y, tile.Z)), unixTimeEncode(time.Now()))
	})
}

func (c *Cache) saveImage(vendor string, tile *tile.Tile) error {
	err := os.MkdirAll(filepath.Join(c.path, vendor, fmt.Sprintf("%d", tile.Z)), 0600)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(c.path, vendor, fmt.Sprintf("%d", tile.Z), fmt.Sprintf("%d_%d.jpeg", tile.X, tile.Y)))
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(tile.Image)
	return err
}
