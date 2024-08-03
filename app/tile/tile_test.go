package tile

import (
	"testing"
)

func TestGetNearby(t *testing.T) {
	tests := []struct {
		tile     Tile
		side     int
		expected []Tile
	}{
		{
			tile: Tile{X: 2, Y: 2, Z: 0},
			side: 3,
			expected: []Tile{
				{X: 1, Y: 1, Z: 0},
				{X: 1, Y: 2, Z: 0},
				{X: 1, Y: 3, Z: 0},
				{X: 2, Y: 1, Z: 0},
				{X: 2, Y: 2, Z: 0},
				{X: 2, Y: 3, Z: 0},
				{X: 3, Y: 1, Z: 0},
				{X: 3, Y: 2, Z: 0},
				{X: 3, Y: 3, Z: 0},
			},
		},
		{
			tile: Tile{X: 0, Y: 0, Z: 1},
			side: 2,
			expected: []Tile{
				{X: -1, Y: -1, Z: 1},
				{X: -1, Y: 0, Z: 1},
				{X: 0, Y: -1, Z: 1},
				{X: 0, Y: 0, Z: 1},
			},
		},
	}

	for _, test := range tests {
		result := test.tile.GetNearby(test.side)

		if len(result) != len(test.expected) {
			t.Errorf("For tile %v with side %d, expected %d tiles but got %d", test.tile, test.side, len(test.expected), len(result))
			continue
		}

		for i := range result {
			if !tilesEqual(result[i], test.expected[i]) {
				t.Errorf("For tile %v with side %d, expected tile %v but got %v", test.tile, test.side, test.expected[i], result[i])
			}
		}
	}
}

// Helper function to compare two tiles
func tilesEqual(a, b Tile) bool {
	return a.X == b.X && a.Y == b.Y && a.Z == b.Z
}
