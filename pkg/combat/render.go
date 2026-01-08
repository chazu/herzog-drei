package combat

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Renderer handles rendering of combat effects
type Renderer struct{}

// NewRenderer creates a new combat renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Draw renders all combat effects
func (r *Renderer) Draw(sys *System) {
	r.drawExplosions(sys)
}

// drawExplosions renders explosion effects
func (r *Renderer) drawExplosions(sys *System) {
	for _, e := range sys.GetExplosions() {
		if !e.Active {
			continue
		}

		// Calculate fade based on time
		t := e.Elapsed / e.Duration
		alpha := uint8(255 * (1.0 - t))

		// Inner core (bright)
		coreColor := rl.Color{R: 255, G: 255, B: 200, A: alpha}
		rl.DrawSphere(e.Position, e.Radius*0.3, coreColor)

		// Outer ring (colored)
		outerColor := rl.Color{R: e.Color.R, G: e.Color.G, B: e.Color.B, A: alpha / 2}
		rl.DrawSphere(e.Position, e.Radius, outerColor)

		// Draw ring on ground
		groundPos := rl.Vector3{X: e.Position.X, Y: 0.05, Z: e.Position.Z}
		rl.DrawCircle3D(groundPos, e.Radius*1.5, rl.Vector3{X: 1, Y: 0, Z: 0}, 90, outerColor)
	}
}

// DrawUI renders combat-related UI elements
func (r *Renderer) DrawUI(sys *System, screenWidth, screenHeight int) {
	// Draw respawn countdown if mech is dead
	if sys.IsMechDead() {
		timer := sys.GetRespawnTimer()

		// Dark overlay
		rl.DrawRectangle(0, 0, int32(screenWidth), int32(screenHeight), rl.Color{R: 0, G: 0, B: 0, A: 150})

		// "DESTROYED" text
		text := "DESTROYED"
		textWidth := rl.MeasureText(text, 60)
		rl.DrawText(text, int32(screenWidth/2)-textWidth/2, int32(screenHeight/2)-60, 60, rl.Red)

		// Respawn countdown
		countdownText := fmt.Sprintf("Respawning in %.1f...", timer)
		countdownWidth := rl.MeasureText(countdownText, 30)
		rl.DrawText(countdownText, int32(screenWidth/2)-countdownWidth/2, int32(screenHeight/2)+20, 30, rl.White)
	}

	// Draw invulnerability indicator
	if sys.IsMechInvulnerable() {
		timer := sys.GetInvulnTimer()
		text := fmt.Sprintf("INVULNERABLE %.1f", timer)
		textWidth := rl.MeasureText(text, 20)

		// Flash effect
		alpha := uint8(200)
		if int(timer*10)%2 == 0 {
			alpha = 100
		}

		rl.DrawText(text, int32(screenWidth/2)-textWidth/2, 50, 20, rl.Color{R: 0, G: 255, B: 255, A: alpha})
	}
}
