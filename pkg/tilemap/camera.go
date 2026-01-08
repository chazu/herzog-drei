package tilemap

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// GameCamera provides a camera that follows a target with smooth scrolling
type GameCamera struct {
	Camera       rl.Camera3D
	Target       rl.Vector3  // Position to follow
	Offset       rl.Vector3  // Offset from target (determines viewing angle)
	SmoothSpeed  float32     // How quickly camera catches up (0-1, higher = faster)
	Bounds       *rl.BoundingBox // Optional bounds to constrain camera
	ZoomLevel    float32     // Zoom multiplier
	MinZoom      float32
	MaxZoom      float32
}

// NewGameCamera creates a new camera configured for Herzog Drei-style viewing
func NewGameCamera() *GameCamera {
	gc := &GameCamera{
		Target:      rl.NewVector3(0, 0, 0),
		Offset:      rl.NewVector3(0, 15, 10), // High above, slightly behind
		SmoothSpeed: 0.1,
		ZoomLevel:   1.0,
		MinZoom:     0.5,
		MaxZoom:     2.0,
	}

	gc.Camera = rl.Camera3D{
		Position:   rl.Vector3Add(gc.Target, gc.Offset),
		Target:     gc.Target,
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	return gc
}

// SetTarget sets the position the camera should follow
func (gc *GameCamera) SetTarget(pos rl.Vector3) {
	gc.Target = pos
}

// SetBounds sets the world bounds to constrain camera movement
func (gc *GameCamera) SetBounds(bounds rl.BoundingBox) {
	gc.Bounds = &bounds
}

// Update smoothly moves the camera toward its target
func (gc *GameCamera) Update() {
	// Calculate desired camera position
	scaledOffset := rl.Vector3Scale(gc.Offset, gc.ZoomLevel)
	desiredPos := rl.Vector3Add(gc.Target, scaledOffset)

	// Apply bounds constraints if set
	if gc.Bounds != nil {
		desiredPos = gc.constrainToBounds(desiredPos)
	}

	// Smooth interpolation toward desired position
	gc.Camera.Position = rl.Vector3Lerp(gc.Camera.Position, desiredPos, gc.SmoothSpeed)
	gc.Camera.Target = rl.Vector3Lerp(gc.Camera.Target, gc.Target, gc.SmoothSpeed)
}

// constrainToBounds keeps the camera view within map bounds
func (gc *GameCamera) constrainToBounds(pos rl.Vector3) rl.Vector3 {
	if gc.Bounds == nil {
		return pos
	}

	// Keep target within bounds (camera can extend past to see edges)
	result := pos

	minX := gc.Bounds.Min.X
	maxX := gc.Bounds.Max.X
	minZ := gc.Bounds.Min.Z
	maxZ := gc.Bounds.Max.Z

	// Constrain the look-at target
	if gc.Target.X < minX {
		gc.Target.X = minX
	}
	if gc.Target.X > maxX {
		gc.Target.X = maxX
	}
	if gc.Target.Z < minZ {
		gc.Target.Z = minZ
	}
	if gc.Target.Z > maxZ {
		gc.Target.Z = maxZ
	}

	return result
}

// Zoom adjusts the camera zoom level
func (gc *GameCamera) Zoom(delta float32) {
	gc.ZoomLevel += delta
	if gc.ZoomLevel < gc.MinZoom {
		gc.ZoomLevel = gc.MinZoom
	}
	if gc.ZoomLevel > gc.MaxZoom {
		gc.ZoomLevel = gc.MaxZoom
	}
}

// HandleInput processes camera-related input (zoom, etc.)
func (gc *GameCamera) HandleInput() {
	// Mouse wheel zoom
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		gc.Zoom(-wheel * 0.1)
	}
}

// Begin3D starts 3D rendering mode with this camera
func (gc *GameCamera) Begin3D() {
	rl.BeginMode3D(gc.Camera)
}

// End3D ends 3D rendering mode
func (gc *GameCamera) End3D() {
	rl.EndMode3D()
}

// ScreenToWorld converts screen coordinates to a world position on a plane at height y
func (gc *GameCamera) ScreenToWorld(screenPos rl.Vector2, groundY float32) rl.Vector3 {
	ray := rl.GetScreenToWorldRay(screenPos, gc.Camera)

	// Calculate intersection with ground plane
	if math.Abs(float64(ray.Direction.Y)) < 0.0001 {
		// Ray is nearly parallel to ground, return camera target
		return gc.Target
	}

	t := (groundY - ray.Position.Y) / ray.Direction.Y
	if t < 0 {
		t = 0 // Intersection is behind camera
	}

	return rl.Vector3{
		X: ray.Position.X + ray.Direction.X*t,
		Y: groundY,
		Z: ray.Position.Z + ray.Direction.Z*t,
	}
}

// GetVisibleTileRange returns the range of tiles currently visible
func (gc *GameCamera) GetVisibleTileRange(tm *TileMap) (minX, minY, maxX, maxY int) {
	// Get corners of visible area at ground level
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())

	corners := []rl.Vector2{
		{X: 0, Y: 0},
		{X: screenW, Y: 0},
		{X: 0, Y: screenH},
		{X: screenW, Y: screenH},
	}

	minX, minY = tm.Width, tm.Height
	maxX, maxY = 0, 0

	for _, corner := range corners {
		worldPos := gc.ScreenToWorld(corner, 0)
		tileX, tileY := tm.WorldToTile(worldPos.X, worldPos.Z)

		if tileX < minX {
			minX = tileX
		}
		if tileX > maxX {
			maxX = tileX
		}
		if tileY < minY {
			minY = tileY
		}
		if tileY > maxY {
			maxY = tileY
		}
	}

	// Add some padding for safety
	minX -= 2
	minY -= 2
	maxX += 2
	maxY += 2

	// Clamp to map bounds
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX >= tm.Width {
		maxX = tm.Width - 1
	}
	if maxY >= tm.Height {
		maxY = tm.Height - 1
	}

	return minX, minY, maxX, maxY
}
