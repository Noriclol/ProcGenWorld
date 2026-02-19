package game

import world "proc-gen-world/internal/world"

type Game struct {
	world.World
}

func NewGame() (Game, error) {
	return Game{}, nil
}
