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
	ready  chan int // for blocking animations
}

func NewShell(screen *tl.Screen) *Shell {
	l := tl.NewBaseLevel(tl.Cell{
		Bg: DarkBG,
	})

	return &Shell{Level: l, ready: make(chan int)}
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

	s.AddEntity(NewStdOut(lines...).BlockUntil(s.ready))

	s.linum = <-s.ready

	cmd := NewCmdLine(s.linum)
	s.AddEntity(cmd)
}

type CmdLine struct {
	Username   string
	Path       string
	x          int
	y          int
	cursorX    int
	prefix     *tl.Text
	cmd        *tl.Text
	Cursor     *tl.Cell
	execdLines []*tl.Text
}

func NewCmdLine(y int) *CmdLine {
	// get currently logged in user
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	prefix := usr.Username + "@☠☠☠☠☠☠☠☠:~ $"

	return &CmdLine{
		Username:   usr.Username,
		Path:       "~",
		x:          0,
		y:          y,
		cursorX:    len(prefix),
		prefix:     tl.NewText(0, y, prefix, tl.ColorGreen, DarkBG),
		cmd:        tl.NewText(len(prefix), y, "", tl.ColorWhite, DarkBG),
		Cursor:     &tl.Cell{Fg: DarkBG, Bg: tl.ColorYellow, Ch: ' '},
		execdLines: []*tl.Text{},
	}
}

func (c *CmdLine) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		switch event.Key {
		case tl.KeyEnter:
			c.Exec()
		}
	}
}

func (c *CmdLine) Draw(s *tl.Screen) {
	// draw all execdlines
	for _, line := range c.execdLines {
		line.Draw(s)
	}

	//draw current line
	c.prefix.Draw(s)
	c.cmd.Draw(s)

	// draw cursor
	s.RenderCell(c.cursorX, c.y, c.Cursor)
}

func (c *CmdLine) Exec() {
	// evalute the cmd

	// copy this line to execdLines
	c.execdLines = append(c.execdLines, c.prefix, c.cmd)

	// feed line
	c.y++

	// make new prefix and cmd
	c.prefix = tl.NewText(0, c.y, c.prefix.Text(), tl.ColorGreen, DarkBG)
	c.cmd = tl.NewText(len(c.prefix.Text()), c.y, "", tl.ColorWhite, DarkBG)
	c.cursorX = len(c.prefix.Text())
}
