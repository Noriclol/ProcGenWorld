package terrain

type BiomeID int

const (
	BiomeOcean             BiomeID = iota
	BiomeBeach
	BiomelPlains
	BiomeForest
	BiomeDesert
	BiomeTundra
	BiomelMountain
	BiomeSnowPeak
	BiomeTropicalRainforest
	BiomeSavanna
)

type Biome struct {
	Name  string
	Color [3]uint8
}

var Biomes = map[BiomeID]Biome{
	BiomeOcean:              {Name: "Ocean", Color: [3]uint8{30, 100, 200}},
	BiomeBeach:              {Name: "Beach", Color: [3]uint8{220, 210, 150}},
	BiomelPlains:            {Name: "Plains", Color: [3]uint8{100, 180, 80}},
	BiomeForest:             {Name: "Forest", Color: [3]uint8{34, 120, 40}},
	BiomeDesert:             {Name: "Desert", Color: [3]uint8{230, 200, 100}},
	BiomeTundra:             {Name: "Tundra", Color: [3]uint8{170, 200, 180}},
	BiomelMountain:          {Name: "Mountain", Color: [3]uint8{130, 110, 100}},
	BiomeSnowPeak:           {Name: "Snow Peak", Color: [3]uint8{240, 240, 255}},
	BiomeTropicalRainforest: {Name: "Tropical Rainforest", Color: [3]uint8{0, 100, 30}},
	BiomeSavanna:            {Name: "Savanna", Color: [3]uint8{180, 170, 80}},
}

func ClassifyBiome(height, rainfall float32) BiomeID {
	if height < 0.25 {
		return BiomeOcean
	}
	if height < 0.28 {
		return BiomeBeach
	}
	if height >= 0.82 {
		return BiomeSnowPeak
	}

	// Elevation tier: 0=Lowland, 1=Midland, 2=Highland
	var elevTier int
	switch {
	case height < 0.40:
		elevTier = 0
	case height < 0.55:
		elevTier = 1
	default:
		elevTier = 2
	}

	// Moisture tier: 0=Arid, 1=Semi-arid, 2=Moderate, 3=Wet
	var moistTier int
	switch {
	case rainfall < 0.15:
		moistTier = 0
	case rainfall < 0.40:
		moistTier = 1
	case rainfall < 0.65:
		moistTier = 2
	default:
		moistTier = 3
	}

	// 2D Whittaker lookup [elevTier][moistTier]
	lookup := [3][4]BiomeID{
		//  Arid          Semi-arid     Moderate      Wet
		{BiomeDesert, BiomelPlains, BiomelPlains, BiomeTropicalRainforest}, // Lowland
		{BiomeDesert, BiomelPlains, BiomeForest, BiomeForest},              // Midland
		{BiomeTundra, BiomeTundra, BiomelMountain, BiomelMountain},         // Highland
	}

	return lookup[elevTier][moistTier]
}
