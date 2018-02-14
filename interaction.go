package gpsdec

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	noCollision = iota
	invalidClick
	buildingCollision
	buttonBuildingCollision
	buttonWeatherCollision
	buttonSignalCollision
	buttonScaleCollision
	buttonLineCollision
	buttonClearCollision
	buttonPerson1Collision
	buttonPerson2Collision
	buttonControlsCollision
	buttonRunCollision
)

var (
	currSatelliteError = 1
	flip               = 0
)

func whereClick(loc pixel.Vec) int {
	if loc.Y > maxY-200 && loc.X < 800 {
		return invalidClick
	}
	if loc.Y > 200 {
		if collisionDetected(loc, staticobjects) {
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
	if loc.X < buttons[7].posX+buttons[7].frame.W() && loc.X > buttons[7].posX-buttons[7].frame.W()/2 {
		if loc.Y < buttons[7].posY+buttons[7].frame.H() {
			return buttonLineCollision
		}
	}
	if loc.X < buttons[8].posX+buttons[8].frame.W() && loc.X > buttons[8].posX-buttons[8].frame.W()/2 {
		if loc.Y < buttons[8].posY+buttons[8].frame.H() {
			return buttonControlsCollision
		}
	}
	if loc.X < buttons[9].posX+buttons[9].frame.W() && loc.X > buttons[9].posX-buttons[9].frame.W()/2 {
		if loc.Y < buttons[9].posY+buttons[9].frame.H() {
			return buttonRunCollision
		}
	}
	if loc.X < buttons[6].posX+buttons[6].frame.W() && loc.X > (buttons[6].posX-buttons[6].frame.W()/2) {
		if loc.Y < buttons[6].posY+buttons[6].frame.H() {
			return buttonScaleCollision
		}
	}
	if loc.X < buttons[3].posX+buttons[3].frame.W() && loc.X > buttons[3].posX-buttons[3].frame.W()/2 {
		if loc.Y < buttons[3].posY+buttons[3].frame.H() {
			return buttonClearCollision
		}
	}
	if loc.X < buttons[2].posX+buttons[2].frame.W() && loc.X > buttons[2].posX-buttons[2].frame.W()/2 {
		if loc.Y < buttons[2].posY+buttons[2].frame.H() {
			return buttonSignalCollision
		}
	}
	if loc.X < buttons[1].posX+buttons[1].frame.W() && loc.X > buttons[1].posX-buttons[1].frame.W()/2 {
		if loc.Y < buttons[1].posY+buttons[1].frame.H() {
			return buttonWeatherCollision
		}
	}
	if loc.X < buttons[0].posX+buttons[0].frame.W() && loc.X > buttons[0].posX-buttons[0].frame.W()/2 {
		if loc.Y < buttons[0].posY+buttons[0].frame.H() {
			return buttonBuildingCollision
		}
	}
	return noCollision
}

func handleCollision(win *pixelgl.Window) {
	if !win.JustPressed(pixelgl.MouseButtonLeft) {
		return
	}
	switch whereClick(win.MousePosition()) {
	case noCollision:
		return
	case buildingCollision:
		newMessage("Cannot place a building on top of another building!", 100, standardFont)
	case buttonBuildingCollision:
		handleBuildingButton(win.MousePosition())
	case buttonWeatherCollision:
		handleWeatherButton()
	case buttonSignalCollision:
		handleSignalButton()
	case buttonClearCollision:
		handleClearButton()
	case buttonPerson1Collision:
		handlePersonPButton()
	case buttonPerson2Collision:
		handlePersonQButton()
	case buttonScaleCollision:
		handleScaleButton()
	case buttonLineCollision:
		handleLineButton()
	case buttonControlsCollision:
		handleControlsButton(win)
	case buttonRunCollision:
		handleRunButton(win)
	}
}

func handleMovementKeyPress(win *pixelgl.Window) {
	if drawStatic {
		handleStaticMovement(win)
		drawPerson(win, &personP)
		drawPerson(win, &personQ)
		return
	}
	p1 := &personP
	p2 := &personQ
	if currPerson == p {
		walkMap = walkingP
	} else {
		p1 = &personQ
		p2 = &personP
		walkMap = walkingQ
	}
	if win.Pressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
		walkMap[directionLeft][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
		if p1.loc.X > minSpriteX {
			p1.loc.X -= walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
		walkMap[directionRight][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
		if p1.loc.X < maxSpriteX {
			p1.loc.X += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
		walkMap[directionUp][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
		if p1.loc.Y < maxSpriteY {
			p1.loc.Y += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
		walkMap[directionDown][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
		if p1.loc.Y > minSpriteY {
			p1.loc.Y -= walkSpeed
		}
	} else {
		drawPerson(win, p1)
	}
	if currPerson == p {
		walkMap = walkingQ
	} else {
		walkMap = walkingP
	}
	if win.Pressed(pixelgl.KeyA) || win.Repeated(pixelgl.KeyA) {
		walkMap[directionLeft][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
		if p2.loc.X > minSpriteX {
			p2.loc.X -= walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyD) || win.Repeated(pixelgl.KeyD) {
		walkMap[directionRight][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
		if p2.loc.X < maxSpriteX {
			p2.loc.X += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyW) || win.Repeated(pixelgl.KeyW) {
		walkMap[directionUp][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
		if p2.loc.Y < maxSpriteY {
			p2.loc.Y += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyS) || win.Repeated(pixelgl.KeyS) {
		walkMap[directionDown][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
		if p2.loc.Y > minSpriteY {
			p2.loc.Y -= walkSpeed
		}
	} else {
		drawPerson(win, p2)
	}
}

func handleLoadingScreenOk(clicked bool, mousepos pixel.Vec) bool {
	if clicked {
		if vectorIntersectionWithObject(mousepos, okbutton) {
			return true
		}
	}
	return false
}

func handleBuildingAdded(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.Key1) {
		if drawAnimation {
			stopAnimation()
			staticobject = object{}
			drawStatic = false
			return
		}
		drawNewBuilding(win.MousePosition())
	}
	if win.JustReleased(pixelgl.Key2) {
		handleWeatherButton()
	}
}

func handleStaticMovement(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
		staticobject.loc.X -= staticSpeed
	}
	if win.Pressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
		staticobject.loc.X += staticSpeed
	}
	if win.Pressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
		staticobject.loc.Y += staticSpeed
	}
	if win.Pressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
		staticobject.loc.Y -= staticSpeed
	}
	if win.JustReleased(pixelgl.KeyA) {
		staticobject = newBuilding(pixel.V(staticobject.loc.X, staticobject.loc.Y), -1)
	}
	if win.JustReleased(pixelgl.KeyD) {
		staticobject = newBuilding(pixel.V(staticobject.loc.X, staticobject.loc.Y), 1)
	}
	if win.Pressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
		staticobject.loc.Y -= staticSpeed
	}
	if win.Pressed(pixelgl.KeyEnter) || win.JustReleased(pixelgl.KeyEnter) {
		if collisionDetected(pixel.V(staticobject.loc.X, staticobject.loc.Y), staticobjects) {
			newMessage("Cannot place a building on top of another building!", 100, standardFont)
		} else {
			staticobjects = append(staticobjects, staticobject)
			staticobject = object{}
			drawStatic = false
			newMessage("Building added!", 100, standardFont)
			stopAnimation()
		}
	}
}

func handleDistanceLineKey(pressed bool) {
	if pressed {
		buttons[7].drawcount = 10
		drawingDistanceLine = !drawingDistanceLine
		if drawingDistanceLine {
			newMessage("Showing distance line", 100, standardFont)
		}
	}
}

func handlePersonKeyPressed(pressed bool) {
	if pressed {
		currPerson = !currPerson
	}
}

func handleRunButton(win *pixelgl.Window) {
	buttons[9].drawcount = 10
}

func handleControlsButton(win *pixelgl.Window) {
	buttons[8].drawcount = 10
	drawControlScreen(win)
}

func handleLineButton() {
	buttons[7].drawcount = 10
	handleDistanceLineKey(true)
}

func handleWeatherButton() {
	buttons[1].drawcount = 10
	drawingRain = !drawingRain
	newMessage("Changing environment", 100, standardFont)
}

func handleSignalButton() {
	buttons[2].drawcount = 10
	newMessage(satelliteError(), 100, standardFont)
}

func handleScaleButton() {
	buttons[6].drawcount = 10
	currScale++
	if currScale == len(scaleNames) {
		currScale = 0
	}
	newMessage("Distance scale changed to "+scaleNames[currScale], 100, standardFont)
}

func handlePersonPButton() {
	currPerson = p
}

func handlePersonQButton() {
	currPerson = q
}

func handleClearButton() {
	buttons[3].drawcount = 10
	clearSprites()
}

func handleBuildingButton(pos pixel.Vec) {
	newBuilding(pixel.V(maxSpriteX/2, maxSpriteY/2), 1)
	drawNewBuilding(pixel.V(maxSpriteX/2, maxSpriteY/2))
	buttons[0].drawcount = 10
}

func collisionDetected(v pixel.Vec, objects []object) bool {
	for _, o := range objects {
		if v.X >= o.loc.X && v.X <= o.frame.W()+o.loc.X {
			if v.Y >= o.loc.Y && v.Y <= o.frame.H()+o.loc.Y {
				return true
			}
		}
		println("new object: [", int(v.X), int(v.Y), "]", "object: [", int(o.loc.X), int(o.loc.Y), "]")
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
