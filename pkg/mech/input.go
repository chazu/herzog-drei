package mech

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// InputHandler processes player input for the mech
type InputHandler struct {
	transformPressed bool // Track transform key state for edge detection
}

// NewInputHandler creates a new input handler
func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

// Update reads input and applies it to the mech
func (h *InputHandler) Update(m *Mech) {
	// Movement input (WASD)
	var moveX, moveZ float32

	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		moveZ = 1
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		moveZ = -1
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		moveX = 1
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		moveX = -1
	}

	// Normalize diagonal movement
	if moveX != 0 && moveZ != 0 {
		invLen := float32(1.0 / math.Sqrt(2))
		moveX *= invLen
		moveZ *= invLen
	}

	m.InputMove = rl.Vector2{X: moveX, Y: moveZ}

	// Shooting input (Space or Left Mouse)
	m.InputShoot = rl.IsKeyDown(rl.KeySpace) || rl.IsMouseButtonDown(rl.MouseLeftButton)

	// Transform input (T key) - edge triggered
	transformDown := rl.IsKeyDown(rl.KeyT)
	m.InputTransform = transformDown && !h.transformPressed
	h.transformPressed = transformDown
}
