package main

import (
	tl "github.com/JoelOtter/termloop"
)

func main() {
	// initialize game
	game := tl.NewGame()
	game.Screen().SetFps(60)
	screen := game.Screen()

	cemetery := BuildCemetery(screen)
	cemetery.Login()
	// cemetery.EnterDreamer()
	screen.SetLevel(cemetery)

	game.Start()
}
