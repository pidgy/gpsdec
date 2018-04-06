package gpsdec

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
)

const (
	gpsDistance = 21000.00
)

var (
	estimatesLoaded = false

	factorPropDelay     = 0
	factorGPSDrift      = 0
	factorEphemeris     = 0
	factorHardware      = 0
	factorMultipathPro  = 0
	factorSatelliteGeom = 0

	zenithDelays map[int]float64
	zenithScales map[int]float64

	pv, qv            pixel.Vec
	ceZenithPathDelay float64

	gpsClockDrift float64
)

func init() {
	zenithDelays = map[int]float64{}
	zenithDelays[WEATHER_DRY] = 290
	zenithDelays[WEATHER_RAIN] = 140
	zenithDelays[WEATHER_SAND] = 18
	zenithDelays[WEATHER_ASH] = 0.01

	zenithScales = map[int]float64{}
	zenithScales[WEATHER_DRY] = 8
	zenithScales[WEATHER_RAIN] = 2.7
	zenithScales[WEATHER_SAND] = 1
	zenithScales[WEATHER_ASH] = 4

}

func insertPositionEstimates() {
	pv = personP.loc
	qv = personQ.loc
	if ceZenithPathDelay != 0 {
		pv.X -= ceZenithPathDelay / 10
		qv.X += ceZenithPathDelay / 10
	}
	if currScale == SCALE_M {
		createEstimatePlot(30, pv, qv)
	} else if currScale == SCALE_KM {
		createEstimatePlot(5, pv, qv)
	}
}

func mapAtmosphericDelay(T float64) (float64, float64) {
	N := 0.0
	ZD := 0.0
	Ns := ""
	F := kelToFah(T)
	printFmtOutput(fmt.Sprintf("Current Temperature (F): %.2f", F))
	F -= tempChangeFromElevation()
	printFmtOutput(fmt.Sprintf("Temperature Adjustment from Elevation: %.2f Degrees.", F))
	C := fahToCel(F)
	printFmtOutput(fmt.Sprintf("Current Temperature (C): %.2f", C))
	switch currWeather {
	case WEATHER_RAIN:
		printFmtOutput("Heavy water vapor levels detected..")
		Ptorr := waterVaporPressureAntoine(C)
		printFmtOutput(fmt.Sprintf("Antoine Equation for Water Vapor Pressure (torr): %.2f", Ptorr))
		Pmb := ptorrToPmb(Ptorr)
		printFmtOutput(fmt.Sprintf("Antoine Equation for Water Vapor Pressure (mb): %.2f", Pmb))
		N = wetRefractivity(celToKel(C), Pmb)
		Ns = "Nvap"
		ZD = zenithDelays[currWeather]
	case WEATHER_ASH:
		printFmtOutput("Heavy distribution of ash detected..")
		ZD = zenithDelays[currWeather]
	case WEATHER_SAND:
		printFmtOutput("Heavy sandstorm detected..")
		ZD = zenithDelays[currWeather]
	case WEATHER_DRY:
		printFmtOutput("Hydrostatic air pressure detected..")
		k1 := 77.6
		P := getAirPressure()
		printFmtOutput(fmt.Sprintf("Current air pressure: %.2f", P))
		N = dryRefractivity(k1, P)
		Ns = "Ndry"
		ZD = zenithDelays[currWeather]
	}
	printFmtOutput(fmt.Sprintf("[Propagation Delay]: Refractivity: %s = %.5f", Ns, N))
	return N, ZD
}

func printDashLine() {
	println("|--------------------------------------------------------------|")
}

func ceEstimateDistance() {
	printDashLine()
	println("|                 gpsdec simulation output                     |")
	printDashLine()
	drawingDistanceLine = true
	weather := getWeather()
	printFmtOutput(fmt.Sprintf("True distance P->Q: %.2f %s", ceTrueDistance(&personP, &personQ), scaleNames[currScale]))
	printFmtOutput(fmt.Sprintf("GPS Distance: %.2f km", gpsDistance))
	printFmtOutput(fmt.Sprintf("GPS Distance given elevation: %.1f km", gpsDistance-elevations[currElevation]))
	printFmtOutput(fmt.Sprintf("Current Temperature (C): %.2f", kelToCel(weather.Main.Temp)))
	printFmtOutput(fmt.Sprintf("Current Time: %s", time.Now().UTC()))
	printDashLine()
	printFmtOutput("Estimating distance..")
	drawingPositionEstimates = true
	printFmtOutput(fmt.Sprintf("Current Elevation: %.1f ft.", elevations[currElevation]))
	printFmtOutput(fmt.Sprintf("Current Temperature (K): %.2f", weather.Main.Temp))
	printFmtOutput(fmt.Sprintf("Checking for measurement error factors.. "))

	if drawingWeather {
		printDashLine()
		printFmtOutput("Propagation Delay Detected!")
		printDashLine()
		_, ZD := mapAtmosphericDelay(weather.Main.Temp)
		ZS := zenithScales[currWeather]
		printFmtOutput(fmt.Sprintf("[Propagation Delay]: Zenith Path Delay: %.2f mm/km", ZD))
		printFmtOutput(fmt.Sprintf("[Propagation Delay]: Scale Height of Constituent : %.2f km", zenithScales[currWeather]))
		printFmtOutput(fmt.Sprintf("[Propagation Delay]: Total Delay: %.2f mm", ZD*ZS))
		ceZenithPathDelay += ZD * ZS
	}
	if gpsClockDrift != 0 {
		printDashLine()
		printFmtOutput("GPS Clock Drift Detected!")
		printDashLine()
		printFmtOutput(fmt.Sprintf("[GPS Clock Drift]: Drift Time: %.2f ns", gpsClockDrift))
		printFmtOutput(fmt.Sprintf("[GPS Clock Drift]: Range Error %.2f metres", 0.3*gpsClockDrift))
		ceZenithPathDelay += gpsClockDrift * 10
	}
	printDashLine()
}

func ceTrueDistance(p, q *object) float64 {
	x, y := distance(p.loc, q.loc)
	angleLen := distanceAngleLength(x, y)
	return angleLen
}

func printFmtOutput(s string) {
	lenBetweenBars := len("|--------------------------------------------------------------|") - 2
	lenStr := len(s)
	diff := lenBetweenBars - lenStr
	if diff < 0 {
		fmt.Printf("%s\n", s)
		return
	}
	outStr := "| "
	outStr += s
	i := 0
	for i < diff-1 {
		outStr += " "
		i++
	}
	outStr += "|"
	println(outStr)
}

func createEstimatePlot(move int, pv, qv pixel.Vec) {
	distance := move
	modY := -1
	modX := -1
	pestimate.loc = pv
	qestimate.loc = qv
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
