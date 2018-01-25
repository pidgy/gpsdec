package gpsdec

import (
	"image"
	"os"
	"sync"

	"github.com/faiface/pixel"
)

const (
	maxX                float64 = 1024
	maxY                float64 = 768
	left                bool    = true
	right               bool    = false
	up                  bool    = true
	down                bool    = false
	p                   bool    = true
	q                   bool    = false
	spritedirectory     string  = "gpsdec/pics/"
	backgrounddirectory string  = "background/"
	buttonsdirectory    string  = "buttons/"
	objectsdirectory    string  = "objects/"
	pdirectory          string  = "person1/"
	qdirectory          string  = "person2/"
	imgtestdirectory    string  = "tests/"

	directionLeft = iota
	directionRight
	directionUp
	directionDown

	maxSpriteX = maxX - 5
	minSpriteX = 5
	maxSpriteY = 381
	minSpriteY = 100
	walkSpeed  = 3
)

var (
	buildings  []object
	satellites []object
	buttons    []object
	rain       []object

	buildingBatch *pixel.Batch

	personP      object
	personQ      object
	distanceLine object
	loadScreen   object

	background object

	buildPic    pixel.Picture
	buildFrames []pixel.Rect

	buildingNames    []string
	currBuildingName int

	scaleNames []string
	currScale  int

	currPerson = q

	walkingP map[int][]object
	walkingQ map[int][]object

	load sync.WaitGroup

	top = 740.0

	numSatellites = 3

	walkMap         map[int][]object
	rainSprites     []*pixel.Sprite
	buildingSprites []*pixel.Sprite
)

type message struct {
	msg       string
	drawcount int
}

type object struct {
	posX            float64
	posY            float64
	angle           float64
	directionX      bool
	directionY      bool
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
	loadingPic, err := loadPicture(spritedirectory + "loadscreen.png")
	if err != nil {
		panic(err)
	}
	loadScreen = object{
		pic:    loadingPic,
		sprite: pixel.NewSprite(loadingPic, loadingPic.Bounds()),
		mat:    pixel.IM.Moved(pixel.V(maxX/2, maxY/2)),
	}
	load.Add(1)
	go loadAllTheThings()
}

func loadAllTheThings() {
	buildPic, buildFrames = buildingFrames()
	buttons = []object{
		object{posX: 40, posY: 30, filename: spritedirectory + buttonsdirectory + "button-buildings.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-buildings.png"},
		object{posX: 120, posY: 30, filename: spritedirectory + buttonsdirectory + "button-weather.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-weather.png"},
		object{posX: 200, posY: 30, filename: spritedirectory + buttonsdirectory + "button-gps.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-gps.png"},
		object{posX: maxX - 40, posY: 30, filename: spritedirectory + buttonsdirectory + "button-clear.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-clear.png"},
		object{posX: maxX - 40, posY: maxY - 100, filename: spritedirectory + buttonsdirectory + "button-person1.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-person1.png"},
		object{posX: maxX - 40, posY: maxY - 170, filename: spritedirectory + buttonsdirectory + "button-person2.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-person2.png"},
		object{posX: 280, posY: 30, filename: spritedirectory + buttonsdirectory + "button-scale.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-scale.png"},
		object{posX: 360, posY: 30, filename: spritedirectory + buttonsdirectory + "button-line.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-line.png"},
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
	loadBackground()
	loadSatelliteFrames()
	loadButtonFrames()
	loadHumans()
	loadRain()
	loadDistanceLine()
	load.Done()
}

func clearSprites() {
	buildings = []object{}
	drawRain = false
	drawDistanceLine = false
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

func loadDistanceLine() {
	sprite, err := loadPicture(spritedirectory + objectsdirectory + "distance-line.png")
	if err != nil {
		panic(err)
	}
	distanceLine = object{
		pic:    sprite,
		sprite: pixel.NewSprite(sprite, sprite.Bounds()),
	}
}

func loadHumans() {
	sprite1, err := loadPicture(spritedirectory + pdirectory + "person-standing.png")
	sprite2, err := loadPicture(spritedirectory + qdirectory + "person-standing.png")
	if err != nil {
		panic(err)
	}
	personP.loc = pixel.V(341, 99)
	personP.frame = sprite1.Bounds()
	personP.mat = pixel.IM.Scaled(pixel.ZV, 0.3)
	personP.pic = sprite1
	personP.sprite = pixel.NewSprite(sprite1, sprite1.Bounds())

	personQ.loc = pixel.V(682, 99)
	personQ.frame = sprite2.Bounds()
	personQ.mat = pixel.IM.Scaled(pixel.ZV, 0.3)
	personQ.pic = sprite2
	personQ.sprite = pixel.NewSprite(sprite2, sprite2.Bounds())

	spriteup1, err := loadPicture(spritedirectory + pdirectory + "person-walking-up1.png")
	spriteup2, err := loadPicture(spritedirectory + pdirectory + "person-walking-up2.png")
	spritedown1, err := loadPicture(spritedirectory + pdirectory + "person-walking-down1.png")
	spritedown2, err := loadPicture(spritedirectory + pdirectory + "person-walking-down2.png")
	spriteleft1, err := loadPicture(spritedirectory + pdirectory + "person-walking-left1.png")
	spriteleft2, err := loadPicture(spritedirectory + pdirectory + "person-walking-left2.png")
	spriteright1, err := loadPicture(spritedirectory + pdirectory + "person-walking-right1.png")
	spriteright2, err := loadPicture(spritedirectory + pdirectory + "person-walking-right2.png")
	if err != nil {
		panic(err)
	}
	walkingP = map[int][]object{
		directionLeft:  []object{object{pic: spriteleft1, frame: spriteleft1.Bounds(), sprite: pixel.NewSprite(spriteleft1, spriteleft1.Bounds())}, object{pic: spriteleft2, frame: spriteleft2.Bounds(), sprite: pixel.NewSprite(spriteleft2, spriteleft2.Bounds())}},
		directionRight: []object{object{pic: spriteright1, frame: spriteright1.Bounds(), sprite: pixel.NewSprite(spriteright1, spriteright1.Bounds())}, object{pic: spriteright2, frame: spriteright2.Bounds(), sprite: pixel.NewSprite(spriteright2, spriteright2.Bounds())}},
		directionUp:    []object{object{pic: spriteup1, frame: spriteup1.Bounds(), sprite: pixel.NewSprite(spriteup1, spriteup1.Bounds())}, object{pic: spriteup2, frame: spriteup2.Bounds(), sprite: pixel.NewSprite(spriteup2, spriteup2.Bounds())}},
		directionDown:  []object{object{pic: spritedown1, frame: spritedown1.Bounds(), sprite: pixel.NewSprite(spritedown1, spritedown1.Bounds())}, object{pic: spritedown2, frame: spritedown2.Bounds(), sprite: pixel.NewSprite(spritedown2, spritedown2.Bounds())}},
	}
	spriteup1, err = loadPicture(spritedirectory + qdirectory + "person-walking-up1.png")
	spriteup2, err = loadPicture(spritedirectory + qdirectory + "person-walking-up2.png")
	spritedown1, err = loadPicture(spritedirectory + qdirectory + "person-walking-down1.png")
	spritedown2, err = loadPicture(spritedirectory + qdirectory + "person-walking-down2.png")
	spriteleft1, err = loadPicture(spritedirectory + qdirectory + "person-walking-left1.png")
	spriteleft2, err = loadPicture(spritedirectory + qdirectory + "person-walking-left2.png")
	spriteright1, err = loadPicture(spritedirectory + qdirectory + "person-walking-right1.png")
	spriteright2, err = loadPicture(spritedirectory + qdirectory + "person-walking-right2.png")
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
	sprite, err := loadPicture(spritedirectory + objectsdirectory + "raining-tr.png")
	if err != nil {
		panic(err)
	}
	rain = append(rain, object{
		frame: sprite.Bounds(),
		batch: pixel.NewBatch(&pixel.TrianglesData{}, sprite),
		pic:   sprite,
		posY:  maxY,
		posX:  maxX / 2})

	for i := range rain {
		rainSprites = append(rainSprites, pixel.NewSprite(rain[i].pic, rain[i].frame))
	}

}

func loadBackground() {
	sprite, err := loadPicture(spritedirectory + backgrounddirectory + "bg-ottawa.png")
	if err != nil {
		panic(err)
	}
	background = object{
		frame:  sprite.Bounds(),
		batch:  pixel.NewBatch(&pixel.TrianglesData{}, sprite),
		pic:    sprite,
		sprite: pixel.NewSprite(sprite, sprite.Bounds()),
		posY:   maxY / 2,
		posX:   maxX / 2,
	}
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
	satellites = []object{
		object{posX: 900, posY: top, angle: 10, directionX: left, directionY: down},
		object{posX: 20, posY: top, angle: -0.45, directionX: right, directionY: down},
		object{posX: maxX / 2, posY: top, angle: -1.5, directionX: left, directionY: down},
	}
	sprite, err := loadPicture(spritedirectory + objectsdirectory + "satellite-pixel.png")
	if err != nil {
		panic(err)
	}
	for i := 0; i < numSatellites; i++ {
		satellites[i].frame = sprite.Bounds()
		satellites[i].pic = sprite
		satellites[i].sprite = pixel.NewSprite(sprite, sprite.Bounds())
	}
}
