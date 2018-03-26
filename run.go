package gpsdec

import (
	_ "image/png" // enforce png image use only
	"time"

	"github.com/faiface/pixel/pixelgl"
)

var (
	second          = time.Tick(time.Second)
	millisecond     = time.Tick(time.Millisecond * 100)
	framespersecond int
	last            = time.Now()
)

// Run this from your program to start simulation
func Run() {
	pixelgl.Run(maincfg.run)
}
