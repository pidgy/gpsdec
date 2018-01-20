package gpsdec

import (
	"github.com/faiface/pixel"
)

const (
	noCollision = iota
	invalidClick
	buildingCollision
	buttonBuildingCollision
	buttonWeatherCollision
	buttonGPSCollision
	buttonScaleCollision
	buttonClearCollision
	buttonPerson1Collision
	buttonPerson2Collision
)

var (
	currSatelliteError = 1
)

func whereClick(loc pixel.Vec) int {
	if loc.Y > maxY-200 && loc.X < 800 {
		return invalidClick
	}
	if loc.Y > 200 {
		if collisionDetected(loc, buildings) {
			return buildingCollision
		}
	}
	if loc.X > 920 {
		if loc.Y > buttons[4].posY-buttons[4].frame.H()/2 && loc.Y < buttons[4].posY+buttons[4].frame.H()/2 {
			return buttonPerson1Collision
		}
		if loc.Y > buttons[5].posY-buttons[5].frame.H()/2 && loc.Y < buttons[5].posY+buttons[5].frame.H()/2 {
			return buttonPerson2Collision
		}
	}

	if loc.X < buttons[0].posX+buttons[0].frame.W() && loc.X > buttons[0].posX-buttons[0].frame.W()/2 {
		if loc.Y < buttons[0].posY+buttons[0].frame.H() {
			return buttonBuildingCollision
		}
	}
	if loc.X < buttons[1].posX+buttons[1].frame.W() && loc.X > buttons[1].posX-buttons[1].frame.W()/2 {
		if loc.Y < buttons[1].posY+buttons[1].frame.H() {
			return buttonWeatherCollision
		}
	}
	if loc.X < buttons[2].posX+buttons[2].frame.W() && loc.X > buttons[2].posX-buttons[2].frame.W()/2 {
		if loc.Y < buttons[2].posY+buttons[2].frame.H() {
			return buttonGPSCollision
		}
	}
	if loc.X < buttons[3].posX+buttons[3].frame.W() && loc.X > buttons[3].posX-buttons[3].frame.W()/2 {
		if loc.Y < buttons[3].posY+buttons[3].frame.H() {
			return buttonClearCollision
		}
	}
	if loc.X < buttons[6].posX+buttons[6].frame.W() && loc.X > buttons[6].posX-buttons[6].frame.W()/2 {
		if loc.Y < buttons[6].posY+buttons[6].frame.H() {
			return buttonScaleCollision
		}
	}

	return noCollision
}

func collisionDetected(v pixel.Vec, objects []object) bool {
	for _, o := range objects {
		if v.X >= o.posX && v.X <= o.frame.W()+o.posX {
			if v.Y >= o.posY && v.Y <= o.frame.H()+o.posY {
				return true
			}
		}
	}
	return false
}

func satelliteError() string {
	returnString := ""
	switch currSatelliteError {
	case 1:
		numSatellites--
		returnString = "Removing Satellite"
	case 2:
		numSatellites--
		returnString = "Removing Satellite"
	case 3:
		numSatellites--
		returnString = "Removing Satellite"
	case 4:
		numSatellites = 3
		returnString = "Readding Satellites"
	case 5:
		// TODO add GPS clock drift
		numSatellites = 3
		returnString = "Adding GPS clock drift"
	}
	if currSatelliteError == 5 {
		currSatelliteError = 1
	} else {
		currSatelliteError++
	}
	return returnString
}
