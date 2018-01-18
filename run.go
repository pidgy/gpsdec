package gpsimulation

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	maxX            float64 = 1024
	maxY            float64 = 768
	left            bool    = true
	right           bool    = false
	spritedirectory string  = "gpssimulation/pics/"
)

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

var (
	buildings       []object
	messages        []message
	second          = time.Tick(time.Second)
	framespersecond int
	buildPic        pixel.Picture
	buildFrames     []pixel.Rect
	satellites      []object
	buttons         []object
	humans          []object
	selectedButton  int
	locP            = pixel.V(maxX/2, maxY/2)
	locQ            = pixel.V(float64(rand.Intn(int(maxX))), float64(rand.Intn(int(maxY))))

	numHumans     = 1
	numButtons    = 3
	numSatellites = 3
)

func init() {
	buildPic, buildFrames = buildingFrames()

	buttons = []object{
		object{posX: 40, posY: 30, filename: spritedirectory + "button-buildings.png", pressedfilename: spritedirectory + "button-pressed-buildings.png"},
		object{posX: 120, posY: 30, filename: spritedirectory + "button-weather.png", pressedfilename: spritedirectory + "button-pressed-weather.png"},
		object{posX: 200, posY: 30, filename: spritedirectory + "button-gps.png", pressedfilename: spritedirectory + "button-pressed-gps.png"},
	}
	satellites = []object{
		object{posX: 1000, angle: 10, direction: left},
		object{posX: 20, angle: -0.45, direction: right},
		object{posX: maxX / 2, angle: -1.5, direction: left},
	}
	humans = []object{
		object{},
	}
	loadSatelliteFrames()
	loadButtonFrames()
	loadHumanFrames()
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func loadHumanFrames() {
	sprite, err := loadPicture(spritedirectory + "person.gif")
	if err != nil {
		panic(err)
	}
	for i := 0; i < numHumans; i++ {
		humans[i].loc = pixel.V(maxX/2, maxY/2)
		humans[i].frame = sprite.Bounds()
		humans[i].mat = pixel.IM.Scaled(pixel.ZV, 0.3)
		humans[i].pic = sprite
	}
}

func buildingFrames() (pixel.Picture, []pixel.Rect) {
	spritesheet, err := loadPicture(spritedirectory + "skyscraper-spritesheet.png")
	if err != nil {
		panic(err)
	}
	var buildingFrames []pixel.Rect

	buildingFrames = append(buildingFrames, pixel.R(0, 0, 100, 100))    // DC
	buildingFrames = append(buildingFrames, pixel.R(100, 0, 170, 100))  // Hospital
	buildingFrames = append(buildingFrames, pixel.R(170, 0, 240, 100))  // Hospital 2
	buildingFrames = append(buildingFrames, pixel.R(0, 100, 75, 190))   // Apartment
	buildingFrames = append(buildingFrames, pixel.R(75, 100, 140, 190)) // Apartment 2

	buildingFrames = append(buildingFrames, pixel.R(0, 385, 120, 700))   // Big Sky
	buildingFrames = append(buildingFrames, pixel.R(515, 385, 600, 700)) // Big Sky 2
	buildingFrames = append(buildingFrames, pixel.R(120, 385, 210, 700)) // Big Sky 3
	buildingFrames = append(buildingFrames, pixel.R(210, 385, 310, 700)) // Big Sky 4
	buildingFrames = append(buildingFrames, pixel.R(310, 385, 410, 700)) // Big Sky 5
	buildingFrames = append(buildingFrames, pixel.R(410, 385, 510, 700)) // Big Sky 6
	buildingFrames = append(buildingFrames, pixel.R(510, 385, 610, 700)) // Big Sky 7

	return spritesheet, buildingFrames
}

func loadButtonFrames() {
	for i := 0; i < numButtons; i++ {
		button1, err := loadPicture(buttons[i].filename)
		buttonpressed, err := loadPicture(buttons[i].pressedfilename)
		if err != nil {
			panic(err)
		}
		buttons[i].frame = button1.Bounds()
		buttons[i].pic = button1
		buttons[i].pressedpic = buttonpressed
		buttons[i].mat = pixel.IM.Moved(pixel.V(buttons[i].posX, buttons[i].posY))
	}
}

func loadSatelliteFrames() {
	sprite, err := loadPicture(spritedirectory + "satellite-pixel.png")
	if err != nil {
		panic(err)
	}
	for i := 0; i < numSatellites; i++ {
		satellites[i].frame = sprite.Bounds()
		satellites[i].pic = sprite
	}
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

	road, err := loadPicture(spritedirectory + "road.png")
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Skyblue)

	satelliteAngle := 0.0
	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(colornames.Skyblue)

		// Draw the satellites
		satelliteAngle += 2 * dt
		for i := range satellites {
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
			if win.MousePosition().Y > 100 {
				if !collisionDetected(win.MousePosition(), buildings) {
					f := buildFrames[rand.Intn(len(buildFrames))]
					buildings = append(buildings, object{
						pic:   buildPic,
						posX:  win.MousePosition().X - f.W()/2,
						posY:  win.MousePosition().Y - f.H()/2,
						frame: f,
						mat:   pixel.IM.Moved(win.MousePosition())})
				} else {
					messages = append(messages, message{msg: "Cannot place a building on top of another building!", drawcount: 200})
				}
				goto exitpress
			}
			if win.MousePosition().X < buttons[0].posX+buttons[0].frame.W() && win.MousePosition().X > buttons[0].posX-buttons[0].frame.W()/2 {
				buttons[0].drawcount = 10
				goto exitpress
			}
			if win.MousePosition().X < buttons[1].posX+buttons[1].frame.W() && win.MousePosition().X > buttons[1].posX-buttons[1].frame.W()/2 {
				buttons[1].drawcount = 10
				goto exitpress
			}
			if win.MousePosition().X < buttons[2].posX+buttons[2].frame.W() && win.MousePosition().X > buttons[2].posX-buttons[2].frame.W()/2 {
				buttons[2].drawcount = 10
				goto exitpress
			}

		exitpress:
		}
		if win.Pressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
			humans[0].loc.X -= 3
		}
		if win.Pressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
			humans[0].loc.X += 3
		}
		if win.Pressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
			humans[0].loc.Y += 3
		}
		if win.Pressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
			humans[0].loc.Y -= 3
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
		for i := range messages {
			if messages[i].drawcount > 0 {
				basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
				basicTxt := text.New(pixel.V(100, 500), basicAtlas)
				basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
				fmt.Fprintln(basicTxt, messages[i].msg)
				basicTxt.Draw(win, pixel.IM)
				messages[i].drawcount--
			}
		}
		// Draw the human last so they are above the road
		for _, human := range humans {
			hum := pixel.NewSprite(human.pic, human.frame)
			hum.Draw(win, human.mat.Moved(human.loc))
		}
		win.Update()

		// FPS Update
		framespersecond++
		select {
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
