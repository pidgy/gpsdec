package gpsdec

import (
	"fmt"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	angle                     = 0.0
	satelliteAngle            = 0.0
	drawingDistanceLine       = false
	drawingRain               = false
	displayMessageAlpha uint8 = 255
	displayMessage      string
	displayMessageCount int
	displayMessageSize  *basicfont.Face
	standardFont        = basicfont.Face7x13
)

func drawBackground(win *pixelgl.Window) {
	background.sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
}

func drawLoadingScreen(win *pixelgl.Window) {
	loadScreen.sprite.Draw(win, loadScreen.mat)
	win.Update()
	load.Wait()
}

func drawOkButton(win *pixelgl.Window) {
	okbutton.sprite.Draw(win, pixel.IM.Moved(pixel.V(okbutton.posX, okbutton.posY)))
}

func drawOkButtonPressed(win *pixelgl.Window) {
	pixel.NewSprite(okbutton.pressedpic, okbutton.pressedpic.Bounds()).Draw(win, pixel.IM.Moved(pixel.V(okbutton.posX, okbutton.posY)))
	okbutton.sprite.Draw(win, pixel.IM.Moved(pixel.V(okbutton.posX, okbutton.posY)))
}

func drawControlScreen(win *pixelgl.Window) {
	loadSats := loadExtraSprites(2, scrambleXPositions(satellites))
	for {
		controlScreen.sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		if win.Closed() {
			win.SetClosed(true)
			break
		}
		if handleLoadingScreenOk(win.JustPressed(pixelgl.MouseButtonLeft), win.MousePosition()) {
			drawOkButtonPressed(win)
			win.Update()
			break
		}
		for _, l := range loadSats {
			drawRandom(win, l, 2)
		}
		drawOkButton(win)
		win.Update()
	}
}

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

func drawButtons(win *pixelgl.Window) {
	for i := range buttons {
		if buttons[i].drawcount == 0 {
			button := pixel.NewSprite(buttons[i].pic, buttons[i].frame)
			button.Draw(win, buttons[i].mat)
		} else {
			button := pixel.NewSprite(buttons[i].pressedpic, buttons[i].frame)
			button.Draw(win, buttons[i].mat)
			buttons[i].drawcount--
		}
	}
	if currPerson == p {
		button := pixel.NewSprite(buttons[4].pressedpic, buttons[4].frame)
		button.Draw(win, buttons[4].mat)
	} else {
		button := pixel.NewSprite(buttons[5].pressedpic, buttons[5].frame)
		button.Draw(win, buttons[5].mat)
	}
}

// Draw a rectangle around an object
func drawRectangle(win *pixelgl.Window, obj object) *imdraw.IMDraw {
	imd := imdraw.New(nil)
	imd.Color = colornames.Honeydew
	imd.EndShape = imdraw.SharpEndShape
	imd.Push(pixel.V(obj.loc.X-obj.frame.W(), obj.loc.Y-obj.frame.H()), pixel.V(obj.loc.X+obj.frame.W(), obj.loc.Y+obj.frame.H()))
	imd.EndShape = imdraw.SharpEndShape
	imd.Rectangle(1)
	imd.Draw(win)
	return imd
}

func drawNewBuilding(pos pixel.Vec) {
	if int(pos.Y) < 100 {
		staticobject = newBuilding(pixel.V(maxSpriteX/2, maxSpriteY/2), 0)
	} else {
		staticobject = newBuilding(pos, 0)
	}
	drawStatic = true
	startAnimation(currentBuilding)
}

func drawStaticObject(win *pixelgl.Window) {
	if drawStatic {
		staticobject.sprite.Draw(win, pixel.IM.Moved(pixel.V(staticobject.loc.X, staticobject.loc.Y)))
		drawRectangle(win, staticobject)
	}
	if drawAnimation {
		for _, o := range getCurrentAnimation() {
			atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
			txt := text.New(pixel.V(o.posX-o.descalphaX, o.posY+o.descalphaY), atlas)
			fmt.Fprintln(txt, o.desc)
			txt.Draw(win, pixel.IM)
			o.sprite.Draw(win, pixel.IM.Moved(pixel.V(o.posX, o.posY)))
		}
	}
}

func drawStaticObjects(win *pixelgl.Window) {
	staticBatch.Clear()
	for i := range staticobjects {
		staticobjects[i].sprite.Draw(staticBatch, pixel.IM.Moved(pixel.V(staticobjects[i].loc.X, staticobjects[i].loc.Y)))
	}
	staticBatch.Draw(win)
}

func drawDistanceLine(win *pixelgl.Window) {
	if !drawingDistanceLine {
		return
	}
	imd := imdraw.New(nil)
	imd.Color = colornames.Red
	imd.EndShape = imdraw.RoundEndShape
	imd.Push(pixel.V(personP.loc.X, personP.loc.Y), pixel.V(personQ.loc.X, personQ.loc.Y))
	imd.EndShape = imdraw.SharpEndShape
	imd.Line(1)
	imd.Draw(win)
}

func drawDistanceLineLength(win *pixelgl.Window) {
	x, y := distance(personP.loc, personQ.loc)
	angleLen := distanceAngleLength(x, y)
	basicAtlas := text.NewAtlas(standardFont, text.ASCII)
	basicTxt := text.New(pixel.V(400, 20), basicAtlas)
	fmt.Fprintln(basicTxt, fmt.Sprintf("P->Q Distance: %.2f %s", angleLen, scaleNames[currScale]))
	basicTxt.Draw(win, pixel.IM)
}

func drawPerson(win *pixelgl.Window, p *object) {
	p.sprite.Draw(win, p.mat.Moved(p.loc))
}

func drawRain(win *pixelgl.Window) {
	if drawingRain {
		for i, rdrop := range rainSprites {
			rain[i].batch.Clear()
			rdrop.Draw(rain[i].batch, pixel.IM.Moved(pixel.V(rain[i].posX, rain[i].posY)))
			rain[i].batch.Draw(win)

			if rain[i].posY < maxY/3 {
				rain[i].posY = maxY
			} else {
				rain[i].posY -= 10
			}
		}
	}
}

func newMessage(m string, c int, s *basicfont.Face) {
	displayMessage = m
	displayMessageCount = c
	displayMessageSize = s
	displayMessageAlpha = 255
}

func drawMessage(win *pixelgl.Window) {
	if displayMessageCount > 0 {
		basicAtlas := text.NewAtlas(displayMessageSize, text.ASCII)
		basicTxt := text.New(pixel.V(10, maxY-70), basicAtlas)
		basicTxt.Color = color.RGBA{0, 0, 0, displayMessageAlpha}
		fmt.Fprintln(basicTxt, displayMessage)
		basicTxt.Draw(win, pixel.IM)
		displayMessageAlpha = uint8(displayMessageAlpha - uint8(int(displayMessageAlpha)/displayMessageCount))
		displayMessageCount--
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
