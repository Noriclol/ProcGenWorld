package main

import (
	"log"
	"proc-gen-world/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("procgenworld")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
