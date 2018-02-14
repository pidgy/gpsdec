package gpsdec

import (
	"math/rand"
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
	animationdirectory  string  = "animations/"

	directionLeft = iota
	directionRight
	directionUp
	directionDown

	maxSpriteX  = maxX - 5
	minSpriteX  = 5
	maxSpriteY  = 381
	minSpriteY  = 100
	walkSpeed   = 3.0
	staticSpeed = 3.0
)

var (
	buildings     []object
	satellites    []object
	buttons       []object
	rain          []object
	staticobjects []object

	buildinghelp objectqueue
	buildingnext objectqueue
	buildingdone objectqueue

	personP       object
	personQ       object
	distanceLine  object
	loadScreen    object
	okbutton      object
	controlScreen object
	background    object
	staticobject  object

	staticBatch *pixel.Batch

	locP       = pixel.V(maxX/2, maxY/2)
	locQ       = pixel.V(float64(rand.Intn(int(maxX))), float64(rand.Intn(int(maxY))))
	numHumans  = 1
	numButtons = 4

	drawStatic = false

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

	buttonY = 30.0
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
	desc            string
	descalphaX      float64
	descalphaY      float64
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

func clearSprites() {
	buildings = []object{}
	staticobjects = []object{}
	staticobject = object{}
	drawingDistanceLine = false
	drawingRain = false
	drawStatic = false
	stopAnimation()
}

func buildingFrames() (pixel.Picture, []pixel.Rect) {
	spritesheet, err := loadPicture(spritedirectory + "skyscraper-spritesheet.png")
	if err != nil {
		panic(err)
	}
	var buildingFrames []pixel.Rect
	staticBatch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)

	buildingFrames = append(buildingFrames, pixel.R(40, 10, 100, 80))   // DC
	buildingFrames = append(buildingFrames, pixel.R(100, 10, 170, 90))  // Hospital
	buildingFrames = append(buildingFrames, pixel.R(170, 0, 240, 100))  // Hospital 2
	buildingFrames = append(buildingFrames, pixel.R(40, 110, 75, 165))  // Apartment
	buildingFrames = append(buildingFrames, pixel.R(75, 100, 140, 180)) // Apartment 2

	buildingFrames = append(buildingFrames, pixel.R(0, 400, 120, 600))   // Big Sky
	buildingFrames = append(buildingFrames, pixel.R(515, 385, 600, 600)) // Big Sky 2
	buildingFrames = append(buildingFrames, pixel.R(120, 385, 210, 600)) // Big Sky 3
	buildingFrames = append(buildingFrames, pixel.R(210, 385, 310, 600)) // Big Sky 4
	buildingFrames = append(buildingFrames, pixel.R(310, 385, 410, 600)) // Big Sky 5
	buildingFrames = append(buildingFrames, pixel.R(410, 385, 510, 610)) // Big Sky 6
	buildingFrames = append(buildingFrames, pixel.R(510, 385, 610, 600)) // Big Sky 7

	return spritesheet, buildingFrames
}

func newBuilding(pos pixel.Vec, incr int) object {
	currBuildingName = (currBuildingName + incr)
	if currBuildingName > len(buildFrames)-1 {
		currBuildingName = 0
	}
	if currBuildingName < 0 {
		currBuildingName = len(buildFrames) - 1
	}
	newMessage(buildingNames[currBuildingName]+" Selected", 200, standardFont)

	f := buildFrames[currBuildingName]
	building := pixel.NewSprite(buildPic, f)
	buildingSprites = append(buildingSprites, building)
	nb := object{
		pic:    buildPic,
		posX:   pos.X - f.W()/2,
		posY:   pos.Y - f.H()/2,
		loc:    pos,
		frame:  f,
		mat:    pixel.IM.Moved(pos),
		sprite: building,
	}

	return nb
}
