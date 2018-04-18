package main

import (
	"math/rand"

	tl "github.com/JoelOtter/termloop"
)

const (
	MaxFeedTime float64 = 1
)

type StdOut struct {
	*tl.Entity
	Lines       []*tl.Text
	idx         int
	elapsedTime float64
	feedTime    float64
}

func NewStdOut(lines ...*tl.Text) *StdOut {
	return &StdOut{
		Lines:    lines,
		feedTime: MaxFeedTime * rand.Float64(),
	}
}

func (l *StdOut) Draw(screen *tl.Screen) {
	l.elapsedTime += screen.TimeDelta()

	// when loader is complete, remove it from screen
	if l.idx == len(l.Lines) && l.elapsedTime > l.feedTime {
		screen.RemoveEntity(l)
		return
	}

	if l.elapsedTime > l.feedTime {
		// set the position of next line feed
		l.Lines[l.idx].SetPosition(0, l.idx)

		// draw all the text lines to the screen
		l.DrawLines(l.idx, screen)

		// reset elapsedTime, re-init feed time, incr idx
		l.elapsedTime = 0
		l.feedTime = MaxFeedTime * rand.Float64()
		l.idx++

		return
	}

	// draw all the text lines to the screen
	l.DrawLines(l.idx, screen)
}

func (l *StdOut) DrawLines(idx int, s *tl.Screen) {
	for _, line := range l.Lines[:idx] {
		line.Draw(s)
	}
}
