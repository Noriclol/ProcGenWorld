package world

type WorldGenResult struct {
	gen_settings WorldGenSettings
	output       World
}

type World struct {
	name string
	//chunks map[Chunk]vec.Vec2
}
