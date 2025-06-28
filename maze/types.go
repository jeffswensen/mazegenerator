package maze

import "image/color"

// Direction represents the four cardinal directions
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// Cell represents a single cell in the maze
type Cell struct {
	X, Y    int
	Visited bool
	Walls   map[Direction]bool
}

// NewCell creates a new cell with all walls intact
func NewCell(x, y int) *Cell {
	return &Cell{
		X:       x,
		Y:       y,
		Visited: false,
		Walls: map[Direction]bool{
			North: true,
			East:  true,
			South: true,
			West:  true,
		},
	}
}

// Point represents a coordinate in the maze
type Point struct {
	X, Y int
}

// Maze represents the entire maze structure
type Maze struct {
	Width, Height int
	Cells         [][]*Cell
	Start, Finish Point
}

// NewMaze creates a new maze with the specified dimensions
func NewMaze(width, height int) *Maze {
	cells := make([][]*Cell, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]*Cell, width)
		for x := 0; x < width; x++ {
			cells[y][x] = NewCell(x, y)
		}
	}

	return &Maze{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// GetCell returns the cell at the given coordinates
func (m *Maze) GetCell(x, y int) *Cell {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return nil
	}
	return m.Cells[y][x]
}

// GetNeighbor returns the neighboring cell in the given direction
func (m *Maze) GetNeighbor(cell *Cell, dir Direction) *Cell {
	switch dir {
	case North:
		return m.GetCell(cell.X, cell.Y-1)
	case East:
		return m.GetCell(cell.X+1, cell.Y)
	case South:
		return m.GetCell(cell.X, cell.Y+1)
	case West:
		return m.GetCell(cell.X-1, cell.Y)
	}
	return nil
}

// RemoveWall removes the wall between two adjacent cells
func (m *Maze) RemoveWall(cell1, cell2 *Cell) {
	dx := cell2.X - cell1.X
	dy := cell2.Y - cell1.Y

	if dx == 1 { // cell2 is to the east of cell1
		cell1.Walls[East] = false
		cell2.Walls[West] = false
	} else if dx == -1 { // cell2 is to the west of cell1
		cell1.Walls[West] = false
		cell2.Walls[East] = false
	} else if dy == 1 { // cell2 is to the south of cell1
		cell1.Walls[South] = false
		cell2.Walls[North] = false
	} else if dy == -1 { // cell2 is to the north of cell1
		cell1.Walls[North] = false
		cell2.Walls[South] = false
	}
}

// CanMove checks if movement is possible from one cell to another
func (m *Maze) CanMove(from, to *Cell) bool {
	if from == nil || to == nil {
		return false
	}

	dx := to.X - from.X
	dy := to.Y - from.Y

	// Check if cells are adjacent
	if (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1)) {
		if dx == 1 { // moving east
			return !from.Walls[East]
		} else if dx == -1 { // moving west
			return !from.Walls[West]
		} else if dy == 1 { // moving south
			return !from.Walls[South]
		} else if dy == -1 { // moving north
			return !from.Walls[North]
		}
	}
	return false
}

// RenderConfig holds configuration for rendering the maze
type RenderConfig struct {
	CellSize       int
	WallThickness  int
	ImageWidth     int
	ImageHeight    int
	Padding        int
	HeaderHeight   int
	LegendFontSize int    // Font size multiplier for legend text
	FontPath       string // Path to TrueType font file (optional)
	WallColor      color.Color
	PathColor      color.Color
	TextColor      color.Color
}

// DefaultRenderConfig returns a default configuration optimized for 8.5"x11" printing
func DefaultRenderConfig() RenderConfig {
	return RenderConfig{
		CellSize:       84,                             // Size of each cell in pixels
		WallThickness:  8,                              // Thickness of walls in pixels
		ImageWidth:     2100,                           // ~7" at 300 DPI
		ImageHeight:    2700,                           // ~9" at 300 DPI
		Padding:        100,                            // Padding around the maze in pixels
		HeaderHeight:   120,                            // Height of header area for legend (increased for larger font)
		LegendFontSize: 3,                              // 3x font size multiplier
		WallColor:      color.RGBA{0, 0, 0, 255},       // Black
		PathColor:      color.RGBA{255, 255, 255, 255}, // White
		TextColor:      color.RGBA{0, 0, 0, 255},       // Black text
	}
}
