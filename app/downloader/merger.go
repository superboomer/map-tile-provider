package downloader

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"sort"

	"github.com/superboomer/map-tile-provider/app/tile"
)

type imageTileSlice []tile.Tile

// Len is part of sort.Interface.
func (d imageTileSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface
func (d imageTileSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less compare two value and return less
func (d imageTileSlice) Less(i, j int) bool {
	if d[i].Y == d[j].Y {
		return d[i].X < d[j].X
	}
	return d[i].Y < d[j].Y
}

// Merge combines multiple tiles into a single image
func Merge(side int, centerTile tile.Tile, tiles ...tile.Tile) ([]byte, error) {
	var coordsData = make(imageTileSlice, 0)

	var centerTileOK bool

	for _, t := range tiles {
		if t.Image == nil {
			return nil, fmt.Errorf("image empty on tile x=%d, y=%d, z=%d", t.X, t.Y, t.Z)
		}
		if t.X == centerTile.X && t.Y == centerTile.Y && t.Z == centerTile.Z {
			centerTileOK = true
		}

		coordsData = append(coordsData, t)
	}

	if !centerTileOK {
		return nil, fmt.Errorf("center tile is not exist in tiles, x=%d, y=%d, z=%d", centerTile.X, centerTile.Y, centerTile.Z)
	}

	sort.Sort(coordsData)

	images := make([][]image.Image, side)
	for i := range images {
		images[i] = make([]image.Image, side)
	}

	for _, file := range coordsData {
		img, _, err := image.Decode(bytes.NewReader(file.Image))
		if err != nil {
			return nil, fmt.Errorf("error occurred with decoding image err=%w", err)
		}

		x := file.X - centerTile.X + (side / 2)
		y := file.Y - centerTile.Y + (side / 2)

		if x >= 0 && x < side && y >= 0 && y < side {
			images[x][y] = img
		}
	}

	totalWidth := images[0][0].Bounds().Dx() * side
	totalHeight := images[0][0].Bounds().Dy() * side
	result := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))

	for x := 0; x < side; x++ {
		for y := 0; y < side; y++ {
			if images[x][y] != nil {
				img := images[x][y]
				r := image.Rect(
					x*img.Bounds().Dx(),
					y*img.Bounds().Dy(),
					(x+1)*img.Bounds().Dx(),
					(y+1)*img.Bounds().Dy(),
				)
				draw.Draw(result, r, img, img.Bounds().Min, draw.Over)
			}
		}
	}

	resultImage := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(resultImage, result, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, fmt.Errorf("error occurred with encoding new image err=%w", err)
	}

	return resultImage.Bytes(), nil
}
