package base

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// PlayerState tracks economy and game state for a player
type PlayerState struct {
	Credits float32
}

// Manager manages all bases in the game
type Manager struct {
	Config Config
	Bases  []*Base
	nextID int

	// Player economies
	Player1 PlayerState
	Player2 PlayerState
}

// NewManager creates a new base manager
func NewManager(cfg Config) *Manager {
	return &Manager{
		Config: cfg,
		Bases:  make([]*Base, 0, 16),
		nextID: 1,
		Player1: PlayerState{Credits: 500}, // Starting credits
		Player2: PlayerState{Credits: 500},
	}
}

// AddBase creates and adds a new base
func (m *Manager) AddBase(baseType Type, position rl.Vector3, owner Owner) *Base {
	base := NewBase(m.nextID, baseType, position, owner, m.Config)
	m.nextID++
	m.Bases = append(m.Bases, base)
	return base
}

// Update updates all bases and collects income
func (m *Manager) Update(dt float32) {
	for _, base := range m.Bases {
		base.Update(dt, m.Config)

		// Collect income for owners
		income := base.CollectIncome()
		switch base.Owner {
		case OwnerPlayer1:
			m.Player1.Credits += income
		case OwnerPlayer2:
			m.Player2.Credits += income
		}
	}
}

// GetBase returns a base by ID
func (m *Manager) GetBase(id int) *Base {
	for _, base := range m.Bases {
		if base.ID == id {
			return base
		}
	}
	return nil
}

// GetBaseAt returns the base at or near the given position
func (m *Manager) GetBaseAt(pos rl.Vector3, radius float32) *Base {
	for _, base := range m.Bases {
		dx := base.Position.X - pos.X
		dz := base.Position.Z - pos.Z
		distSq := dx*dx + dz*dz
		if distSq <= radius*radius {
			return base
		}
	}
	return nil
}

// GetBasesOwnedBy returns all bases owned by a specific owner
func (m *Manager) GetBasesOwnedBy(owner Owner) []*Base {
	result := make([]*Base, 0)
	for _, base := range m.Bases {
		if base.Owner == owner {
			result = append(result, base)
		}
	}
	return result
}

// GetHQ returns the HQ for a specific owner (nil if destroyed or not found)
func (m *Manager) GetHQ(owner Owner) *Base {
	for _, base := range m.Bases {
		if base.Type == TypeHQ && base.Owner == owner && !base.IsDestroyed() {
			return base
		}
	}
	return nil
}

// IsGameOver checks if either player has lost their HQ
// Returns the losing owner, or OwnerNeutral if game continues
func (m *Manager) IsGameOver() Owner {
	player1HQ := m.GetHQ(OwnerPlayer1)
	player2HQ := m.GetHQ(OwnerPlayer2)

	if player1HQ == nil {
		return OwnerPlayer1 // Player 1 lost
	}
	if player2HQ == nil {
		return OwnerPlayer2 // Player 2 lost
	}
	return OwnerNeutral // Game continues
}

// SpendCredits attempts to spend credits for a player
// Returns true if successful, false if insufficient funds
func (m *Manager) SpendCredits(owner Owner, amount float32) bool {
	var player *PlayerState
	switch owner {
	case OwnerPlayer1:
		player = &m.Player1
	case OwnerPlayer2:
		player = &m.Player2
	default:
		return false
	}

	if player.Credits >= amount {
		player.Credits -= amount
		return true
	}
	return false
}

// GetCredits returns credits for a player
func (m *Manager) GetCredits(owner Owner) float32 {
	switch owner {
	case OwnerPlayer1:
		return m.Player1.Credits
	case OwnerPlayer2:
		return m.Player2.Credits
	default:
		return 0
	}
}

// CreateDefaultMap creates a standard symmetric map layout
func (m *Manager) CreateDefaultMap() {
	// Player 1 HQ (bottom of map)
	m.AddBase(TypeHQ, rl.NewVector3(0, 0, -15), OwnerPlayer1)

	// Player 2 HQ (top of map)
	m.AddBase(TypeHQ, rl.NewVector3(0, 0, 15), OwnerPlayer2)

	// Neutral outposts in a symmetric pattern
	// Center outpost
	m.AddBase(TypeOutpost, rl.NewVector3(0, 0, 0), OwnerNeutral)

	// Side outposts
	m.AddBase(TypeOutpost, rl.NewVector3(-10, 0, -5), OwnerNeutral)
	m.AddBase(TypeOutpost, rl.NewVector3(10, 0, -5), OwnerNeutral)
	m.AddBase(TypeOutpost, rl.NewVector3(-10, 0, 5), OwnerNeutral)
	m.AddBase(TypeOutpost, rl.NewVector3(10, 0, 5), OwnerNeutral)

	// Corner outposts
	m.AddBase(TypeOutpost, rl.NewVector3(-8, 0, -10), OwnerPlayer1) // Near P1
	m.AddBase(TypeOutpost, rl.NewVector3(8, 0, -10), OwnerPlayer1)
	m.AddBase(TypeOutpost, rl.NewVector3(-8, 0, 10), OwnerPlayer2) // Near P2
	m.AddBase(TypeOutpost, rl.NewVector3(8, 0, 10), OwnerPlayer2)
}
