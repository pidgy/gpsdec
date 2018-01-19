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
	second              = time.Tick(time.Second)
	millisecond         = time.Tick(time.Millisecond * 100)
	framespersecond     int
	selectedButton      int
	locP                = pixel.V(maxX/2, maxY/2)
	locQ                = pixel.V(float64(rand.Intn(int(maxX))), float64(rand.Intn(int(maxY))))
	drawRain            = false
	numHumans           = 1
	numButtons          = 4
	numSatellites       = 3
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

func run() {
	units := drawUnits()
	cfg := pixelgl.WindowConfig{
		Title:  "GPS Error Simulations",
		Bounds: pixel.R(0, 0, maxX, maxY),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	road, err := loadPicture(spritedirectory + "road.png")
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Skyblue)

	satelliteAngle := 0.0
	last := time.Now()
	flip := 0

	var walkMap map[int][]object

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Skyblue)

		for i := range units {
			units[i].Draw(win, pixel.IM)
		}

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
				buildings = append(buildings, object{
					pic:   buildPic,
					posX:  win.MousePosition().X - f.W()/2,
					posY:  win.MousePosition().Y - f.H()/2,
					frame: f,
					mat:   pixel.IM.Moved(win.MousePosition())})
			case buildingCollision:
				displayMessage = "Cannot place a building on top of another building!"
				displayMessageCount = 100
			case buttonBuildingCollision:
				currBuildingName = (currBuildingName + 1) % len(buildFrames)
				displayMessage = buildingNames[currBuildingName] + " Selected"
				displayMessageCount = 100
				buttons[0].drawcount = 10
			case buttonWeatherCollision:
				buttons[1].drawcount = 10
				drawRain = !drawRain
				displayMessage = "Changing environment"
				displayMessageCount = 100
			case buttonGPSCollision:
				buttons[2].drawcount = 10
				displayMessage = satelliteError()
				displayMessageCount = 100
			case buttonClearCollision:
				buttons[3].drawcount = 10
				clearSprites()
			case buttonPerson1Collision:
				currPerson = p
			case buttonPerson2Collision:
				currPerson = q
			}
		}
		obj := &personP
		if currPerson == p {
			walkMap = walkingP
		} else {
			obj = &personQ
			walkMap = walkingQ
		}
		if win.JustReleased(pixelgl.KeyTab) {
			currPerson = !currPerson
		}
		if win.Pressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
			hum := pixel.NewSprite(walkMap[directionLeft][flip].pic, walkMap[directionLeft][flip].frame)
			hum.Draw(win, obj.mat.Moved(obj.loc))
			if obj.loc.X > 150 {
				obj.loc.X -= 3
			}
		} else if win.Pressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
			hum := pixel.NewSprite(walkMap[directionRight][flip].pic, walkMap[directionRight][flip].frame)
			hum.Draw(win, obj.mat.Moved(obj.loc))
			if obj.loc.X < maxX-150 {
				obj.loc.X += 3
			}
		} else if win.Pressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
			hum := pixel.NewSprite(walkMap[directionUp][flip].pic, walkMap[directionUp][flip].frame)
			hum.Draw(win, obj.mat.Moved(obj.loc))
			if obj.loc.Y < 632 {
				obj.loc.Y += 3
			}
		} else if win.Pressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
			hum := pixel.NewSprite(walkMap[directionDown][flip].pic, walkMap[directionDown][flip].frame)
			hum.Draw(win, obj.mat.Moved(obj.loc))
			if obj.loc.Y > 100 {
				obj.loc.Y -= 3
			}
		} else {
			hum := pixel.NewSprite(obj.pic, obj.frame)
			hum.Draw(win, obj.mat.Moved(obj.loc))
		}
		if currPerson == p {
			hum := pixel.NewSprite(personQ.pic, personQ.frame)
			hum.Draw(win, personQ.mat.Moved(personQ.loc))
		} else {
			hum := pixel.NewSprite(personP.pic, personP.frame)
			hum.Draw(win, personP.mat.Moved(personP.loc))
		}

		for i := range buildings {
			building := pixel.NewSprite(buildings[i].pic, buildings[i].frame)
			building.Draw(win, buildings[i].mat)
		}
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
			basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
			basicTxt := text.New(pixel.V(maxX/3, maxY-200), basicAtlas)
			fmt.Fprintln(basicTxt, displayMessage)
			basicTxt.Draw(win, pixel.IM)
			displayMessageCount--
		}
		if drawRain {
			for i := range rain {
				rdrop := pixel.NewSprite(rain[i].pic, rain[i].frame)
				rdrop.Draw(win, pixel.IM.Moved(pixel.V(rain[i].posX, rain[i].posY)))
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

func Run() {
	pixelgl.Run(run)
}
