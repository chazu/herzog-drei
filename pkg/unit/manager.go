package unit

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Manager handles all units in the game
type Manager struct {
	Units      []*Unit
	nextID     int
	pickupDist float32 // Distance for mech to pick up units
}

// NewManager creates a new unit manager
func NewManager() *Manager {
	return &Manager{
		Units:      make([]*Unit, 0, 64),
		nextID:     1,
		pickupDist: 1.5, // Mech must be within 1.5 units to pick up
	}
}

// Spawn creates a new unit at the given position
func (m *Manager) Spawn(unitType Type, position rl.Vector3, owner Owner) *Unit {
	unit := New(m.nextID, unitType, position, owner)
	m.nextID++
	m.Units = append(m.Units, unit)
	return unit
}

// Update updates all units
func (m *Manager) Update(dt float32) {
	for _, u := range m.Units {
		u.Update(dt)
	}

	// Remove dead units
	m.cleanup()
}

func (m *Manager) cleanup() {
	alive := m.Units[:0]
	for _, u := range m.Units {
		if !u.IsDead() {
			alive = append(alive, u)
		}
	}
	m.Units = alive
}

// GetNearestPickupTarget finds the nearest friendly unit that can be picked up
func (m *Manager) GetNearestPickupTarget(pos rl.Vector3, owner Owner) *Unit {
	var nearest *Unit
	minDist := float32(math.MaxFloat32)

	for _, u := range m.Units {
		// Skip if not owned by picker, already carried, or dead
		if u.Owner != owner || u.IsCarried() || u.IsDead() {
			continue
		}

		dist := distance(pos, u.Position)
		if dist < minDist && dist < m.pickupDist {
			minDist = dist
			nearest = u
		}
	}

	return nearest
}

// GetUnitsInRange returns all units within range of a position
func (m *Manager) GetUnitsInRange(pos rl.Vector3, radius float32, owner Owner, enemiesOnly bool) []*Unit {
	result := make([]*Unit, 0)

	for _, u := range m.Units {
		if u.IsCarried() || u.IsDead() {
			continue
		}

		if enemiesOnly && u.Owner == owner {
			continue
		}

		dist := distance(pos, u.Position)
		if dist <= radius {
			result = append(result, u)
		}
	}

	return result
}

// GetNearestEnemy returns the nearest enemy unit to a position
func (m *Manager) GetNearestEnemy(pos rl.Vector3, owner Owner) *Unit {
	var nearest *Unit
	minDist := float32(math.MaxFloat32)

	for _, u := range m.Units {
		if u.Owner == owner || u.Owner == OwnerNeutral || u.IsCarried() || u.IsDead() {
			continue
		}

		dist := distance(pos, u.Position)
		if dist < minDist {
			minDist = dist
			nearest = u
		}
	}

	return nearest
}

// UpdateTargets assigns attack targets based on orders and proximity
func (m *Manager) UpdateTargets() {
	for _, u := range m.Units {
		if u.IsCarried() || u.IsDead() {
			continue
		}

		switch u.Order {
		case OrderAttackNearest, OrderDefendPosition, OrderPatrolArea:
			// Find nearest enemy in attack range
			enemy := m.GetNearestEnemy(u.Position, u.Owner)
			if enemy != nil {
				dist := distance(u.Position, enemy.Position)
				if dist <= u.AttackRange*1.5 {
					u.Target = enemy
					u.OrderTarget = enemy.Position
				} else {
					u.Target = nil
				}
			} else {
				u.Target = nil
			}
		}
	}
}

// Count returns the number of units for an owner
func (m *Manager) Count(owner Owner) int {
	count := 0
	for _, u := range m.Units {
		if u.Owner == owner && !u.IsCarried() && !u.IsDead() {
			count++
		}
	}
	return count
}

// CountByType returns the number of units of a type for an owner
func (m *Manager) CountByType(owner Owner, unitType Type) int {
	count := 0
	for _, u := range m.Units {
		if u.Owner == owner && u.Type == unitType && !u.IsCarried() && !u.IsDead() {
			count++
		}
	}
	return count
}
