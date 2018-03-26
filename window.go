package gpsdec

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	title             = "GPS Distance Error Simulation"
	maintitle         = title + " - Main"
	simtitle          = title + " - Data"
	maxX      float64 = 1024
	maxY      float64 = 768
)

type gpsdecWindow struct {
	cfg          pixelgl.WindowConfig
	winnotclosed func(*pixelgl.Window)
	showcontrols bool
}

var (
	maincfg *gpsdecWindow
	simcfg  *gpsdecWindow
	blink   = 0
	flip    = 0
)

func initMainWindow() {
	sprite, err := loadPicture(spritedirectory + objectsdirectory + "satellite-pixel.png")
	if err != nil {
		panic(err)
	}
	maincfg = &gpsdecWindow{
		cfg: pixelgl.WindowConfig{
			Icon:   []pixel.Picture{sprite},
			Title:  maintitle,
			Bounds: pixel.R(0, 0, maxX, maxY),
			VSync:  true,
		},
		winnotclosed: mainNotClosed,
		showcontrols: true,
	}
}

func initSimWindow() {
	simcfg = &gpsdecWindow{
		cfg: pixelgl.WindowConfig{
			Title:  simtitle,
			Bounds: pixel.R(0, 0, maxX/2, maxY/2),
			VSync:  true,
		},
		winnotclosed: simNotClosed,
		showcontrols: false,
	}
}

func init() {
	initMainWindow()
	initSimWindow()
}

func (w *gpsdecWindow) run() {

	win, err := pixelgl.NewWindow(w.cfg)
	if err != nil {
		panic(err)
	}

	if w.showcontrols {
		drawLoadingScreen(win)
		drawControlScreen(win)
	}
	w.winnotclosed(win)
}

func mainNotClosed(win *pixelgl.Window) {
	for !win.Closed() {
		drawBackground(win)
		drawSatellites(win)
		handleCollision(win)
		handleMovementKeyPress(win)
		handleBuildingAdded(win)
		handleDistanceLineKey(win.JustReleased(pixelgl.Key4))
		handlePersonKeyPressed(win.JustReleased(pixelgl.KeyTab))
		drawStaticObjects(win)
		drawButtons(win)
		drawRain(win)
		drawMessage(win)
		drawDistanceLine(win, &personP, &personQ)
		drawDistanceLineLength(win, &personP, &personQ)
		drawRectangle(win, object{})
		drawStaticObject(win)
		drawPositionEstimates(win)
		handleMouseHover(win)
		drawTip(win)
		drawOnTipMessage(win, currentTipMessageByte)
		win.Update()

		if firstRun {
			newTipMessage(tipMessages[currTipMessage], standardFont)
			drawingTip = true
		}
		framespersecond++
		select {
		case <-millisecond:
			flip = (flip + 1) % 2
			blink++
			if currentTipMessageByte != len(tipMessage) {
				currentTipMessageByte++
			}
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", maintitle, framespersecond))
			framespersecond = 0
			//println("Mouse X->", int(win.MousePosition().X), "Y->", int(win.MousePosition().Y))
		default:
		}
	}
}

func simNotClosed(win *pixelgl.Window) {
	for !win.Closed() {
		win.Update()
	}
}
