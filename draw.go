package gpsdec

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

var (
	angle          = 0.0
	satelliteAngle = 0.0
)

func drawUnits() []*text.Text {
	j := 0
	maxYInt := int(maxY)
	var texts []*text.Text
	for i := 100; i < maxYInt; i += int(maxYInt / 10) {
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(10, float64(i)), basicAtlas)
		fmt.Fprintln(basicTxt, fmt.Sprintf("-----= %d", i))
		texts = append(texts, basicTxt)
		j++
	}
	return texts
}

func drawMessage(m string, c int, s *basicfont.Face) {
	displayMessage = m
	displayMessageCount = c
	displayMessageSize = s
}

func drawSatellites(win *pixelgl.Window) {
	dt := time.Since(last).Seconds()
	last = time.Now()
	satelliteAngle += 2 * dt
	speed := 1.0

	for i := 0; i < numSatellites; i++ {
		mat := pixel.IM
		mat = mat.Scaled(pixel.ZV, 0.1).Rotated(pixel.ZV, satelliteAngle)
		if satellites[i].posX < speed || satellites[i].posX > maxX-speed {
			satellites[i].directionX = !satellites[i].directionX
		}
		if satellites[i].directionX == left {
			satellites[i].posX -= speed
		} else {
			satellites[i].posX += speed
		}
		mat = mat.Moved(pixel.V(satellites[i].posX, 740))
		satellites[i].sprite.Draw(win, mat)
	}
}

func drawRandom(win *pixelgl.Window, o []object, speed float64) {
	dt := time.Since(last).Seconds()
	last = time.Now()
	angle += 2 * dt
	for i := 0; i < len(o); i++ {
		mat := pixel.IM
		mat = mat.Scaled(pixel.ZV, 0.1).Rotated(pixel.ZV, angle)
		if o[i].posX < speed || o[i].posX > maxX-speed {
			o[i].directionX = !o[i].directionX
		}
		if o[i].posY < speed || o[i].posY > maxY-speed {
			o[i].directionY = !o[i].directionY
		}
		if o[i].directionX == left {
			o[i].posX -= speed
		} else {
			o[i].posX += speed
		}
		if o[i].directionY == up {
			o[i].posY += speed
		} else {
			o[i].posY -= speed
		}
		mat = mat.Moved(pixel.V(o[i].posX, o[i].posY))
		o[i].sprite.Draw(win, mat)
	}
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
