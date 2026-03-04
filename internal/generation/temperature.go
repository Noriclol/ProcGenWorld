package generation

import "math"

func generateTemperature(width, height int, heightmap [][]float32) [][]float32 {
	temp := make([][]float32, height)
	for y := range temp {
		temp[y] = make([]float32, width)
		for x := 0; x < width; x++ {
			ny := float64(y) / float64(height)
			latitudeFactor := float32(1.0 - 2.0*math.Abs(ny-0.5))

			elevationChill := heightmap[y][x] * 0.8

			t := latitudeFactor - elevationChill
			if t < 0 {
				t = 0
			}
			if t > 1 {
				t = 1
			}
			temp[y][x] = t
		}
	}
	return temp
}
