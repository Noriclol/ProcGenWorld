package generation

import (
	"math/rand"

	opensimplex "github.com/ojrac/opensimplex-go"
)

func generateWeather(width, height int, rainfall [][]float32, seed int64) [][]float32 {
	noise := opensimplex.New(seed)
	weather := make([][]float32, height)
	for y := range weather {
		weather[y] = make([]float32, width)
		for x := 0; x < width; x++ {
			perturbation := float32(noise.Eval2(float64(x)*0.05, float64(y)*0.05)) * 0.1
			v := rainfall[y][x] + perturbation
			if v < 0 {
				v = 0
			} else if v > 1 {
				v = 1
			}
			weather[y][x] = v
		}
	}
	return weather
}

func RegenerateWeather(rainfall [][]float32, seed int64) [][]float32 {
	if len(rainfall) == 0 {
		return rainfall
	}
	rng := rand.New(rand.NewSource(seed))
	newSeed := rng.Int63()
	return generateWeather(len(rainfall[0]), len(rainfall), rainfall, newSeed)
}
