package tile

import (
	"math"
)

// Elips contains eccentrcity for calculating
type Elips struct {
	Eccentricity float64
}

var (
	// ElipsWGS84 mercator WGS84 Eccentricity
	ElipsWGS84 = Elips{Eccentricity: 0.08181919084262157}
	// ElipsSpherical mercator for spherical Eccentricity
	ElipsSpherical = Elips{Eccentricity: 0}
)

// ConvertToTile convert latitude and longtitude to XYZ tile for specified mercator projection
func ConvertToTile(lat, long, zoom float64, proj *Elips) (x, y int) {
	rho := math.Pow(2, zoom+8) / 2
	beta := lat * math.Pi / 180

	phi := (1 - proj.Eccentricity*math.Sin(beta)) / (1 + proj.Eccentricity*math.Sin(beta))
	theta := math.Tan(math.Pi/4+beta/2) * math.Pow(phi, proj.Eccentricity/2)

	xP := rho * (1 + long/180)
	yP := rho * (1 - math.Log(theta)/math.Pi)

	x = int(math.Floor(xP / 256))
	y = int(math.Floor(yP / 256))

	return
}
