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

	"math/rand"
)

var (
	angle                    = 0.0
	satelliteAngle           = 0.0
	drawingDistanceLine      = false
	drawingWeather           = false
	drawingPositionEstimates = false
	drawingTip               = false
	drawingTipMessage        = false
	drawingUserSelectionWin  = false
	tipMessage               string
	tipMessageAlpha          uint8 = 255
	tipMessageSize           *basicfont.Face
	currentTipMessageByte    = 0
	tipMaxScaleX             = 0.0
	tipCurrScaleX            = 0.0
	tipMaxScaleY             = 0.0
	tipCurrScaleY            = 0.0
	helpMessage              string
	helpMessageCount         int
	helpMessageAlpha         uint8 = 255
	helpMessageSize          *basicfont.Face
	displayMessageAlpha      uint8 = 255
	displayMessage           string
	displayMessageCount      int
	displayMessageSize       *basicfont.Face
	standardFont             = basicfont.Face7x13
	firstRun                 = true
	elevations               = []float64{230, 1000, 3000, 6800}
	currElevation            = 0
)

func drawBackground(win *pixelgl.Window) {
	background.sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
}

func drawLoadingScreen(win *pixelgl.Window) {
	loadScreen.sprite.Draw(win, loadScreen.mat)
	win.Update()
	load.Wait()
}

func drawOkButton(win *pixelgl.Window, scale float64) {
	okbutton.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scale).Moved(pixel.V(okbutton.posX, okbutton.posY)))
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
		if handleOKButtonClicked(win.JustPressed(pixelgl.MouseButtonLeft), win.MousePosition()) {
			drawOkButtonPressed(win)
			win.Update()
			break
		}
		for _, l := range loadSats {
			drawRandom(win, l, 2)
		}
		drawOkButton(win, 1)
		win.Update()
	}
}

func drawTip(win *pixelgl.Window) {
	if !drawingTip {
		return
	}
	okbutton.posY = 284
	if tipCurrScaleX < tipMaxScaleX {
		tipCurrScaleX += 5
	}
	if tipCurrScaleY < tipMaxScaleY {
		tipCurrScaleY += 5
	}
	scale := tipCurrScaleX / tipMaxScaleX
	tipmessage.sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, scale).Moved(win.Bounds().Center()))
	if scale > 0.99 {
		drawOkButton(win, scale)
		drawingTipMessage = true
	}
	if handleOKButtonClicked(win.JustPressed(pixelgl.MouseButtonLeft), win.MousePosition()) {
		drawOkButtonPressed(win)
		drawingTipMessage = false
		drawingTip = false
		tipCurrScaleX = 0.0
		tipCurrScaleY = 0.0
		if firstRun {
			firstRun = false
		}
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
		if currScale == SCALE_KM {
			mat = mat.Scaled(pixel.ZV, 0.5)
		}
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
	if drawingAnimation {
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

func drawDistanceLine(win *pixelgl.Window, p, q *object) {
	if !drawingDistanceLine {
		return
	}
	getDistanceLine(pixel.V(p.loc.X, p.loc.Y), pixel.V(q.loc.X, q.loc.Y)).Draw(win)
}

func drawDistanceLineLengthWithError(win *pixelgl.Window, p, q *object) {
	angleLen := 0.0
	if ceZenithPathDelay != 0 {
		pp := object{loc: personP.loc}
		qq := object{loc: personQ.loc}
		pp.loc.X -= ceZenithPathDelay / 10
		pp.loc.Y -= ceZenithPathDelay / 10
		qq.loc.X += ceZenithPathDelay / 10
		qq.loc.Y += ceZenithPathDelay / 10
		angleLen = ceTrueDistance(&pp, &qq)
		basicAtlas := text.NewAtlas(standardFont, text.ASCII)
		basicTxt := text.New(pixel.V(400, 10), basicAtlas)
		fmt.Fprintln(basicTxt, fmt.Sprintf("P->Q With Err: %.2f %s", angleLen, scaleNames[currScale]))
		basicTxt.Draw(win, pixel.IM)
		return
	}
	basicAtlas := text.NewAtlas(standardFont, text.ASCII)
	basicTxt := text.New(pixel.V(400, 10), basicAtlas)
	fmt.Fprintln(basicTxt, fmt.Sprintf("P->Q With Err: NK"))
	basicTxt.Draw(win, pixel.IM)
}

func drawDistanceLineLength(win *pixelgl.Window, p, q *object) {
	angleLen := ceTrueDistance(p, q)
	basicAtlas := text.NewAtlas(standardFont, text.ASCII)
	basicTxt := text.New(pixel.V(400, 30), basicAtlas)
	fmt.Fprintln(basicTxt, fmt.Sprintf("P->Q Distance: %.2f %s", angleLen, scaleNames[currScale]))
	basicTxt.Draw(win, pixel.IM)
}

func drawPerson(win *pixelgl.Window, p *object) {
	if currScale == SCALE_M {
		p.sprite.Draw(win, p.mat.Moved(p.loc))
		walkSpeed = 3.0
	} else if currScale == SCALE_KM {
		p.sprite.Draw(win, p.mat.Scaled(pixel.ZV, 0.5).Moved(p.loc))
		walkSpeed = 1.0
	}
}

func drawMovingPerson(win *pixelgl.Window, direction, flip int, p *object) {
	if currScale == SCALE_M {
		walkMap[direction][flip].sprite.Draw(win, p.mat.Moved(p.loc))
		walkSpeed = 3.0
	} else if currScale == SCALE_KM {
		walkMap[direction][flip].sprite.Draw(win, p.mat.Scaled(pixel.ZV, 0.5).Moved(p.loc))
		walkSpeed = 1.0
	}
}

func drawPositionEstimates(win *pixelgl.Window) {
	if !drawingPositionEstimates {
		return
	}
	if blink%4 != 0 {
		return
	}
	insertPositionEstimates()
	prand := &pestimates[rand.Int()%len(pestimates)]
	qrand := &qestimates[rand.Int()%len(qestimates)]
	getDistanceLine(pixel.V(prand.loc.X, prand.loc.Y), pixel.V(qrand.loc.X, qrand.loc.Y)).Draw(win)
	for _, p := range pestimates {
		p.sprite.Draw(win, p.mat)
	}
	for _, q := range qestimates {
		q.sprite.Draw(win, q.mat)
	}
}

func drawWeather(win *pixelgl.Window) {
	if !drawingWeather {
		return
	}
	wea := weatherSprites[currWeather]
	robj := weatherObjects[currWeather]
	for i, rdrop := range wea {
		robj[i].batch.Clear()
		rdrop.Draw(robj[i].batch, pixel.IM.Moved(pixel.V(robj[i].posX, robj[i].posY)))
		robj[i].batch.Draw(win)
		if robj[i].posY < maxY/3 {
			robj[i].posY = maxY
		} else {
			switch currWeather {
			case WEATHER_RAIN:
				robj[i].posY -= 10
			case WEATHER_ASH:
				robj[i].posY -= 5
			case WEATHER_DRY:
				robj[i].posY -= 2
			case WEATHER_SAND:
				robj[i].posY -= 30
			}
		}
	}
}

func drawOnTipMessage(win *pixelgl.Window, c int) {
	if !drawingTipMessage {
		return
	}
	basicAtlas := text.NewAtlas(tipMessageSize, text.ASCII)
	basicTxt := text.New(pixel.V(365, 465), basicAtlas)
	basicTxt.Color = color.RGBA{0, 0, 0, tipMessageAlpha}
	fmt.Fprintln(basicTxt, tipMessage[:c])
	basicTxt.Draw(win, pixel.IM)
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
	if helpMessageCount > 0 {
		basicAtlas := text.NewAtlas(helpMessageSize, text.ASCII)
		basicTxt := text.New(pixel.V(10, maxY-100), basicAtlas)
		basicTxt.Color = color.RGBA{0, 0, 0, helpMessageAlpha}
		fmt.Fprintln(basicTxt, helpMessage)
		basicTxt.Draw(win, pixel.IM)
		helpMessageAlpha = uint8(helpMessageAlpha - uint8(int(helpMessageAlpha)/helpMessageCount))
		helpMessageCount--
	}
}

func drawUserSelectionWindow(win *pixelgl.Window) {
	if !drawingUserSelectionWin {
		return
	}
	if drawingTipMessage && !drawingUserSelectionWin {
		newMessage("Cannot select options while tip message is displayed!", 100, standardFont)
		drawingUserSelectionWin = false
		return
	}
	currentAnimation = currentUserSelect
	tipmessage.sprite.Draw(win, pixel.IM.Moved(pixel.V(maxX/2, maxY/2)))
	newTipMessage("Select a clock drift in nanoseconds", standardFont)
	drawingTipMessage = true
	selection := handleUserSelectInput(win)
	if selection != 0 {
		drawingTipMessage = false
		stopAnimation()
		gpsClockDrift = selection
		drawingUserSelectionWin = false
		newMessage(fmt.Sprintf("GPS Clock drift set to %.2f nanoseconds", selection), 100, standardFont)
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
