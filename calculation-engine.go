package gpsdec

import (
	"fmt"

	"github.com/faiface/pixel"
)

var (
	estimatesLoaded = false

	factorPropDelay     = 0
	factorGPSDrift      = 0
	factorEphemeris     = 0
	factorHardware      = 0
	factorMultipathPro  = 0
	factorSatelliteGeom = 0
)

func insertPositionEstimates() {
	if currScale == SCALE_M {
		createEstimatePlot(30)
	} else if currScale == SCALE_KM {
		createEstimatePlot(5)
	}
}

func mapAtmosphericDelay() {
	switch currWeather {
	case WEATHER_RAIN:
	case WEATHER_ASH:
	case WEATHER_SAND:
	case WEATHER_DRY:
	}
}

func estimateDistance() {
	println("|--------------------------------------------------------------|")
	println("|                 gpsdec simulation output                     |")
	println("|--------------------------------------------------------------|")
	fmt.Printf("| True distance P->Q: %.2f %s\n", trueDistance(&personP, &personQ), scaleNames[currScale])
	fmt.Printf("| Current Elevation: %.1f ft.\n", elevations[currElevation])
	fmt.Printf("| Current Temperature: %.1f ft.\n", elevations[currElevation])
	if drawingWeather {
		println("| Propagation Delay: ")
		mapAtmosphericDelay()
	}
	println("|--------------------------------------------------------------|")
}

func trueDistance(p, q *object) float64 {
	x, y := distance(p.loc, q.loc)
	angleLen := distanceAngleLength(x, y)
	return angleLen
}

func createEstimatePlot(move int) {
	distance := move
	modY := -1
	modX := -1
	pestimate.loc = personP.loc
	pestimate.mat = personP.mat
	qestimate.loc = personQ.loc
	qestimate.mat = personQ.mat
	pestimates = []object{pestimate}
	qestimates = []object{qestimate}
	for modX < 2 {
		for modY < 2 {
			i := 0
			distance = move
			for i < 2 {
				pm := pestimate
				qm := qestimate
				if modX != 0 {
					pm.loc.X += float64((distance + 10) * modX)
					qm.loc.X += float64((distance + 10) * modX)
				}
				if modY != 0 {
					pm.loc.Y += float64((distance + 10) * modY)
					qm.loc.Y += float64((distance + 10) * modY)
				}
				pm.mat = pixel.IM.Moved(pixel.V(pm.loc.X, pm.loc.Y))
				qm.mat = pixel.IM.Moved(pixel.V(qm.loc.X, qm.loc.Y))
				pestimates = append(pestimates, pm)
				qestimates = append(qestimates, qm)
				i++
				distance += 30
			}
			modY++
		}
		modY = -1
		modX++
	}
}
