package world

import (
	"proc-gen-world/internal/calendar"
	"proc-gen-world/internal/world/entities"
	"proc-gen-world/internal/world/terrain"

	"github.com/eihigh/vec"
)

type WorldGenResult struct {
	GenSettings WorldGenSettings
	Output      World
	Heightmap   [][]float32
	Rainfall    [][]float32
	Weather     [][]float32
	Flow        [][]float32
	Temperature [][]float32
}

type World struct {
	Name     string
	Chunks   map[vec.Vec2i]*terrain.Chunk
	Entities []entities.Entity
	Calendar calendar.Calendar
}
