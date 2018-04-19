package main

import (
	"log"
	"os/user"
	"strings"
	"unicode/utf8"

	tl "github.com/JoelOtter/termloop"
)

type Shell struct {
	tl.Level
	screen *tl.Screen
	ready  chan []*tl.Text // for blocking animations
}

func NewShell(screen *tl.Screen) *Shell {
	l := tl.NewBaseLevel(tl.Cell{
		Bg: DarkBG,
	})

	return &Shell{Level: l, ready: make(chan []*tl.Text), screen: screen}
}

func (s *Shell) Login() {

	lines := GetStdOutLinesFromString(linuxLoginStdout, tl.ColorWhite, DarkBG, 0.1, func() float64 { return 1 })

	stdoutLines := NewStdOut(lines...).BlockUntil(s.ready)
	s.AddEntity(stdoutLines)

	printedLines := <-s.ready

	cmd := NewCmdLine(len(printedLines), s.screen)
	cmd.execdLines = append(cmd.execdLines, printedLines...)
	s.AddEntity(cmd)

	s.RemoveEntity(stdoutLines) // idk why i have to do this lol -___-
}

type CmdLine struct {
	screen       *tl.Screen
	Username     string
	Path         string
	x            int
	y            int
	cursorX      int
	cursorXInCmd int
	prefix       *tl.Text
	cmd          *tl.Text
	results      []*tl.Text
	Cursor       *tl.Cell
	execdLines   []*tl.Text
	cursorRate   float64
	elapsedTime  float64
	cursorOn     bool
}

func NewCmdLine(y int, s *tl.Screen) *CmdLine {
	// get currently logged in user
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	prefix := usr.Username + "@☠☠☠☠☠☠☠☠:~ $"

	prefixLen := utf8.RuneCountInString(prefix) + 1

	return &CmdLine{
		screen:       s,
		Username:     usr.Username,
		Path:         "~",
		x:            0,
		y:            y,
		cursorX:      prefixLen,
		cursorXInCmd: 0,
		prefix:       tl.NewText(0, y, prefix, tl.ColorGreen, DarkBG),
		cmd:          tl.NewText(prefixLen, y, "", tl.ColorWhite, DarkBG),
		results:      nil,
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
			c.results = []*tl.Text{
				tl.NewText(0, 0, "hehe sorry", tl.ColorRed, DarkBG),
			}

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

func (c *CmdLine) mvExecdLinesUp(n int) {
	newExecdLines := []*tl.Text{}
	for _, t := range c.execdLines {
		_, y := t.Position()
		if y < n-1 {
			c.screen.RemoveEntity(t)
		} else {
			x := 0
			if tx, _ := t.Position(); tx != 0 {
				x = tx
			}
			t.SetPosition(x, y-n)
			newExecdLines = append(newExecdLines, t)
		}
	}

	c.execdLines = newExecdLines
}

func (c *CmdLine) Eval() []*tl.Text {
	var results = []*tl.Text{}

	// if there is already a result use that.
	if c.results != nil {
		// copy ney array
		for _, r := range c.results {
			results = append(results, r)
		}
		c.results = nil
	} else {
		args := strings.Split(c.cmd.Text(), " ")

		// evaluate based on cmd ...
		results = c.EvalCmd(args)
	}

	// move stuff up after result (make room for new results in line feed)
	_, h := c.screen.Size()
	n := len(results)
	if c.y+n > h-2 {
		c.mvExecdLinesUp(n)
		c.y -= n
	}

	// set results position in line feed
	for _, r := range results {
		c.y++
		r.SetPosition(0, c.y)
	}

	return results
}

func (c *CmdLine) Exec() {
	// copy this line to execdLines
	c.execdLines = append(c.execdLines, c.prefix, c.cmd)

	// evalute the cmd
	results := c.Eval()
	if results != nil {
		// append result to execd lines
		c.execdLines = append(c.execdLines, results...)
	}

	// feed line
	_, h := c.screen.Size()
	if c.y == h-2 {
		c.mvExecdLinesUp(1)

	} else {
		c.y++ // TODO (cw|4.18.2018) results could be multi-line
	}

	// make new prefix and cmd
	prefixLen := utf8.RuneCountInString(c.prefix.Text()) + 1
	c.prefix = tl.NewText(0, c.y, c.prefix.Text(), tl.ColorGreen, DarkBG)
	c.cmd = tl.NewText(prefixLen, c.y, "", tl.ColorWhite, DarkBG)
	c.cursorX = prefixLen
	c.Cursor.Ch = ' '
	c.cursorXInCmd = 0
}
