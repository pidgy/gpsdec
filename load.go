package gpsdec

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/faiface/pixel"
)

func loadAllTheThings() {
	buildPic, buildFrames = buildingFrames()

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
	loadAnimations()
	loadUserSelectionAnimations()
	loadControlScreen()
	loadBackground()
	loadSatelliteFrames()
	loadButtonFrames()
	loadHumans()
	loadWeather()
	loadDistanceLine()
	loadOkButton()
	loadPositionEstimates()
	loadTipScreen()
	load.Done()
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

func setPictureColor(im pixel.Picture, color []color.RGBA) pixel.Picture {
	img := pixel.PictureDataFromPicture(im)
	img.Pix = color
	return img
}

func loadAnimations() {
	buildinghelp = newObjectQueue("Move the building", 50, 60)
	i := 1
	for i < 5 {
		sprite, err := loadPicture(spritedirectory + animationdirectory + fmt.Sprintf("arrow-help-%d.png", i))
		if err != nil {
			panic(err)
		}
		buildinghelp.push(object{
			pic:    sprite,
			frame:  sprite.Bounds(),
			sprite: pixel.NewSprite(sprite, sprite.Bounds()),
			posX:   (maxX / 3) + 150,
			posY:   maxY / 1.25,
		})
		i++
	}
	buildingnext = newObjectQueue("Choose another building", 85, 40)
	i = 1
	for i < 3 {
		sprite, err := loadPicture(spritedirectory + animationdirectory + fmt.Sprintf("ad-help-%d.png", i))
		if err != nil {
			panic(err)
		}
		buildingnext.push(object{
			pic:    sprite,
			frame:  sprite.Bounds(),
			sprite: pixel.NewSprite(sprite, sprite.Bounds()),
			posX:   maxX / 3,
			posY:   maxY / 1.3,
		})
		i++
	}

	buildingdone = newObjectQueue("Finish placing the building", 85, 40)
	i = 1
	for i < 3 {
		sprite, err := loadPicture(spritedirectory + animationdirectory + fmt.Sprintf("done-help-%d.png", i))
		if err != nil {
			panic(err)
		}
		buildingdone.push(object{
			pic:    sprite,
			frame:  sprite.Bounds(),
			sprite: pixel.NewSprite(sprite, sprite.Bounds()),
			posX:   (maxX / 3) + 300,
			posY:   maxY / 1.3,
		})
		i++
	}
}

func loadUserSelectionAnimations() {
	userselectionbuttons1 = newObjectQueue("1 ns", maxX/2, maxY/2)
	userselectionbuttons2 = newObjectQueue("2 ns", maxX/2, maxY/2)
	userselectionbuttons3 = newObjectQueue("3 ns", maxX/2, maxY/2)
	userselectionbuttons4 = newObjectQueue("4 ns", maxX/2, maxY/2)
	allobjects := []object{}
	i := 1
	posX := maxX/2 - 110
	for i < 5 {
		sprite1, err := loadPicture(spritedirectory + animationdirectory + fmt.Sprintf("user-select-button-%d-up.png", i))
		sprite2, err := loadPicture(spritedirectory + animationdirectory + fmt.Sprintf("user-select-button-%d-down.png", i))
		if err != nil {
			panic(err)
		}
		if i != 1 {
			posX = posX + 75
		}
		butup := object{
			pic:    sprite1,
			frame:  sprite1.Bounds(),
			sprite: pixel.NewSprite(sprite1, sprite1.Bounds()),
			posX:   posX,
			posY:   maxY / 2,
			mat:    pixel.IM.Moved(pixel.V(posX, maxY/2)),
		}
		butdown := object{
			pic:    sprite2,
			frame:  sprite2.Bounds(),
			sprite: pixel.NewSprite(sprite2, sprite2.Bounds()),
			posX:   posX,
			posY:   maxY / 2,
			mat:    pixel.IM.Moved(pixel.V(posX, maxY/2)),
		}
		allobjects = append(allobjects, butup)
		allobjects = append(allobjects, butdown)
		i++
	}
	userselectionbuttons1.push(allobjects[0])
	userselectionbuttons1.push(allobjects[1])

	userselectionbuttons2.push(allobjects[3])
	userselectionbuttons2.push(allobjects[2])

	userselectionbuttons3.push(allobjects[4])
	userselectionbuttons3.push(allobjects[5])

	userselectionbuttons4.push(allobjects[7])
	userselectionbuttons4.push(allobjects[6])
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

func loadWeather() {
	weatherSprites = map[int][]*pixel.Sprite{}
	weatherObjects = map[int][]object{}

	for i, png := range []string{"NONE", "raining-tr.png", "ash-tr.png", "dry-tr.png", "sand-tr.png"} {
		if i == 0 {
			continue
		}
		sprite, err := loadPicture(spritedirectory + objectsdirectory + png)
		if err != nil {
			panic(err)
		}
		weatherObjects[i] = append(weatherObjects[i], object{
			frame: sprite.Bounds(),
			batch: pixel.NewBatch(&pixel.TrianglesData{}, sprite),
			pic:   sprite,
			posY:  maxY,
			posX:  maxX / 2,
		})
		weatherSprites[i] = append(weatherSprites[i], pixel.NewSprite(sprite, sprite.Bounds()))
	}
}

func loadOkButton() {
	sprite1, err := loadPicture(spritedirectory + buttonsdirectory + "button-ok.png")
	sprite2, err := loadPicture(spritedirectory + buttonsdirectory + "button-pressed-ok.png")
	if err != nil {
		panic(err)
	}
	okbutton = &object{
		frame:      sprite1.Bounds(),
		pic:        sprite1,
		pressedpic: sprite2,
		sprite:     pixel.NewSprite(sprite1, sprite1.Bounds()),
		posY:       maxY / 2,
		posX:       maxX / 2,
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
	}
}

func loadDataBackground() {
	sprite, err := loadPicture(spritedirectory + backgrounddirectory + "bg-ottawa.png")
	if err != nil {
		panic(err)
	}
	background = object{
		frame:  sprite.Bounds(),
		batch:  pixel.NewBatch(&pixel.TrianglesData{}, sprite),
		pic:    sprite,
		sprite: pixel.NewSprite(sprite, sprite.Bounds()),
	}
}

func loadButtonFrames() {
	buttons = []object{
		object{posX: 40, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-buildings.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-buildings.png"},
		object{posX: 120, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-weather.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-weather.png"},
		object{posX: 200, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-gps.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-gps.png"},
		object{posX: maxX - 40, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-clear.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-clear.png"},
		object{posX: maxX - 40, posY: maxY - 100, filename: spritedirectory + buttonsdirectory + "button-tip.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-tip.png"},
		object{posX: maxX - 40, posY: maxY - 170, filename: spritedirectory + buttonsdirectory + "button-controls.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-controls.png"},
		object{posX: 280, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-scale.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-scale.png"},
		object{posX: 360, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-line.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-line.png"},
		object{posX: maxX - 115, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-controls.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-controls.png"},
		object{posX: maxX - 195, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-run.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-run.png"},
		object{posX: maxX - 270, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-estimate.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-estimate.png"},
		object{posX: maxX - 345, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-ephemeris.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-ephemeris.png"},
		object{posX: maxX - 420, posY: buttonY, filename: spritedirectory + buttonsdirectory + "button-elevation.png", pressedfilename: spritedirectory + buttonsdirectory + "button-pressed-elevation.png"},
	}

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
		object{loc: pixel.V(900, top), posX: 900, posY: top, angle: 10, directionX: left, directionY: down},
		object{loc: pixel.V(20, top), posX: 20, posY: top, angle: -0.45, directionX: right, directionY: down},
		object{loc: pixel.V(maxX-100, top), posX: maxX - 100, posY: top, angle: -0.25, directionX: right, directionY: down},
		object{loc: pixel.V(maxX/2, top), posX: maxX / 2, posY: top, angle: -1.5, directionX: left, directionY: down},
	}
	sprite, err := loadPicture(spritedirectory + objectsdirectory + "satellite-pixel.png")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(satellites); i++ {
		satellites[i].frame = sprite.Bounds()
		satellites[i].pic = sprite
		satellites[i].sprite = pixel.NewSprite(sprite, sprite.Bounds())
	}
	numSatellites = len(satellites)
}

func loadControlScreen() {
	controlPic, err := loadPicture(spritedirectory + "controls.png")
	if err != nil {
		panic(err)
	}
	controlScreen = object{
		sprite: pixel.NewSprite(controlPic, controlPic.Bounds()),
	}
}

func loadExtraSprites(n int, obj []object) [][]object {
	var extras [][]object
	i := 0
	for i < n {
		extras = append(extras, scramblePositions(obj))
		i++
	}
	return extras
}

func loadPositionEstimates() {
	pmsprite, err := loadPicture(spritedirectory + objectsdirectory + "pm.png")
	qmsprite, err := loadPicture(spritedirectory + objectsdirectory + "qm.png")
	if err != nil {
		panic(err)
	}
	pmmain := personP
	qmmain := personQ
	pmmain.sprite = pixel.NewSprite(pmsprite, pmsprite.Bounds())
	qmmain.sprite = pixel.NewSprite(qmsprite, qmsprite.Bounds())
	pestimate = pmmain
	qestimate = qmmain
}

func loadTipScreen() {
	hmsprite, err := loadPicture(spritedirectory + objectsdirectory + "help-message.png")
	if err != nil {
		panic(err)
	}
	tipmessage = object{
		sprite: pixel.NewSprite(hmsprite, hmsprite.Bounds()),
		loc:    pixel.V(maxX/2, maxY/2),
		posY:   maxY / 2,
		posX:   maxX / 2,
		mat:    pixel.IM.Scaled(pixel.ZV, 0).Moved(pixel.V(maxX/2, maxY/2)),
	}
	tipMaxScaleX = hmsprite.Bounds().H()
	tipMaxScaleY = hmsprite.Bounds().W()
}
