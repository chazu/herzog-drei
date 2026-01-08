package base

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Renderer handles base rendering
type Renderer struct{}

// NewRenderer creates a new base renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Draw renders all bases
func (r *Renderer) Draw(mgr *Manager) {
	for _, base := range mgr.Bases {
		if base.IsDestroyed() {
			r.drawDestroyed(base)
		} else if base.Type == TypeHQ {
			r.drawHQ(base)
		} else {
			r.drawOutpost(base)
		}
	}
}

func (r *Renderer) drawHQ(b *Base) {
	pos := b.Position
	ownerColor := b.GetOwnerColor()

	// HQ is a larger, more elaborate structure
	// Main building
	rl.DrawCube(pos, 4.0, 3.0, 4.0, ownerColor)
	rl.DrawCubeWires(pos, 4.0, 3.0, 4.0, rl.Black)

	// Roof/tower
	roofPos := rl.Vector3{X: pos.X, Y: pos.Y + 2.5, Z: pos.Z}
	rl.DrawCube(roofPos, 2.5, 2.0, 2.5, darkenColor(ownerColor))
	rl.DrawCubeWires(roofPos, 2.5, 2.0, 2.5, rl.Black)

	// Antenna/flag pole
	antennaPos := rl.Vector3{X: pos.X, Y: pos.Y + 4.5, Z: pos.Z}
	rl.DrawCylinder(antennaPos, 0.1, 0.1, 2.0, 8, rl.DarkGray)

	// Flag at top (colored by owner)
	flagPos := rl.Vector3{X: pos.X + 0.3, Y: pos.Y + 5.0, Z: pos.Z}
	rl.DrawCube(flagPos, 0.6, 0.4, 0.05, ownerColor)

	// Door
	doorPos := rl.Vector3{X: pos.X, Y: pos.Y - 0.5, Z: pos.Z + 2.01}
	rl.DrawCube(doorPos, 1.0, 2.0, 0.1, rl.DarkGray)

	// Windows
	windowColor := rl.Color{R: 150, G: 200, B: 255, A: 200}
	for _, wx := range []float32{-1.2, 1.2} {
		windowPos := rl.Vector3{X: pos.X + wx, Y: pos.Y + 0.5, Z: pos.Z + 2.01}
		rl.DrawCube(windowPos, 0.6, 0.8, 0.05, windowColor)
	}

	// Health bar above
	r.drawHealthBar(b, 4.0)

	// Draw spawn point indicator
	r.drawSpawnPoint(b)
}

func (r *Renderer) drawOutpost(b *Base) {
	pos := b.Position
	ownerColor := b.GetOwnerColor()

	// Outpost is a smaller structure
	// Main building
	rl.DrawCube(pos, 2.0, 1.5, 2.0, ownerColor)
	rl.DrawCubeWires(pos, 2.0, 1.5, 2.0, rl.Black)

	// Small roof
	roofPos := rl.Vector3{X: pos.X, Y: pos.Y + 1.25, Z: pos.Z}
	rl.DrawCube(roofPos, 2.2, 0.5, 2.2, darkenColor(ownerColor))
	rl.DrawCubeWires(roofPos, 2.2, 0.5, 2.2, rl.Black)

	// Flag pole (smaller than HQ)
	polePos := rl.Vector3{X: pos.X, Y: pos.Y + 2.0, Z: pos.Z}
	rl.DrawCylinder(polePos, 0.05, 0.05, 1.0, 8, rl.DarkGray)

	// Small flag
	flagPos := rl.Vector3{X: pos.X + 0.2, Y: pos.Y + 2.3, Z: pos.Z}
	rl.DrawCube(flagPos, 0.4, 0.25, 0.03, ownerColor)

	// Health bar
	r.drawHealthBar(b, 2.5)

	// Capture progress bar (if being captured)
	if b.CaptureProgress > 0 {
		r.drawCaptureBar(b)
	}

	// Draw spawn point indicator
	r.drawSpawnPoint(b)
}

func (r *Renderer) drawDestroyed(b *Base) {
	pos := b.Position

	// Draw rubble
	rubbleColor := rl.Color{R: 80, G: 80, B: 80, A: 255}

	// Scattered debris cubes
	for i := 0; i < 5; i++ {
		offset := rl.Vector3{
			X: float32(i%3-1) * 0.8,
			Y: 0.2,
			Z: float32(i/3) * 0.6 - 0.3,
		}
		debrisPos := rl.Vector3{
			X: pos.X + offset.X,
			Y: pos.Y + offset.Y,
			Z: pos.Z + offset.Z,
		}
		size := 0.3 + float32(i%3)*0.2
		rl.DrawCube(debrisPos, size, size*0.5, size, rubbleColor)
	}
}

func (r *Renderer) drawHealthBar(b *Base, yOffset float32) {
	pos := b.Position
	barWidth := float32(2.0)
	barHeight := float32(0.15)

	barPos := rl.Vector3{
		X: pos.X,
		Y: pos.Y + yOffset,
		Z: pos.Z,
	}

	// Background
	rl.DrawCube(barPos, barWidth, barHeight, 0.1, rl.DarkGray)

	// Health fill
	healthPct := b.Health / b.MaxHealth
	fillWidth := barWidth * healthPct

	var healthColor rl.Color
	if healthPct > 0.6 {
		healthColor = rl.Green
	} else if healthPct > 0.3 {
		healthColor = rl.Yellow
	} else {
		healthColor = rl.Red
	}

	fillPos := rl.Vector3{
		X: pos.X - (barWidth-fillWidth)/2,
		Y: pos.Y + yOffset,
		Z: pos.Z + 0.05,
	}
	rl.DrawCube(fillPos, fillWidth, barHeight, 0.05, healthColor)
}

func (r *Renderer) drawCaptureBar(b *Base) {
	pos := b.Position
	barWidth := float32(2.0)
	barHeight := float32(0.1)

	barPos := rl.Vector3{
		X: pos.X,
		Y: pos.Y + 2.0,
		Z: pos.Z,
	}

	// Background
	rl.DrawCube(barPos, barWidth, barHeight, 0.1, rl.DarkGray)

	// Capture progress fill
	fillWidth := barWidth * b.CaptureProgress

	// Color based on who is capturing
	var captureColor rl.Color
	switch b.CapturingOwner {
	case OwnerPlayer1:
		captureColor = rl.Blue
	case OwnerPlayer2:
		captureColor = rl.Red
	default:
		captureColor = rl.White
	}

	fillPos := rl.Vector3{
		X: pos.X - (barWidth-fillWidth)/2,
		Y: pos.Y + 2.0,
		Z: pos.Z + 0.05,
	}
	rl.DrawCube(fillPos, fillWidth, barHeight, 0.05, captureColor)
}

func (r *Renderer) drawSpawnPoint(b *Base) {
	if b.Owner == OwnerNeutral {
		return // Neutral bases don't show spawn points
	}

	sp := b.SpawnPoint
	ownerColor := b.GetOwnerColor()

	// Draw a small marker at spawn point
	rl.DrawCylinder(sp, 0.3, 0.3, 0.05, 16, lightenColor(ownerColor))
	rl.DrawCylinderWires(sp, 0.3, 0.3, 0.05, 16, ownerColor)
}

// DrawUI renders base-related UI elements
func (r *Renderer) DrawUI(mgr *Manager, screenWidth, screenHeight int) {
	// Credits display for Player 1
	creditsText := fmt.Sprintf("Credits: $%.0f", mgr.Player1.Credits)
	rl.DrawText(creditsText, int32(screenWidth)-150, 10, 20, rl.Yellow)

	// Base count
	p1Bases := len(mgr.GetBasesOwnedBy(OwnerPlayer1))
	p2Bases := len(mgr.GetBasesOwnedBy(OwnerPlayer2))
	neutralBases := len(mgr.GetBasesOwnedBy(OwnerNeutral))

	baseText := fmt.Sprintf("Bases: P1:%d  N:%d  P2:%d", p1Bases, neutralBases, p2Bases)
	rl.DrawText(baseText, int32(screenWidth)-200, 35, 15, rl.White)

	// Game over check
	loser := mgr.IsGameOver()
	if loser != OwnerNeutral {
		var winText string
		if loser == OwnerPlayer1 {
			winText = "PLAYER 2 WINS!"
		} else {
			winText = "PLAYER 1 WINS!"
		}
		textWidth := rl.MeasureText(winText, 40)
		rl.DrawText(winText, int32(screenWidth/2)-textWidth/2, int32(screenHeight/2)-20, 40, rl.Gold)
	}
}

// Helper color functions

func darkenColor(c rl.Color) rl.Color {
	return rl.Color{
		R: uint8(float32(c.R) * 0.6),
		G: uint8(float32(c.G) * 0.6),
		B: uint8(float32(c.B) * 0.6),
		A: c.A,
	}
}

func lightenColor(c rl.Color) rl.Color {
	return rl.Color{
		R: uint8(min(255, int(float32(c.R)*1.3))),
		G: uint8(min(255, int(float32(c.G)*1.3))),
		B: uint8(min(255, int(float32(c.B)*1.3))),
		A: c.A,
	}
}
