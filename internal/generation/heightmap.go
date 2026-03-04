package generation

import (
	"math"

	opensimplex "github.com/ojrac/opensimplex-go"
)

func generateHeightmap(width, height int, seed int64) [][]float32 {
	noise := opensimplex.New(seed)
	hm := make([][]float32, height)
	for y := range hm {
		hm[y] = make([]float32, width)
	}

	octaves := []struct {
		freq, amp float64
	}{
		{freq: 2.0, amp: 0.5},
		{freq: 4.0, amp: 0.3},
		{freq: 8.0, amp: 0.2},
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nx := float64(x) / float64(width)
			ny := float64(y) / float64(height)
			var v float64
			for _, oct := range octaves {
				v += oct.amp * noise.Eval2(nx*oct.freq, ny*oct.freq)
			}
			// Island falloff: depress edges, elevate center
			nx2 := 2*nx - 1
			ny2 := 2*ny - 1
			d := 2.0 * math.Max(math.Abs(nx2), math.Abs(ny2))
			v = (1.0 + v - d) / 2.0
			if v < 0 {
				v = 0
			} else if v > 1 {
				v = 1
			}
			hm[y][x] = float32(v)
		}
	}
	return hm
}
