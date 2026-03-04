package game

import (
	"fmt"
	"image"
	"proc-gen-world/internal/camera"
	"proc-gen-world/internal/generation"
	"proc-gen-world/internal/render"
	"proc-gen-world/internal/world"
	"proc-gen-world/internal/world/terrain"
	"strings"

	"github.com/eihigh/vec"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	world         world.World
	settings      world.WorldGenSettings
	images        [5]*ebiten.Image
	viewMode      int // 1=biome 2=elevation 3=rainfall 4=weather 5=voronoi
	inspectedTile *terrain.Tile
	inspectPos    image.Point
	cam           camera.Camera
}

func NewGame() (*Game, error) {
	settings := world.DefaultSettings()
	result, err := generation.GenWithSettings(settings)
	if err != nil {
		return nil, err
	}
	var images [5]*ebiten.Image
	images[0] = render.BuildMapImage(&result.Output, settings.WorldSize)
	images[1] = render.BuildElevationImage(result.Heightmap)
	images[2] = render.BuildRainfallImage(result.Rainfall)
	images[3] = render.BuildWeatherImage(result.Weather)
	images[4] = render.BuildVoronoiImage(&result.Output, settings.WorldSize)
	return &Game{
		world:    result.Output,
		settings: settings,
		images:   images,
		viewMode: 1,
		cam:      camera.New(float64(settings.WorldSize.X), float64(settings.WorldSize.Y)),
	}, nil
}

func (g *Game) Update() error {
	tps := ebiten.ActualTPS()
	if tps == 0 {
		tps = 60
	}
	g.cam.Update(1.0 / tps)

	newDay, newWeek := g.world.Calendar.Tick()
	if newWeek {
		_ = newDay
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.viewMode = 1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.viewMode = 2
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.viewMode = 3
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		g.viewMode = 4
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		g.viewMode = 5
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		sx, sy := ebiten.CursorPosition()
		sw, sh := ebiten.WindowSize()
		wfx, wfy := g.cam.ScreenToWorld(sx, sy, sw, sh)
		wx, wy := int(wfx), int(wfy)
		if wx >= 0 && wy >= 0 && wx < g.settings.WorldSize.X && wy < g.settings.WorldSize.Y {
			cx, cy := wx/terrain.ChunkSize, wy/terrain.ChunkSize
			if chunk, ok := g.world.Chunks[vec.Vec2i{X: cx, Y: cy}]; ok {
				t := chunk.Tiles[wy%terrain.ChunkSize][wx%terrain.ChunkSize]
				g.inspectedTile = &t
				g.inspectPos = image.Point{X: sx, Y: sy}
			}
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.inspectedTile = nil
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	src := g.images[g.viewMode-1]
	mapW, mapH := src.Bounds().Dx(), src.Bounds().Dy()

	srcRect, offsetX, offsetY := g.cam.DrawParams(sw, sh)
	clampedRect := srcRect.Intersect(image.Rect(0, 0, mapW, mapH))

	if !clampedRect.Empty() {
		sub := src.SubImage(clampedRect).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(offsetX, offsetY)
		op.GeoM.Scale(g.cam.Zoom, g.cam.Zoom)
		// Compensate for left/top clamping when camera is near world edge
		op.GeoM.Translate(
			float64(clampedRect.Min.X-srcRect.Min.X)*g.cam.Zoom,
			float64(clampedRect.Min.Y-srcRect.Min.Y)*g.cam.Zoom,
		)
		screen.DrawImage(sub, op)
	}

	date := g.world.Calendar.CurrentDate()
	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"Year %d  Month %d  Day %d  (%s)", date.Year, date.Month, date.Day, date.Season()))

	labels := []string{"[F1] Biome", "[F2] Elevation", "[F3] Rainfall", "[F4] Weather", "[F5] Voronoi"}
	var sb strings.Builder
	for i, label := range labels {
		if g.viewMode == i+1 {
			sb.WriteString("> ")
		} else {
			sb.WriteString("  ")
		}
		sb.WriteString(label)
		if i < len(labels)-1 {
			sb.WriteByte('\n')
		}
	}
	ebitenutil.DebugPrintAt(screen, sb.String(), sw-130, 4)

	if g.inspectedTile != nil {
		t := g.inspectedTile
		biome := terrain.Biomes[t.Biome]
		tempC := int(t.Temperature*60) - 20 // remap [0,1] → [-20°C, 40°C]
		info := fmt.Sprintf(
			"Biome:     %s\nElevation: %.2f\nTemp:      %d°C",
			biome.Name, t.Height, tempC,
		)
		ebitenutil.DebugPrintAt(screen, info, g.inspectPos.X+8, g.inspectPos.Y+8)
	}
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ow, oh
}
