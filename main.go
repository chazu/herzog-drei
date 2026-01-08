package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/chazu/herzog-drei/pkg/base"
	"github.com/chazu/herzog-drei/pkg/combat"
	"github.com/chazu/herzog-drei/pkg/mech"
	"github.com/chazu/herzog-drei/pkg/tilemap"
	"github.com/chazu/herzog-drei/pkg/unit"
)

const (
	screenWidth  = 1280
	screenHeight = 720
	targetFPS    = 60
	gameTitle    = "Herzog Drei"

	mapWidth  = 64
	mapHeight = 48
)

// Game holds the game state
type Game struct {
	// Map and camera
	tileMap *tilemap.TileMap
	camera  *tilemap.GameCamera
	minimap *tilemap.Minimap

	// Player mech
	playerMech   *mech.Mech
	mechInput    *mech.InputHandler
	mechRenderer *mech.Renderer

	// Units
	unitManager    *unit.Manager
	unitRenderer   *unit.Renderer
	unitPathfinder *unit.Pathfinder

	// Bases
	baseManager  *base.Manager
	baseRenderer *base.Renderer

	// Combat
	combatSystem   *combat.System
	combatRenderer *combat.Renderer
}

// NewGame creates and initializes a new game instance
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

// init sets up initial game state
func (g *Game) init() {
	// Create tile map with test terrain
	g.tileMap = tilemap.GenerateTestMap(mapWidth, mapHeight)

	// Set up game camera
	g.camera = tilemap.NewGameCamera()
	g.camera.SetBounds(g.tileMap.GetWorldBounds())

	// Create player mech at center of map
	centerX, centerZ := g.tileMap.TileToWorld(mapWidth/2, mapHeight/2)
	startPos := rl.NewVector3(centerX, 3, centerZ)
	g.playerMech = mech.New(startPos, mech.DefaultConfig())
	g.mechInput = mech.NewInputHandler()
	g.mechRenderer = mech.NewRenderer()

	// Set camera to follow mech
	g.camera.SetTarget(g.playerMech.Position)

	// Set up minimap in top-right corner
	g.minimap = tilemap.NewMinimap()
	g.minimap.SetPosition(screenWidth-210, 10)
	g.minimap.SetSize(200, 150)

	// Initialize unit system
	g.unitManager = unit.NewManager(100) // Max 100 units
	g.unitRenderer = unit.NewRenderer()
	g.unitPathfinder = unit.NewPathfinder(mapWidth, mapHeight, 1.0)
	g.unitManager.Pathfinder = g.unitPathfinder

	// Initialize base system
	g.baseManager = base.NewManager(base.DefaultConfig())
	g.baseRenderer = base.NewRenderer()
	g.baseManager.CreateDefaultMap()

	// Initialize combat system
	g.combatSystem = combat.NewSystem(combat.DefaultConfig())
	g.combatRenderer = combat.NewRenderer()
	g.combatSystem.SetRespawnPosition(startPos) // Respawn at start position

	// Spawn test units for demonstration
	g.spawnTestUnits()
}

// Update handles game logic each frame
func (g *Game) Update() {
	dt := rl.GetFrameTime()

	// Handle camera input (zoom)
	g.camera.HandleInput()

	// Process player input
	g.mechInput.Update(g.playerMech)

	// Update mech
	g.playerMech.Update(dt)

	// Check terrain collision for ground (robot) mode
	if g.playerMech.Mode == mech.ModeRobot {
		if !g.tileMap.IsPassableAt(g.playerMech.Position.X, g.playerMech.Position.Z) {
			// Push mech back if on impassable terrain
			g.playerMech.Position.X -= g.playerMech.Velocity.X * dt
			g.playerMech.Position.Z -= g.playerMech.Velocity.Z * dt
		}
		// Adjust height based on terrain
		g.playerMech.Position.Y = g.tileMap.GetHeightAt(g.playerMech.Position.X, g.playerMech.Position.Z)
	}

	// Handle transport (pickup/drop units)
	g.handleTransport()

	// Update units
	g.unitManager.Update(dt)

	// Update bases (income, capture progress, spawns)
	g.baseManager.Update(dt)

	// Update combat (hit detection, damage, respawn)
	g.combatSystem.Update(dt, g.playerMech, g.unitManager)

	// Process base spawn queues - spawn units from bases
	g.processBaseSpawns()

	// Handle unit purchasing (press 1-6 to buy units at nearest owned base)
	g.handleUnitPurchaseInput()

	// Update camera to follow mech
	g.camera.SetTarget(g.playerMech.Position)
	g.camera.Update()
}

// handleTransport handles picking up and dropping units
func (g *Game) handleTransport() {
	// Handle pickup
	if g.playerMech.InputPickup && g.playerMech.CanPickup() {
		pickupRadius := float32(2.0)
		nearUnit := g.unitManager.GetNearestPickupableUnit(
			g.playerMech.Position,
			pickupRadius,
			g.playerMech.Team,
		)
		if nearUnit != nil {
			g.playerMech.PickupUnit(nearUnit)
		}
	}

	// Handle drop
	if g.playerMech.InputDrop && g.playerMech.CanDrop() {
		g.playerMech.DropUnit()
	}
}

// Render draws the game each frame
func (g *Game) Render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.SkyBlue)

	// 3D rendering
	g.camera.Begin3D()

	// Render tile map
	g.tileMap.Render()

	// Draw bases
	g.baseRenderer.Draw(g.baseManager)

	// Draw units
	g.unitRenderer.Draw(g.unitManager)

	// Draw player mech
	g.mechRenderer.Draw(g.playerMech)

	// Draw combat effects (explosions)
	g.combatRenderer.Draw(g.combatSystem)

	g.camera.End3D()

	// Draw minimap with player marker
	markers := []tilemap.MinimapMarker{
		tilemap.NewMarker(g.playerMech.Position.X, g.playerMech.Position.Z, tilemap.MarkerPlayer, rl.Red),
	}
	g.minimap.RenderWithMarkers(g.tileMap, g.camera, markers)

	// Draw UI overlay
	rl.DrawText(gameTitle, 10, 10, 20, rl.DarkGray)
	rl.DrawFPS(screenWidth-100, 10)

	// Draw mech UI (health bar, mode indicator)
	g.mechRenderer.DrawUI(g.playerMech, screenWidth, screenHeight)

	// Draw unit UI
	g.unitRenderer.DrawUI(g.unitManager, screenWidth, screenHeight)

	// Draw base UI (credits, base counts)
	g.baseRenderer.DrawUI(g.baseManager, screenWidth, screenHeight)

	// Draw combat UI (respawn timer, invulnerability)
	g.combatRenderer.DrawUI(g.combatSystem, screenWidth, screenHeight)

	// Show current terrain info
	terrain := g.tileMap.GetTerrainAt(g.playerMech.Position.X, g.playerMech.Position.Z)
	info := tilemap.GetTerrainInfo(terrain)
	rl.DrawText("Terrain: "+info.Name, 10, screenHeight-100, 15, rl.DarkGray)

	// Show transport info
	if g.playerMech.IsCarrying() {
		carriedInfo := "Carrying: " + g.playerMech.CarriedUnit.Config.Type.String()
		rl.DrawText(carriedInfo, 10, screenHeight-80, 15, rl.Green)
	}
	orderInfo := "Order: " + g.playerMech.GetSelectedOrderName() + " (R/F to cycle)"
	rl.DrawText(orderInfo, 10, screenHeight-60, 15, rl.DarkGray)

	rl.DrawText("T: Transform | E: Pickup | Q: Drop | R/F: Cycle Order | Scroll: Zoom", 10, screenHeight-40, 12, rl.DarkGray)
	rl.DrawText("1-6: Spawn units | 1:Infantry 2:Tank 3:Bike 4:SAM 5:Boat 6:Supply", 10, screenHeight-20, 12, rl.DarkGray)

	rl.EndDrawing()
}

// spawnTestUnits creates initial units for testing
func (g *Game) spawnTestUnits() {
	centerX, centerZ := g.tileMap.TileToWorld(mapWidth/2, mapHeight/2)

	// Spawn player units on left side
	g.unitManager.SpawnWithObjective(
		unit.TypeInfantry, unit.TeamPlayer,
		rl.NewVector3(centerX-10, 0, centerZ+5),
		rl.NewVector3(centerX+10, 0, centerZ+5),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeTank, unit.TeamPlayer,
		rl.NewVector3(centerX-10, 0, centerZ),
		rl.NewVector3(centerX+10, 0, centerZ),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeMotorcycle, unit.TeamPlayer,
		rl.NewVector3(centerX-10, 0, centerZ-5),
		rl.NewVector3(centerX+10, 0, centerZ-5),
	)

	// Spawn enemy units on right side
	g.unitManager.SpawnWithObjective(
		unit.TypeInfantry, unit.TeamEnemy,
		rl.NewVector3(centerX+10, 0, centerZ+5),
		rl.NewVector3(centerX-10, 0, centerZ+5),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeTank, unit.TeamEnemy,
		rl.NewVector3(centerX+10, 0, centerZ),
		rl.NewVector3(centerX-10, 0, centerZ),
	)
	g.unitManager.SpawnWithObjective(
		unit.TypeSAM, unit.TeamEnemy,
		rl.NewVector3(centerX+10, 0, centerZ-5),
		rl.NewVector3(centerX-10, 0, centerZ-5),
	)
}

// processBaseSpawns handles spawning units from base queues
func (g *Game) processBaseSpawns() {
	for _, b := range g.baseManager.Bases {
		unitType, spawned := b.TrySpawn(g.baseManager.Config)
		if !spawned {
			continue
		}

		// Map base owner to unit team
		var team unit.Team
		switch b.Owner {
		case base.OwnerPlayer1:
			team = unit.TeamPlayer
		case base.OwnerPlayer2:
			team = unit.TeamEnemy
		default:
			continue // Neutral bases shouldn't spawn
		}

		// Spawn the unit at the base's spawn point
		g.unitManager.Spawn(unitType, team, b.SpawnPoint)
	}
}

// handleUnitPurchaseInput purchases units based on number key presses
func (g *Game) handleUnitPurchaseInput() {
	// Find nearest owned base to purchase from
	nearestBase := g.findNearestOwnedBase(base.OwnerPlayer1)
	if nearestBase == nil {
		return // No owned bases to purchase from
	}

	// Map keys 1-6 to unit types
	type keyMapping struct {
		key      int32
		unitType unit.UnitType
	}
	mappings := []keyMapping{
		{rl.KeyOne, unit.TypeInfantry},
		{rl.KeyTwo, unit.TypeTank},
		{rl.KeyThree, unit.TypeMotorcycle},
		{rl.KeyFour, unit.TypeSAM},
		{rl.KeyFive, unit.TypeBoat},
		{rl.KeySix, unit.TypeSupply},
	}

	for _, m := range mappings {
		if rl.IsKeyPressed(m.key) {
			// Try to purchase - this checks credits and queues at the base
			g.baseManager.TryPurchaseUnit(nearestBase.ID, m.unitType, base.OwnerPlayer1)
		}
	}
}

// findNearestOwnedBase finds the player's nearest owned base
func (g *Game) findNearestOwnedBase(owner base.Owner) *base.Base {
	ownedBases := g.baseManager.GetBasesOwnedBy(owner)
	if len(ownedBases) == 0 {
		return nil
	}

	var nearest *base.Base
	nearestDist := float32(1e9)

	for _, b := range ownedBases {
		dx := b.Position.X - g.playerMech.Position.X
		dz := b.Position.Z - g.playerMech.Position.Z
		dist := dx*dx + dz*dz // squared distance is fine for comparison
		if dist < nearestDist {
			nearestDist = dist
			nearest = b
		}
	}

	return nearest
}

func main() {
	// Initialize window
	rl.InitWindow(screenWidth, screenHeight, gameTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(targetFPS)

	// Create game instance
	game := NewGame()

	// Main game loop
	for !rl.WindowShouldClose() {
		game.Update()
		game.Render()
	}
}
