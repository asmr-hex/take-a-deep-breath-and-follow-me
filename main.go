package main

import (
	tl "github.com/JoelOtter/termloop"
)

func main() {
	// initialize game
	game := tl.NewGame()
	game.Screen().SetFps(60)
	game.SetEndKey(tl.KeyCtrlD) // 0xFF for an unreachable key
	screen := game.Screen()

	cemetery := BuildCemetery(screen)
	go cemetery.Login()
	// cemetery.EnterDreamer()
	screen.SetLevel(cemetery)

	game.Start()
}
