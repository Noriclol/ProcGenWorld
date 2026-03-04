package terrain

const ChunkSize = 16

type Tile struct {
	Height      float32
	Rainfall    float32
	Flow        float32 // accumulated river flow 0.0–1.0
	Temperature float32 // 0.0 = freezing, 1.0 = tropical
	Biome       BiomeID
	Region      int
}

type Chunk struct {
	Tiles [ChunkSize][ChunkSize]Tile
}
