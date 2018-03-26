package gpsdec

import (
	"math"

	"github.com/faiface/pixel"
)

const (
	SCALE_M = iota
	SCALE_KM
)

var (
	scaleNames [2]string
	currScale  int
)

func init() {
	scaleNames[SCALE_M] = "M"
	scaleNames[SCALE_KM] = "Km"
	currScale = SCALE_M
}

func distance(p, q pixel.Vec) (float64, float64) {
	return (personP.loc.X - personQ.loc.X), (personP.loc.Y - personQ.loc.Y)
}

func distanceAngle(x, y float64, p, q object) float64 {
	tan := math.Tan(y / x)
	arctan := math.Atan(tan)
	return arctan * (math.Pi / 180)
}

func distanceAngleLength(x, y float64) float64 {
	var pyth float64
	pyth = x*x + y*y
	pyth = math.Sqrt(pyth)
	return pyth
}
