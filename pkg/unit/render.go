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

// Draw renders all units from a manager
func (r *Renderer) Draw(m *Manager) {
	for _, u := range m.GetUnits() {
		r.DrawUnit(u)
	}
}

// DrawUnit renders a single unit
func (r *Renderer) DrawUnit(u *Unit) {
	if u.IsDead() {
		r.drawDeadUnit(u)
		return
	}

	// Get colors based on team
	mainColor, trimColor := r.getTeamColors(u.Team)

	// Draw based on unit type
	switch u.Config.Type {
	case TypeInfantry:
		r.drawInfantry(u, mainColor, trimColor)
	case TypeTank:
		r.drawTank(u, mainColor, trimColor)
	case TypeMotorcycle:
		r.drawMotorcycle(u, mainColor, trimColor)
	case TypeSAM:
		r.drawSAM(u, mainColor, trimColor)
	case TypeBoat:
		r.drawBoat(u, mainColor, trimColor)
	case TypeSupply:
		r.drawSupply(u, mainColor, trimColor)
	}

	// Draw health bar
	r.drawHealthBar(u)

	// Draw attack effect if attacking
	if u.State == StateAttacking && u.Target != nil {
		r.drawAttackEffect(u)
	}
}

func (r *Renderer) getTeamColors(team Team) (main, trim rl.Color) {
	if team == TeamPlayer {
		return rl.Blue, rl.DarkBlue
	}
	return rl.Red, rl.Maroon
}

func (r *Renderer) drawInfantry(u *Unit, main, trim rl.Color) {
	pos := u.Position
	rot := u.Rotation * 180.0 / math.Pi

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Body
	rl.DrawCube(rl.NewVector3(0, 0.15, 0), 0.2, 0.3, 0.15, main)

	// Head
	rl.DrawSphere(rl.NewVector3(0, 0.35, 0), 0.08, main)

	// Legs
	rl.DrawCube(rl.NewVector3(0.05, 0.05, 0), 0.06, 0.15, 0.06, trim)
	rl.DrawCube(rl.NewVector3(-0.05, 0.05, 0), 0.06, 0.15, 0.06, trim)

	// Gun
	rl.DrawCube(rl.NewVector3(0.1, 0.15, 0.1), 0.03, 0.03, 0.15, rl.Gray)

	rl.PopMatrix()
}

func (r *Renderer) drawTank(u *Unit, main, trim rl.Color) {
	pos := u.Position
	rot := u.Rotation * 180.0 / math.Pi

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Track base
	rl.DrawCube(rl.NewVector3(0, 0.1, 0), 0.6, 0.2, 0.8, trim)

	// Hull
	rl.DrawCube(rl.NewVector3(0, 0.25, 0), 0.5, 0.15, 0.6, main)

	// Turret
	rl.DrawCube(rl.NewVector3(0, 0.4, -0.05), 0.35, 0.15, 0.35, main)

	// Barrel
	rl.DrawCube(rl.NewVector3(0, 0.4, 0.35), 0.08, 0.08, 0.5, rl.DarkGray)

	rl.PopMatrix()
}

func (r *Renderer) drawMotorcycle(u *Unit, main, trim rl.Color) {
	pos := u.Position
	rot := u.Rotation * 180.0 / math.Pi

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Frame
	rl.DrawCube(rl.NewVector3(0, 0.15, 0), 0.15, 0.1, 0.5, main)

	// Front wheel
	rl.DrawCylinder(rl.NewVector3(0, 0.08, 0.2), 0.08, 0.08, 0.05, 8, trim)

	// Back wheel
	rl.DrawCylinder(rl.NewVector3(0, 0.08, -0.2), 0.08, 0.08, 0.05, 8, trim)

	// Rider
	rl.DrawCube(rl.NewVector3(0, 0.25, -0.05), 0.12, 0.15, 0.15, main)
	rl.DrawSphere(rl.NewVector3(0, 0.38, -0.05), 0.06, main)

	rl.PopMatrix()
}

func (r *Renderer) drawSAM(u *Unit, main, trim rl.Color) {
	pos := u.Position
	rot := u.Rotation * 180.0 / math.Pi

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Vehicle base
	rl.DrawCube(rl.NewVector3(0, 0.1, 0), 0.5, 0.2, 0.7, trim)

	// Cab
	rl.DrawCube(rl.NewVector3(0, 0.25, -0.2), 0.4, 0.15, 0.25, main)

	// Launcher platform
	rl.DrawCube(rl.NewVector3(0, 0.25, 0.15), 0.35, 0.08, 0.3, main)

	// Missile tubes
	rl.DrawCylinder(rl.NewVector3(0.1, 0.35, 0.15), 0.05, 0.05, 0.25, 6, rl.Gray)
	rl.DrawCylinder(rl.NewVector3(-0.1, 0.35, 0.15), 0.05, 0.05, 0.25, 6, rl.Gray)

	rl.PopMatrix()
}

func (r *Renderer) drawBoat(u *Unit, main, trim rl.Color) {
	pos := u.Position
	rot := u.Rotation * 180.0 / math.Pi

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Hull
	rl.DrawCube(rl.NewVector3(0, 0, 0), 0.4, 0.15, 0.8, main)

	// Cabin
	rl.DrawCube(rl.NewVector3(0, 0.15, -0.1), 0.25, 0.15, 0.3, trim)

	// Gun turret
	rl.DrawCylinder(rl.NewVector3(0, 0.15, 0.2), 0.1, 0.1, 0.1, 8, rl.Gray)
	rl.DrawCube(rl.NewVector3(0, 0.2, 0.35), 0.04, 0.04, 0.2, rl.DarkGray)

	rl.PopMatrix()
}

func (r *Renderer) drawSupply(u *Unit, main, trim rl.Color) {
	pos := u.Position
	rot := u.Rotation * 180.0 / math.Pi

	rl.PushMatrix()
	rl.Translatef(pos.X, pos.Y, pos.Z)
	rl.Rotatef(rot, 0, 1, 0)

	// Truck chassis
	rl.DrawCube(rl.NewVector3(0, 0.1, 0), 0.4, 0.15, 0.8, trim)

	// Cab
	rl.DrawCube(rl.NewVector3(0, 0.25, -0.25), 0.35, 0.2, 0.25, main)

	// Cargo box
	rl.DrawCube(rl.NewVector3(0, 0.3, 0.15), 0.35, 0.25, 0.45, rl.Color{R: 80, G: 80, B: 60, A: 255})

	// Cross symbol (supply)
	rl.DrawCube(rl.NewVector3(0, 0.43, 0.15), 0.2, 0.05, 0.05, rl.White)
	rl.DrawCube(rl.NewVector3(0, 0.43, 0.15), 0.05, 0.05, 0.2, rl.White)

	rl.PopMatrix()
}

func (r *Renderer) drawDeadUnit(u *Unit) {
	// Draw wreckage
	pos := u.Position
	rl.DrawCube(pos, 0.3, 0.1, 0.3, rl.DarkGray)
	// Smoke effect would go here
}

func (r *Renderer) drawHealthBar(u *Unit) {
	// Position health bar above unit
	pos := u.Position
	pos.Y += 0.7

	// Health bar dimensions in world space
	barWidth := float32(0.5)
	barHeight := float32(0.08)

	// Background
	rl.DrawCube(pos, barWidth, barHeight, 0.01, rl.DarkGray)

	// Health fill
	healthPct := u.Health / u.MaxHealth
	if healthPct > 0 {
		fillWidth := barWidth * healthPct
		fillPos := pos
		fillPos.X -= (barWidth - fillWidth) / 2

		var healthColor rl.Color
		if healthPct > 0.6 {
			healthColor = rl.Green
		} else if healthPct > 0.3 {
			healthColor = rl.Yellow
		} else {
			healthColor = rl.Red
		}

		rl.DrawCube(fillPos, fillWidth, barHeight*0.8, 0.02, healthColor)
	}
}

func (r *Renderer) drawAttackEffect(u *Unit) {
	// Draw a line from unit to target when attacking
	if u.Target == nil {
		return
	}

	startPos := u.Position
	startPos.Y += 0.3

	endPos := u.Target.Position
	endPos.Y += 0.3

	// Flash color
	var color rl.Color
	if u.AttackCooldown > 0.9/u.Config.AttackRate {
		color = rl.Yellow
	} else {
		color = rl.Orange
	}

	rl.DrawLine3D(startPos, endPos, color)
}

// DrawUI renders unit-related UI elements
func (r *Renderer) DrawUI(m *Manager, screenWidth, screenHeight int) {
	// Unit count display
	playerCount := m.CountByTeam(TeamPlayer)
	enemyCount := m.CountByTeam(TeamEnemy)

	unitText := fmt.Sprintf("Units - Player: %d | Enemy: %d", playerCount, enemyCount)
	rl.DrawText(unitText, int32(screenWidth-200), 40, 15, rl.White)
}

// DrawDebugPath draws a unit's current path (for debugging)
func (r *Renderer) DrawDebugPath(u *Unit) {
	if len(u.Path) == 0 {
		return
	}

	// Draw path waypoints
	for i, wp := range u.Path {
		pos := rl.Vector3{X: wp.X, Y: 0.1, Z: wp.Y}

		// Current waypoint is highlighted
		if i == u.PathIndex {
			rl.DrawSphere(pos, 0.15, rl.Green)
		} else if i > u.PathIndex {
			rl.DrawSphere(pos, 0.1, rl.Yellow)
		}

		// Draw lines between waypoints
		if i < len(u.Path)-1 {
			nextWp := u.Path[i+1]
			nextPos := rl.Vector3{X: nextWp.X, Y: 0.1, Z: nextWp.Y}
			rl.DrawLine3D(pos, nextPos, rl.Yellow)
		}
	}

	// Draw objective
	if u.HasObjective {
		objPos := rl.Vector3{X: u.Objective.X, Y: 0.2, Z: u.Objective.Z}
		rl.DrawSphere(objPos, 0.2, rl.Magenta)
	}
}
