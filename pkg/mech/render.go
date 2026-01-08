package mech

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Renderer handles mech and projectile rendering
type Renderer struct{}

// NewRenderer creates a new mech renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Draw renders the mech using placeholder geometry
func (r *Renderer) Draw(m *Mech) {
	// Draw based on mode and transformation state
	if m.State == StateTransforming {
		r.drawTransforming(m)
	} else if m.Mode == ModeJet {
		r.drawJetMode(m)
	} else {
		r.drawRobotMode(m)
	}

	// Draw projectiles
	r.drawProjectiles(m)
}

func (r *Renderer) drawJetMode(m *Mech) {
	pos := m.Position
	rot := m.Rotation * 180.0 / math.Pi // Convert to degrees

	// Jet body (elongated box)
	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Main fuselage
	rl.DrawCube(rl.NewVector3(0, 0, 0), 0.4, 0.3, 1.2, rl.Blue)
	rl.DrawCubeWires(rl.NewVector3(0, 0, 0), 0.4, 0.3, 1.2, rl.DarkBlue)

	// Wings
	rl.DrawCube(rl.NewVector3(0, 0, 0.1), 1.4, 0.05, 0.5, rl.Blue)
	rl.DrawCubeWires(rl.NewVector3(0, 0, 0.1), 1.4, 0.05, 0.5, rl.DarkBlue)

	// Tail fins
	rl.DrawCube(rl.NewVector3(0.15, 0.15, -0.5), 0.05, 0.3, 0.2, rl.Blue)
	rl.DrawCube(rl.NewVector3(-0.15, 0.15, -0.5), 0.05, 0.3, 0.2, rl.Blue)

	// Cockpit
	rl.DrawCube(rl.NewVector3(0, 0.2, 0.3), 0.25, 0.15, 0.3, rl.SkyBlue)

	rl.PopMatrix()

	// Draw shadow on ground
	r.drawShadow(pos)
}

func (r *Renderer) drawRobotMode(m *Mech) {
	pos := m.Position
	rot := m.Rotation * 180.0 / math.Pi

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Legs
	rl.DrawCube(rl.NewVector3(0.2, 0.3, 0), 0.15, 0.6, 0.2, rl.Blue)
	rl.DrawCube(rl.NewVector3(-0.2, 0.3, 0), 0.15, 0.6, 0.2, rl.Blue)

	// Feet
	rl.DrawCube(rl.NewVector3(0.2, 0.05, 0.1), 0.18, 0.1, 0.35, rl.DarkBlue)
	rl.DrawCube(rl.NewVector3(-0.2, 0.05, 0.1), 0.18, 0.1, 0.35, rl.DarkBlue)

	// Torso
	rl.DrawCube(rl.NewVector3(0, 0.8, 0), 0.5, 0.4, 0.3, rl.Blue)
	rl.DrawCubeWires(rl.NewVector3(0, 0.8, 0), 0.5, 0.4, 0.3, rl.DarkBlue)

	// Head
	rl.DrawCube(rl.NewVector3(0, 1.1, 0), 0.25, 0.2, 0.2, rl.Blue)
	rl.DrawCube(rl.NewVector3(0, 1.1, 0.12), 0.2, 0.1, 0.05, rl.Red) // Visor

	// Arms
	rl.DrawCube(rl.NewVector3(0.35, 0.75, 0), 0.1, 0.35, 0.12, rl.Blue)
	rl.DrawCube(rl.NewVector3(-0.35, 0.75, 0), 0.1, 0.35, 0.12, rl.Blue)

	// Shoulder pads
	rl.DrawCube(rl.NewVector3(0.35, 0.95, 0), 0.2, 0.1, 0.2, rl.DarkBlue)
	rl.DrawCube(rl.NewVector3(-0.35, 0.95, 0), 0.2, 0.1, 0.2, rl.DarkBlue)

	// Gun on right arm
	rl.DrawCube(rl.NewVector3(0.35, 0.6, 0.15), 0.08, 0.08, 0.25, rl.Gray)

	rl.PopMatrix()
}

func (r *Renderer) drawTransforming(m *Mech) {
	pos := m.Position
	rot := m.Rotation * 180.0 / math.Pi
	t := m.TransformProgress

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Interpolate between forms
	var height, width, length float32
	var color rl.Color

	if m.Mode == ModeJet {
		// Jet -> Robot: compact and rise
		height = lerp(0.3, 1.0, t)
		width = lerp(1.4, 0.5, t)
		length = lerp(1.2, 0.3, t)
	} else {
		// Robot -> Jet: stretch and lower
		height = lerp(1.0, 0.3, t)
		width = lerp(0.5, 1.4, t)
		length = lerp(0.3, 1.2, t)
	}

	// Flash during transformation
	if int(t*10)%2 == 0 {
		color = rl.White
	} else {
		color = rl.Blue
	}

	// Draw morphing shape
	rl.DrawCube(rl.NewVector3(0, height/2, 0), width, height, length, color)
	rl.DrawCubeWires(rl.NewVector3(0, height/2, 0), width, height, length, rl.DarkBlue)

	rl.PopMatrix()

	// Draw shadow
	r.drawShadow(pos)
}

func (r *Renderer) drawShadow(pos rl.Vector3) {
	// Simple circular shadow on ground
	shadowY := float32(0.01) // Slightly above ground to avoid z-fighting
	shadowSize := float32(0.8)

	rl.DrawCylinder(
		rl.NewVector3(pos.X, shadowY, pos.Z),
		shadowSize, shadowSize, 0.01,
		16,
		rl.Color{R: 0, G: 0, B: 0, A: 64},
	)
}

func (r *Renderer) drawProjectiles(m *Mech) {
	for _, p := range m.Projectiles {
		if !p.Alive {
			continue
		}

		// Different colors based on damage (jet vs robot)
		var color rl.Color
		if p.Damage <= 15 {
			color = rl.Yellow // Jet mode projectiles
		} else {
			color = rl.Orange // Robot mode projectiles
		}

		// Draw projectile as a small sphere
		rl.DrawSphere(p.Position, 0.1, color)

		// Draw a small trail
		trailLen := float32(0.3)
		speed := float32(math.Sqrt(float64(p.Velocity.X*p.Velocity.X + p.Velocity.Z*p.Velocity.Z)))
		if speed > 0 {
			trailEnd := rl.Vector3{
				X: p.Position.X - p.Velocity.X/speed*trailLen,
				Y: p.Position.Y,
				Z: p.Position.Z - p.Velocity.Z/speed*trailLen,
			}
			rl.DrawLine3D(p.Position, trailEnd, color)
		}
	}
}

// DrawUI renders mech-related UI elements
func (r *Renderer) DrawUI(m *Mech, screenWidth, screenHeight int) {
	// Health bar
	barWidth := float32(200)
	barHeight := float32(20)
	barX := float32(10)
	barY := float32(screenHeight) - barHeight - 40

	// Background
	rl.DrawRectangle(int32(barX), int32(barY), int32(barWidth), int32(barHeight), rl.DarkGray)

	// Health fill
	healthPct := m.Health / m.MaxHealth
	fillWidth := barWidth * healthPct

	var healthColor rl.Color
	if healthPct > 0.6 {
		healthColor = rl.Green
	} else if healthPct > 0.3 {
		healthColor = rl.Yellow
	} else {
		healthColor = rl.Red
	}
	rl.DrawRectangle(int32(barX), int32(barY), int32(fillWidth), int32(barHeight), healthColor)

	// Border
	rl.DrawRectangleLines(int32(barX), int32(barY), int32(barWidth), int32(barHeight), rl.Black)

	// Health text
	healthText := fmt.Sprintf("HP: %.0f/%.0f", m.Health, m.MaxHealth)
	rl.DrawText(healthText, int32(barX), int32(barY-20), 15, rl.White)

	// Mode indicator
	var modeText string
	var modeColor rl.Color
	if m.Mode == ModeJet {
		modeText = "JET MODE"
		modeColor = rl.SkyBlue
	} else {
		modeText = "ROBOT MODE"
		modeColor = rl.Orange
	}

	if m.State == StateTransforming {
		modeText = "TRANSFORMING..."
		modeColor = rl.White
	}

	rl.DrawText(modeText, int32(barX), int32(barY-40), 20, modeColor)

	// Controls hint
	rl.DrawText("WASD: Move | SPACE: Shoot | T: Transform", 10, int32(screenHeight)-20, 15, rl.Gray)
}

func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}
