package main

import (
	tl "github.com/JoelOtter/termloop"
)

var (
	// (◟'')◜
	DreamerSprite = []rune{'(', '◟', '\'', '\'', ')', '◜'}
)

type Dreamer struct {
	*tl.Entity
	Screen *tl.Screen
}

func NewDreamer(screen *tl.Screen) *Dreamer {
	e := tl.NewEntity(1, 1, 6, 1)

	d := &Dreamer{Entity: e, Screen: screen}
	d.FaceRight(true)

	return d
}

func (d *Dreamer) FaceRight(forward bool) {
	for i := 0; i < len(DreamerSprite); i++ {
		fg := tl.ColorBlue
		idx := i
		if !forward {
			idx = len(DreamerSprite) - 1 - i
			switch i {
			case 1:
				idx = 0
			case len(DreamerSprite) - 1:
				idx = i - 1
			}
		}
		if i == 2 || i == 3 {
			fg = tl.ColorWhite
		}
		d.SetCell(i, 0, &tl.Cell{
			Fg: fg,
			Ch: DreamerSprite[idx],
		})
	}
}

func (d *Dreamer) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		x, y := d.Position()
		w, h := d.Screen.Size()
		switch event.Key {
		case tl.KeyArrowRight:
			x += 1
			if x > w {
				x = 0
			}
			d.FaceRight(true)
		case tl.KeyArrowLeft:
			x -= 1
			if x < 0 {
				x = w
			}
			d.FaceRight(false)
		case tl.KeyArrowUp:
			y -= 1
			if y < 0 {
				y = h
			}
		case tl.KeyArrowDown:
			y += 1
			if y > h {
				y = 0
			}
		}
		d.SetPosition(x, y)
	}
}
