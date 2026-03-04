package render

import (
	"image"
	"image/color"
	"proc-gen-world/internal/world"
	"proc-gen-world/internal/world/terrain"

	"github.com/eihigh/vec"
	"github.com/hajimehoshi/ebiten/v2"
)

func buildFloat32Image(data [][]float32, colorFn func(float32) color.RGBA) *ebiten.Image {
	h := len(data)
	if h == 0 {
		return ebiten.NewImage(1, 1)
	}
	w := len(data[0])
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for ty := 0; ty < h; ty++ {
		for tx := 0; tx < w; tx++ {
			img.SetRGBA(tx, ty, colorFn(data[ty][tx]))
		}
	}
	return ebiten.NewImageFromImage(img)
}

func BuildElevationImage(data [][]float32) *ebiten.Image {
	return buildFloat32Image(data, func(v float32) color.RGBA {
		c := uint8(v * 255)
		return color.RGBA{c, c, c, 255}
	})
}

func BuildRainfallImage(data [][]float32) *ebiten.Image {
	return buildFloat32Image(data, func(v float32) color.RGBA {
		return color.RGBA{uint8(v * 80), uint8(v * 140), uint8(v * 255), 255}
	})
}

func BuildWeatherImage(data [][]float32) *ebiten.Image {
	return buildFloat32Image(data, func(v float32) color.RGBA {
		return color.RGBA{uint8(v * 80), uint8(v * 220), uint8(v * 200), 255}
	})
}

var voronoiPalette = []color.RGBA{
	{228, 26, 28, 255},
	{55, 126, 184, 255},
	{77, 175, 74, 255},
	{152, 78, 163, 255},
	{255, 127, 0, 255},
	{255, 255, 51, 255},
	{166, 86, 40, 255},
	{247, 129, 191, 255},
	{153, 153, 153, 255},
	{0, 210, 213, 255},
}

func BuildVoronoiImage(w *world.World, size vec.Vec2i) *ebiten.Image {
	mapW, mapH := size.X, size.Y
	img := image.NewRGBA(image.Rect(0, 0, mapW, mapH))
	for ty := 0; ty < mapH; ty++ {
		for tx := 0; tx < mapW; tx++ {
			cx, cy := tx/terrain.ChunkSize, ty/terrain.ChunkSize
			chunk, ok := w.Chunks[vec.Vec2i{X: cx, Y: cy}]
			if !ok {
				img.SetRGBA(tx, ty, color.RGBA{0, 0, 0, 255})
				continue
			}
			lx, ly := tx%terrain.ChunkSize, ty%terrain.ChunkSize
			region := chunk.Tiles[ly][lx].Region
			img.SetRGBA(tx, ty, voronoiPalette[region%len(voronoiPalette)])
		}
	}
	return ebiten.NewImageFromImage(img)
}

const riverThreshold = 0.15

func BuildMapImage(w *world.World, size vec.Vec2i) *ebiten.Image {
	mapW := size.X
	mapH := size.Y
	img := image.NewRGBA(image.Rect(0, 0, mapW, mapH))

	for ty := 0; ty < mapH; ty++ {
		for tx := 0; tx < mapW; tx++ {
			cx := tx / terrain.ChunkSize
			cy := ty / terrain.ChunkSize
			key := vec.Vec2i{X: cx, Y: cy}
			chunk, ok := w.Chunks[key]
			if !ok {
				img.SetRGBA(tx, ty, color.RGBA{0, 0, 0, 255})
				continue
			}
			lx := tx % terrain.ChunkSize
			ly := ty % terrain.ChunkSize
			tile := chunk.Tiles[ly][lx]
			biome := terrain.Biomes[tile.Biome]
			c := color.RGBA{biome.Color[0], biome.Color[1], biome.Color[2], 255}
			if tile.Flow > riverThreshold {
				alpha := tile.Flow
				c.R = uint8(float32(c.R)*(1-alpha) + 60*alpha)
				c.G = uint8(float32(c.G)*(1-alpha) + 100*alpha)
				c.B = uint8(float32(c.B)*(1-alpha) + 200*alpha)
			}
			img.SetRGBA(tx, ty, c)
		}
	}

	return ebiten.NewImageFromImage(img)
}
