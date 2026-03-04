package camera

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MinZoom     = 1.0
	MaxZoom     = 4.0
	DefaultZoom = 2.0   // 256 tiles × 2px = 512px — world fits in 640×480 window
	PanSpeed    = 200.0 // screen pixels per second; divided by Zoom → consistent screen-space feel
)

type Camera struct {
	X, Y           float64 // world center in tile coords
	Zoom           float64 // screen pixels per tile
	WorldW, WorldH float64
}

func New(worldW, worldH float64) Camera {
	return Camera{X: worldW / 2, Y: worldH / 2, Zoom: DefaultZoom, WorldW: worldW, WorldH: worldH}
}

func (c *Camera) Update(dt float64) {
	delta := PanSpeed / c.Zoom * dt
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		c.X -= delta
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		c.X += delta
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Y -= delta
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Y += delta
	}

	if ebiten.IsKeyPressed(ebiten.KeyI) {
		c.Zoom *= 1.15
	}
	if ebiten.IsKeyPressed(ebiten.KeyU) {
		c.Zoom /= 1.15
	}
	c.Zoom = clamp(c.Zoom, MinZoom, MaxZoom)

	c.X = clamp(c.X, 0, c.WorldW)
	c.Y = clamp(c.Y, 0, c.WorldH)
}

// DrawParams returns the source SubImage rect and the screen-space translation offset
// for sub-pixel smooth scrolling.
func (c *Camera) DrawParams(sw, sh int) (srcRect image.Rectangle, offsetX, offsetY float64) {
	visW := float64(sw) / c.Zoom
	visH := float64(sh) / c.Zoom
	left := c.X - visW/2
	top := c.Y - visH/2

	srcLeft := int(math.Floor(left))
	srcTop := int(math.Floor(top))
	fracX := left - float64(srcLeft) // sub-tile remainder [0, 1)
	fracY := top - float64(srcTop)

	srcRect = image.Rect(srcLeft, srcTop,
		int(math.Ceil(left+visW)),
		int(math.Ceil(top+visH)))

	offsetX = -fracX * c.Zoom // shift image left by fractional pixel amount
	offsetY = -fracY * c.Zoom
	return
}

// ScreenToWorld converts a screen pixel to a world tile coordinate.
func (c *Camera) ScreenToWorld(sx, sy, sw, sh int) (wx, wy float64) {
	visW := float64(sw) / c.Zoom
	visH := float64(sh) / c.Zoom
	wx = (c.X - visW/2) + float64(sx)/c.Zoom
	wy = (c.Y - visH/2) + float64(sy)/c.Zoom
	return
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
