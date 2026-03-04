package generation

import (
	"fmt"
	"math/rand"
	"proc-gen-world/internal/world/entities"

	"github.com/eihigh/vec"
)

func scatterEntities(width, height int, heightmap [][]float32, count uint32, rng *rand.Rand) []entities.Entity {
	// Collect candidate land tiles
	type tile struct{ x, y int }
	var land []tile
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if heightmap[y][x] >= 0.38 {
				land = append(land, tile{x, y})
			}
		}
	}

	result := make([]entities.Entity, 0, count)
	for i := uint32(0); i < count && len(land) > 0; i++ {
		t := land[rng.Intn(len(land))]
		pos := vec.Vec2{X: float64(t.x), Y: float64(t.y)}
		uid := uint64(i + 1)
		id := fmt.Sprintf("e%d", uid)

		var e entities.Entity
		switch i % 3 {
		case 0:
			e = &entities.Person{BaseEntity: entities.BaseEntity{UID: uid, ID: id, Position: pos, Kind: entities.KindPerson}}
		case 1:
			e = &entities.Animal{BaseEntity: entities.BaseEntity{UID: uid, ID: id, Position: pos, Kind: entities.KindAnimal}}
		default:
			e = &entities.Monster{BaseEntity: entities.BaseEntity{UID: uid, ID: id, Position: pos, Kind: entities.KindMonster}}
		}
		result = append(result, e)
	}
	return result
}
