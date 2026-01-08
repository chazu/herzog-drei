package combat

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/mech"
	"github.com/chazu/herzog-drei/pkg/unit"
)

// Config holds combat system configuration
type Config struct {
	// Collision
	ProjectileRadius float32 // Collision radius for projectiles
	UnitHitboxRadius float32 // Default hitbox radius for units
	MechHitboxRadius float32 // Hitbox radius for mech

	// Respawn
	MechRespawnDelay float32 // Seconds before mech respawns
	MechSpawnInvuln  float32 // Seconds of invulnerability after spawn

	// Effects
	ExplosionDuration float32
}

// DefaultConfig returns default combat configuration
func DefaultConfig() Config {
	return Config{
		ProjectileRadius: 0.15,
		UnitHitboxRadius: 0.5,
		MechHitboxRadius: 0.6,
		MechRespawnDelay: 3.0,
		MechSpawnInvuln:  2.0,
		ExplosionDuration: 0.5,
	}
}

// Explosion represents a visual explosion effect
type Explosion struct {
	Position rl.Vector3
	Radius   float32
	MaxRadius float32
	Duration float32
	Elapsed  float32
	Color    rl.Color
	Active   bool
}

// System manages combat interactions
type System struct {
	Config Config

	// Effects
	explosions []Explosion

	// Mech respawn
	mechDead        bool
	respawnTimer    float32
	invulnTimer     float32
	respawnPosition rl.Vector3
}

// NewSystem creates a new combat system
func NewSystem(cfg Config) *System {
	return &System{
		Config:     cfg,
		explosions: make([]Explosion, 0, 32),
	}
}

// SetRespawnPosition sets where the mech will respawn
func (s *System) SetRespawnPosition(pos rl.Vector3) {
	s.respawnPosition = pos
}

// Update runs combat checks and updates effects
func (s *System) Update(dt float32, playerMech *mech.Mech, unitMgr *unit.Manager) {
	// Handle mech respawn
	s.updateMechRespawn(dt, playerMech)

	// Skip combat checks if mech is dead or invulnerable
	if playerMech.IsDead() {
		return
	}

	// Check mech projectiles vs enemy units
	s.checkProjectileUnitCollisions(playerMech, unitMgr)

	// Check unit attacks vs mech (if not invulnerable)
	if s.invulnTimer <= 0 {
		s.checkUnitMechCollisions(playerMech, unitMgr)
	}

	// Update effects
	s.updateExplosions(dt)
}

// checkProjectileUnitCollisions checks mech projectiles hitting units
func (s *System) checkProjectileUnitCollisions(playerMech *mech.Mech, unitMgr *unit.Manager) {
	enemies := unitMgr.GetUnitsByTeam(unit.TeamEnemy)

	for i := range playerMech.Projectiles {
		proj := &playerMech.Projectiles[i]
		if !proj.Alive {
			continue
		}

		for _, enemy := range enemies {
			if enemy.IsDead() {
				continue
			}

			// Check collision
			dist := distance3D(proj.Position, enemy.Position)
			hitRadius := s.Config.ProjectileRadius + s.Config.UnitHitboxRadius

			if dist <= hitRadius {
				// Hit! Apply damage
				enemy.TakeDamage(proj.Damage)
				proj.Alive = false

				// Spawn hit effect
				s.spawnHitEffect(proj.Position)

				// Spawn explosion if enemy died
				if enemy.IsDead() {
					s.spawnExplosion(enemy.Position, 1.0, rl.Orange)
				}
				break
			}
		}
	}
}

// checkUnitMechCollisions checks if units are attacking the mech
func (s *System) checkUnitMechCollisions(playerMech *mech.Mech, unitMgr *unit.Manager) {
	enemies := unitMgr.GetEnemiesInRadius(playerMech.Position, 10.0, unit.TeamPlayer)

	for _, enemy := range enemies {
		if enemy.IsDead() {
			continue
		}

		// Check if enemy can attack air (jet mode) or ground (robot mode)
		isAir := playerMech.Mode == mech.ModeJet
		if isAir && !enemy.Config.CanAttackAir {
			continue
		}
		if !isAir && !enemy.Config.CanAttackGround {
			continue
		}

		// Check range
		dist := distance3D(enemy.Position, playerMech.Position)
		if dist > enemy.Config.AttackRange {
			continue
		}

		// Attack if cooldown ready (using existing unit attack rate)
		if enemy.AttackCooldown <= 0 {
			playerMech.TakeDamage(enemy.Config.AttackDamage)
			enemy.AttackCooldown = 1.0 / enemy.Config.AttackRate

			// Spawn small hit effect
			s.spawnHitEffect(playerMech.Position)

			// Check if mech died
			if playerMech.IsDead() {
				s.onMechDeath(playerMech)
				return
			}
		}
	}
}

// onMechDeath handles mech death
func (s *System) onMechDeath(playerMech *mech.Mech) {
	s.mechDead = true
	s.respawnTimer = s.Config.MechRespawnDelay

	// Big explosion
	s.spawnExplosion(playerMech.Position, 2.0, rl.Red)
}

// updateMechRespawn handles mech respawn timing
func (s *System) updateMechRespawn(dt float32, playerMech *mech.Mech) {
	// Update invulnerability timer
	if s.invulnTimer > 0 {
		s.invulnTimer -= dt
	}

	// Handle respawn timer
	if !s.mechDead {
		return
	}

	s.respawnTimer -= dt
	if s.respawnTimer <= 0 {
		s.respawnMech(playerMech)
	}
}

// respawnMech respawns the mech at the respawn position
func (s *System) respawnMech(playerMech *mech.Mech) {
	playerMech.Position = s.respawnPosition
	playerMech.Velocity = rl.Vector3{}
	playerMech.Health = playerMech.MaxHealth
	playerMech.Mode = mech.ModeJet
	playerMech.State = mech.StateIdle
	playerMech.Projectiles = playerMech.Projectiles[:0]

	s.mechDead = false
	s.invulnTimer = s.Config.MechSpawnInvuln
}

// IsMechDead returns true if mech is waiting to respawn
func (s *System) IsMechDead() bool {
	return s.mechDead
}

// GetRespawnTimer returns time until respawn
func (s *System) GetRespawnTimer() float32 {
	return s.respawnTimer
}

// IsMechInvulnerable returns true if mech has spawn protection
func (s *System) IsMechInvulnerable() bool {
	return s.invulnTimer > 0
}

// GetInvulnTimer returns remaining invulnerability time
func (s *System) GetInvulnTimer() float32 {
	return s.invulnTimer
}

// spawnExplosion creates an explosion effect
func (s *System) spawnExplosion(pos rl.Vector3, size float32, color rl.Color) {
	s.explosions = append(s.explosions, Explosion{
		Position:  pos,
		Radius:    0.1,
		MaxRadius: size,
		Duration:  s.Config.ExplosionDuration,
		Elapsed:   0,
		Color:     color,
		Active:    true,
	})
}

// spawnHitEffect creates a small hit particle effect
func (s *System) spawnHitEffect(pos rl.Vector3) {
	s.explosions = append(s.explosions, Explosion{
		Position:  pos,
		Radius:    0.05,
		MaxRadius: 0.3,
		Duration:  0.2,
		Elapsed:   0,
		Color:     rl.Yellow,
		Active:    true,
	})
}

// updateExplosions updates explosion animations
func (s *System) updateExplosions(dt float32) {
	active := s.explosions[:0]
	for i := range s.explosions {
		e := &s.explosions[i]
		if !e.Active {
			continue
		}

		e.Elapsed += dt
		if e.Elapsed >= e.Duration {
			e.Active = false
			continue
		}

		// Expand radius
		t := e.Elapsed / e.Duration
		e.Radius = e.MaxRadius * t

		active = append(active, *e)
	}
	s.explosions = active
}

// GetExplosions returns active explosions for rendering
func (s *System) GetExplosions() []Explosion {
	return s.explosions
}

// Helper functions

func distance3D(a, b rl.Vector3) float32 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	dz := b.Z - a.Z
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}
