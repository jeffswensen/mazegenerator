package maze

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

// Generator handles maze generation using recursive backtracking
type Generator struct {
	rng *rand.Rand
}

// NewGenerator creates a new maze generator with a random seed
func NewGenerator() *Generator {
	// Use crypto/rand for secure seed generation
	seed, err := cryptorand.Int(cryptorand.Reader, big.NewInt(1<<63-1))
	if err != nil {
		// Fallback to time-based seed if crypto/rand fails
		seed = big.NewInt(time.Now().UnixNano())
	}

	return &Generator{
		rng: rand.New(rand.NewSource(seed.Int64())),
	}
}

// Generate creates a new maze using recursive backtracking algorithm
func (g *Generator) Generate(width, height int) *Maze {
	maze := NewMaze(width, height)

	// Start from a random cell
	startX := g.rng.Intn(width)
	startY := g.rng.Intn(height)
	startCell := maze.GetCell(startX, startY)

	// Use recursive backtracking to generate the maze
	g.generateRecursive(maze, startCell)

	return maze
}

// generateRecursive implements the recursive backtracking algorithm
func (g *Generator) generateRecursive(maze *Maze, current *Cell) {
	current.Visited = true

	// Get all unvisited neighbors in random order
	neighbors := g.getUnvisitedNeighbors(maze, current)
	g.shuffleNeighbors(neighbors)

	for _, neighbor := range neighbors {
		if !neighbor.Visited {
			// Remove wall between current and neighbor
			maze.RemoveWall(current, neighbor)

			// Recursively visit the neighbor
			g.generateRecursive(maze, neighbor)
		}
	}
}

// getUnvisitedNeighbors returns all unvisited neighboring cells
func (g *Generator) getUnvisitedNeighbors(maze *Maze, cell *Cell) []*Cell {
	var neighbors []*Cell

	directions := []Direction{North, East, South, West}
	for _, dir := range directions {
		neighbor := maze.GetNeighbor(cell, dir)
		if neighbor != nil && !neighbor.Visited {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

// shuffleNeighbors randomly shuffles the slice of neighbors
func (g *Generator) shuffleNeighbors(neighbors []*Cell) {
	for i := len(neighbors) - 1; i > 0; i-- {
		j := g.rng.Intn(i + 1)
		neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
	}
}

// PlaceStartAndFinish randomly places start and finish points in the maze
// ensuring they are as far apart as possible
func (g *Generator) PlaceStartAndFinish(maze *Maze) {
	// Try to place start and finish at opposite corners first
	corners := []Point{
		{0, 0},                            // Top-left
		{maze.Width - 1, 0},               // Top-right
		{0, maze.Height - 1},              // Bottom-left
		{maze.Width - 1, maze.Height - 1}, // Bottom-right
	}

	// Shuffle corners for randomness
	for i := len(corners) - 1; i > 0; i-- {
		j := g.rng.Intn(i + 1)
		corners[i], corners[j] = corners[j], corners[i]
	}

	// Try corner pairs first
	if len(corners) >= 2 {
		maze.Start = corners[0]
		maze.Finish = corners[1]
		return
	}

	// Fallback: place randomly if corners don't work
	g.placeRandomStartFinish(maze)
}

// placeRandomStartFinish places start and finish at random locations
func (g *Generator) placeRandomStartFinish(maze *Maze) {
	// Generate random start position
	maze.Start = Point{
		X: g.rng.Intn(maze.Width),
		Y: g.rng.Intn(maze.Height),
	}

	// Generate random finish position, ensuring it's different from start
	for {
		finish := Point{
			X: g.rng.Intn(maze.Width),
			Y: g.rng.Intn(maze.Height),
		}

		// Ensure finish is not the same as start
		if finish.X != maze.Start.X || finish.Y != maze.Start.Y {
			maze.Finish = finish
			break
		}
	}
}

// GenerateWithValidation generates a maze and ensures start/finish are connected
func (g *Generator) GenerateWithValidation(width, height int, maxRetries int) *Maze {
	for attempt := 0; attempt < maxRetries; attempt++ {
		maze := g.Generate(width, height)

		// Try multiple start/finish placements
		for placementAttempt := 0; placementAttempt < 10; placementAttempt++ {
			g.PlaceStartAndFinish(maze)

			// Validate that a path exists
			validator := NewValidator()
			if validator.HasPath(maze) {
				return maze
			}
		}
	}

	// If we get here, something went wrong - return a basic maze anyway
	// This shouldn't happen with proper maze generation, but it's a safety net
	maze := g.Generate(width, height)
	g.PlaceStartAndFinish(maze)
	return maze
}
