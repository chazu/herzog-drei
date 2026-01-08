package unit

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Renderer handles unit rendering
type Renderer struct{}

// NewRenderer creates a new unit renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Draw renders all units
func (r *Renderer) Draw(m *Manager) {
	for _, u := range m.Units {
		if u.IsCarried() || u.IsDead() {
			continue
		}
		r.drawUnit(u)
	}
}

func (r *Renderer) drawUnit(u *Unit) {
	pos := u.Position
	rot := u.Rotation * 180.0 / math.Pi
	color := u.GetOwnerColor()
	darkColor := darken(color)

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	switch u.Type {
	case TypeInfantry:
		r.drawInfantry(color, darkColor)
	case TypeTank:
		r.drawTank(color, darkColor)
	case TypeAA:
		r.drawAA(color, darkColor)
	case TypeArtillery:
		r.drawArtillery(color, darkColor)
	case TypeMotorcycle:
		r.drawMotorcycle(color, darkColor)
	case TypeSupplyTruck:
		r.drawSupplyTruck(color, darkColor)
	}

	rl.PopMatrix()

	// Draw health bar if damaged
	if u.Health < u.MaxHealth {
		r.drawHealthBar(u)
	}

	// Draw order indicator
	r.drawOrderIndicator(u)
}

func (r *Renderer) drawInfantry(color, darkColor rl.Color) {
	// Small humanoid shape
	// Body
	rl.DrawCube(rl.NewVector3(0, 0.25, 0), 0.15, 0.3, 0.1, color)
	// Head
	rl.DrawSphere(rl.NewVector3(0, 0.45, 0), 0.08, color)
	// Legs
	rl.DrawCube(rl.NewVector3(0.04, 0.08, 0), 0.04, 0.16, 0.04, darkColor)
	rl.DrawCube(rl.NewVector3(-0.04, 0.08, 0), 0.04, 0.16, 0.04, darkColor)
}

func (r *Renderer) drawTank(color, darkColor rl.Color) {
	// Tank body
	rl.DrawCube(rl.NewVector3(0, 0.15, 0), 0.6, 0.2, 0.4, darkColor)
	// Tank treads
	rl.DrawCube(rl.NewVector3(0.25, 0.08, 0), 0.1, 0.16, 0.5, rl.DarkGray)
	rl.DrawCube(rl.NewVector3(-0.25, 0.08, 0), 0.1, 0.16, 0.5, rl.DarkGray)
	// Turret
	rl.DrawCube(rl.NewVector3(0, 0.3, -0.05), 0.3, 0.15, 0.25, color)
	// Gun barrel
	rl.DrawCube(rl.NewVector3(0, 0.32, 0.25), 0.06, 0.06, 0.35, rl.Gray)
}

func (r *Renderer) drawAA(color, darkColor rl.Color) {
	// AA platform
	rl.DrawCube(rl.NewVector3(0, 0.1, 0), 0.4, 0.15, 0.4, darkColor)
	// AA gun mount
	rl.DrawCube(rl.NewVector3(0, 0.25, 0), 0.2, 0.1, 0.2, color)
	// Twin barrels pointing up
	rl.DrawCube(rl.NewVector3(0.05, 0.45, 0.1), 0.04, 0.3, 0.04, rl.Gray)
	rl.DrawCube(rl.NewVector3(-0.05, 0.45, 0.1), 0.04, 0.3, 0.04, rl.Gray)
}

func (r *Renderer) drawArtillery(color, darkColor rl.Color) {
	// Artillery base
	rl.DrawCube(rl.NewVector3(0, 0.12, 0), 0.5, 0.15, 0.35, darkColor)
	// Wheels
	rl.DrawCylinder(rl.NewVector3(0.2, 0.1, 0.2), 0.1, 0.1, 0.05, 8, rl.DarkGray)
	rl.DrawCylinder(rl.NewVector3(-0.2, 0.1, 0.2), 0.1, 0.1, 0.05, 8, rl.DarkGray)
	rl.DrawCylinder(rl.NewVector3(0.2, 0.1, -0.2), 0.1, 0.1, 0.05, 8, rl.DarkGray)
	rl.DrawCylinder(rl.NewVector3(-0.2, 0.1, -0.2), 0.1, 0.1, 0.05, 8, rl.DarkGray)
	// Gun platform
	rl.DrawCube(rl.NewVector3(0, 0.25, -0.05), 0.25, 0.1, 0.2, color)
	// Long barrel
	rl.DrawCube(rl.NewVector3(0, 0.3, 0.3), 0.08, 0.08, 0.5, rl.Gray)
}

func (r *Renderer) drawMotorcycle(color, darkColor rl.Color) {
	// Motorcycle body
	rl.DrawCube(rl.NewVector3(0, 0.15, 0), 0.15, 0.15, 0.4, color)
	// Wheels
	rl.DrawCylinder(rl.NewVector3(0, 0.1, 0.2), 0.1, 0.1, 0.08, 8, rl.DarkGray)
	rl.DrawCylinder(rl.NewVector3(0, 0.1, -0.15), 0.08, 0.08, 0.08, 8, rl.DarkGray)
	// Rider
	rl.DrawCube(rl.NewVector3(0, 0.28, -0.05), 0.1, 0.15, 0.1, darkColor)
	rl.DrawSphere(rl.NewVector3(0, 0.4, -0.05), 0.06, darkColor)
}

func (r *Renderer) drawSupplyTruck(color, darkColor rl.Color) {
	// Truck cab
	rl.DrawCube(rl.NewVector3(0, 0.2, 0.2), 0.35, 0.25, 0.25, color)
	// Truck bed
	rl.DrawCube(rl.NewVector3(0, 0.15, -0.15), 0.4, 0.2, 0.4, darkColor)
	// Wheels
	rl.DrawCylinder(rl.NewVector3(0.18, 0.08, 0.2), 0.08, 0.08, 0.06, 8, rl.DarkGray)
	rl.DrawCylinder(rl.NewVector3(-0.18, 0.08, 0.2), 0.08, 0.08, 0.06, 8, rl.DarkGray)
	rl.DrawCylinder(rl.NewVector3(0.18, 0.08, -0.25), 0.08, 0.08, 0.06, 8, rl.DarkGray)
	rl.DrawCylinder(rl.NewVector3(-0.18, 0.08, -0.25), 0.08, 0.08, 0.06, 8, rl.DarkGray)
	// Cargo
	rl.DrawCube(rl.NewVector3(0, 0.3, -0.15), 0.3, 0.1, 0.3, rl.Brown)
}

func (r *Renderer) drawHealthBar(u *Unit) {
	// Small health bar above unit
	barWidth := float32(0.6)
	barHeight := float32(0.08)
	y := float32(0.8) // Height above unit

	// Background
	rl.DrawCube(
		rl.NewVector3(u.Position.X, u.Position.Y+y, u.Position.Z),
		barWidth, barHeight, 0.02,
		rl.DarkGray,
	)

	// Health fill
	healthPct := u.Health / u.MaxHealth
	fillWidth := barWidth * healthPct

	var healthColor rl.Color
	if healthPct > 0.6 {
		healthColor = rl.Green
	} else if healthPct > 0.3 {
		healthColor = rl.Yellow
	} else {
		healthColor = rl.Red
	}

	offsetX := (barWidth - fillWidth) / 2
	rl.DrawCube(
		rl.NewVector3(u.Position.X-offsetX, u.Position.Y+y, u.Position.Z+0.01),
		fillWidth, barHeight, 0.02,
		healthColor,
	)
}

func (r *Renderer) drawOrderIndicator(u *Unit) {
	if u.Order == OrderNone {
		return
	}

	// Draw a small icon/indicator above the unit showing its order
	indicatorY := u.Position.Y + 1.0

	var indicatorColor rl.Color
	switch u.Order {
	case OrderAttackHQ:
		indicatorColor = rl.Red
	case OrderAttackNearest:
		indicatorColor = rl.Orange
	case OrderCaptureOutpost:
		indicatorColor = rl.Green
	case OrderDefendPosition:
		indicatorColor = rl.SkyBlue
	case OrderPatrolArea:
		indicatorColor = rl.Yellow
	}

	// Small sphere indicator
	rl.DrawSphere(
		rl.NewVector3(u.Position.X, indicatorY, u.Position.Z),
		0.08,
		indicatorColor,
	)
}

// DrawUI renders unit-related UI elements
func (r *Renderer) DrawUI(m *Manager, screenWidth, screenHeight int) {
	// Unit counts
	y := int32(100)
	rl.DrawText("UNITS:", 10, y, 15, rl.White)

	p1Count := m.Count(OwnerPlayer1)
	p2Count := m.Count(OwnerPlayer2)

	rl.DrawText(fmt.Sprintf("Player 1: %d", p1Count), 10, y+20, 15, rl.Blue)
	rl.DrawText(fmt.Sprintf("Player 2: %d", p2Count), 10, y+40, 15, rl.Red)
}

func darken(c rl.Color) rl.Color {
	return rl.Color{
		R: uint8(float32(c.R) * 0.6),
		G: uint8(float32(c.G) * 0.6),
		B: uint8(float32(c.B) * 0.6),
		A: c.A,
	}
}
