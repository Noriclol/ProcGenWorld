# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run ./cmd/ebiten/        # Run the game
go build ./cmd/ebiten/      # Build
go test ./...               # Run all tests
```

## Architecture

This is a procedurally generated world using the [Ebiten](https://ebitengine.org/) 2D game engine.

**Entry point:** `cmd/ebiten/main.go` — sets up the Ebiten window and implements the `ebiten.Game` interface (`Update`, `Draw`, `Layout`).

**Game loop flow:**
1. `cmd/ebiten/main.go` initializes Ebiten and creates a `game.Game`
2. `internal/game/game.go` embeds `world.World` and holds game state
3. World is generated via `internal/generation/gen.go` (`Gen()` returns `WorldGenResult`)
4. The world is chunk-based: `World` stores a `map[vec.Vec2]terrain.Chunk`

**Key packages:**
- `internal/world` — core types: `World`, `WorldGenResult`, `WorldGenSettings`
- `internal/world/terrain` — `Chunk` and `biome` (terrain primitives)
- `internal/world/common` — `Location` with Vec2 position
- `internal/world/entities` — `Entity` interface
- `internal/generation` — `Gen()` function that produces a `WorldGenResult`

**Coordinate system:** Uses `github.com/eihigh/vec` (`vec.Vec2`) for 2D positions throughout.

**Chunk-based world:** The world is divided into chunks indexed by `vec.Vec2`. This is the primary spatial data structure for the terrain.
