package gpsdec

import (
	"fmt"

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
	buttonEstimateCollision
	buttonTipCollision
	buttonEphemerisCollision
	buttonElevationCollision
)

var (
	currSatelliteError = 1
	currWeather        = 0

	userSelectSignal = false
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
			return buttonTipCollision
		}
		if loc.Y > buttons[5].posY-buttons[5].frame.H()/2 && loc.Y < buttons[5].posY+buttons[5].frame.H()/2 {
			return buttonControlsCollision
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
	if loc.X < buttons[10].posX+buttons[10].frame.W() && loc.X > buttons[10].posX-buttons[10].frame.W()/2 {
		if loc.Y < buttons[10].posY+buttons[10].frame.H() {
			return buttonEstimateCollision
		}
	}
	if loc.X < buttons[11].posX+buttons[11].frame.W() && loc.X > buttons[11].posX-buttons[11].frame.W()/2 {
		if loc.Y < buttons[11].posY+buttons[11].frame.H() {
			return buttonEphemerisCollision
		}
	}
	if loc.X < buttons[12].posX+buttons[12].frame.W() && loc.X > buttons[12].posX-buttons[12].frame.W()/2 {
		if loc.Y < buttons[12].posY+buttons[12].frame.H() {
			return buttonElevationCollision
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
	case buttonEstimateCollision:
		handleEstimateButton(win)
	case buttonTipCollision:
		handleTipButton(win)
	case buttonElevationCollision:
		handleElevationButton(win)
	case buttonEphemerisCollision:
		handleEphemerisButton(win)
	}
}

func handleMouseHover(win *pixelgl.Window) {
	switch whereClick(win.MousePosition()) {
	case noCollision:
		return
	case buildingCollision:
		newHelpMessage("Building", 100, standardFont)
	case buttonBuildingCollision:
		newHelpMessage("Add a Building!", 100, standardFont)
	case buttonWeatherCollision:
		newHelpMessage("Add Harsh Weather Elements", 100, standardFont)
	case buttonSignalCollision:
		newHelpMessage("Add Satellite Error", 100, standardFont)
	case buttonClearCollision:
		newHelpMessage("Remove All Elements", 100, standardFont)
	case buttonPerson1Collision:
		newHelpMessage("Switch Control For Red", 100, standardFont)
	case buttonPerson2Collision:
		newHelpMessage("Switch Control For Purple", 100, standardFont)
	case buttonScaleCollision:
		newHelpMessage("Change Scale For P->Q Distance", 100, standardFont)
	case buttonLineCollision:
		newHelpMessage("Show Distance Line From P->Q", 100, standardFont)
	case buttonControlsCollision:
		newHelpMessage("Show HotKey Controls", 100, standardFont)
	case buttonRunCollision:
		newHelpMessage("Run Simulation", 100, standardFont)
	case buttonEstimateCollision:
		newHelpMessage("Estimate Positions and Distance", 100, standardFont)
	case buttonTipCollision:
		newHelpMessage("Show Last Tip", 100, standardFont)
	case buttonElevationCollision:
		newHelpMessage(fmt.Sprintf("Change Elevation (Currently: %.1f ft.)", elevations[currElevation]), 100, standardFont)
	case buttonEphemerisCollision:
		newHelpMessage("Add/Remove Ephemeris Error (Hint: Elevation will affect the error!)", 100, standardFont)
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
		drawMovingPerson(win, directionLeft, flip, p1)
		if p1.loc.X > minSpriteX {
			p1.loc.X -= walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
		drawMovingPerson(win, directionRight, flip, p1)
		if p1.loc.X < maxSpriteX {
			p1.loc.X += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
		drawMovingPerson(win, directionUp, flip, p1)
		if p1.loc.Y < maxSpriteY {
			p1.loc.Y += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
		drawMovingPerson(win, directionDown, flip, p1)
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
		drawMovingPerson(win, directionLeft, flip, p2)
		if p2.loc.X > minSpriteX {
			p2.loc.X -= walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyD) || win.Repeated(pixelgl.KeyD) {
		drawMovingPerson(win, directionRight, flip, p2)
		if p2.loc.X < maxSpriteX {
			p2.loc.X += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyW) || win.Repeated(pixelgl.KeyW) {
		drawMovingPerson(win, directionUp, flip, p2)
		if p2.loc.Y < maxSpriteY {
			p2.loc.Y += walkSpeed
		}
	} else if win.Pressed(pixelgl.KeyS) || win.Repeated(pixelgl.KeyS) {
		drawMovingPerson(win, directionDown, flip, p2)
		if p2.loc.Y > minSpriteY {
			p2.loc.Y -= walkSpeed
		}
	} else {
		drawPerson(win, p2)
	}
}

func handleOKButtonClicked(clicked bool, mousepos pixel.Vec) bool {
	if clicked && vectorIntersectionWithObject(mousepos, okbutton) {
		return true
	}
	return false
}

func handleBuildingAdded(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.Key1) {
		if drawingAnimation {
			stopAnimation()
			staticobject = object{}
			drawStatic = false
			return
		}
		drawNewBuilding(win.MousePosition())
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

func handleUserSelectInput(win *pixelgl.Window) float64 {
	switch true {
	case win.JustPressed(pixelgl.Key1) || win.JustReleased(pixelgl.Key1):
		return 1
	case win.JustPressed(pixelgl.Key2) || win.JustReleased(pixelgl.Key2):
		return 2
	case win.JustPressed(pixelgl.Key3) || win.JustReleased(pixelgl.Key3):
		return 3
	case win.JustPressed(pixelgl.Key4) || win.JustReleased(pixelgl.Key4):
		return 4
	}
	return 0
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
	ceEstimateDistance()
}

func handleEstimateButton(win *pixelgl.Window) {
	buttons[10].drawcount = 10
	drawingPositionEstimates = !drawingPositionEstimates
	if drawingPositionEstimates {
		if currScale == SCALE_M {
			newMessage("Big Overstimation of Distance", 100, standardFont)
		} else if currScale == SCALE_KM {
			newMessage("Small Overstimation of Distance", 100, standardFont)
		}
	}
}

func handleTipButton(win *pixelgl.Window) {
	buttons[4].drawcount = 10
	drawingTip = true
}

func handleEphemerisButton(win *pixelgl.Window) {
	buttons[11].drawcount = 10
	if !ephemerisError {
		ephemerisError = true
		newMessage("Ephemeris Error Has Been Added", 100, standardFont)
		return
	}
	ephemerisError = false
	newMessage("Ephemeris Error Has Been Removed", 100, standardFont)
}

func handleElevationButton(win *pixelgl.Window) {
	buttons[12].drawcount = 10
	currElevation = (currElevation + 1) % len(elevations)
	newMessage(fmt.Sprintf("Elevation changed to: %.1f ft.", elevations[currElevation]), 100, standardFont)
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
	currWeather = (currWeather + 1) % (len(weatherObjects) + 1)
	switch currWeather {
	case WEATHER_NONE:
		newMessage("Removed all weather elements", 100, standardFont)
		drawingWeather = false
	case WEATHER_RAIN:
		newMessage("Adding rain to the environment", 100, standardFont)
		drawingWeather = true
	case WEATHER_ASH:
		newMessage("Adding volcanic ash to the environment", 100, standardFont)
		drawingWeather = true
	case WEATHER_DRY:
		newMessage("Adding dry air to the environment", 100, standardFont)
		drawingWeather = true
	case WEATHER_SAND:
		newMessage("Adding a sandstorm to the environment", 100, standardFont)
		drawingWeather = true
	}
}

func handleSignalButton() {
	buttons[2].drawcount = 10
	newErr := satelliteError()
	newMessage(newErr, 100, standardFont)
}

func handleScaleButton() {
	buttons[6].drawcount = 10
	currScale = (currScale + 1) % 2
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
	if currentAnimation == currentUserSelect {
		newMessage("Cannot add a building right now!", 100, standardFont)
		return
	}
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
	}
	return false
}

func satelliteError() string {
	returnString := ""
	switch currSatelliteError {
	case 1:
		numSatellites--
		dopValue = DOP_RANGE_GOOD
		returnString = "DOP value set to range 2-5 (Good)"
	case 2:
		numSatellites--
		dopValue = DOP_RANGE_MODERATE
		returnString = "DOP value set to range 5-10 (Moderate)"
	case 3:
		numSatellites--
		dopValue = DOP_RANGE_FAIR
		returnString = "DOP value set to range 10-20 (Fair)"
	case 4:
		numSatellites--
		dopValue = DOP_RANGE_POOR
		returnString = "DOP value set to > 20 (Poor)"
	case 5:
		numSatellites = len(satellites)
		dopValue = DOP_RANGE_IDEAL
		returnString = "DOP value set to < 1 (Ideal)"
	case 6:
		if currentAnimation == currentBuilding {
			newMessage("Cannot select clock drift while adding a building!", 100, standardFont)
		}
		drawingUserSelectionWin = true
		returnString = "Select GPS clock drift!"
		currentTipMessageByte = 0
		userSelectSignal = true
		startAnimation(currentUserSelect)
	}
	if currSatelliteError == 6 {
		currSatelliteError = 1
	} else {
		currSatelliteError++
	}
	return returnString
}
