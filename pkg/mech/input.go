package mech

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// InputHandler processes player input for the mech
type InputHandler struct {
	transformPressed bool // Track transform key state for edge detection
	pickupPressed    bool // Track pickup key state for edge detection
	dropPressed      bool // Track drop key state for edge detection
	orderNextPressed bool // Track order cycle next key state
	orderPrevPressed bool // Track order cycle prev key state
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

	// Pickup input (E key) - edge triggered
	pickupDown := rl.IsKeyDown(rl.KeyE)
	m.InputPickup = pickupDown && !h.pickupPressed
	h.pickupPressed = pickupDown

	// Drop input (Q key) - edge triggered
	dropDown := rl.IsKeyDown(rl.KeyQ)
	m.InputDrop = dropDown && !h.dropPressed
	h.dropPressed = dropDown

	// Order cycling (R = next, F = previous) - edge triggered
	orderNextDown := rl.IsKeyDown(rl.KeyR)
	m.InputOrderNext = orderNextDown && !h.orderNextPressed
	h.orderNextPressed = orderNextDown

	orderPrevDown := rl.IsKeyDown(rl.KeyF)
	m.InputOrderPrev = orderPrevDown && !h.orderPrevPressed
	h.orderPrevPressed = orderPrevDown

	// Handle order cycling immediately
	if m.InputOrderNext {
		m.CycleOrderNext()
	}
	if m.InputOrderPrev {
		m.CycleOrderPrev()
	}
}
