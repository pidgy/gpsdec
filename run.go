package gpsdec

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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
	standardFont        = basicfont.Face7x13
	drawRain            = false
	drawDistanceLine    = false
	last                = time.Now()
)

func remove(s []message, i int) []message {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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
	win.Clear(colornames.Blue)

	flip := 0

	loadScreen.sprite.Draw(win, loadScreen.mat)
	win.Update()
	load.Wait()

	controlPic, err := loadPicture(spritedirectory + "controls.png")
	if err != nil {
		panic(err)
	}
	controls := pixel.NewSprite(controlPic, controlPic.Bounds())
	var loadSats [][]object
	i := 0
	for i < 2 {
		loadSats = append(loadSats, scramblePositions(satellites))
		i++
	}
	for {
		controls.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		if win.Closed() {
			goto exit
		}
		if win.JustReleased(pixelgl.MouseButtonLeft) {
			break
		}
		if win.JustReleased(pixelgl.KeyEnter) {
			break
		}
		for _, l := range loadSats {
			drawRandom(win, l, 2)
		}
		win.Update()
	}

	for !win.Closed() {
		win.Clear(colornames.Blue)
		background.sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		drawSatellites(win)
		handleLeftClick(win)
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
			p1.sprite.Draw(win, p1.mat.Moved(p1.loc))
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
			p2.sprite.Draw(win, p2.mat.Moved(p2.loc))
		}
		if drawDistanceLine {
			x, y := distance(personP.loc, personQ.loc)
			angleLen := angleLength(x, y)
			imd := imdraw.New(nil)
			imd.Color = colornames.Red
			imd.EndShape = imdraw.RoundEndShape
			imd.Push(pixel.V(personP.loc.X, personP.loc.Y), pixel.V(personQ.loc.X, personQ.loc.Y))
			imd.EndShape = imdraw.SharpEndShape
			imd.Line(1)
			imd.Draw(win)
			basicAtlas := text.NewAtlas(standardFont, text.ASCII)
			basicTxt := text.New(win.Bounds().Center(), basicAtlas)
			fmt.Fprintln(basicTxt, fmt.Sprintf("%.2f %s", angleLen*1000, scaleNames[currScale]))
			basicTxt.Draw(win, pixel.IM)
		}
		if win.JustReleased(pixelgl.KeyL) {
			buttons[7].drawcount = 10
			drawDistanceLine = !drawDistanceLine
			if drawDistanceLine {
				drawMessage("Showing distance line", 100, standardFont)
			}
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
exit:
}

// Run this from your program to start simulation
func Run() {
	pixelgl.Run(run)
}
