package gpsdec

import (
	"fmt"
	_ "image/png" // enforce png image use only
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	second          = time.Tick(time.Second)
	millisecond     = time.Tick(time.Millisecond * 100)
	framespersecond int
	last            = time.Now()
	debugMode       = false
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "GPS Error Simulations",
		Bounds: pixel.R(0, 0, maxX, maxY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	if !debugMode {
		drawLoadingScreen(win)
		drawControlScreen(win)
	}

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
		drawDistanceLine(win)
		drawDistanceLineLength(win)
		drawRectangle(win, object{})
		drawStaticObject(win)
		win.Update()

		// FPS Update
		framespersecond++
		select {
		case <-millisecond:
			flip = (flip + 1) % 2
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, framespersecond))
			framespersecond = 0
		default:
		}
	}
}

// Run this from your program to start simulation
func Run(debug bool) {
	debugMode = debug
	pixelgl.Run(run)
}
