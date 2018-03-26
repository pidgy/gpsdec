package gpsdec

import (
	"time"
)

const (
	currentNothing = iota
	currentBuilding
)

var (
	animationTickers []*time.Ticker
	drawAnimation    = false
	currentAnimation int
)

func startAnimation(c int) {
	stopAnimation()
	currentAnimation = c
	switch currentAnimation {
	case currentBuilding:
		t1 := time.NewTicker(time.Millisecond * 500)
		go startTimer(&buildinghelp, t1)
		animationTickers = append(animationTickers, t1)
		t2 := time.NewTicker(time.Millisecond * 500)
		go startTimer(&buildingnext, t2)
		animationTickers = append(animationTickers, t2)
		t3 := time.NewTicker(time.Millisecond * 500)
		go startTimer(&buildingdone, t3)
		animationTickers = append(animationTickers, t3)
	}
	drawAnimation = true
}

func stopAnimation() {
	for _, t := range animationTickers {
		t.Stop()
	}
	currentAnimation = currentNothing
	drawAnimation = false
}

func buildingHelp() []*object {
	return []*object{buildinghelp.next(), buildingnext.next(), buildingdone.next()}
}

func getCurrentAnimation() []*object {
	switch currentAnimation {
	case currentBuilding:
		return buildingHelp()
	}
	return nil
}

func startTimer(q *objectqueue, event *time.Ticker) {
	for _ = range event.C {
		q.roll()
	}
}

func slideAnimation(o *object) {

}
