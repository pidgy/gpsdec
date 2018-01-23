package gpsdec

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var (
	displayMessage      string
	displayMessageCount int
	displayMessageSize  *basicfont.Face
	second              = time.Tick(time.Second)
	millisecond         = time.Tick(time.Millisecond * 100)
	framespersecond     int
	selectedButton      int
	locP                = pixel.V(maxX/2, maxY/2)
	locQ                = pixel.V(float64(rand.Intn(int(maxX))), float64(rand.Intn(int(maxY))))
	numHumans           = 1
	numButtons          = 4
	numSatellites       = 3
	standardFont        = basicfont.Face7x13
	drawRain            = false
	drawDistanceLine    = false
)

func remove(s []message, i int) []message {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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

func drawMessage(m string, c int, s *basicfont.Face) {
	displayMessage = m
	displayMessageCount = c
	displayMessageSize = s
}

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
	units := drawUnits()
	road, err := loadPicture(spritedirectory + "road.png")
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Skyblue)

	satelliteAngle := 0.0
	last := time.Now()
	flip := 0

	var walkMap map[int][]object
	var rainSprites []*pixel.Sprite
	var buildingSprites []*pixel.Sprite

	for i := range rain {
		rainSprites = append(rainSprites, pixel.NewSprite(rain[i].pic, rain[i].frame))
	}

	loadScreen.sprite.Draw(win, loadScreen.mat)
	win.Update()
	load.Wait()

	controlPic, err := loadPicture(spritedirectory + "controls.png")
	if err != nil {
		panic(err)
	}
	controls := pixel.NewSprite(controlPic, controlPic.Bounds())
	controls.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	for {
		if win.JustReleased(pixelgl.MouseButtonLeft) {
			break
		}
		if win.JustReleased(pixelgl.KeyEnter) {
			break
		}
		time.Sleep(1000)
		win.Update()
	}

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Skyblue)

		// Draw the satellites
		satelliteAngle += 2 * dt
		for i := 0; i < numSatellites; i++ {
			mat := pixel.IM
			mat = mat.Scaled(pixel.ZV, 0.1).Rotated(pixel.ZV, satelliteAngle)
			if satellites[i].posX == 1 || satellites[i].posX == maxX-1 {
				satellites[i].direction = !satellites[i].direction
			}
			if satellites[i].direction == left {
				satellites[i].posX--
			}
			if satellites[i].direction == right {
				satellites[i].posX++
			}
			mat = mat.Moved(pixel.V(satellites[i].posX, 740))
			sat := pixel.NewSprite(satellites[i].pic, satellites[i].frame)
			sat.Draw(win, mat)
		}
		// Draw a road under the person
		roadPix := pixel.NewSprite(road, road.Bounds())
		roadPix.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(win.Bounds().Max.X/2, (win.Bounds().Max.Y/2)-50)))
		roadPix.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(0, (win.Bounds().Max.Y/2)-50)))
		roadPix.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(pixel.V(win.Bounds().Max.X, (win.Bounds().Max.Y/2)-50)))

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			switch whereClick(win.MousePosition()) {
			case noCollision:
				f := buildFrames[currBuildingName]
				building := pixel.NewSprite(buildPic, f)
				buildingSprites = append(buildingSprites, building)
				buildings = append(buildings, object{
					pic:   buildPic,
					posX:  win.MousePosition().X - f.W()/2,
					posY:  win.MousePosition().Y - f.H()/2,
					frame: f,
					mat:   pixel.IM.Moved(win.MousePosition())})
			case buildingCollision:
				drawMessage("Cannot place a building on top of another building!", 100, standardFont)
			case buttonBuildingCollision:
				currBuildingName = (currBuildingName + 1) % len(buildFrames)
				drawMessage(buildingNames[currBuildingName]+" Selected", 100, standardFont)
				buttons[0].drawcount = 10
			case buttonWeatherCollision:
				buttons[1].drawcount = 10
				drawRain = !drawRain
				drawMessage("Changing environment", 100, standardFont)
			case buttonGPSCollision:
				buttons[2].drawcount = 10
				drawMessage(satelliteError(), 100, standardFont)
			case buttonClearCollision:
				buttons[3].drawcount = 10
				clearSprites()
			case buttonPerson1Collision:
				currPerson = p
			case buttonPerson2Collision:
				currPerson = q
			case buttonScaleCollision:
				buttons[6].drawcount = 10
				currScale++
				if currScale == len(scaleNames) {
					currScale = 0
				}
				drawMessage("Distance scale changed to "+scaleNames[currScale], 100, standardFont)
			case buttonLineCollision:
				buttons[7].drawcount = 10
				drawDistanceLine = !drawDistanceLine
				if drawDistanceLine {
					drawMessage("Showing distance line", 100, standardFont)
				}
			}
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
		if win.JustReleased(pixelgl.KeyTab) {
			currPerson = !currPerson
		}
		if win.Pressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
			walkMap[directionLeft][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
			if p1.loc.X > 150 {
				p1.loc.X -= 3
			}
		} else if win.Pressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
			walkMap[directionRight][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
			if p1.loc.X < maxX-150 {
				p1.loc.X += 3
			}
		} else if win.Pressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
			walkMap[directionUp][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
			if p1.loc.Y < 632 {
				p1.loc.Y += 3
			}
		} else if win.Pressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
			walkMap[directionDown][flip].sprite.Draw(win, p1.mat.Moved(p1.loc))
			if p1.loc.Y > 100 {
				p1.loc.Y -= 3
			}
		} else {
			p1.sprite.Draw(win, p1.mat.Moved(p1.loc))
		}
		if currPerson == p {
			walkMap = walkingQ
		} else {
			walkMap = walkingP
		}
		if win.Pressed(pixelgl.KeyA) || win.Repeated(pixelgl.KeyA) {
			walkMap[directionLeft][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
			if p2.loc.X > 150 {
				p2.loc.X -= 3
			}
		} else if win.Pressed(pixelgl.KeyD) || win.Repeated(pixelgl.KeyD) {
			walkMap[directionRight][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
			if p2.loc.X < maxX-150 {
				p2.loc.X += 3
			}
		} else if win.Pressed(pixelgl.KeyW) || win.Repeated(pixelgl.KeyW) {
			walkMap[directionUp][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
			if p2.loc.Y < 632 {
				p2.loc.Y += 3
			}
		} else if win.Pressed(pixelgl.KeyS) || win.Repeated(pixelgl.KeyS) {
			walkMap[directionDown][flip].sprite.Draw(win, p2.mat.Moved(p2.loc))
			if p2.loc.Y > 100 {
				p2.loc.Y -= 3
			}
		} else {
			p2.sprite.Draw(win, p2.mat.Moved(p2.loc))
		}
		if !drawDistanceLine {
			x, y := distance(personP.loc, personQ.loc)
			angleLen := angleLength(x, y)
			angle := distanceAngle(x, y, personP, personQ)
			distanceLine.mat = pixel.IM.ScaledXY(pixel.ZV, pixel.V(angleLen, 1)).Moved(pixel.V(personP.loc.X+angleLen*100, personP.loc.Y)).Rotated(pixel.V(personP.loc.X, personP.loc.Y), angle)
			distanceLine.sprite.Draw(win, distanceLine.mat)
		}
		for _, u := range units {
			u.Draw(win, pixel.IM)
		}

		basicAtlas := text.NewAtlas(standardFont, text.ASCII)
		basicTxt := text.New(pixel.V(personQ.loc.X, personQ.loc.Y), basicAtlas)
		fmt.Fprintln(basicTxt, fmt.Sprintf("  [%f|%f]", personQ.loc.X, personQ.loc.Y))
		basicTxt.Draw(win, pixel.IM)

		basicAtlas = text.NewAtlas(standardFont, text.ASCII)
		basicTxt = text.New(pixel.V(personP.loc.X, personP.loc.Y), basicAtlas)
		fmt.Fprintln(basicTxt, fmt.Sprintf("  [%f|%f]", personP.loc.X, personP.loc.Y))
		basicTxt.Draw(win, pixel.IM)

		if win.JustReleased(pixelgl.KeyL) {
			buttons[7].drawcount = 10
			drawDistanceLine = !drawDistanceLine
			if drawDistanceLine {
				drawMessage("Showing distance line", 100, standardFont)
			}
		}

		if win.JustReleased(pixelgl.KeyEnter) {
			println(int(personP.loc.X), int(personP.loc.Y))
			println(int(personQ.loc.X), int(personQ.loc.Y))
			x, y := distance(personP.loc, personQ.loc)
			println("X distance:", x)
			println("Y distance:", y)
			println("angle distance:", distanceAngle(x, y, personP, personQ))

			println(fmt.Sprintf("%f,%f,%f", distanceAngle(x, y, personP, personQ), x, y))

		}

		buildingBatch.Clear()
		for i, building := range buildings {
			buildingSprites[i].Draw(buildingBatch, building.mat)
		}
		buildingBatch.Draw(win)

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

		if displayMessageCount > 0 {
			basicAtlas := text.NewAtlas(displayMessageSize, text.ASCII)
			basicTxt := text.New(pixel.V(maxX/3, maxY-200), basicAtlas)
			fmt.Fprintln(basicTxt, displayMessage)
			basicTxt.Draw(win, pixel.IM)
			displayMessageCount--
		}
		if drawRain {
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
func Run() {
	pixelgl.Run(run)
}
