package unit

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Team represents which side a unit belongs to
type Team int

const (
	TeamPlayer Team = iota
	TeamEnemy
)

// UnitType identifies the kind of unit
type UnitType int

const (
	TypeInfantry UnitType = iota
	TypeTank
	TypeMotorcycle
	TypeSAM
	TypeBoat
	TypeSupply
)

// State represents what the unit is currently doing
type State int

const (
	StateIdle State = iota
	StateMoving
	StateAttacking
	StateDead
	StateCapturing    // Infantry only
	StateBeingCarried // Being transported by mech
)

// Order represents the unit's assigned behavior
type Order int

const (
	OrderNone Order = iota
	OrderAttackHQ       // Attack enemy HQ
	OrderAttackNearest  // Attack nearest enemy
	OrderCaptureOutpost // Capture nearest outpost
	OrderDefendPosition // Hold current position
	OrderPatrolArea     // Patrol around drop point
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

// Config holds unit configuration values
type Config struct {
	Type UnitType

	// Movement
	Speed         float32
	TurnSpeed     float32 // radians per second
	CanTraverseWater bool

	// Combat
	AttackRange   float32
	AttackDamage  float32
	AttackRate    float32 // attacks per second
	CanAttackAir  bool
	CanAttackGround bool

	// Health
	MaxHealth float32
	Armor     float32 // damage reduction 0-1

	// Special
	CanCapture bool // Infantry only
	Cost       int  // Resource cost to spawn
}

// Unit represents a deployable combat unit
type Unit struct {
	ID     uint32
	Config Config
	Team   Team

	// Position and movement
	Position rl.Vector3
	Velocity rl.Vector3
	Rotation float32 // Y-axis rotation in radians

	// State
	State State
	Order Order

	// Health
	Health    float32
	MaxHealth float32

	// Combat
	AttackCooldown float32
	Target         *Unit // Current attack target

	// AI
	Objective    rl.Vector3   // Where the unit is trying to go
	HasObjective bool
	Path         []rl.Vector2 // Pathfinding result (X, Z)
	PathIndex    int

	// Order-specific data
	OrderTarget  rl.Vector3 // Target position for orders
	PatrolCenter rl.Vector3 // Center of patrol area
	PatrolRadius float32
}

// New creates a new unit of the specified type
func New(id uint32, unitType UnitType, team Team, pos rl.Vector3) *Unit {
	cfg := GetConfig(unitType)
	return &Unit{
		ID:        id,
		Config:    cfg,
		Team:      team,
		Position:  pos,
		Velocity:  rl.Vector3{},
		Rotation:  0,
		State:     StateIdle,
		Health:    cfg.MaxHealth,
		MaxHealth: cfg.MaxHealth,
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

	// Execute order-based behavior if we have an order
	if u.Order != OrderNone {
		u.executeOrder(dt)
	} else if u.HasObjective && len(u.Path) > 0 && u.PathIndex < len(u.Path) {
		// Movement along path
		u.updateMovement(dt)
	} else if u.HasObjective {
		// Direct movement to objective (fallback when no path)
		u.moveToward(u.Objective, dt)
	}

	// Update state based on movement
	speed := float32(math.Sqrt(float64(u.Velocity.X*u.Velocity.X + u.Velocity.Z*u.Velocity.Z)))
	if u.State != StateAttacking && u.State != StateCapturing {
		if speed > 0.1 {
			u.State = StateMoving
		} else {
			u.State = StateIdle
		}
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
	if u.moveTowardOrder(u.OrderTarget, dt) {
		// Reached target, attack if we have a target
		if u.Target != nil && u.AttackCooldown <= 0 {
			u.State = StateAttacking
			u.Target.TakeDamage(u.Config.AttackDamage)
			u.AttackCooldown = 1.0 / u.Config.AttackRate
		}
	}
}

func (u *Unit) executeCaptureOrder(dt float32) {
	// Move toward capture target
	u.moveTowardOrder(u.OrderTarget, dt)
	// Capture logic handled by base system checking infantry in range
}

func (u *Unit) executeDefendOrder(dt float32) {
	// Stay near order target, attack enemies in range
	dist := u.DistanceToPoint(u.OrderTarget)
	if dist > 2.0 {
		u.moveTowardOrder(u.OrderTarget, dt)
	} else {
		u.Velocity = rl.Vector3{}
		u.State = StateIdle
		// Attack logic handled externally
	}
}

func (u *Unit) executePatrolOrder(dt float32) {
	// Move around patrol center
	dist := u.DistanceToPoint(u.PatrolCenter)

	if dist > u.PatrolRadius {
		// Move back toward center
		u.moveTowardOrder(u.PatrolCenter, dt)
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
		u.moveTowardOrder(u.OrderTarget, dt)
	}
}

// moveTowardOrder moves the unit toward a target position for order execution
// Returns true if within attack range
func (u *Unit) moveTowardOrder(target rl.Vector3, dt float32) bool {
	dx := target.X - u.Position.X
	dz := target.Z - u.Position.Z
	dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))

	if dist < u.Config.AttackRange {
		u.Velocity = rl.Vector3{}
		u.State = StateIdle
		return true
	}

	if dist > 0.1 {
		u.State = StateMoving
		u.Velocity = rl.Vector3{
			X: (dx / dist) * u.Config.Speed,
			Y: 0,
			Z: (dz / dist) * u.Config.Speed,
		}

		// Update position
		u.Position.X += u.Velocity.X * dt
		u.Position.Z += u.Velocity.Z * dt

		// Update rotation to face movement direction
		u.Rotation = float32(math.Atan2(float64(dx), float64(dz)))
	}
	return false
}

func (u *Unit) updateMovement(dt float32) {
	if u.PathIndex >= len(u.Path) {
		u.HasObjective = false
		u.Path = nil
		return
	}

	// Get current waypoint
	wp := u.Path[u.PathIndex]
	target := rl.Vector3{X: wp.X, Y: 0, Z: wp.Y}

	// Check if we reached the waypoint
	dx := target.X - u.Position.X
	dz := target.Z - u.Position.Z
	dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))

	if dist < 0.5 {
		u.PathIndex++
		if u.PathIndex >= len(u.Path) {
			u.HasObjective = false
			u.Path = nil
			u.Velocity = rl.Vector3{}
		}
		return
	}

	// Move toward waypoint
	u.moveToward(target, dt)
}

func (u *Unit) moveToward(target rl.Vector3, dt float32) {
	dx := target.X - u.Position.X
	dz := target.Z - u.Position.Z
	dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))

	if dist < 0.1 {
		u.Velocity = rl.Vector3{}
		return
	}

	// Calculate target rotation
	targetRot := float32(math.Atan2(float64(dx), float64(dz)))

	// Smoothly rotate toward target
	u.Rotation = lerpAngle(u.Rotation, targetRot, u.Config.TurnSpeed*dt)

	// Move forward if roughly facing the target
	angleDiff := math.Abs(float64(normalizeAngle(targetRot - u.Rotation)))
	if angleDiff < math.Pi/4 { // Within 45 degrees
		// Set velocity in facing direction
		u.Velocity.X = float32(math.Sin(float64(u.Rotation))) * u.Config.Speed
		u.Velocity.Z = float32(math.Cos(float64(u.Rotation))) * u.Config.Speed
	} else {
		// Slow down while turning
		u.Velocity.X *= 0.9
		u.Velocity.Z *= 0.9
	}

	// Apply velocity
	u.Position.X += u.Velocity.X * dt
	u.Position.Z += u.Velocity.Z * dt
}

// SetObjective sets a destination for the unit
func (u *Unit) SetObjective(pos rl.Vector3) {
	u.Objective = pos
	u.HasObjective = true
	u.Path = nil
	u.PathIndex = 0
}

// ClearObjective stops the unit from moving
func (u *Unit) ClearObjective() {
	u.HasObjective = false
	u.Path = nil
	u.PathIndex = 0
	u.Velocity = rl.Vector3{}
}

// SetPath sets a pathfinding result for the unit to follow
func (u *Unit) SetPath(path []rl.Vector2) {
	u.Path = path
	u.PathIndex = 0
}

// TakeDamage applies damage to the unit
func (u *Unit) TakeDamage(amount float32) {
	// Apply armor reduction
	actualDamage := amount * (1.0 - u.Config.Armor)
	u.Health -= actualDamage
	if u.Health <= 0 {
		u.Health = 0
		u.State = StateDead
	}
}

// Heal restores health to the unit
func (u *Unit) Heal(amount float32) {
	u.Health += amount
	if u.Health > u.MaxHealth {
		u.Health = u.MaxHealth
	}
}

// IsDead returns true if the unit has no health
func (u *Unit) IsDead() bool {
	return u.Health <= 0 || u.State == StateDead
}

// DistanceTo returns the distance to another unit
func (u *Unit) DistanceTo(other *Unit) float32 {
	dx := other.Position.X - u.Position.X
	dz := other.Position.Z - u.Position.Z
	return float32(math.Sqrt(float64(dx*dx + dz*dz)))
}

// DistanceToPoint returns the distance to a point
func (u *Unit) DistanceToPoint(pos rl.Vector3) float32 {
	dx := pos.X - u.Position.X
	dz := pos.Z - u.Position.Z
	return float32(math.Sqrt(float64(dx*dx + dz*dz)))
}

// CanAttack returns true if this unit can attack the target
func (u *Unit) CanAttack(target *Unit) bool {
	if target == nil || target.IsDead() {
		return false
	}
	if u.Team == target.Team {
		return false
	}
	// For now, all units are ground units
	return u.Config.CanAttackGround
}

// IsInRange returns true if the target is within attack range
func (u *Unit) IsInRange(target *Unit) bool {
	return u.DistanceTo(target) <= u.Config.AttackRange
}

// GetForward returns the forward direction vector
func (u *Unit) GetForward() rl.Vector3 {
	return rl.Vector3{
		X: float32(math.Sin(float64(u.Rotation))),
		Y: 0,
		Z: float32(math.Cos(float64(u.Rotation))),
	}
}

// Transport methods

// PickUp marks the unit as being carried
func (u *Unit) PickUp() {
	u.State = StateBeingCarried
	u.Velocity = rl.Vector3{}
	u.ClearObjective()
}

// Drop places the unit at a position with an order
func (u *Unit) Drop(position rl.Vector3, order Order) {
	u.Position = position
	u.State = StateIdle
	u.SetOrder(order, position)
}

// IsCarried returns true if the unit is being carried
func (u *Unit) IsCarried() bool {
	return u.State == StateBeingCarried
}

// SetOrder sets the unit's order with a target position
func (u *Unit) SetOrder(order Order, target rl.Vector3) {
	u.Order = order
	u.OrderTarget = target
	u.HasObjective = true
	u.Objective = target

	if order == OrderPatrolArea {
		u.PatrolCenter = target
		u.PatrolRadius = 5.0 // Default patrol radius
	}
}

// GetOrderName returns the display name for the unit's current order
func (u *Unit) GetOrderName() string {
	names := OrderNames()
	if int(u.Order) < len(names) {
		return names[u.Order]
	}
	return "Unknown"
}

// Helper functions

func lerpAngle(a, b, t float32) float32 {
	a = normalizeAngle(a)
	b = normalizeAngle(b)

	diff := b - a
	if diff > math.Pi {
		diff -= 2 * math.Pi
	}
	if diff < -math.Pi {
		diff += 2 * math.Pi
	}

	result := a + diff*t
	return normalizeAngle(result)
}

func normalizeAngle(a float32) float32 {
	for a > math.Pi {
		a -= 2 * math.Pi
	}
	for a < -math.Pi {
		a += 2 * math.Pi
	}
	return a
}
