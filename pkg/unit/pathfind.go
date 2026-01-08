package unit

import (
	"container/heap"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Pathfinder implements A* pathfinding on a grid
type Pathfinder struct {
	width, height int
	cellSize      float32
	blocked       []bool // true if cell is blocked
}

// NewPathfinder creates a new pathfinder for the given map size
func NewPathfinder(width, height int, cellSize float32) *Pathfinder {
	return &Pathfinder{
		width:    width,
		height:   height,
		cellSize: cellSize,
		blocked:  make([]bool, width*height),
	}
}

// SetBlocked marks a cell as blocked or unblocked
func (p *Pathfinder) SetBlocked(x, y int, blocked bool) {
	if x >= 0 && x < p.width && y >= 0 && y < p.height {
		p.blocked[y*p.width+x] = blocked
	}
}

// IsBlocked returns true if a cell is blocked
func (p *Pathfinder) IsBlocked(x, y int) bool {
	if x < 0 || x >= p.width || y < 0 || y >= p.height {
		return true // Out of bounds is blocked
	}
	return p.blocked[y*p.width+x]
}

// WorldToGrid converts world coordinates to grid coordinates
func (p *Pathfinder) WorldToGrid(pos rl.Vector2) (int, int) {
	// Center the grid on the world origin
	offsetX := float32(p.width) * p.cellSize / 2
	offsetY := float32(p.height) * p.cellSize / 2

	x := int((pos.X + offsetX) / p.cellSize)
	y := int((pos.Y + offsetY) / p.cellSize)
	return x, y
}

// GridToWorld converts grid coordinates to world coordinates (center of cell)
func (p *Pathfinder) GridToWorld(x, y int) rl.Vector2 {
	offsetX := float32(p.width) * p.cellSize / 2
	offsetY := float32(p.height) * p.cellSize / 2

	return rl.Vector2{
		X: float32(x)*p.cellSize + p.cellSize/2 - offsetX,
		Y: float32(y)*p.cellSize + p.cellSize/2 - offsetY,
	}
}

// FindPath finds a path from start to goal using A*
// Returns nil if no path is found
func (p *Pathfinder) FindPath(start, goal rl.Vector2) []rl.Vector2 {
	startX, startY := p.WorldToGrid(start)
	goalX, goalY := p.WorldToGrid(goal)

	// If start or goal is blocked, return nil
	if p.IsBlocked(startX, startY) || p.IsBlocked(goalX, goalY) {
		return nil
	}

	// If start == goal, return single point
	if startX == goalX && startY == goalY {
		return []rl.Vector2{p.GridToWorld(goalX, goalY)}
	}

	// A* implementation
	openSet := &nodeHeap{}
	heap.Init(openSet)

	startNode := &pathNode{
		x: startX, y: startY,
		g: 0,
		h: heuristic(startX, startY, goalX, goalY),
	}
	startNode.f = startNode.g + startNode.h
	heap.Push(openSet, startNode)

	cameFrom := make(map[int]*pathNode)
	gScore := make(map[int]float32)
	gScore[startY*p.width+startX] = 0

	// Direction vectors for 8-directional movement
	dirs := [][2]int{
		{0, -1}, {0, 1}, {-1, 0}, {1, 0},   // Cardinal
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1}, // Diagonal
	}
	costs := []float32{1, 1, 1, 1, 1.41, 1.41, 1.41, 1.41}

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*pathNode)

		// Check if we reached the goal
		if current.x == goalX && current.y == goalY {
			return p.reconstructPath(cameFrom, current)
		}

		// Explore neighbors
		for i, dir := range dirs {
			nx, ny := current.x+dir[0], current.y+dir[1]

			// Skip if blocked or out of bounds
			if p.IsBlocked(nx, ny) {
				continue
			}

			// For diagonal movement, check if both adjacent cells are free
			if i >= 4 { // Diagonal
				if p.IsBlocked(current.x+dir[0], current.y) || p.IsBlocked(current.x, current.y+dir[1]) {
					continue
				}
			}

			tentativeG := gScore[current.y*p.width+current.x] + costs[i]
			neighborKey := ny*p.width + nx

			existingG, exists := gScore[neighborKey]
			if !exists || tentativeG < existingG {
				neighbor := &pathNode{
					x: nx, y: ny,
					g: tentativeG,
					h: heuristic(nx, ny, goalX, goalY),
				}
				neighbor.f = neighbor.g + neighbor.h

				cameFrom[neighborKey] = current
				gScore[neighborKey] = tentativeG
				heap.Push(openSet, neighbor)
			}
		}
	}

	// No path found
	return nil
}

// reconstructPath builds the path from the cameFrom map
func (p *Pathfinder) reconstructPath(cameFrom map[int]*pathNode, current *pathNode) []rl.Vector2 {
	path := []rl.Vector2{p.GridToWorld(current.x, current.y)}

	key := current.y*p.width + current.x
	for {
		prev, exists := cameFrom[key]
		if !exists {
			break
		}
		path = append([]rl.Vector2{p.GridToWorld(prev.x, prev.y)}, path...)
		key = prev.y*p.width + prev.x
	}

	// Simplify path by removing redundant waypoints
	return p.simplifyPath(path)
}

// simplifyPath removes unnecessary waypoints that are in a straight line
func (p *Pathfinder) simplifyPath(path []rl.Vector2) []rl.Vector2 {
	if len(path) <= 2 {
		return path
	}

	result := []rl.Vector2{path[0]}

	for i := 1; i < len(path)-1; i++ {
		prev := result[len(result)-1]
		curr := path[i]
		next := path[i+1]

		// Check if direction changes
		dx1 := curr.X - prev.X
		dy1 := curr.Y - prev.Y
		dx2 := next.X - curr.X
		dy2 := next.Y - curr.Y

		// Normalize directions
		len1 := float32(math.Sqrt(float64(dx1*dx1 + dy1*dy1)))
		len2 := float32(math.Sqrt(float64(dx2*dx2 + dy2*dy2)))

		if len1 > 0 && len2 > 0 {
			dx1 /= len1
			dy1 /= len1
			dx2 /= len2
			dy2 /= len2

			// If direction changed significantly, keep this waypoint
			dot := dx1*dx2 + dy1*dy2
			if dot < 0.99 { // Not parallel
				result = append(result, curr)
			}
		}
	}

	result = append(result, path[len(path)-1])
	return result
}

// heuristic calculates the estimated cost from (x,y) to (gx,gy)
// Using octile distance for 8-directional movement
func heuristic(x, y, gx, gy int) float32 {
	dx := abs(gx - x)
	dy := abs(gy - y)
	// Octile distance
	return float32(max(dx, dy)) + 0.41*float32(min(dx, dy))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// pathNode represents a node in the A* search
type pathNode struct {
	x, y  int
	f, g, h float32
	index int // heap index
}

// nodeHeap implements heap.Interface for A* priority queue
type nodeHeap []*pathNode

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].f < h[j].f }
func (h nodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *nodeHeap) Push(x interface{}) {
	n := len(*h)
	node := x.(*pathNode)
	node.index = n
	*h = append(*h, node)
}

func (h *nodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*h = old[0 : n-1]
	return node
}
