package tilemap

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Minimap renders a small overview of the tilemap
type Minimap struct {
	X, Y          int32   // Screen position (top-left)
	Width, Height int32   // Size in pixels
	BorderColor   rl.Color
	BorderWidth   int32
	ShowViewport  bool    // Draw rectangle showing current camera view
	Alpha         uint8   // Transparency (0-255)
}

// NewMinimap creates a new minimap with default settings
func NewMinimap() *Minimap {
	return &Minimap{
		X:            10,
		Y:            10,
		Width:        200,
		Height:       150,
		BorderColor:  rl.DarkGray,
		BorderWidth:  2,
		ShowViewport: true,
		Alpha:        220,
	}
}

// SetPosition sets the screen position of the minimap
func (mm *Minimap) SetPosition(x, y int32) {
	mm.X = x
	mm.Y = y
}

// SetSize sets the size of the minimap in pixels
func (mm *Minimap) SetSize(width, height int32) {
	mm.Width = width
	mm.Height = height
}

// Render draws the minimap
func (mm *Minimap) Render(tm *TileMap, camera *GameCamera) {
	// Draw background
	bgColor := rl.NewColor(20, 20, 20, mm.Alpha)
	rl.DrawRectangle(mm.X-mm.BorderWidth, mm.Y-mm.BorderWidth,
		mm.Width+mm.BorderWidth*2, mm.Height+mm.BorderWidth*2, mm.BorderColor)
	rl.DrawRectangle(mm.X, mm.Y, mm.Width, mm.Height, bgColor)

	// Calculate scale factors
	scaleX := float32(mm.Width) / float32(tm.Width)
	scaleY := float32(mm.Height) / float32(tm.Height)

	// Draw terrain tiles
	for y := 0; y < tm.Height; y++ {
		for x := 0; x < tm.Width; x++ {
			tile := tm.Tiles[y][x]
			info := GetTerrainInfo(tile.Terrain)

			// Apply alpha to terrain color
			color := rl.NewColor(info.Color.R, info.Color.G, info.Color.B, mm.Alpha)

			pixelX := mm.X + int32(float32(x)*scaleX)
			pixelY := mm.Y + int32(float32(y)*scaleY)
			pixelW := int32(scaleX) + 1
			pixelH := int32(scaleY) + 1

			rl.DrawRectangle(pixelX, pixelY, pixelW, pixelH, color)
		}
	}

	// Draw viewport indicator
	if mm.ShowViewport && camera != nil {
		mm.drawViewport(tm, camera, scaleX, scaleY)
	}
}

// drawViewport draws a rectangle showing the current camera view on the minimap
func (mm *Minimap) drawViewport(tm *TileMap, camera *GameCamera, scaleX, scaleY float32) {
	minX, minY, maxX, maxY := camera.GetVisibleTileRange(tm)

	vpX := mm.X + int32(float32(minX)*scaleX)
	vpY := mm.Y + int32(float32(minY)*scaleY)
	vpW := int32(float32(maxX-minX+1) * scaleX)
	vpH := int32(float32(maxY-minY+1) * scaleY)

	// Draw viewport rectangle outline
	viewportColor := rl.NewColor(255, 255, 255, 200)
	rl.DrawRectangleLines(vpX, vpY, vpW, vpH, viewportColor)
}

// RenderWithMarkers draws the minimap with additional markers (units, bases, etc.)
func (mm *Minimap) RenderWithMarkers(tm *TileMap, camera *GameCamera, markers []MinimapMarker) {
	// First render the base minimap
	mm.Render(tm, camera)

	// Calculate scale factors
	scaleX := float32(mm.Width) / float32(tm.Width)
	scaleY := float32(mm.Height) / float32(tm.Height)

	// Draw markers
	for _, marker := range markers {
		tileX, tileY := tm.WorldToTile(marker.WorldX, marker.WorldZ)

		pixelX := mm.X + int32(float32(tileX)*scaleX)
		pixelY := mm.Y + int32(float32(tileY)*scaleY)

		switch marker.Type {
		case MarkerUnit:
			rl.DrawCircle(pixelX, pixelY, 3, marker.Color)
		case MarkerBase:
			rl.DrawRectangle(pixelX-3, pixelY-3, 6, 6, marker.Color)
		case MarkerObjective:
			// Draw a diamond shape
			rl.DrawTriangle(
				rl.NewVector2(float32(pixelX), float32(pixelY-4)),
				rl.NewVector2(float32(pixelX-4), float32(pixelY)),
				rl.NewVector2(float32(pixelX+4), float32(pixelY)),
				marker.Color,
			)
			rl.DrawTriangle(
				rl.NewVector2(float32(pixelX-4), float32(pixelY)),
				rl.NewVector2(float32(pixelX), float32(pixelY+4)),
				rl.NewVector2(float32(pixelX+4), float32(pixelY)),
				marker.Color,
			)
		case MarkerPlayer:
			// Draw player indicator (triangle pointing up)
			rl.DrawTriangle(
				rl.NewVector2(float32(pixelX), float32(pixelY-5)),
				rl.NewVector2(float32(pixelX-4), float32(pixelY+3)),
				rl.NewVector2(float32(pixelX+4), float32(pixelY+3)),
				marker.Color,
			)
		}
	}
}

// MarkerType defines different types of minimap markers
type MarkerType int

const (
	MarkerUnit MarkerType = iota
	MarkerBase
	MarkerObjective
	MarkerPlayer
)

// MinimapMarker represents an icon on the minimap
type MinimapMarker struct {
	WorldX, WorldZ float32
	Type           MarkerType
	Color          rl.Color
}

// NewMarker creates a new minimap marker
func NewMarker(worldX, worldZ float32, markerType MarkerType, color rl.Color) MinimapMarker {
	return MinimapMarker{
		WorldX: worldX,
		WorldZ: worldZ,
		Type:   markerType,
		Color:  color,
	}
}
