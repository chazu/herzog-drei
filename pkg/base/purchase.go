package base

import (
	"github.com/chazu/herzog-drei/pkg/unit"
)

// AllUnitTypes lists all purchasable unit types
var AllUnitTypes = []unit.UnitType{
	unit.TypeInfantry,
	unit.TypeTank,
	unit.TypeMotorcycle,
	unit.TypeSAM,
	unit.TypeBoat,
	unit.TypeSupply,
}

// UnitCost returns the credit cost for a unit type
func UnitCost(unitType unit.UnitType) float32 {
	return float32(unit.GetConfig(unitType).Cost)
}

// UnitName returns a display name for a unit type
func UnitName(unitType unit.UnitType) string {
	return unit.TypeName(unitType)
}

// PurchaseRequest represents a request to purchase a unit
type PurchaseRequest struct {
	UnitType unit.UnitType
	BaseID   int
	Owner    Owner
}

// TryPurchaseUnit attempts to purchase and queue a unit at a base
// Returns true if successful
func (m *Manager) TryPurchaseUnit(baseID int, unitType unit.UnitType, owner Owner) bool {
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
func (m *Manager) GetPurchasableUnits(owner Owner) []unit.UnitType {
	credits := m.GetCredits(owner)
	available := make([]unit.UnitType, 0, len(AllUnitTypes))

	for _, ut := range AllUnitTypes {
		if UnitCost(ut) <= credits {
			available = append(available, ut)
		}
	}

	return available
}
