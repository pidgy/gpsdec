package gpsdec

import "github.com/faiface/pixel"

var (
	estimatesLoaded = false
)

func insertPositionEstimates() {
	if currScale == SCALE_M {
		createEstimatePlot(30)
	} else if currScale == SCALE_KM {
		createEstimatePlot(5)
	}
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
