package generation

import (
	"math/rand"
	"proc-gen-world/internal/world"
	"proc-gen-world/internal/world/terrain"
	"sort"

	"github.com/eihigh/vec"
)

func Gen() (world.WorldGenResult, error) {
	return GenWithSettings(world.DefaultSettings())
}

func GenWithSettings(settings world.WorldGenSettings) (world.WorldGenResult, error) {
	w := settings.WorldSize.X
	h := settings.WorldSize.Y
	rng := rand.New(rand.NewSource(settings.Seed))

	heightmap := generateHeightmap(w, h, settings.Seed)
	voronoi := generateVoronoi(w, h, settings.NumRegions, rng)

	// Voronoi elevation quantization: blend pixel heights toward region average
	regionSum := make([]float64, settings.NumRegions)
	regionCount := make([]int, settings.NumRegions)
	for y := range heightmap {
		for x, hv := range heightmap[y] {
			r := voronoi.Regions[y][x]
			regionSum[r] += float64(hv)
			regionCount[r]++
		}
	}
	regionElev := make([]float32, settings.NumRegions)
	for r := range regionElev {
		if regionCount[r] > 0 {
			regionElev[r] = float32(regionSum[r] / float64(regionCount[r]))
		}
	}
	const quantizeBlend = 0.7
	for y := range heightmap {
		for x := range heightmap[y] {
			r := voronoi.Regions[y][x]
			heightmap[y][x] = heightmap[y][x]*(1-quantizeBlend) + regionElev[r]*quantizeBlend
		}
	}

	rainfall := generateRainfall(w, h, heightmap, settings.WindDir)
	weather := generateWeather(w, h, rainfall, rng.Int63())
	temperature := generateTemperature(w, h, heightmap)

	// River flow: build region adjacency
	adj := make([]map[int]bool, settings.NumRegions)
	for i := range adj {
		adj[i] = make(map[int]bool)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r := voronoi.Regions[y][x]
			if x+1 < w {
				if n := voronoi.Regions[y][x+1]; n != r {
					adj[r][n] = true
					adj[n][r] = true
				}
			}
			if y+1 < h {
				if n := voronoi.Regions[y+1][x]; n != r {
					adj[r][n] = true
					adj[n][r] = true
				}
			}
		}
	}

	// Find downslope neighbor for each region
	downslope := make([]int, settings.NumRegions)
	for r := range downslope {
		downslope[r] = -1
		lowest := regionElev[r]
		for nb := range adj[r] {
			if regionElev[nb] < lowest {
				lowest = regionElev[nb]
				downslope[r] = nb
			}
		}
	}

	// Sort regions high→low, accumulate flow downstream
	order := make([]int, settings.NumRegions)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(a, b int) bool {
		return regionElev[order[a]] > regionElev[order[b]]
	})

	regionRainfall := computeRegionAvg(settings.NumRegions, rainfall, voronoi.Regions)
	flow := make([]float32, settings.NumRegions)
	for _, r := range order {
		flow[r] += regionRainfall[r]
		if ds := downslope[r]; ds >= 0 {
			flow[ds] += flow[r]
		}
	}

	// Normalize flow to [0,1]
	var maxFlow float32
	for _, f := range flow {
		if f > maxFlow {
			maxFlow = f
		}
	}
	if maxFlow > 0 {
		for r := range flow {
			flow[r] /= maxFlow
		}
	}

	// Build flow 2D map
	flowMap := make([][]float32, h)
	for y := range flowMap {
		flowMap[y] = make([]float32, w)
		for x := range flowMap[y] {
			flowMap[y][x] = flow[voronoi.Regions[y][x]]
		}
	}

	// Build chunks
	chunkCountX := (w + terrain.ChunkSize - 1) / terrain.ChunkSize
	chunkCountY := (h + terrain.ChunkSize - 1) / terrain.ChunkSize
	chunks := make(map[vec.Vec2i]*terrain.Chunk, chunkCountX*chunkCountY)

	for ty := 0; ty < h; ty++ {
		for tx := 0; tx < w; tx++ {
			cx := tx / terrain.ChunkSize
			cy := ty / terrain.ChunkSize
			key := vec.Vec2i{X: cx, Y: cy}
			chunk, ok := chunks[key]
			if !ok {
				chunk = &terrain.Chunk{}
				chunks[key] = chunk
			}
			lx := tx % terrain.ChunkSize
			ly := ty % terrain.ChunkSize
			ht := heightmap[ty][tx]
			rf := rainfall[ty][tx]
			chunk.Tiles[ly][lx] = terrain.Tile{
				Height:      ht,
				Rainfall:    rf,
				Flow:        flow[voronoi.Regions[ty][tx]],
				Temperature: temperature[ty][tx],
				Biome:       terrain.ClassifyBiome(ht, rf),
				Region:      voronoi.Regions[ty][tx],
			}
		}
	}

	ents := scatterEntities(w, h, heightmap, settings.AmountLocations, rng)

	output := world.World{
		Name:     "World",
		Chunks:   chunks,
		Entities: ents,
	}

	return world.WorldGenResult{
		GenSettings: settings,
		Output:      output,
		Heightmap:   heightmap,
		Rainfall:    rainfall,
		Weather:     weather,
		Flow:        flowMap,
		Temperature: temperature,
	}, nil
}

func computeRegionAvg(numRegions int, data [][]float32, regions [][]int) []float32 {
	sum := make([]float64, numRegions)
	count := make([]int, numRegions)
	for y := range data {
		for x, v := range data[y] {
			r := regions[y][x]
			sum[r] += float64(v)
			count[r]++
		}
	}
	avg := make([]float32, numRegions)
	for r := range avg {
		if count[r] > 0 {
			avg[r] = float32(sum[r] / float64(count[r]))
		}
	}
	return avg
}
