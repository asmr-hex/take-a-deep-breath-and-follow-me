package main

import (
	"log"
	"os/user"

	tl "github.com/JoelOtter/termloop"
)

type Shell struct {
	tl.Level
	screen *tl.Screen
	linum  int
}

func NewShell(screen *tl.Screen) *Shell {
	l := tl.NewBaseLevel(tl.Cell{
		Bg: DarkBG,
	})

	return &Shell{Level: l}
}

func (s *Shell) Login() {
	// get currently logged in user
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	lines := append(
		GetStdOutLinesFromString(linuxLoginStdout, tl.ColorWhite, DarkBG, 0.1, func() float64 { return 1 }),
		NewStdOutLine(usr.Username+"@☠☠☠☠☠☠☠☠:~ $", tl.ColorGreen, DarkBG, 0.3),
	)

	s.AddEntity(NewStdOut(lines...))
}
