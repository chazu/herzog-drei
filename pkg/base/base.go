package base

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/unit"
)

// Owner represents who controls a base
type Owner int

const (
	OwnerNeutral Owner = iota
	OwnerPlayer1
	OwnerPlayer2
)

// Type represents the kind of base
type Type int

const (
	TypeHQ      Type = iota // Main base, losing it = game over
	TypeOutpost             // Capturable, generates income, spawns units
)

// Config holds configuration for base behavior
type Config struct {
	// Income
	OutpostIncomeRate float32 // Credits per second for outposts
	HQIncomeRate      float32 // Credits per second for HQ

	// Capture
	CaptureTime float32 // Seconds to capture when fully occupied

	// Health
	HQMaxHealth      float32
	OutpostMaxHealth float32

	// Spawn
	SpawnCooldown float32 // Minimum time between spawns
}

// DefaultConfig returns the default base configuration
func DefaultConfig() Config {
	return Config{
		OutpostIncomeRate: 15.0, // Credits per second for outposts
		HQIncomeRate:      5.0,  // Credits per second for HQ
		CaptureTime:       5.0,
		HQMaxHealth:       500.0,
		OutpostMaxHealth:  200.0,
		SpawnCooldown:     2.0, // Slightly faster spawns
	}
}

// Base represents a capturable structure on the map
type Base struct {
	// Identity
	ID   int
	Type Type

	// Position
	Position rl.Vector3

	// Ownership
	Owner           Owner
	CaptureProgress float32 // 0.0 to 1.0, who is capturing
	CapturingOwner  Owner   // Who is currently capturing (if any)

	// State
	Health    float32
	MaxHealth float32

	// Economy
	IncomeRate     float32
	AccumulatedIncome float32

	// Spawning
	SpawnPoint    rl.Vector3      // Where units spawn
	SpawnCooldown float32         // Time until next spawn allowed
	SpawnQueue    []unit.UnitType // Units waiting to spawn

	// Infantry occupying this base (for capture mechanic)
	OccupyingInfantry int   // Count of infantry inside
	OccupyingOwner    Owner // Owner of occupying infantry
}

// NewBase creates a new base at the given position
func NewBase(id int, baseType Type, position rl.Vector3, owner Owner, cfg Config) *Base {
	var maxHealth, incomeRate float32
	if baseType == TypeHQ {
		maxHealth = cfg.HQMaxHealth
		incomeRate = cfg.HQIncomeRate
	} else {
		maxHealth = cfg.OutpostMaxHealth
		incomeRate = cfg.OutpostIncomeRate
	}

	// Spawn point is slightly in front of the base
	spawnPoint := rl.Vector3{
		X: position.X,
		Y: 0,
		Z: position.Z + 2.0,
	}

	return &Base{
		ID:            id,
		Type:          baseType,
		Position:      position,
		Owner:         owner,
		Health:        maxHealth,
		MaxHealth:     maxHealth,
		IncomeRate:    incomeRate,
		SpawnPoint:    spawnPoint,
		SpawnQueue:    make([]unit.UnitType, 0, 8),
	}
}

// Update updates the base state for the frame
func (b *Base) Update(dt float32, cfg Config) {
	// Generate income if owned
	if b.Owner != OwnerNeutral {
		b.AccumulatedIncome += b.IncomeRate * dt
	}

	// Update capture progress
	b.updateCapture(dt, cfg)

	// Update spawn cooldown
	if b.SpawnCooldown > 0 {
		b.SpawnCooldown -= dt
	}
}

func (b *Base) updateCapture(dt float32, cfg Config) {
	// Only outposts can be captured
	if b.Type == TypeHQ {
		return
	}

	// Check if infantry are occupying
	if b.OccupyingInfantry <= 0 {
		// No infantry, capture progress decays
		if b.CaptureProgress > 0 {
			b.CaptureProgress -= dt / cfg.CaptureTime
			if b.CaptureProgress < 0 {
				b.CaptureProgress = 0
				b.CapturingOwner = OwnerNeutral
			}
		}
		return
	}

	// Infantry present - update capture
	if b.OccupyingOwner == b.Owner {
		// Friendly infantry, no capture (but could heal?)
		return
	}

	// Enemy or neutral capturing
	if b.CapturingOwner != b.OccupyingOwner {
		// New capturer, reset progress
		b.CapturingOwner = b.OccupyingOwner
		b.CaptureProgress = 0
	}

	// Progress capture based on infantry count
	captureSpeed := float32(b.OccupyingInfantry) * dt / cfg.CaptureTime
	b.CaptureProgress += captureSpeed

	if b.CaptureProgress >= 1.0 {
		// Capture complete!
		b.Owner = b.CapturingOwner
		b.CaptureProgress = 0
		b.CapturingOwner = OwnerNeutral
		b.OccupyingInfantry = 0 // Infantry "merge" into garrison
		b.OccupyingOwner = OwnerNeutral
	}
}

// QueueUnit adds a unit to the spawn queue
func (b *Base) QueueUnit(unitType unit.UnitType) {
	if b.Owner == OwnerNeutral {
		return // Can't spawn from neutral bases
	}
	b.SpawnQueue = append(b.SpawnQueue, unitType)
}

// TrySpawn attempts to spawn the next unit in queue
// Returns the unit type and true if a spawn occurred
func (b *Base) TrySpawn(cfg Config) (unit.UnitType, bool) {
	if len(b.SpawnQueue) == 0 {
		return 0, false
	}
	if b.SpawnCooldown > 0 {
		return 0, false
	}

	// Spawn the unit
	unitType := b.SpawnQueue[0]
	b.SpawnQueue = b.SpawnQueue[1:]
	b.SpawnCooldown = cfg.SpawnCooldown

	return unitType, true
}

// CollectIncome collects and resets accumulated income
func (b *Base) CollectIncome() float32 {
	income := b.AccumulatedIncome
	b.AccumulatedIncome = 0
	return income
}

// TakeDamage applies damage to the base
func (b *Base) TakeDamage(amount float32) {
	b.Health -= amount
	if b.Health < 0 {
		b.Health = 0
	}
}

// IsDestroyed returns true if the base has no health
func (b *Base) IsDestroyed() bool {
	return b.Health <= 0
}

// SetOccupyingInfantry sets the infantry occupying this base for capture
func (b *Base) SetOccupyingInfantry(count int, owner Owner) {
	b.OccupyingInfantry = count
	b.OccupyingOwner = owner
}

// GetOwnerColor returns the color associated with the base's owner
func (b *Base) GetOwnerColor() rl.Color {
	switch b.Owner {
	case OwnerPlayer1:
		return rl.Blue
	case OwnerPlayer2:
		return rl.Red
	default:
		return rl.Gray
	}
}
