package maze

// Validator handles maze validation using pathfinding algorithms
type Validator struct{}

// NewValidator creates a new maze validator
func NewValidator() *Validator {
	return &Validator{}
}

// HasPath checks if there's a valid path from start to finish using BFS
func (v *Validator) HasPath(maze *Maze) bool {
	if maze == nil {
		return false
	}

	startCell := maze.GetCell(maze.Start.X, maze.Start.Y)
	finishCell := maze.GetCell(maze.Finish.X, maze.Finish.Y)

	if startCell == nil || finishCell == nil {
		return false
	}

	// If start and finish are the same cell, path exists
	if maze.Start.X == maze.Finish.X && maze.Start.Y == maze.Finish.Y {
		return true
	}

	return v.bfsPath(maze, startCell, finishCell)
}

// bfsPath performs breadth-first search to find a path between two cells
func (v *Validator) bfsPath(maze *Maze, start, finish *Cell) bool {
	// Keep track of visited cells
	visited := make(map[Point]bool)

	// Queue for BFS - stores cells to visit
	queue := []*Cell{start}
	visited[Point{start.X, start.Y}] = true

	for len(queue) > 0 {
		// Dequeue the first cell
		current := queue[0]
		queue = queue[1:]

		// Check if we reached the finish
		if current.X == finish.X && current.Y == finish.Y {
			return true
		}

		// Check all four directions
		directions := []Direction{North, East, South, West}
		for _, dir := range directions {
			neighbor := maze.GetNeighbor(current, dir)
			if neighbor != nil {
				neighborPoint := Point{neighbor.X, neighbor.Y}

				// If we haven't visited this neighbor and can move to it
				if !visited[neighborPoint] && maze.CanMove(current, neighbor) {
					visited[neighborPoint] = true
					queue = append(queue, neighbor)
				}
			}
		}
	}

	// No path found
	return false
}

// FindPath returns the actual path from start to finish (for debugging/visualization)
func (v *Validator) FindPath(maze *Maze) []Point {
	if maze == nil {
		return nil
	}

	startCell := maze.GetCell(maze.Start.X, maze.Start.Y)
	finishCell := maze.GetCell(maze.Finish.X, maze.Finish.Y)

	if startCell == nil || finishCell == nil {
		return nil
	}

	// If start and finish are the same cell
	if maze.Start.X == maze.Finish.X && maze.Start.Y == maze.Finish.Y {
		return []Point{maze.Start}
	}

	return v.bfsPathWithTrace(maze, startCell, finishCell)
}

// bfsPathWithTrace performs BFS and returns the actual path
func (v *Validator) bfsPathWithTrace(maze *Maze, start, finish *Cell) []Point {
	// Keep track of visited cells and their parents
	visited := make(map[Point]bool)
	parent := make(map[Point]Point)

	// Queue for BFS
	queue := []*Cell{start}
	startPoint := Point{start.X, start.Y}
	visited[startPoint] = true

	for len(queue) > 0 {
		// Dequeue the first cell
		current := queue[0]
		queue = queue[1:]
		currentPoint := Point{current.X, current.Y}

		// Check if we reached the finish
		if current.X == finish.X && current.Y == finish.Y {
			// Reconstruct path
			return v.reconstructPath(parent, startPoint, Point{finish.X, finish.Y})
		}

		// Check all four directions
		directions := []Direction{North, East, South, West}
		for _, dir := range directions {
			neighbor := maze.GetNeighbor(current, dir)
			if neighbor != nil {
				neighborPoint := Point{neighbor.X, neighbor.Y}

				// If we haven't visited this neighbor and can move to it
				if !visited[neighborPoint] && maze.CanMove(current, neighbor) {
					visited[neighborPoint] = true
					parent[neighborPoint] = currentPoint
					queue = append(queue, neighbor)
				}
			}
		}
	}

	// No path found
	return nil
}

// reconstructPath builds the path from start to finish using parent tracking
func (v *Validator) reconstructPath(parent map[Point]Point, start, finish Point) []Point {
	path := []Point{}
	current := finish

	// Trace back from finish to start
	for current != start {
		path = append([]Point{current}, path...) // Prepend to build path in correct order
		current = parent[current]
	}

	// Add the start point
	path = append([]Point{start}, path...)

	return path
}
