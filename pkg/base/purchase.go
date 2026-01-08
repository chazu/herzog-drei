package base

// UnitCost returns the credit cost for a unit type
func UnitCost(unitType UnitType) float32 {
	switch unitType {
	case UnitInfantry:
		return 25
	case UnitTank:
		return 100
	case UnitAA:
		return 75
	case UnitArtillery:
		return 150
	default:
		return 0
	}
}

// UnitName returns a display name for a unit type
func UnitName(unitType UnitType) string {
	switch unitType {
	case UnitInfantry:
		return "Infantry"
	case UnitTank:
		return "Tank"
	case UnitAA:
		return "AA Gun"
	case UnitArtillery:
		return "Artillery"
	default:
		return "Unknown"
	}
}

// PurchaseRequest represents a request to purchase a unit
type PurchaseRequest struct {
	UnitType UnitType
	BaseID   int
	Owner    Owner
}

// TryPurchaseUnit attempts to purchase and queue a unit at a base
// Returns true if successful
func (m *Manager) TryPurchaseUnit(baseID int, unitType UnitType, owner Owner) bool {
	base := m.GetBase(baseID)
	if base == nil {
		return false
	}

	// Verify ownership
	if base.Owner != owner {
		return false
	}

	// Check cost
	cost := UnitCost(unitType)
	if !m.SpendCredits(owner, cost) {
		return false
	}

	// Queue the unit
	base.QueueUnit(unitType)
	return true
}

// GetPurchasableUnits returns units that can be purchased with current credits
func (m *Manager) GetPurchasableUnits(owner Owner) []UnitType {
	credits := m.GetCredits(owner)
	available := make([]UnitType, 0, 4)

	for _, ut := range []UnitType{UnitInfantry, UnitTank, UnitAA, UnitArtillery} {
		if UnitCost(ut) <= credits {
			available = append(available, ut)
		}
	}

	return available
}
