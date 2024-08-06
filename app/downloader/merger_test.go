package downloader

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/superboomer/map-tile-provider/app/tile"
)

func createTestImage(c color.Color) []byte {
	width, height := 100, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, c)
		}
	}
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestMerge_Success(t *testing.T) {
	side := 3
	centerTile := tile.Tile{X: 1, Y: 1}

	tiles := []tile.Tile{
		{X: 0, Y: 0, Image: createTestImage(color.RGBA{255, 0, 0, 255})},     // Red
		{X: 1, Y: 0, Image: createTestImage(color.RGBA{0, 255, 0, 255})},     // Green
		{X: 2, Y: 0, Image: createTestImage(color.RGBA{0, 0, 255, 255})},     // Blue
		{X: 0, Y: 1, Image: createTestImage(color.RGBA{255, 255, 0, 255})},   // Yellow
		{X: 1, Y: 1, Image: createTestImage(color.RGBA{255, 0, 255, 255})},   // Magenta (center)
		{X: 1, Y: 2, Image: createTestImage(color.RGBA{128, 128, 128, 255})}, // Gray
		{X: 2, Y: 1, Image: createTestImage(color.RGBA{0, 255, 255, 255})},   // Cyan
	}

	downloader := NewMapDownloader(http.DefaultClient)

	resultBytes, err := downloader.Merge(side, centerTile, tiles...)
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	resultImg, _, err := image.Decode(bytes.NewReader(resultBytes))
	if err != nil {
		t.Fatalf("Failed to decode result image: %v", err)
	}

	expectedWidth := side * 100
	expectedHeight := side * 100
	if resultImg.Bounds().Dx() != expectedWidth || resultImg.Bounds().Dy() != expectedHeight {
		t.Errorf("Expected size (%d,%d), got (%d,%d)", expectedWidth, expectedHeight,
			resultImg.Bounds().Dx(), resultImg.Bounds().Dy())
	}
}

func TestMerge_FailInvalidTile(t *testing.T) {
	side := 3
	centerTile := tile.Tile{X: 1, Y: 1}

	tiles := []tile.Tile{
		{X: 0, Y: 0, Image: createTestImage(color.RGBA{255, 0, 0, 255})},
		{X: 1, Y: 0, Image: nil},
	}

	downloader := NewMapDownloader(http.DefaultClient)

	_, err := downloader.Merge(side, centerTile, tiles...)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "center tile is not exist") // because function return only one hardcoded error message :(
}

func TestMerge_ErrorDecoding(t *testing.T) {
	side := 3
	centerTile := tile.Tile{X: 1, Y: 1}

	tiles := []tile.Tile{
		{X: 1, Y: 1, Image: []byte{0x00}}, // Invalid image data
	}

	downloader := NewMapDownloader(http.DefaultClient)

	resultBytes, err := downloader.Merge(side, centerTile, tiles...)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
	if resultBytes != nil {
		t.Fatalf("Expected nil result bytes but got some data")
	}
}

func TestMerge_ErrorEncoding(t *testing.T) {
	side := 3
	centerTile := tile.Tile{X: 1, Y: 1}

	tiles := []tile.Tile{
		{X: 0, Y: 0, Image: createTestImage(color.RGBA{255, 0, 0, 255})},
	}

	downloader := NewMapDownloader(http.DefaultClient)

	resultBytes, err := downloader.Merge(side, centerTile, tiles...)
	if err == nil {
		t.Fatal("Expected an error but got none")
	}
	if resultBytes != nil {
		t.Fatalf("Expected nil result bytes but got some data")
	}
}
