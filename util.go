package gpsdec

import (
	"math/rand"

	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/faiface/pixel"
)

var (
	tipMessages    = []string{"Welcome to gpsdec! \n\n\nAs you can see you have 2 characters below\n\n\nWe will call them P and Q!"}
	currTipMessage = 0
)

func remove(s []message, i int) []message {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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
