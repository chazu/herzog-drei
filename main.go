package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/mech"
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
	// Map and camera
	tileMap *tilemap.TileMap
	camera  *tilemap.GameCamera
	minimap *tilemap.Minimap

	// Player mech
	playerMech   *mech.Mech
	mechInput    *mech.InputHandler
	mechRenderer *mech.Renderer
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

	// Create player mech at center of map
	centerX, centerZ := g.tileMap.TileToWorld(mapWidth/2, mapHeight/2)
	startPos := rl.NewVector3(centerX, 3, centerZ)
	g.playerMech = mech.New(startPos, mech.DefaultConfig())
	g.mechInput = mech.NewInputHandler()
	g.mechRenderer = mech.NewRenderer()

	// Set camera to follow mech
	g.camera.SetTarget(g.playerMech.Position)

	// Set up minimap in top-right corner
	g.minimap = tilemap.NewMinimap()
	g.minimap.SetPosition(screenWidth-210, 10)
	g.minimap.SetSize(200, 150)
}

// Update handles game logic each frame
func (g *Game) Update() {
	dt := rl.GetFrameTime()

	// Handle camera input (zoom)
	g.camera.HandleInput()

	// Process player input
	g.mechInput.Update(g.playerMech)

	// Update mech
	g.playerMech.Update(dt)

	// Check terrain collision for ground (robot) mode
	if g.playerMech.Mode == mech.ModeRobot {
		if !g.tileMap.IsPassableAt(g.playerMech.Position.X, g.playerMech.Position.Z) {
			// Push mech back if on impassable terrain
			g.playerMech.Position.X -= g.playerMech.Velocity.X * dt
			g.playerMech.Position.Z -= g.playerMech.Velocity.Z * dt
		}
		// Adjust height based on terrain
		g.playerMech.Position.Y = g.tileMap.GetHeightAt(g.playerMech.Position.X, g.playerMech.Position.Z)
	}

	// Update camera to follow mech
	g.camera.SetTarget(g.playerMech.Position)
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

	// Draw player mech
	g.mechRenderer.Draw(g.playerMech)

	g.camera.End3D()

	// Draw minimap with player marker
	markers := []tilemap.MinimapMarker{
		tilemap.NewMarker(g.playerMech.Position.X, g.playerMech.Position.Z, tilemap.MarkerPlayer, rl.Red),
	}
	g.minimap.RenderWithMarkers(g.tileMap, g.camera, markers)

	// Draw UI overlay
	rl.DrawText(gameTitle, 10, 10, 20, rl.DarkGray)
	rl.DrawFPS(screenWidth-100, 10)

	// Draw mech UI (health bar, mode indicator)
	g.mechRenderer.DrawUI(g.playerMech, screenWidth, screenHeight)

	// Show current terrain info
	terrain := g.tileMap.GetTerrainAt(g.playerMech.Position.X, g.playerMech.Position.Z)
	info := tilemap.GetTerrainInfo(terrain)
	rl.DrawText("Terrain: "+info.Name, 10, screenHeight-60, 15, rl.DarkGray)
	rl.DrawText("Space: Transform | Mouse: Aim | Click: Shoot | Scroll: Zoom", 10, screenHeight-35, 12, rl.DarkGray)

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
