package tilemap

import rl "github.com/gen2brain/raylib-go/raylib"

// TerrainType represents different terrain categories
type TerrainType int

const (
	TerrainGround TerrainType = iota
	TerrainWater
	TerrainMountain
	TerrainForest
	TerrainRoad
)

// TerrainInfo holds properties for a terrain type
type TerrainInfo struct {
	Type       TerrainType
	Name       string
	Color      rl.Color
	Height     float32 // Base height for 3D rendering
	Passable   bool    // Can ground units traverse this?
	Flyable    bool    // Can air units fly over this?
	SpeedMod   float32 // Movement speed modifier (1.0 = normal)
	DefenseMod float32 // Defense bonus modifier (1.0 = normal)
}

// TerrainRegistry maps terrain types to their info
var TerrainRegistry = map[TerrainType]TerrainInfo{
	TerrainGround: {
		Type:       TerrainGround,
		Name:       "Ground",
		Color:      rl.NewColor(139, 119, 101, 255), // Tan/brown
		Height:     0.0,
		Passable:   true,
		Flyable:    true,
		SpeedMod:   1.0,
		DefenseMod: 1.0,
	},
	TerrainWater: {
		Type:       TerrainWater,
		Name:       "Water",
		Color:      rl.NewColor(64, 164, 223, 255), // Blue
		Height:     -0.2,
		Passable:   false,
		Flyable:    true,
		SpeedMod:   0.0,
		DefenseMod: 0.0,
	},
	TerrainMountain: {
		Type:       TerrainMountain,
		Name:       "Mountain",
		Color:      rl.NewColor(128, 128, 128, 255), // Gray
		Height:     1.5,
		Passable:   false,
		Flyable:    false,
		SpeedMod:   0.0,
		DefenseMod: 0.0,
	},
	TerrainForest: {
		Type:       TerrainForest,
		Name:       "Forest",
		Color:      rl.NewColor(34, 139, 34, 255), // Forest green
		Height:     0.3,
		Passable:   true,
		Flyable:    true,
		SpeedMod:   0.6,
		DefenseMod: 1.3,
	},
	TerrainRoad: {
		Type:       TerrainRoad,
		Name:       "Road",
		Color:      rl.NewColor(160, 160, 160, 255), // Light gray
		Height:     0.05,
		Passable:   true,
		Flyable:    true,
		SpeedMod:   1.5,
		DefenseMod: 0.8,
	},
}

// GetTerrainInfo returns the info for a terrain type
func GetTerrainInfo(t TerrainType) TerrainInfo {
	if info, ok := TerrainRegistry[t]; ok {
		return info
	}
	return TerrainRegistry[TerrainGround]
}

// IsPassable checks if a terrain type can be traversed by ground units
func (t TerrainType) IsPassable() bool {
	return GetTerrainInfo(t).Passable
}

// IsFlyable checks if a terrain type can be flown over
func (t TerrainType) IsFlyable() bool {
	return GetTerrainInfo(t).Flyable
}
