package gpsdec

import (
	"math"
	"math/rand"
	"os"

	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"

	owm "github.com/briandowns/openweathermap"
)

const (
	ownAPIKey = "2301adeee7376aae473f2e88288708f6"
)

var (
	tipMessages    = []string{"Welcome to gpsdec! \n\n\nAs you can see you have 2 characters below\n\n\nWe will call them P and Q!"}
	currTipMessage = 0
)

func init() {
	os.Setenv("OWM_API_KEY", ownAPIKey)
}

func remove(s []message, i int) []message {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func celToKel(C float64) float64 {
	return C + 273.0
}

func kelToCel(K float64) float64 {
	return K - 273.0
}

func kelToFah(K float64) float64 {
	return K*9/5 - 459.67
}

func fahToCel(F float64) float64 {
	return (F - 32) * 5 / 9
}

func tempChangeFromElevation() float64 {
	if currWeather == WEATHER_RAIN {
		return (3.3 * (elevations[currElevation] - elevations[0])) / 1000
	}
	return (5.4 * (elevations[currElevation] - elevations[0])) / 1000
}

func waterVaporPressureAntoine(T float64) float64 {
	A := 8.07131
	B := 1730.63
	C := 233.426
	if T > 99 {
		A = 8.07131
		B = 1730.63
		C = 233.426
	}
	Ptorr := math.Pow(10, (A - (B / (C + T))))
	return Ptorr
}

func dryRefractivity(T, P float64) float64 {
	k1 := 77.6
	return k1 * (P / T)

}

func wetRefractivity(T, P float64) float64 {
	k2 := 64.8
	k3 := 3.776 * (math.Pow(10, 5))
	return (k2 * (P / T)) + (k3 * (P / math.Pow(T, 2)))
}

func ptorrToPmb(p float64) float64 {
	return 1.33322 * p
}

func getAirPressure() float64 {
	ap := 1015.0
	switch elevations[currElevation] {
	case 6000:
		ap = 980.0
	case 3000:
		ap = 1005.0
	case 1000:
		ap = 1025.0
	case 230:
		ap = 1050.0
	}
	return ap
}

func scramblePositions(o []object) []object {
	var ret []object
	for _, obj := range o {
		ret = append(ret, obj)
	}
	for i := 0; i < len(ret); i++ {
		ret[i].posX = float64(rand.Int() % int(maxX))
		ret[i].posY = float64(rand.Int() % int(maxY))
	}
	return ret
}

func scrambleXPositions(o []object) []object {
	var ret []object
	for _, obj := range o {
		ret = append(ret, obj)
	}
	for i := 0; i < len(ret); i++ {
		ret[i].posX = float64(rand.Int() % int(maxX))
	}
	return ret
}

func scrambleYPositions(o []object) []object {
	var ret []object
	for _, obj := range o {
		ret = append(ret, obj)
	}
	for i := 0; i < len(ret); i++ {
		ret[i].posY = float64(rand.Int() % int(maxY))
	}
	return ret
}

func getWeather() *owm.CurrentWeatherData {
	w, err := owm.NewCurrent("K", "EN", ownAPIKey)
	if err != nil {
		m := owm.Main{Temp: 303.15}
		return &owm.CurrentWeatherData{Main: m}
	}
	w.CurrentByName("Ottawa")
	return w
}

func mmToKm(mm float64) float64 {
	return mm / 1000000
}

// Check if a vector lands in the frame of a manually created object
// Note: this function assumes the objects posX and posY fields exist
// and they lay in the center and not the origin
func vectorIntersectionWithObject(src pixel.Vec, dst *object) bool {
	if src.X > dst.posX-dst.frame.W()/2 && src.X < dst.posX+dst.frame.W()/2 {
		if src.Y > dst.posY-dst.frame.H()/2 && src.Y < dst.posY+dst.frame.H()/2 {
			return true
		}
	}
	return false
}

func getDistanceLine(pvec, qvec pixel.Vec) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = colornames.Red
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pvec, qvec)
	imd.EndShape = imdraw.SharpEndShape
	imd.Line(1)
	return imd
}

func newTipMessage(m string, s *basicfont.Face) {
	tipMessage = m
	tipMessageSize = s
	tipMessageAlpha = 255
}

func newMessage(m string, c int, s *basicfont.Face) {
	displayMessage = m
	displayMessageCount = c
	displayMessageSize = s
	displayMessageAlpha = 255
}

func newHelpMessage(m string, c int, s *basicfont.Face) {
	helpMessage = m
	helpMessageCount = c
	helpMessageSize = s
	helpMessageAlpha = 255
}
