package main

import (
	tl "github.com/JoelOtter/termloop"
)

func (c *CmdLine) clear(args []string) []*tl.Text {
	for _, l := range c.execdLines {
		c.screen.RemoveEntity(l)
	}
	c.y = -1
	c.execdLines = []*tl.Text{}

	return []*tl.Text{}
}

func (c *CmdLine) ls(args []string) []*tl.Text {
	// get working dir

	return []*tl.Text{}
}

func (c *CmdLine) EvalCmd(args []string) []*tl.Text {
	var (
		// add commands here for support
		CMDS = map[string]func([]string) []*tl.Text{
			"clear": c.clear,
			"ls":    c.ls,
		}
		results = []*tl.Text{}
	)

	if len(args) == 0 {
		return results
	}

	if fn, ok := CMDS[args[0]]; ok {
		return fn(args)
	}

	return results
}
