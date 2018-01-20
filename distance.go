package gpsdec

import (
	"github.com/faiface/pixel"
)

func distance(p, q pixel.Vec) {
	xdiff := 0.0
	ydiff := 0.0

	if p.X > q.X {
		xdiff = p.X - q.X
	} else {
		xdiff = q.X - p.X
	}
	if p.Y > q.Y {
		ydiff = p.Y - q.Y
	} else {
		ydiff = q.Y - p.Y
	}

	println("X diff:", int(xdiff), "Y diff:", int(ydiff))
}
