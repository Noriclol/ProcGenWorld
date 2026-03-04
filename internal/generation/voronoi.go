package generation

import "math/rand"

type VoronoiResult struct {
	Regions [][]int  // [y][x] → region ID
	Seeds   [][2]int // [regionID] → {x, y} seed point
}

func generateVoronoi(width, height, numRegions int, rng *rand.Rand) VoronoiResult {
	seeds := make([][2]int, numRegions)
	for i := range seeds {
		seeds[i] = [2]int{rng.Intn(width), rng.Intn(height)}
	}

	regions := make([][]int, height)
	for y := range regions {
		regions[y] = make([]int, width)
		for x := 0; x < width; x++ {
			best := -1
			bestDist := int(^uint(0) >> 1)
			for i, s := range seeds {
				dx := x - s[0]
				if dx < 0 {
					dx = -dx
				}
				if dx > width/2 {
					dx = width - dx
				}
				dy := y - s[1]
				if dy < 0 {
					dy = -dy
				}
				if dy > height/2 {
					dy = height - dy
				}
				d := dx*dx + dy*dy
				if d < bestDist {
					bestDist = d
					best = i
				}
			}
			regions[y][x] = best
		}
	}
	return VoronoiResult{Regions: regions, Seeds: seeds}
}
