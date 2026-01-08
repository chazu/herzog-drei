package mech

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Mode represents the mech's current form
type Mode int

const (
	ModeJet   Mode = iota // Fast flight, pick up/drop units, limited weapons
	ModeRobot             // Ground movement, full weapons, can enter buildings
)

// State represents what the mech is currently doing
type State int

const (
	StateIdle State = iota
	StateMoving
	StateTransforming
	StateShooting
	StateDead
)

// Config holds mech configuration values
type Config struct {
	// Movement
	JetSpeed         float32
	JetAcceleration  float32
	RobotSpeed       float32
	RobotAcceleration float32
	FlightHeight     float32

	// Combat
	JetFireRate      float32 // shots per second
	RobotFireRate    float32
	JetDamage        float32
	RobotDamage      float32
	ProjectileSpeed  float32

	// Health
	MaxHealth float32

	// Transformation
	TransformDuration float32 // seconds
}

// DefaultConfig returns the default mech configuration
func DefaultConfig() Config {
	return Config{
		JetSpeed:          15.0,
		JetAcceleration:   30.0,
		RobotSpeed:        5.0,
		RobotAcceleration: 20.0,
		FlightHeight:      3.0,

		JetFireRate:     2.0,
		RobotFireRate:   5.0,
		JetDamage:       10.0,
		RobotDamage:     25.0,
		ProjectileSpeed: 30.0,

		MaxHealth: 100.0,

		TransformDuration: 0.5,
	}
}

// Projectile represents a bullet/missile fired by the mech
type Projectile struct {
	Position  rl.Vector3
	Velocity  rl.Vector3
	Damage    float32
	Alive     bool
	LifeTime  float32
	MaxLife   float32
}

// Mech represents the player's transforming mech
type Mech struct {
	Config Config

	// Position and movement
	Position rl.Vector3
	Velocity rl.Vector3
	Rotation float32 // Y-axis rotation in radians

	// State
	Mode  Mode
	State State

	// Health
	Health    float32
	MaxHealth float32

	// Combat
	FireCooldown float32
	Projectiles  []Projectile

	// Transformation
	TransformProgress float32 // 0.0 to 1.0, used for animation

	// Input state (set by player input)
	InputMove  rl.Vector2 // normalized movement input (x, z)
	InputShoot bool
	InputTransform bool
}

// New creates a new mech at the given position
func New(pos rl.Vector3, cfg Config) *Mech {
	return &Mech{
		Config:    cfg,
		Position:  pos,
		Velocity:  rl.Vector3{},
		Rotation:  0,
		Mode:      ModeJet,
		State:     StateIdle,
		Health:    cfg.MaxHealth,
		MaxHealth: cfg.MaxHealth,
		Projectiles: make([]Projectile, 0, 32),
	}
}

// Update updates the mech state for the frame
func (m *Mech) Update(dt float32) {
	if m.State == StateDead {
		return
	}

	// Update transformation
	if m.State == StateTransforming {
		m.updateTransformation(dt)
		return
	}

	// Check for transform input
	if m.InputTransform && m.State != StateTransforming {
		m.startTransformation()
		return
	}

	// Update movement based on mode
	if m.Mode == ModeJet {
		m.updateJetMovement(dt)
	} else {
		m.updateRobotMovement(dt)
	}

	// Update shooting
	m.updateShooting(dt)

	// Update projectiles
	m.updateProjectiles(dt)

	// Update state
	m.updateState()
}

func (m *Mech) updateJetMovement(dt float32) {
	// Jet mode: 8-directional flight at fixed height
	targetVelX := m.InputMove.X * m.Config.JetSpeed
	targetVelZ := m.InputMove.Y * m.Config.JetSpeed

	// Smooth acceleration
	accel := m.Config.JetAcceleration * dt
	m.Velocity.X = approach(m.Velocity.X, targetVelX, accel)
	m.Velocity.Z = approach(m.Velocity.Z, targetVelZ, accel)

	// Maintain flight height
	targetY := m.Config.FlightHeight
	m.Velocity.Y = (targetY - m.Position.Y) * 5.0

	// Apply velocity
	m.Position.X += m.Velocity.X * dt
	m.Position.Y += m.Velocity.Y * dt
	m.Position.Z += m.Velocity.Z * dt

	// Update rotation to face movement direction
	if m.InputMove.X != 0 || m.InputMove.Y != 0 {
		targetRotation := float32(math.Atan2(float64(m.InputMove.X), float64(m.InputMove.Y)))
		m.Rotation = lerpAngle(m.Rotation, targetRotation, 10.0*dt)
	}
}

func (m *Mech) updateRobotMovement(dt float32) {
	// Robot mode: ground movement
	targetVelX := m.InputMove.X * m.Config.RobotSpeed
	targetVelZ := m.InputMove.Y * m.Config.RobotSpeed

	// Smooth acceleration
	accel := m.Config.RobotAcceleration * dt
	m.Velocity.X = approach(m.Velocity.X, targetVelX, accel)
	m.Velocity.Z = approach(m.Velocity.Z, targetVelZ, accel)

	// Stay on ground
	m.Velocity.Y = -m.Position.Y * 5.0 // Smoothly approach ground

	// Apply velocity
	m.Position.X += m.Velocity.X * dt
	m.Position.Y += m.Velocity.Y * dt
	m.Position.Z += m.Velocity.Z * dt

	// Clamp to ground
	if m.Position.Y < 0 {
		m.Position.Y = 0
		m.Velocity.Y = 0
	}

	// Update rotation to face movement direction
	if m.InputMove.X != 0 || m.InputMove.Y != 0 {
		targetRotation := float32(math.Atan2(float64(m.InputMove.X), float64(m.InputMove.Y)))
		m.Rotation = lerpAngle(m.Rotation, targetRotation, 8.0*dt)
	}
}

func (m *Mech) startTransformation() {
	m.State = StateTransforming
	m.TransformProgress = 0.0
	m.Velocity = rl.Vector3{} // Stop movement during transform
}

func (m *Mech) updateTransformation(dt float32) {
	m.TransformProgress += dt / m.Config.TransformDuration

	if m.TransformProgress >= 1.0 {
		m.TransformProgress = 0.0
		m.State = StateIdle

		// Toggle mode
		if m.Mode == ModeJet {
			m.Mode = ModeRobot
		} else {
			m.Mode = ModeJet
		}
	}
}

func (m *Mech) updateShooting(dt float32) {
	// Decrease fire cooldown
	if m.FireCooldown > 0 {
		m.FireCooldown -= dt
	}

	// Check if we can fire
	if !m.InputShoot || m.FireCooldown > 0 {
		return
	}

	// Get fire rate based on mode
	var fireRate, damage float32
	if m.Mode == ModeJet {
		fireRate = m.Config.JetFireRate
		damage = m.Config.JetDamage
	} else {
		fireRate = m.Config.RobotFireRate
		damage = m.Config.RobotDamage
	}

	// Fire projectile
	m.FireCooldown = 1.0 / fireRate

	// Calculate projectile direction (forward)
	direction := rl.Vector3{
		X: float32(math.Sin(float64(m.Rotation))),
		Y: 0,
		Z: float32(math.Cos(float64(m.Rotation))),
	}

	// Spawn projectile slightly in front of mech
	spawnOffset := float32(0.5)
	proj := Projectile{
		Position: rl.Vector3{
			X: m.Position.X + direction.X*spawnOffset,
			Y: m.Position.Y + 0.5,
			Z: m.Position.Z + direction.Z*spawnOffset,
		},
		Velocity: rl.Vector3{
			X: direction.X * m.Config.ProjectileSpeed,
			Y: 0,
			Z: direction.Z * m.Config.ProjectileSpeed,
		},
		Damage:   damage,
		Alive:    true,
		LifeTime: 0,
		MaxLife:  3.0, // 3 seconds before despawn
	}

	m.Projectiles = append(m.Projectiles, proj)
}

func (m *Mech) updateProjectiles(dt float32) {
	for i := range m.Projectiles {
		if !m.Projectiles[i].Alive {
			continue
		}

		// Update position
		m.Projectiles[i].Position.X += m.Projectiles[i].Velocity.X * dt
		m.Projectiles[i].Position.Y += m.Projectiles[i].Velocity.Y * dt
		m.Projectiles[i].Position.Z += m.Projectiles[i].Velocity.Z * dt

		// Update lifetime
		m.Projectiles[i].LifeTime += dt
		if m.Projectiles[i].LifeTime >= m.Projectiles[i].MaxLife {
			m.Projectiles[i].Alive = false
		}
	}

	// Remove dead projectiles (compact slice)
	alive := m.Projectiles[:0]
	for _, p := range m.Projectiles {
		if p.Alive {
			alive = append(alive, p)
		}
	}
	m.Projectiles = alive
}

func (m *Mech) updateState() {
	if m.Health <= 0 {
		m.State = StateDead
		return
	}

	if m.FireCooldown > 0 {
		m.State = StateShooting
		return
	}

	speed := float32(math.Sqrt(float64(m.Velocity.X*m.Velocity.X + m.Velocity.Z*m.Velocity.Z)))
	if speed > 0.1 {
		m.State = StateMoving
	} else {
		m.State = StateIdle
	}
}

// TakeDamage applies damage to the mech
func (m *Mech) TakeDamage(amount float32) {
	m.Health -= amount
	if m.Health < 0 {
		m.Health = 0
	}
}

// Heal restores health to the mech
func (m *Mech) Heal(amount float32) {
	m.Health += amount
	if m.Health > m.MaxHealth {
		m.Health = m.MaxHealth
	}
}

// IsDead returns true if the mech has no health
func (m *Mech) IsDead() bool {
	return m.Health <= 0
}

// GetForward returns the forward direction vector
func (m *Mech) GetForward() rl.Vector3 {
	return rl.Vector3{
		X: float32(math.Sin(float64(m.Rotation))),
		Y: 0,
		Z: float32(math.Cos(float64(m.Rotation))),
	}
}

// Helper functions

func approach(current, target, delta float32) float32 {
	diff := target - current
	if diff > delta {
		return current + delta
	}
	if diff < -delta {
		return current - delta
	}
	return target
}

func lerpAngle(a, b, t float32) float32 {
	// Normalize angles to -PI to PI
	for a > math.Pi {
		a -= 2 * math.Pi
	}
	for a < -math.Pi {
		a += 2 * math.Pi
	}
	for b > math.Pi {
		b -= 2 * math.Pi
	}
	for b < -math.Pi {
		b += 2 * math.Pi
	}

	// Find shortest path
	diff := b - a
	if diff > math.Pi {
		diff -= 2 * math.Pi
	}
	if diff < -math.Pi {
		diff += 2 * math.Pi
	}

	return a + diff*t
}
