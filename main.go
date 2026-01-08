package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60
	gameTitle    = "Herzog Drei"
)

// Game holds the game state
type Game struct {
	camera rl.Camera3D
}

// NewGame creates and initializes a new game instance
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

// init sets up initial game state
func (g *Game) init() {
	// Set up 3D camera
	g.camera = rl.Camera3D{
		Position:   rl.NewVector3(10.0, 10.0, 10.0),
		Target:     rl.NewVector3(0.0, 0.0, 0.0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}
}

// Update handles game logic each frame
func (g *Game) Update() {
	// Update camera controls (orbital camera for now)
	rl.UpdateCamera(&g.camera, rl.CameraOrbital)
}

// Render draws the game each frame
func (g *Game) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	rl.BeginMode3D(g.camera)

	// Draw ground plane
	rl.DrawPlane(rl.NewVector3(0, 0, 0), rl.NewVector2(20, 20), rl.LightGray)

	// Draw reference cube at origin
	rl.DrawCube(rl.NewVector3(0, 0.5, 0), 1.0, 1.0, 1.0, rl.Red)
	rl.DrawCubeWires(rl.NewVector3(0, 0.5, 0), 1.0, 1.0, 1.0, rl.Maroon)

	// Draw grid for reference
	rl.DrawGrid(20, 1.0)

	rl.EndMode3D()

	// Draw UI overlay
	rl.DrawText(gameTitle, 10, 10, 20, rl.DarkGray)
	rl.DrawFPS(screenWidth-100, 10)
	rl.DrawText("Use mouse to orbit camera", 10, screenHeight-25, 15, rl.Gray)

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
