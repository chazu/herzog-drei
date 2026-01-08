package unit

// GetConfig returns the configuration for a unit type
func GetConfig(t UnitType) Config {
	switch t {
	case TypeInfantry:
		return Config{
			Type:            TypeInfantry,
			Speed:           2.0,
			TurnSpeed:       4.0,
			CanTraverseWater: false,
			AttackRange:     3.0,
			AttackDamage:    5.0,
			AttackRate:      1.5,
			CanAttackAir:    false,
			CanAttackGround: true,
			MaxHealth:       30.0,
			Armor:           0.0,
			CanCapture:      true,
			Cost:            100,
		}

	case TypeTank:
		return Config{
			Type:            TypeTank,
			Speed:           3.0,
			TurnSpeed:       2.0,
			CanTraverseWater: false,
			AttackRange:     6.0,
			AttackDamage:    20.0,
			AttackRate:      0.8,
			CanAttackAir:    false,
			CanAttackGround: true,
			MaxHealth:       100.0,
			Armor:           0.3,
			CanCapture:      false,
			Cost:            400,
		}

	case TypeMotorcycle:
		return Config{
			Type:            TypeMotorcycle,
			Speed:           6.0,
			TurnSpeed:       5.0,
			CanTraverseWater: false,
			AttackRange:     4.0,
			AttackDamage:    8.0,
			AttackRate:      2.0,
			CanAttackAir:    false,
			CanAttackGround: true,
			MaxHealth:       40.0,
			Armor:           0.0,
			CanCapture:      false,
			Cost:            200,
		}

	case TypeSAM:
		return Config{
			Type:            TypeSAM,
			Speed:           2.5,
			TurnSpeed:       3.0,
			CanTraverseWater: false,
			AttackRange:     8.0,
			AttackDamage:    25.0,
			AttackRate:      1.0,
			CanAttackAir:    true,
			CanAttackGround: false,
			MaxHealth:       50.0,
			Armor:           0.1,
			CanCapture:      false,
			Cost:            350,
		}

	case TypeBoat:
		return Config{
			Type:            TypeBoat,
			Speed:           4.0,
			TurnSpeed:       2.5,
			CanTraverseWater: true,
			AttackRange:     5.0,
			AttackDamage:    15.0,
			AttackRate:      1.2,
			CanAttackAir:    false,
			CanAttackGround: true,
			MaxHealth:       60.0,
			Armor:           0.2,
			CanCapture:      false,
			Cost:            300,
		}

	case TypeSupply:
		return Config{
			Type:            TypeSupply,
			Speed:           3.5,
			TurnSpeed:       2.0,
			CanTraverseWater: false,
			AttackRange:     0.0,
			AttackDamage:    0.0,
			AttackRate:      0.0,
			CanAttackAir:    false,
			CanAttackGround: false,
			MaxHealth:       80.0,
			Armor:           0.1,
			CanCapture:      false,
			Cost:            250,
		}

	default:
		// Default to infantry if unknown type
		return GetConfig(TypeInfantry)
	}
}

// TypeName returns the display name for a unit type
func TypeName(t UnitType) string {
	switch t {
	case TypeInfantry:
		return "Infantry"
	case TypeTank:
		return "Tank"
	case TypeMotorcycle:
		return "Motorcycle"
	case TypeSAM:
		return "SAM Launcher"
	case TypeBoat:
		return "Boat"
	case TypeSupply:
		return "Supply Truck"
	default:
		return "Unknown"
	}
}
