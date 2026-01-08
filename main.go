package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/tilemap"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60
	gameTitle    = "Herzog Drei"

	mapWidth  = 64
	mapHeight = 48
)

// Game holds the game state
type Game struct {
	tileMap   *tilemap.TileMap
	camera    *tilemap.GameCamera
	minimap   *tilemap.Minimap
	playerPos rl.Vector3 // Simulated player position for camera follow
}

// NewGame creates and initializes a new game instance
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

// init sets up initial game state
func (g *Game) init() {
	// Create tile map with test terrain
	g.tileMap = tilemap.GenerateTestMap(mapWidth, mapHeight)

	// Set up game camera
	g.camera = tilemap.NewGameCamera()
	g.camera.SetBounds(g.tileMap.GetWorldBounds())

	// Start player in center of map
	centerX, centerZ := g.tileMap.TileToWorld(mapWidth/2, mapHeight/2)
	g.playerPos = rl.NewVector3(centerX, 0.5, centerZ)
	g.camera.SetTarget(g.playerPos)

	// Set up minimap in top-right corner
	g.minimap = tilemap.NewMinimap()
	g.minimap.SetPosition(screenWidth-210, 10)
	g.minimap.SetSize(200, 150)
}

// Update handles game logic each frame
func (g *Game) Update() {
	// Handle camera input (zoom)
	g.camera.HandleInput()

	// Simulate player movement with WASD keys
	speed := float32(0.2)
	moved := false

	if rl.IsKeyDown(rl.KeyW) {
		g.playerPos.Z -= speed
		moved = true
	}
	if rl.IsKeyDown(rl.KeyS) {
		g.playerPos.Z += speed
		moved = true
	}
	if rl.IsKeyDown(rl.KeyA) {
		g.playerPos.X -= speed
		moved = true
	}
	if rl.IsKeyDown(rl.KeyD) {
		g.playerPos.X += speed
		moved = true
	}

	// Check terrain collision for ground movement
	if moved && !g.tileMap.IsPassableAt(g.playerPos.X, g.playerPos.Z) {
		// Revert to previous position if terrain is impassable
		if rl.IsKeyDown(rl.KeyW) {
			g.playerPos.Z += speed
		}
		if rl.IsKeyDown(rl.KeyS) {
			g.playerPos.Z -= speed
		}
		if rl.IsKeyDown(rl.KeyA) {
			g.playerPos.X += speed
		}
		if rl.IsKeyDown(rl.KeyD) {
			g.playerPos.X -= speed
		}
	}

	// Adjust player height based on terrain
	g.playerPos.Y = g.tileMap.GetHeightAt(g.playerPos.X, g.playerPos.Z) + 0.5

	// Update camera to follow player
	g.camera.SetTarget(g.playerPos)
	g.camera.Update()
}

// Render draws the game each frame
func (g *Game) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.SkyBlue)

	// 3D rendering
	g.camera.Begin3D()

	// Render tile map
	g.tileMap.Render()

	// Draw player indicator (simple cube for now)
	rl.DrawCube(g.playerPos, 0.6, 0.6, 0.6, rl.Red)
	rl.DrawCubeWires(g.playerPos, 0.6, 0.6, 0.6, rl.Maroon)

	g.camera.End3D()

	// Draw minimap with player marker
	markers := []tilemap.MinimapMarker{
		tilemap.NewMarker(g.playerPos.X, g.playerPos.Z, tilemap.MarkerPlayer, rl.Red),
	}
	g.minimap.RenderWithMarkers(g.tileMap, g.camera, markers)

	// Draw UI overlay
	rl.DrawText(gameTitle, 10, screenHeight-60, 20, rl.DarkGray)
	rl.DrawFPS(10, screenHeight-35)
	rl.DrawText("WASD: Move | Mouse wheel: Zoom", 10, screenHeight-15, 12, rl.DarkGray)

	// Show current terrain info
	terrain := g.tileMap.GetTerrainAt(g.playerPos.X, g.playerPos.Z)
	info := tilemap.GetTerrainInfo(terrain)
	rl.DrawText("Terrain: "+info.Name, 10, 170, 15, rl.DarkGray)

	rl.EndDrawing()
}

func main() {
	// Initialize window
	rl.InitWindow(screenWidth, screenHeight, gameTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(targetFPS)

	// Create game instance
	game := NewGame()

	// Main game loop
	for !rl.WindowShouldClose() {
		game.Update()
		game.Render()
	}
}
