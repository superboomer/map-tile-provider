package downloader

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"sort"

	"github.com/superboomer/maptile/app/tile"
)

type imageTileSlice []tile.Tile

func (d imageTileSlice) Len() int {
	return len(d)
}

func (d imageTileSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d imageTileSlice) Less(i, j int) bool {
	if d[i].Y == d[j].Y {
		return d[i].X < d[j].X
	}
	return d[i].Y < d[j].Y
}

// Merge combines multiple tiles into a single image.
func (m *MapDownloader) Merge(side int, centerTile tile.Tile, tiles ...tile.Tile) ([]byte, error) {
	coordsData := filterAndSortTiles(tiles, centerTile)
	if coordsData == nil {
		return nil, fmt.Errorf("center tile is not exist in tiles, x=%d, y=%d, z=%d", centerTile.X, centerTile.Y, centerTile.Z)
	}

	mergedImage, err := createMergedImage(side, centerTile, coordsData)
	if err != nil {
		return nil, fmt.Errorf("error occurred during image merging: %w", err)
	}

	return mergedImage, nil
}

// filterAndSortTiles filters out tiles to ensure the center tile exists and sorts them.
func filterAndSortTiles(tiles []tile.Tile, centerTile tile.Tile) *imageTileSlice {
	var coordsData = make(imageTileSlice, 0)
	centerTileOK := false

	for _, t := range tiles {
		if t.Image == nil {
			return nil // Early return if any tile has no image data
		}
		if t.X == centerTile.X && t.Y == centerTile.Y && t.Z == centerTile.Z {
			centerTileOK = true
		}
		coordsData = append(coordsData, t)
	}

	if !centerTileOK {
		return nil // Center tile not found among provided tiles
	}

	sort.Sort(coordsData)
	return &coordsData
}

// createMergedImage creates a merged image from sorted tiles around a center tile.
func createMergedImage(side int, centerTile tile.Tile, coordsData *imageTileSlice) ([]byte, error) {
	images := prepareImageGrid(side)
	for _, file := range *coordsData {
		img, _, err := image.Decode(bytes.NewReader(file.Image))
		if err != nil {
			return nil, fmt.Errorf("error occurred with decoding image: %w", err)
		}

		x := file.X - centerTile.X + (side / 2)
		y := file.Y - centerTile.Y + (side / 2)
		if x >= 0 && x < side && y >= 0 && y < side {
			images[x][y] = img
		}
	}

	return mergeImagesIntoResult(images, side)
}

// prepareImageGrid initializes a grid to hold images based on the side length.
func prepareImageGrid(side int) [][]image.Image {
	images := make([][]image.Image, side)
	for i := range images {
		images[i] = make([]image.Image, side)
	}
	return images
}

// mergeImagesIntoResult merges individual tile images into a single image.
func mergeImagesIntoResult(images [][]image.Image, side int) ([]byte, error) {
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
		return nil, fmt.Errorf("error occurred with encoding new image: %w", err)
	}

	return resultImage.Bytes(), nil
}
