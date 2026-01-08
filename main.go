package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/mech"
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

	// Draw player mech
	g.mechRenderer.Draw(g.playerMech)

	rl.EndMode3D()

	// Draw UI overlay
	rl.DrawText(gameTitle, 10, 10, 20, rl.White)
	rl.DrawFPS(screenWidth-100, 10)

	// Draw mech UI
	g.mechRenderer.DrawUI(g.playerMech, screenWidth, screenHeight)

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
