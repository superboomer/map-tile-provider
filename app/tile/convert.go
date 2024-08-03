package tile

import (
	"math"
)

type Elips struct {
	Eccentricity float64
}

var (
	ElipsWGS84     = Elips{Eccentricity: 0.08181919084262157} // WGS84 Eccentricity
	ElipsSpherical = Elips{Eccentricity: 0}                   // Spherical Eccentricity
)

func ConvertToTile(lat, long float64, zoom float64, proj *Elips) (int, int) {
	rho := math.Pow(2, zoom+8) / 2
	beta := lat * math.Pi / 180

	phi := (1 - proj.Eccentricity*math.Sin(beta)) / (1 + proj.Eccentricity*math.Sin(beta))
	theta := math.Tan(math.Pi/4+beta/2) * math.Pow(phi, proj.Eccentricity/2)

	xP := rho * (1 + long/180)
	yP := rho * (1 - math.Log(theta)/math.Pi)

	return int(math.Floor(xP / 256)), int(math.Floor(yP / 256))
}
