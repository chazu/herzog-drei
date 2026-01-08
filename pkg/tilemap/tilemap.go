package tilemap

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	DefaultTileSize = 1.0 // World units per tile
)

// Tile represents a single tile in the map
type Tile struct {
	Terrain TerrainType
}

// TileMap holds the game world map data
type TileMap struct {
	Width    int
	Height   int
	TileSize float32
	Tiles    [][]Tile
}

// NewTileMap creates a new tile map with the given dimensions
func NewTileMap(width, height int) *TileMap {
	tm := &TileMap{
		Width:    width,
		Height:   height,
		TileSize: DefaultTileSize,
		Tiles:    make([][]Tile, height),
	}

	// Initialize all tiles as ground
	for y := 0; y < height; y++ {
		tm.Tiles[y] = make([]Tile, width)
		for x := 0; x < width; x++ {
			tm.Tiles[y][x] = Tile{Terrain: TerrainGround}
		}
	}

	return tm
}

// InBounds checks if coordinates are within map bounds
func (tm *TileMap) InBounds(x, y int) bool {
	return x >= 0 && x < tm.Width && y >= 0 && y < tm.Height
}

// GetTile returns the tile at the given coordinates
func (tm *TileMap) GetTile(x, y int) *Tile {
	if !tm.InBounds(x, y) {
		return nil
	}
	return &tm.Tiles[y][x]
}

// SetTerrain sets the terrain type at the given coordinates
func (tm *TileMap) SetTerrain(x, y int, terrain TerrainType) {
	if tm.InBounds(x, y) {
		tm.Tiles[y][x].Terrain = terrain
	}
}

// WorldToTile converts world coordinates to tile coordinates
func (tm *TileMap) WorldToTile(worldX, worldZ float32) (int, int) {
	tileX := int(worldX / tm.TileSize)
	tileY := int(worldZ / tm.TileSize)
	return tileX, tileY
}

// TileToWorld converts tile coordinates to world coordinates (center of tile)
func (tm *TileMap) TileToWorld(tileX, tileY int) (float32, float32) {
	worldX := (float32(tileX) + 0.5) * tm.TileSize
	worldZ := (float32(tileY) + 0.5) * tm.TileSize
	return worldX, worldZ
}

// GetTerrainAt returns the terrain type at world coordinates
func (tm *TileMap) GetTerrainAt(worldX, worldZ float32) TerrainType {
	tileX, tileY := tm.WorldToTile(worldX, worldZ)
	tile := tm.GetTile(tileX, tileY)
	if tile == nil {
		return TerrainGround
	}
	return tile.Terrain
}

// GetHeightAt returns the terrain height at world coordinates
func (tm *TileMap) GetHeightAt(worldX, worldZ float32) float32 {
	terrain := tm.GetTerrainAt(worldX, worldZ)
	return GetTerrainInfo(terrain).Height
}

// IsPassableAt checks if ground units can traverse the given world position
func (tm *TileMap) IsPassableAt(worldX, worldZ float32) bool {
	terrain := tm.GetTerrainAt(worldX, worldZ)
	return terrain.IsPassable()
}

// IsFlyableAt checks if air units can fly over the given world position
func (tm *TileMap) IsFlyableAt(worldX, worldZ float32) bool {
	terrain := tm.GetTerrainAt(worldX, worldZ)
	return terrain.IsFlyable()
}

// GetWorldBounds returns the world-space bounds of the map
func (tm *TileMap) GetWorldBounds() rl.BoundingBox {
	return rl.BoundingBox{
		Min: rl.NewVector3(0, -1, 0),
		Max: rl.NewVector3(float32(tm.Width)*tm.TileSize, 5, float32(tm.Height)*tm.TileSize),
	}
}

// Render draws the tile map in 3D
func (tm *TileMap) Render() {
	for y := 0; y < tm.Height; y++ {
		for x := 0; x < tm.Width; x++ {
			tile := tm.Tiles[y][x]
			info := GetTerrainInfo(tile.Terrain)

			worldX, worldZ := tm.TileToWorld(x, y)

			// Draw tile as a cube with appropriate height
			tileHeight := info.Height
			if tileHeight < 0.1 {
				tileHeight = 0.1 // Minimum visual height
			}

			pos := rl.NewVector3(worldX, info.Height/2, worldZ)
			size := rl.NewVector3(tm.TileSize*0.98, tileHeight, tm.TileSize*0.98)

			rl.DrawCubeV(pos, size, info.Color)

			// Draw water with transparency effect
			if tile.Terrain == TerrainWater {
				waterColor := rl.NewColor(64, 164, 223, 180)
				rl.DrawCubeV(pos, size, waterColor)
			}

			// Draw mountain peaks
			if tile.Terrain == TerrainMountain {
				peakPos := rl.NewVector3(worldX, info.Height, worldZ)
				rl.DrawCube(peakPos, tm.TileSize*0.4, 0.5, tm.TileSize*0.4, rl.DarkGray)
			}

			// Draw trees for forest
			if tile.Terrain == TerrainForest {
				treePos := rl.NewVector3(worldX, info.Height+0.3, worldZ)
				rl.DrawCube(treePos, 0.2, 0.6, 0.2, rl.Brown)
				rl.DrawSphere(rl.NewVector3(worldX, info.Height+0.7, worldZ), 0.3, rl.DarkGreen)
			}
		}
	}
}

// FillRect fills a rectangular area with the specified terrain
func (tm *TileMap) FillRect(x1, y1, x2, y2 int, terrain TerrainType) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			tm.SetTerrain(x, y, terrain)
		}
	}
}

// GenerateTestMap creates a test map with various terrain types
func GenerateTestMap(width, height int) *TileMap {
	tm := NewTileMap(width, height)

	// Add some water (river)
	riverX := width / 3
	for y := 0; y < height; y++ {
		tm.SetTerrain(riverX, y, TerrainWater)
		tm.SetTerrain(riverX+1, y, TerrainWater)
	}

	// Add some mountains
	tm.FillRect(width*2/3, height/4, width*2/3+3, height/4+3, TerrainMountain)

	// Add some forest
	tm.FillRect(5, 5, 8, 8, TerrainForest)
	tm.FillRect(width-10, height-10, width-6, height-6, TerrainForest)

	// Add a road
	roadY := height / 2
	for x := 0; x < width; x++ {
		if tm.GetTile(x, roadY).Terrain == TerrainGround {
			tm.SetTerrain(x, roadY, TerrainRoad)
		}
	}

	return tm
}
