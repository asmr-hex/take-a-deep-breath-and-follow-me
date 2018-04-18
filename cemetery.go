package main

import (
	tl "github.com/JoelOtter/termloop"
)

type Cemetery struct {
	tl.Level
	Screen *tl.Screen
}

func BuildCemetery(screen *tl.Screen) *Cemetery {
	l := tl.NewBaseLevel(tl.Cell{
		Bg: DarkBG,
	})

	return &Cemetery{Level: l, Screen: screen}
}

func (c *Cemetery) Greet(s string) {
	txt := NewStdOut(
		tl.NewText(0, 0, "scanning...", tl.ColorBlue, tl.ColorBlack),
		tl.NewText(0, 0, "ok...", tl.ColorBlue, tl.ColorBlack),
		tl.NewText(0, 0, "whoops...", tl.ColorBlue, tl.ColorBlack),
		tl.NewText(0, 0, "error...", tl.ColorRed, tl.ColorBlack),
	)
	c.AddEntity(txt)
}

func (c *Cemetery) EnterDreamer() {
	dreamer := NewDreamer(c.Screen)
	c.AddEntity(dreamer)
}
