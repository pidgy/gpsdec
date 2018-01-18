package gpssimulation

import "github.com/faiface/pixel"

type message struct {
	msg       string
	drawcount int
}

type object struct {
	posX            float64
	posY            float64
	angle           float64
	direction       bool
	frame           pixel.Rect
	pic             pixel.Picture
	pressedpic      pixel.Picture
	height          float64
	width           float64
	mat             pixel.Matrix
	loc             pixel.Vec
	filename        string
	pressedfilename string
	drawcount       int
}
