package main

import (
	tl "github.com/JoelOtter/termloop"
)

func main() {
	game := tl.NewGame()
	game.Screen().SetFps(60)

	level := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})

	dreamer := NewDreamer(game.Screen())
	level.AddEntity(dreamer)

	game.Screen().SetLevel(level)

	txt := tl.NewText(
		0, 0,
		"...",
		tl.ColorBlack,
		tl.ColorBlue,
	)

	level.AddEntity(txt)

	game.Start()
}

type Dreamer struct {
	*tl.Entity
	GameScreen *tl.Screen
}

// (◟'')◜
func NewDreamer(screen *tl.Screen) *Dreamer {
	e := tl.NewEntity(1, 1, 6, 1)
	e.SetCell(0, 0, &tl.Cell{
		Fg: tl.ColorBlue,
		Ch: '(',
	})
	e.SetCell(1, 0, &tl.Cell{
		Fg: tl.ColorBlue,
		Ch: '◟',
	})
	e.SetCell(2, 0, &tl.Cell{
		Fg: tl.ColorWhite,
		Ch: '\'',
	})
	e.SetCell(3, 0, &tl.Cell{
		Fg: tl.ColorWhite,
		Ch: '\'',
	})
	e.SetCell(4, 0, &tl.Cell{
		Fg: tl.ColorBlue,
		Ch: ')',
	})
	e.SetCell(5, 0, &tl.Cell{
		Fg: tl.ColorBlue,
		Ch: '◜',
	})

	return &Dreamer{Entity: e, GameScreen: screen}
}

func (d *Dreamer) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		x, y := d.Position()
		w, h := d.GameScreen.Size()
		switch event.Key {
		case tl.KeyArrowRight:
			x += 1
			if x > w {
				x = 0
			}
		case tl.KeyArrowLeft:
			x -= 1
			if x < 0 {
				x = w
			}
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
