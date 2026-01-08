package unit

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Type represents the kind of unit
type Type int

const (
	TypeInfantry Type = iota
	TypeTank
	TypeAA       // Anti-air
	TypeArtillery
	TypeMotorcycle
	TypeSupplyTruck
)

// Owner represents who controls a unit
type Owner int

const (
	OwnerNeutral Owner = iota
	OwnerPlayer1
	OwnerPlayer2
)

// Order represents what the unit should do
type Order int

const (
	OrderNone Order = iota
	OrderAttackHQ        // Attack enemy HQ
	OrderAttackNearest   // Attack nearest enemy
	OrderCaptureOutpost  // Capture nearest outpost
	OrderDefendPosition  // Hold current position
	OrderPatrolArea      // Patrol around drop point
)

// OrderNames returns human-readable order names
func OrderNames() []string {
	return []string{
		"None",
		"Attack HQ",
		"Attack Nearest",
		"Capture Outpost",
		"Defend Position",
		"Patrol Area",
	}
}

// State represents the unit's current activity
type State int

const (
	StateIdle State = iota
	StateMoving
	StateAttacking
	StateBeingCarried
	StateDead
)

// Config holds unit configuration
type Config struct {
	Speed       float32
	Health      float32
	Damage      float32
	AttackRange float32
	AttackRate  float32 // attacks per second
	Size        float32 // collision radius
}

// DefaultConfigs returns configurations for each unit type
func DefaultConfigs() map[Type]Config {
	return map[Type]Config{
		TypeInfantry: {
			Speed:       2.0,
			Health:      30.0,
			Damage:      5.0,
			AttackRange: 3.0,
			AttackRate:  1.0,
			Size:        0.3,
		},
		TypeTank: {
			Speed:       3.0,
			Health:      100.0,
			Damage:      20.0,
			AttackRange: 8.0,
			AttackRate:  0.5,
			Size:        0.6,
		},
		TypeAA: {
			Speed:       2.5,
			Health:      50.0,
			Damage:      30.0, // High damage vs air
			AttackRange: 10.0,
			AttackRate:  1.5,
			Size:        0.5,
		},
		TypeArtillery: {
			Speed:       1.5,
			Health:      40.0,
			Damage:      40.0,
			AttackRange: 15.0,
			AttackRate:  0.3,
			Size:        0.6,
		},
		TypeMotorcycle: {
			Speed:       6.0,
			Health:      20.0,
			Damage:      8.0,
			AttackRange: 4.0,
			AttackRate:  2.0,
			Size:        0.4,
		},
		TypeSupplyTruck: {
			Speed:       3.0,
			Health:      60.0,
			Damage:      0.0, // No attack
			AttackRange: 0.0,
			AttackRate:  0.0,
			Size:        0.5,
		},
	}
}

// Unit represents a deployable combat unit
type Unit struct {
	ID   int
	Type Type

	// Ownership
	Owner Owner

	// Position and movement
	Position  rl.Vector3
	Velocity  rl.Vector3
	Rotation  float32 // Y-axis rotation in radians

	// State
	State State
	Order Order

	// Order-specific data
	OrderTarget   rl.Vector3 // Target position for orders
	PatrolCenter  rl.Vector3 // Center of patrol area
	PatrolRadius  float32

	// Stats (from config)
	Speed       float32
	MaxSpeed    float32
	Health      float32
	MaxHealth   float32
	Damage      float32
	AttackRange float32
	AttackRate  float32
	Size        float32

	// Combat
	AttackCooldown float32
	Target         *Unit // Current attack target
}

// New creates a new unit at the given position
func New(id int, unitType Type, position rl.Vector3, owner Owner) *Unit {
	configs := DefaultConfigs()
	cfg := configs[unitType]

	return &Unit{
		ID:          id,
		Type:        unitType,
		Owner:       owner,
		Position:    position,
		Velocity:    rl.Vector3{},
		Rotation:    0,
		State:       StateIdle,
		Order:       OrderNone,
		Speed:       cfg.Speed,
		MaxSpeed:    cfg.Speed,
		Health:      cfg.Health,
		MaxHealth:   cfg.Health,
		Damage:      cfg.Damage,
		AttackRange: cfg.AttackRange,
		AttackRate:  cfg.AttackRate,
		Size:        cfg.Size,
	}
}

// Update updates the unit state for the frame
func (u *Unit) Update(dt float32) {
	if u.State == StateDead || u.State == StateBeingCarried {
		return
	}

	// Update attack cooldown
	if u.AttackCooldown > 0 {
		u.AttackCooldown -= dt
	}

	// Execute current order
	u.executeOrder(dt)

	// Apply movement
	u.Position.X += u.Velocity.X * dt
	u.Position.Z += u.Velocity.Z * dt

	// Keep on ground
	u.Position.Y = 0

	// Update rotation to face movement direction
	if u.Velocity.X != 0 || u.Velocity.Z != 0 {
		u.Rotation = float32(math.Atan2(float64(u.Velocity.X), float64(u.Velocity.Z)))
	}
}

func (u *Unit) executeOrder(dt float32) {
	switch u.Order {
	case OrderNone:
		u.State = StateIdle
		u.Velocity = rl.Vector3{}

	case OrderAttackHQ, OrderAttackNearest:
		u.executeAttackOrder(dt)

	case OrderCaptureOutpost:
		u.executeCaptureOrder(dt)

	case OrderDefendPosition:
		u.executeDefendOrder(dt)

	case OrderPatrolArea:
		u.executePatrolOrder(dt)
	}
}

func (u *Unit) executeAttackOrder(dt float32) {
	// Move toward target position
	if u.moveToward(u.OrderTarget, dt) {
		// Reached target, attack if we have a target
		if u.Target != nil && u.AttackCooldown <= 0 {
			u.State = StateAttacking
			u.Target.TakeDamage(u.Damage)
			u.AttackCooldown = 1.0 / u.AttackRate
		}
	}
}

func (u *Unit) executeCaptureOrder(dt float32) {
	// Move toward capture target
	u.moveToward(u.OrderTarget, dt)
	// Capture logic handled by base system checking infantry in range
}

func (u *Unit) executeDefendOrder(dt float32) {
	// Stay near order target, attack enemies in range
	dist := distance(u.Position, u.OrderTarget)
	if dist > 2.0 {
		u.moveToward(u.OrderTarget, dt)
	} else {
		u.Velocity = rl.Vector3{}
		u.State = StateIdle
		// Attack logic handled externally
	}
}

func (u *Unit) executePatrolOrder(dt float32) {
	// Move around patrol center
	dist := distance(u.Position, u.PatrolCenter)

	if dist > u.PatrolRadius {
		// Move back toward center
		u.moveToward(u.PatrolCenter, dt)
	} else {
		// Wander within patrol area
		if u.State == StateIdle {
			// Pick a new random point in patrol area
			angle := float32(math.Pi * 2.0 * float64(rl.GetRandomValue(0, 100)) / 100.0)
			radius := float32(rl.GetRandomValue(0, int32(u.PatrolRadius*100))) / 100.0
			u.OrderTarget = rl.Vector3{
				X: u.PatrolCenter.X + radius*float32(math.Cos(float64(angle))),
				Y: 0,
				Z: u.PatrolCenter.Z + radius*float32(math.Sin(float64(angle))),
			}
		}
		u.moveToward(u.OrderTarget, dt)
	}
}

// moveToward moves the unit toward a target position
// Returns true if within attack range
func (u *Unit) moveToward(target rl.Vector3, dt float32) bool {
	dx := target.X - u.Position.X
	dz := target.Z - u.Position.Z
	dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))

	if dist < u.AttackRange {
		u.Velocity = rl.Vector3{}
		u.State = StateIdle
		return true
	}

	if dist > 0.1 {
		u.State = StateMoving
		u.Velocity = rl.Vector3{
			X: (dx / dist) * u.Speed,
			Y: 0,
			Z: (dz / dist) * u.Speed,
		}
	}
	return false
}

// SetOrder sets the unit's order with a target position
func (u *Unit) SetOrder(order Order, target rl.Vector3) {
	u.Order = order
	u.OrderTarget = target

	if order == OrderPatrolArea {
		u.PatrolCenter = target
		u.PatrolRadius = 5.0 // Default patrol radius
	}
}

// PickUp marks the unit as being carried
func (u *Unit) PickUp() {
	u.State = StateBeingCarried
	u.Velocity = rl.Vector3{}
}

// Drop places the unit at a position with an order
func (u *Unit) Drop(position rl.Vector3, order Order) {
	u.Position = position
	u.State = StateIdle
	u.SetOrder(order, position)
}

// TakeDamage applies damage to the unit
func (u *Unit) TakeDamage(amount float32) {
	u.Health -= amount
	if u.Health <= 0 {
		u.Health = 0
		u.State = StateDead
	}
}

// IsDead returns true if the unit has no health
func (u *Unit) IsDead() bool {
	return u.Health <= 0 || u.State == StateDead
}

// IsCarried returns true if the unit is being carried
func (u *Unit) IsCarried() bool {
	return u.State == StateBeingCarried
}

// GetOwnerColor returns the color for this unit's owner
func (u *Unit) GetOwnerColor() rl.Color {
	switch u.Owner {
	case OwnerPlayer1:
		return rl.Blue
	case OwnerPlayer2:
		return rl.Red
	default:
		return rl.Gray
	}
}

// GetTypeName returns the name of this unit type
func (u *Unit) GetTypeName() string {
	names := []string{
		"Infantry",
		"Tank",
		"AA",
		"Artillery",
		"Motorcycle",
		"Supply Truck",
	}
	if int(u.Type) < len(names) {
		return names[u.Type]
	}
	return "Unknown"
}

// Helper functions

func distance(a, b rl.Vector3) float32 {
	dx := b.X - a.X
	dz := b.Z - a.Z
	return float32(math.Sqrt(float64(dx*dx + dz*dz)))
}
