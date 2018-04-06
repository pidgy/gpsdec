package gpsdec

import (
	"time"
)

const (
	currentNothing = iota
	currentBuilding
	currentUserSelect
)

var (
	animationTickers []*time.Ticker
	drawingAnimation = false
	currentAnimation int
)

func startAnimation(c int) {
	stopAnimation()
	currentAnimation = c
	t1 := time.NewTicker(time.Millisecond * 500)
	t2 := time.NewTicker(time.Millisecond * 500)
	t3 := time.NewTicker(time.Millisecond * 500)
	t4 := time.NewTicker(time.Millisecond * 500)
	switch currentAnimation {
	case currentBuilding:
		go startTimer(&buildinghelp, t1)
		animationTickers = append(animationTickers, t1)
		go startTimer(&buildingnext, t2)
		animationTickers = append(animationTickers, t2)
		go startTimer(&buildingdone, t3)
		animationTickers = append(animationTickers, t3)
	case currentUserSelect:
		go startTimer(&userselectionbuttons1, t1)
		animationTickers = append(animationTickers, t1)
		go startTimer(&userselectionbuttons2, t2)
		animationTickers = append(animationTickers, t2)
		go startTimer(&userselectionbuttons3, t3)
		animationTickers = append(animationTickers, t3)
		go startTimer(&userselectionbuttons4, t4)
		animationTickers = append(animationTickers, t4)
	}
	drawingAnimation = true
}

func stopAnimation() {
	for _, t := range animationTickers {
		t.Stop()
	}
	currentAnimation = currentNothing
	drawingAnimation = false
}

func buildingHelp() []*object {
	return []*object{buildinghelp.next(), buildingnext.next(), buildingdone.next()}
}

func userSelect() []*object {
	return []*object{userselectionbuttons1.next(), userselectionbuttons2.next(), userselectionbuttons3.next(), userselectionbuttons4.next()}
}

func getCurrentAnimation() []*object {
	switch currentAnimation {
	case currentBuilding:
		return buildingHelp()
	case currentUserSelect:
		return userSelect()
	}
	return nil
}

func startTimer(q *objectqueue, event *time.Ticker) {
	for range event.C {
		q.roll()
	}
}

func slideAnimation(o *object) {

}
