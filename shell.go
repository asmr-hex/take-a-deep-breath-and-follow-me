package main

import (
	"log"
	"os/user"
	"unicode/utf8"

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

	lines := GetStdOutLinesFromString(linuxLoginStdout, tl.ColorWhite, DarkBG, 0.1, func() float64 { return 1 })

	s.AddEntity(NewStdOut(lines...).BlockUntil(s.ready))

	s.linum = <-s.ready

	cmd := NewCmdLine(s.linum)
	s.AddEntity(cmd)
}

type CmdLine struct {
	Username     string
	Path         string
	x            int
	y            int
	cursorX      int
	cursorXInCmd int
	prefix       *tl.Text
	cmd          *tl.Text
	result       *tl.Text
	Cursor       *tl.Cell
	execdLines   []*tl.Text
	cursorRate   float64
	elapsedTime  float64
	cursorOn     bool
}

func NewCmdLine(y int) *CmdLine {
	// get currently logged in user
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	prefix := usr.Username + "@☠☠☠☠☠☠☠☠:~ $"

	prefixLen := utf8.RuneCountInString(prefix) + 1

	return &CmdLine{
		Username:     usr.Username,
		Path:         "~",
		x:            0,
		y:            y,
		cursorX:      prefixLen,
		cursorXInCmd: 0,
		prefix:       tl.NewText(0, y, prefix, tl.ColorGreen, DarkBG),
		cmd:          tl.NewText(prefixLen, y, "", tl.ColorWhite, DarkBG),
		result:       nil,
		Cursor:       &tl.Cell{Fg: DarkBG, Bg: tl.ColorBlue, Ch: ' '},
		execdLines:   []*tl.Text{},
		cursorRate:   1,
		elapsedTime:  0,
		cursorOn:     false,
	}
}

func (c *CmdLine) GetInput(ch string) {
	c.cmd.SetText(c.cmd.Text() + ch)
	c.cursorXInCmd++
	c.Cursor.Ch = ' '
}

func (c *CmdLine) BackSpaceInput() {
	var spliced = []rune{}
	chars := []rune(c.cmd.Text())
	i := c.cursorXInCmd - 1
	if i+1 >= len(chars) {
		spliced = chars[:i]
	} else {
		c.Cursor.Ch = spliced[i]
		spliced = append(chars[:i], chars[i+1:]...)
	}

	c.cmd.SetText(string(spliced))
	c.cursorXInCmd = i
}

func (c *CmdLine) Tick(event tl.Event) {
	if event.Type == tl.EventKey {
		switch event.Key {
		case tl.KeyEnter:
			c.Exec()
		case tl.KeySpace:
			c.GetInput(" ")
		case tl.KeyBackspace2:
			if c.cursorXInCmd != 0 {
				c.BackSpaceInput()
			}
		case tl.KeyCtrlC:
			c.result = tl.NewText(0, 0, "hehe sorry", tl.ColorRed, DarkBG)
			c.Exec()
		}

		if event.Ch != 0 {
			c.GetInput(string(event.Ch))
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
	c.elapsedTime += s.TimeDelta()
	if c.elapsedTime > c.cursorRate {
		if c.cursorOn {
			c.Cursor.Bg = tl.ColorBlue
			c.Cursor.Fg = DarkBG
			// c.Cursor.Ch = []rune(c.cmd.Text())[c.cursorXInCmd] // do this elsewhere...
		} else {
			c.Cursor.Bg = DarkBG
			c.Cursor.Fg = tl.ColorWhite
		}

		c.cursorOn = !c.cursorOn
		c.elapsedTime = 0
	}

	s.RenderCell(c.cursorX+c.cursorXInCmd, c.y, c.Cursor)
}

func (c *CmdLine) Eval() *tl.Text {
	// if there is already a result use that.
	if c.result != nil {
		c.y++
		result := *c.result
		result.SetPosition(0, c.y)
		c.result = nil
		return &result
	}

	// evaluate based on cmd ...

	return nil
}

func (c *CmdLine) Exec() {
	// evalute the cmd
	result := c.Eval()
	if result != nil {
		// append result to execd lines
		c.execdLines = append(c.execdLines, result)
	}

	// copy this line to execdLines
	c.execdLines = append(c.execdLines, c.prefix, c.cmd)

	// feed line
	c.y++

	// make new prefix and cmd
	prefixLen := utf8.RuneCountInString(c.prefix.Text()) + 1
	c.prefix = tl.NewText(0, c.y, c.prefix.Text(), tl.ColorGreen, DarkBG)
	c.cmd = tl.NewText(prefixLen, c.y, "", tl.ColorWhite, DarkBG)
	c.cursorX = prefixLen
	c.Cursor.Ch = ' '
	c.cursorXInCmd = 0
}
