package tile

import (
	"testing"
)

func TestConvertToTile(t *testing.T) {
	tests := []struct {
		lat       float64
		long      float64
		zoom      float64
		proj      *Elips
		expectedX int
		expectedY int
	}{
		{
			lat:       0,
			long:      0,
			zoom:      0,
			proj:      &ElipsWGS84,
			expectedX: 0,
			expectedY: 0,
		},
		{
			lat:       45,
			long:      45,
			zoom:      1,
			proj:      &ElipsWGS84,
			expectedX: 1,
			expectedY: 0,
		},
		{
			lat:       -45,
			long:      -45,
			zoom:      2,
			proj:      &ElipsWGS84,
			expectedX: 1,
			expectedY: 2,
		},
		{
			lat:       0,
			long:      180,
			zoom:      4,
			proj:      &ElipsSpherical,
			expectedX: 16,
			expectedY: 8,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			x, y := ConvertToTile(test.lat, test.long, test.zoom, test.proj)
			if x != test.expectedX || y != test.expectedY {
				t.Errorf("ConvertToTile(%v, %v, %v, %v) = (%v, %v); expected (%v, %v)",
					test.lat, test.long, test.zoom, test.proj, x, y, test.expectedX, test.expectedY)
			}
		})
	}
}
