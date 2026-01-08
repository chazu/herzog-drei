package unit

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Manager handles unit spawning, updates, and cleanup
type Manager struct {
	units    []*Unit
	nextID   uint32
	maxUnits int

	// Pathfinder reference (set externally)
	Pathfinder *Pathfinder
}

// NewManager creates a new unit manager
func NewManager(maxUnits int) *Manager {
	return &Manager{
		units:    make([]*Unit, 0, maxUnits),
		nextID:   1,
		maxUnits: maxUnits,
	}
}

// Spawn creates a new unit at the given position
func (m *Manager) Spawn(unitType UnitType, team Team, pos rl.Vector3) *Unit {
	if len(m.units) >= m.maxUnits {
		return nil
	}

	u := New(m.nextID, unitType, team, pos)
	m.nextID++
	m.units = append(m.units, u)
	return u
}

// SpawnWithObjective creates a new unit and sets an objective
func (m *Manager) SpawnWithObjective(unitType UnitType, team Team, pos, objective rl.Vector3) *Unit {
	u := m.Spawn(unitType, team, pos)
	if u != nil {
		u.SetObjective(objective)
		// Try to pathfind if available
		if m.Pathfinder != nil {
			path := m.Pathfinder.FindPath(
				rl.Vector2{X: pos.X, Y: pos.Z},
				rl.Vector2{X: objective.X, Y: objective.Z},
			)
			if path != nil {
				u.SetPath(path)
			}
		}
	}
	return u
}

// Update updates all units
func (m *Manager) Update(dt float32) {
	for _, u := range m.units {
		u.Update(dt)
	}

	// Run AI for all units
	m.updateAI(dt)

	// Run combat for all units
	m.updateCombat(dt)

	// Cleanup dead units
	m.cleanup()
}

// updateAI handles basic AI behaviors for all units
func (m *Manager) updateAI(dt float32) {
	for _, u := range m.units {
		if u.IsDead() {
			continue
		}

		// Skip if unit already has a target and is attacking
		if u.Target != nil && !u.Target.IsDead() && u.IsInRange(u.Target) {
			continue
		}

		// Find nearest enemy to attack
		var nearest *Unit
		nearestDist := float32(1000000)

		for _, other := range m.units {
			if other == u || other.IsDead() {
				continue
			}
			if other.Team == u.Team {
				continue
			}
			if !u.CanAttack(other) {
				continue
			}

			dist := u.DistanceTo(other)
			if dist < nearestDist {
				nearest = other
				nearestDist = dist
			}
		}

		// Set target if enemy found within aggro range
		aggroRange := u.Config.AttackRange * 2
		if nearest != nil && nearestDist <= aggroRange {
			u.Target = nearest
		} else {
			u.Target = nil
		}
	}
}

// updateCombat handles unit attacking
func (m *Manager) updateCombat(dt float32) {
	for _, u := range m.units {
		if u.IsDead() || u.Target == nil {
			continue
		}

		// Check if target is still valid
		if u.Target.IsDead() {
			u.Target = nil
			u.State = StateIdle
			continue
		}

		// Check if in range
		if !u.IsInRange(u.Target) {
			// Move toward target
			if !u.HasObjective {
				u.SetObjective(u.Target.Position)
			}
			continue
		}

		// Attack if cooldown ready
		if u.AttackCooldown <= 0 {
			u.State = StateAttacking
			u.Target.TakeDamage(u.Config.AttackDamage)
			u.AttackCooldown = 1.0 / u.Config.AttackRate
		}
	}
}

// cleanup removes dead units from the manager
func (m *Manager) cleanup() {
	alive := m.units[:0]
	for _, u := range m.units {
		// Keep unit for a short time after death for death animation
		if !u.IsDead() {
			alive = append(alive, u)
		}
	}
	m.units = alive
}

// GetUnits returns all units (including dead ones pending cleanup)
func (m *Manager) GetUnits() []*Unit {
	return m.units
}

// GetAliveUnits returns only living units
func (m *Manager) GetAliveUnits() []*Unit {
	result := make([]*Unit, 0, len(m.units))
	for _, u := range m.units {
		if !u.IsDead() {
			result = append(result, u)
		}
	}
	return result
}

// GetUnitsByTeam returns units belonging to a specific team
func (m *Manager) GetUnitsByTeam(team Team) []*Unit {
	result := make([]*Unit, 0)
	for _, u := range m.units {
		if u.Team == team && !u.IsDead() {
			result = append(result, u)
		}
	}
	return result
}

// GetUnitByID returns a unit by its ID
func (m *Manager) GetUnitByID(id uint32) *Unit {
	for _, u := range m.units {
		if u.ID == id {
			return u
		}
	}
	return nil
}

// GetUnitsInRadius returns all units within a radius of a point
func (m *Manager) GetUnitsInRadius(center rl.Vector3, radius float32) []*Unit {
	result := make([]*Unit, 0)
	for _, u := range m.units {
		if u.IsDead() {
			continue
		}
		if u.DistanceToPoint(center) <= radius {
			result = append(result, u)
		}
	}
	return result
}

// GetEnemiesInRadius returns enemy units within a radius
func (m *Manager) GetEnemiesInRadius(center rl.Vector3, radius float32, myTeam Team) []*Unit {
	result := make([]*Unit, 0)
	for _, u := range m.units {
		if u.IsDead() || u.Team == myTeam {
			continue
		}
		if u.DistanceToPoint(center) <= radius {
			result = append(result, u)
		}
	}
	return result
}

// Count returns the total number of units
func (m *Manager) Count() int {
	return len(m.units)
}

// CountByTeam returns the number of units on a team
func (m *Manager) CountByTeam(team Team) int {
	count := 0
	for _, u := range m.units {
		if u.Team == team && !u.IsDead() {
			count++
		}
	}
	return count
}

// Clear removes all units
func (m *Manager) Clear() {
	m.units = m.units[:0]
}

// SetPathfinderForUnit calculates and sets a path for a specific unit
func (m *Manager) SetPathfinderForUnit(u *Unit, goal rl.Vector3) {
	if m.Pathfinder == nil {
		return
	}

	path := m.Pathfinder.FindPath(
		rl.Vector2{X: u.Position.X, Y: u.Position.Z},
		rl.Vector2{X: goal.X, Y: goal.Z},
	)
	if path != nil {
		u.SetPath(path)
	}
}
