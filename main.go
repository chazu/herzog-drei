package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/mech"
	"github.com/chazu/herzog-drei/pkg/tilemap"
	"github.com/chazu/herzog-drei/pkg/unit"
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

	// Units
	unitManager    *unit.Manager
	unitRenderer   *unit.Renderer
	unitPathfinder *unit.Pathfinder
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

	// Initialize unit system
	g.unitManager = unit.NewManager(100) // Max 100 units
	g.unitRenderer = unit.NewRenderer()
	g.unitPathfinder = unit.NewPathfinder(mapWidth, mapHeight, 1.0)
	g.unitManager.Pathfinder = g.unitPathfinder

	// Spawn test units for demonstration
	g.spawnTestUnits()
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

	// Update units
	g.unitManager.Update(dt)

	// Handle unit spawning (press 1-6 to spawn player units)
	g.handleUnitSpawnInput()

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

	// Draw units
	g.unitRenderer.Draw(g.unitManager)

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

	// Draw unit UI
	g.unitRenderer.DrawUI(g.unitManager, screenWidth, screenHeight)

	// Show current terrain info
	terrain := g.tileMap.GetTerrainAt(g.playerMech.Position.X, g.playerMech.Position.Z)
	info := tilemap.GetTerrainInfo(terrain)
	rl.DrawText("Terrain: "+info.Name, 10, screenHeight-80, 15, rl.DarkGray)
	rl.DrawText("Space: Transform | Mouse: Aim | Click: Shoot | Scroll: Zoom", 10, screenHeight-55, 12, rl.DarkGray)
	rl.DrawText("1-6: Spawn units | 1:Infantry 2:Tank 3:Bike 4:SAM 5:Boat 6:Supply", 10, screenHeight-35, 12, rl.DarkGray)

	rl.EndDrawing()
}

// spawnTestUnits creates initial units for testing
func (g *Game) spawnTestUnits() {
	centerX, centerZ := g.tileMap.TileToWorld(mapWidth/2, mapHeight/2)

	// Spawn player units on left side
	g.unitManager.SpawnWithObjective(
		unit.TypeInfantry, unit.TeamPlayer,
		rl.NewVector3(centerX-10, 0, centerZ+5),
		rl.NewVector3(centerX+10, 0, centerZ+5),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeTank, unit.TeamPlayer,
		rl.NewVector3(centerX-10, 0, centerZ),
		rl.NewVector3(centerX+10, 0, centerZ),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeMotorcycle, unit.TeamPlayer,
		rl.NewVector3(centerX-10, 0, centerZ-5),
		rl.NewVector3(centerX+10, 0, centerZ-5),
	)

	// Spawn enemy units on right side
	g.unitManager.SpawnWithObjective(
		unit.TypeInfantry, unit.TeamEnemy,
		rl.NewVector3(centerX+10, 0, centerZ+5),
		rl.NewVector3(centerX-10, 0, centerZ+5),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeTank, unit.TeamEnemy,
		rl.NewVector3(centerX+10, 0, centerZ),
		rl.NewVector3(centerX-10, 0, centerZ),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeSAM, unit.TeamEnemy,
		rl.NewVector3(centerX+10, 0, centerZ-5),
		rl.NewVector3(centerX-10, 0, centerZ-5),
	)
}

// handleUnitSpawnInput spawns units based on number key presses
func (g *Game) handleUnitSpawnInput() {
	// Spawn near player mech position
	spawnPos := g.playerMech.Position
	spawnPos.Y = 0 // Ground level

	// Offset spawn slightly behind mech
	spawnPos.X -= g.playerMech.GetForward().X * 2
	spawnPos.Z -= g.playerMech.GetForward().Z * 2

	if rl.IsKeyPressed(rl.KeyOne) {
		g.unitManager.Spawn(unit.TypeInfantry, unit.TeamPlayer, spawnPos)
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		g.unitManager.Spawn(unit.TypeTank, unit.TeamPlayer, spawnPos)
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		g.unitManager.Spawn(unit.TypeMotorcycle, unit.TeamPlayer, spawnPos)
	}
	if rl.IsKeyPressed(rl.KeyFour) {
		g.unitManager.Spawn(unit.TypeSAM, unit.TeamPlayer, spawnPos)
	}
	if rl.IsKeyPressed(rl.KeyFive) {
		g.unitManager.Spawn(unit.TypeBoat, unit.TeamPlayer, spawnPos)
	}
	if rl.IsKeyPressed(rl.KeySix) {
		g.unitManager.Spawn(unit.TypeSupply, unit.TeamPlayer, spawnPos)
	}
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
