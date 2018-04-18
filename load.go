package main

import (
	"math/rand"
	"strings"

	tl "github.com/JoelOtter/termloop"
)

const (
	MaxFeedTime float64 = 1
)

type StdOut struct {
	*tl.Entity
	Lines       []*StdOutLine
	idx         int
	elapsedTime float64
	feedTime    float64
	disappear   bool
}

func NewStdOut(lines ...*StdOutLine) *StdOut {
	return &StdOut{
		Lines:    lines,
		feedTime: MaxFeedTime * rand.Float64(),
	}
}

func (l *StdOut) DisappearOnEnd() {
	l.disappear = true
}

func (l *StdOut) Draw(screen *tl.Screen) {
	l.elapsedTime += screen.TimeDelta()

	// when loader is complete, remove it from screen
	if l.idx == len(l.Lines) && l.elapsedTime > l.feedTime {
		if l.disappear {
			screen.RemoveEntity(l)
		} else {
			// draw all the text lines to the screen
			l.DrawLines(l.idx, screen)
		}

		return
	}

	if l.elapsedTime > l.feedTime {
		// set the position of next line feed
		l.Lines[l.idx].SetPosition(0, l.idx)

		// draw all the text lines to the screen
		l.DrawLines(l.idx, screen)

		// reset elapsedTime, re-init feed time, incr idx
		l.elapsedTime = 0
		l.feedTime = l.Lines[l.idx].feedTime
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

type StdOutLine struct {
	*tl.Text
	feedTime float64
}

func NewStdOutLine(txt string, fg, bg tl.Attr, feedTime float64) *StdOutLine {
	return &StdOutLine{
		Text:     tl.NewText(0, 0, txt, fg, bg),
		feedTime: feedTime,
	}
}

func GetStdOutLinesFromString(s string, fg, bg tl.Attr, feedTime float64, p func() float64) []*StdOutLine {
	lines := []*StdOutLine{}
	for _, line := range strings.Split(s, "\n") {
		lines = append(
			lines,
			NewStdOutLine(line, fg, bg, feedTime*p()),
		)
	}

	return lines
}
