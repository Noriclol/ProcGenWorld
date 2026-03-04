package world

import "github.com/eihigh/vec"

type WorldGenSettings struct {
	WorldSize       vec.Vec2i
	AmountLocations uint32
	Seed            int64
	WindDir         vec.Vec2
	NumRegions      int
}

func DefaultSettings() WorldGenSettings {
	return WorldGenSettings{
		WorldSize:       vec.Vec2i{X: 1024, Y: 1024},
		AmountLocations: 64,
		Seed:            42,
		WindDir:         vec.Vec2{X: 1, Y: 0},
		NumRegions:      1640,
	}
}
