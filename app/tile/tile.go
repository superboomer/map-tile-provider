package tile

type Tile struct {
	X     int
	Y     int
	Z     int
	Image []byte
}

func (t Tile) GetNearby(side int) []Tile {
	var tiles []Tile

	// Calculate the starting point for the square
	startX := t.X - side/2
	startY := t.Y - side/2

	// Generate tiles for the square
	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			tiles = append(tiles, Tile{
				X: startX + i,
				Y: startY + j,
				Z: t.Z,
			})
		}
	}

	return tiles
}
