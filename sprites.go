package gpsdec

import (
	"image"
	"os"

	"github.com/faiface/pixel"
)

const (
	maxX            float64 = 1024
	maxY            float64 = 768
	left            bool    = true
	right           bool    = false
	p               bool    = true
	q               bool    = false
	spritedirectory string  = "gpsdec/pics/"

	directionLeft = iota
	directionRight
	directionUp
	directionDown
)

var (
	buildings  []object
	satellites []object
	buttons    []object
	rain       []object

	buildingBatch *pixel.Batch

	personP object
	personQ object

	buildPic    pixel.Picture
	buildFrames []pixel.Rect

	buildingNames    []string
	currBuildingName int

	scaleNames []string
	currScale  int

	currPerson = q

	walkingP map[int][]object
	walkingQ map[int][]object
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
	batch           *pixel.Batch
	sprite          *pixel.Sprite
}

func init() {
	buildPic, buildFrames = buildingFrames()
	buttons = []object{
		object{posX: 40, posY: 30, filename: spritedirectory + "button-buildings.png", pressedfilename: spritedirectory + "button-pressed-buildings.png"},
		object{posX: 120, posY: 30, filename: spritedirectory + "button-weather.png", pressedfilename: spritedirectory + "button-pressed-weather.png"},
		object{posX: 200, posY: 30, filename: spritedirectory + "button-gps.png", pressedfilename: spritedirectory + "button-pressed-gps.png"},
		object{posX: maxX - 40, posY: 30, filename: spritedirectory + "button-clear.png", pressedfilename: spritedirectory + "button-pressed-clear.png"},
		object{posX: maxX - 40, posY: maxY - 100, filename: spritedirectory + "button-person1.png", pressedfilename: spritedirectory + "button-pressed-person1.png"},
		object{posX: maxX - 40, posY: maxY - 170, filename: spritedirectory + "button-person2.png", pressedfilename: spritedirectory + "button-pressed-person2.png"},
		object{posX: 285, posY: 30, filename: spritedirectory + "button-scale.png", pressedfilename: spritedirectory + "button-pressed-scale.png"},
	}
	satellites = []object{
		object{posX: 1000, angle: 10, direction: left},
		object{posX: 20, angle: -0.45, direction: right},
		object{posX: maxX / 2, angle: -1.5, direction: left},
	}

	scaleNames = []string{
		"M",
		"Km",
	}
	currScale = 0

	buildingNames = []string{
		"Government Building",
		"Hospital 1",
		"Hospital 2",
		"Apartment 1",
		"Apartment 2",
		"Skyscraper 1",
		"Skyscraper 2",
		"Skyscraper 3",
		"Skyscraper 4",
		"Skyscraper 5",
		"Skyscraper 6",
		"Skyscraper 7",
		"Skyscraper 8",
	}
	loadSatelliteFrames()
	loadButtonFrames()
	loadHumanFrames()
	loadRain()
}

func clearSprites() {
	buildings = []object{}
	drawRain = false
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
	sprite1, err := loadPicture(spritedirectory + "person1/person-standing.png")
	sprite2, err := loadPicture(spritedirectory + "person2/person-standing.png")

	if err != nil {
		panic(err)
	}
	personP.loc = pixel.V(maxX/2, maxY/2)
	personP.frame = sprite1.Bounds()
	personP.mat = pixel.IM.Scaled(pixel.ZV, 0.3)
	personP.pic = sprite1
	personP.sprite = pixel.NewSprite(sprite1, sprite1.Bounds())

	personQ.loc = pixel.V(maxX/1.5, maxY/1.5)
	personQ.frame = sprite2.Bounds()
	personQ.mat = pixel.IM.Scaled(pixel.ZV, 0.3)
	personQ.pic = sprite2
	personQ.sprite = pixel.NewSprite(sprite2, sprite2.Bounds())

	spriteup1, err := loadPicture(spritedirectory + "person1/person-walking-up1.png")
	spriteup2, err := loadPicture(spritedirectory + "person1/person-walking-up2.png")
	spritedown1, err := loadPicture(spritedirectory + "person1/person-walking-down1.png")
	spritedown2, err := loadPicture(spritedirectory + "person1/person-walking-down2.png")
	spriteleft1, err := loadPicture(spritedirectory + "person1/person-walking-left1.png")
	spriteleft2, err := loadPicture(spritedirectory + "person1/person-walking-left2.png")
	spriteright1, err := loadPicture(spritedirectory + "person1/person-walking-right1.png")
	spriteright2, err := loadPicture(spritedirectory + "person1/person-walking-right2.png")
	if err != nil {
		panic(err)
	}
	walkingP = map[int][]object{
		directionLeft:  []object{object{pic: spriteleft1, frame: spriteleft1.Bounds(), sprite: pixel.NewSprite(spriteleft1, spriteleft1.Bounds())}, object{pic: spriteleft2, frame: spriteleft2.Bounds(), sprite: pixel.NewSprite(spriteleft2, spriteleft2.Bounds())}},
		directionRight: []object{object{pic: spriteright1, frame: spriteright1.Bounds(), sprite: pixel.NewSprite(spriteright1, spriteright1.Bounds())}, object{pic: spriteright2, frame: spriteright2.Bounds(), sprite: pixel.NewSprite(spriteright2, spriteright2.Bounds())}},
		directionUp:    []object{object{pic: spriteup1, frame: spriteup1.Bounds(), sprite: pixel.NewSprite(spriteup1, spriteup1.Bounds())}, object{pic: spriteup2, frame: spriteup2.Bounds(), sprite: pixel.NewSprite(spriteup2, spriteup2.Bounds())}},
		directionDown:  []object{object{pic: spritedown1, frame: spritedown1.Bounds(), sprite: pixel.NewSprite(spritedown1, spritedown1.Bounds())}, object{pic: spritedown2, frame: spritedown2.Bounds(), sprite: pixel.NewSprite(spritedown2, spritedown2.Bounds())}},
	}
	spriteup1, err = loadPicture(spritedirectory + "person2/person-walking-up1.png")
	spriteup2, err = loadPicture(spritedirectory + "person2/person-walking-up2.png")
	spritedown1, err = loadPicture(spritedirectory + "person2/person-walking-down1.png")
	spritedown2, err = loadPicture(spritedirectory + "person2/person-walking-down2.png")
	spriteleft1, err = loadPicture(spritedirectory + "person2/person-walking-left1.png")
	spriteleft2, err = loadPicture(spritedirectory + "person2/person-walking-left2.png")
	spriteright1, err = loadPicture(spritedirectory + "person2/person-walking-right1.png")
	spriteright2, err = loadPicture(spritedirectory + "person2/person-walking-right2.png")
	if err != nil {
		panic(err)
	}
	walkingQ = map[int][]object{
		directionLeft:  []object{object{pic: spriteleft1, frame: spriteleft1.Bounds(), sprite: pixel.NewSprite(spriteleft1, spriteleft1.Bounds())}, object{pic: spriteleft2, frame: spriteleft2.Bounds(), sprite: pixel.NewSprite(spriteleft2, spriteleft2.Bounds())}},
		directionRight: []object{object{pic: spriteright1, frame: spriteright1.Bounds(), sprite: pixel.NewSprite(spriteright1, spriteright1.Bounds())}, object{pic: spriteright2, frame: spriteright2.Bounds(), sprite: pixel.NewSprite(spriteright2, spriteright2.Bounds())}},
		directionUp:    []object{object{pic: spriteup1, frame: spriteup1.Bounds(), sprite: pixel.NewSprite(spriteup1, spriteup1.Bounds())}, object{pic: spriteup2, frame: spriteup2.Bounds(), sprite: pixel.NewSprite(spriteup2, spriteup2.Bounds())}},
		directionDown:  []object{object{pic: spritedown1, frame: spritedown1.Bounds(), sprite: pixel.NewSprite(spritedown1, spritedown1.Bounds())}, object{pic: spritedown2, frame: spritedown2.Bounds(), sprite: pixel.NewSprite(spritedown2, spritedown2.Bounds())}},
	}
}

func buildingFrames() (pixel.Picture, []pixel.Rect) {
	spritesheet, err := loadPicture(spritedirectory + "skyscraper-spritesheet.png")
	if err != nil {
		panic(err)
	}
	var buildingFrames []pixel.Rect
	buildingBatch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

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

func loadRain() {
	sprite, err := loadPicture(spritedirectory + "raining_tr.png")
	if err != nil {
		panic(err)
	}
	rain = append(rain, object{
		frame: sprite.Bounds(),
		batch: pixel.NewBatch(&pixel.TrianglesData{}, sprite),
		pic:   sprite,
		posY:  maxY,
		posX:  maxX / 2})
}

func loadButtonFrames() {
	for i := 0; i < len(buttons); i++ {
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
