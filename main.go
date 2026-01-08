package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/mech"
	"github.com/chazu/herzog-drei/pkg/unit"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60
	gameTitle    = "Herzog Drei"

	// Camera follow settings
	cameraHeight   = 12.0
	cameraDistance = 8.0
	cameraLerp     = 5.0
)

// Game holds the game state
type Game struct {
	camera rl.Camera3D

	// Player mech
	playerMech    *mech.Mech
	mechInput     *mech.InputHandler
	mechRenderer  *mech.Renderer

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
	// Create player mech at origin
	g.playerMech = mech.New(rl.NewVector3(0, 3, 0), mech.DefaultConfig())
	g.mechInput = mech.NewInputHandler()
	g.mechRenderer = mech.NewRenderer()

	// Initialize unit system
	g.unitManager = unit.NewManager(100) // Max 100 units
	g.unitRenderer = unit.NewRenderer()
	g.unitPathfinder = unit.NewPathfinder(80, 80, 1.0) // 80x80 grid, 1 unit cells
	g.unitManager.Pathfinder = g.unitPathfinder

	// Spawn test units for demonstration
	g.spawnTestUnits()

	// Set up 3D camera (will follow player)
	g.camera = rl.Camera3D{
		Position:   rl.NewVector3(0, cameraHeight, -cameraDistance),
		Target:     rl.NewVector3(0.0, 0.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}
}

// Update handles game logic each frame
func (g *Game) Update() {
	dt := rl.GetFrameTime()

	// Process player input
	g.mechInput.Update(g.playerMech)

	// Update mech
	g.playerMech.Update(dt)

	// Update units
	g.unitManager.Update(dt)

	// Handle unit spawning (press 1-6 to spawn player units)
	g.handleUnitSpawnInput()

	// Update camera to follow mech
	g.updateCamera(dt)
}

func (g *Game) updateCamera(dt float32) {
	mechPos := g.playerMech.Position

	// Camera looks at mech position
	targetPos := rl.Vector3{
		X: mechPos.X,
		Y: mechPos.Y,
		Z: mechPos.Z,
	}

	// Camera positioned above and behind
	desiredCamPos := rl.Vector3{
		X: mechPos.X,
		Y: mechPos.Y + cameraHeight,
		Z: mechPos.Z - cameraDistance,
	}

	// Smooth camera movement
	lerpFactor := cameraLerp * dt
	if lerpFactor > 1 {
		lerpFactor = 1
	}

	g.camera.Position.X += (desiredCamPos.X - g.camera.Position.X) * lerpFactor
	g.camera.Position.Y += (desiredCamPos.Y - g.camera.Position.Y) * lerpFactor
	g.camera.Position.Z += (desiredCamPos.Z - g.camera.Position.Z) * lerpFactor

	g.camera.Target.X += (targetPos.X - g.camera.Target.X) * lerpFactor
	g.camera.Target.Y += (targetPos.Y - g.camera.Target.Y) * lerpFactor
	g.camera.Target.Z += (targetPos.Z - g.camera.Target.Z) * lerpFactor
}

// Render draws the game each frame
func (g *Game) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Color{R: 40, G: 44, B: 52, A: 255}) // Dark background

	rl.BeginMode3D(g.camera)

	// Draw ground plane
	rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(40, 40), rl.Color{R: 60, G: 64, B: 72, A: 255})

	// Draw grid for reference
	rl.DrawGrid(40, 1.0)

	// Draw some reference objects scattered around
	g.drawEnvironment()

	// Draw units
	g.unitRenderer.Draw(g.unitManager)

	// Draw player mech
	g.mechRenderer.Draw(g.playerMech)

	rl.EndMode3D()

	// Draw UI overlay
	rl.DrawText(gameTitle, 10, 10, 20, rl.White)
	rl.DrawFPS(screenWidth-100, 10)

	// Draw mech UI
	g.mechRenderer.DrawUI(g.playerMech, screenWidth, screenHeight)

	// Draw unit UI
	g.unitRenderer.DrawUI(g.unitManager, screenWidth, screenHeight)

	// Draw spawn hint
	rl.DrawText("1-6: Spawn units | 1:Infantry 2:Tank 3:Bike 4:SAM 5:Boat 6:Supply", 10, 35, 15, rl.Gray)

	rl.EndDrawing()
}

func (g *Game) drawEnvironment() {
	// Draw some buildings/obstacles for reference
	buildings := []struct {
		pos    rl.Vector3
		size   rl.Vector3
		color  rl.Color
	}{
		{rl.NewVector3(8, 1, 8), rl.NewVector3(2, 2, 2), rl.DarkGray},
		{rl.NewVector3(-8, 1.5, 5), rl.NewVector3(3, 3, 3), rl.DarkGray},
		{rl.NewVector3(5, 0.75, -8), rl.NewVector3(1.5, 1.5, 1.5), rl.DarkGray},
		{rl.NewVector3(-6, 1, -6), rl.NewVector3(2, 2, 2), rl.DarkGray},
		{rl.NewVector3(10, 2, -3), rl.NewVector3(4, 4, 2), rl.DarkGray},
		{rl.NewVector3(-10, 1.25, 10), rl.NewVector3(2.5, 2.5, 2.5), rl.DarkGray},
	}

	for _, b := range buildings {
		rl.DrawCube(b.pos, b.size.X, b.size.Y, b.size.Z, b.color)
		rl.DrawCubeWires(b.pos, b.size.X, b.size.Y, b.size.Z, rl.Black)
	}
}

// spawnTestUnits creates initial units for testing
func (g *Game) spawnTestUnits() {
	// Spawn player units on left side
	g.unitManager.SpawnWithObjective(
		unit.TypeInfantry, unit.TeamPlayer,
		rl.NewVector3(-10, 0, 5),
		rl.NewVector3(10, 0, 5),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeTank, unit.TeamPlayer,
		rl.NewVector3(-10, 0, 0),
		rl.NewVector3(10, 0, 0),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeMotorcycle, unit.TeamPlayer,
		rl.NewVector3(-10, 0, -5),
		rl.NewVector3(10, 0, -5),
	)

	// Spawn enemy units on right side
	g.unitManager.SpawnWithObjective(
		unit.TypeInfantry, unit.TeamEnemy,
		rl.NewVector3(10, 0, 5),
		rl.NewVector3(-10, 0, 5),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeTank, unit.TeamEnemy,
		rl.NewVector3(10, 0, 0),
		rl.NewVector3(-10, 0, 0),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeSAM, unit.TeamEnemy,
		rl.NewVector3(10, 0, -5),
		rl.NewVector3(-10, 0, -5),
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
