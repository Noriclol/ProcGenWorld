package generation

import "github.com/eihigh/vec"

func generateRainfall(width, height int, heightmap [][]float32, windDir vec.Vec2) [][]float32 {
	rainfall := make([][]float32, height)
	for y := range rainfall {
		rainfall[y] = make([]float32, width)
	}

	// Determine entry edges based on wind direction
	// windDir {1,0} means wind blows W→E, so moisture enters from left (x=0)
	// windDir {0,1} means wind blows N→S, so moisture enters from top (y=0)
	// We march along the wind direction for each entry point

	stepX := windDir.X
	stepY := windDir.Y
	if stepX == 0 && stepY == 0 {
		stepX = 1
	}

	// Build entry points on the upwind edge
	type entry struct{ x, y float64 }
	var entries []entry

	if stepX > 0 {
		for y := 0; y < height; y++ {
			entries = append(entries, entry{0, float64(y)})
		}
	} else if stepX < 0 {
		for y := 0; y < height; y++ {
			entries = append(entries, entry{float64(width - 1), float64(y)})
		}
	}
	if stepY > 0 {
		for x := 0; x < width; x++ {
			entries = append(entries, entry{float64(x), 0})
		}
	} else if stepY < 0 {
		for x := 0; x < width; x++ {
			entries = append(entries, entry{float64(x), float64(height - 1)})
		}
	}

	// March each ray from entry, depositing moisture
	for _, e := range entries {
		moisture := float32(1.0)
		cx, cy := e.x, e.y
		for {
			ix := int(cx + 0.5)
			iy := int(cy + 0.5)
			if ix < 0 || ix >= width || iy < 0 || iy >= height {
				break
			}
			h := heightmap[iy][ix]
			// Orographic: high terrain blocks moisture
			if h > 0.5 {
				moisture *= 1.0 - (h-0.5)*0.8
			}
			if rainfall[iy][ix] < moisture {
				rainfall[iy][ix] = moisture
			}
			cx += float64(stepX)
			cy += float64(stepY)
		}
	}
	return rainfall
}
