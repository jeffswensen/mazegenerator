# Maze Generator

A Go application that generates random mazes and saves them as PNG files optimized for printing on 8.5"x11" paper.

## Features

- **Random Maze Generation**: Uses recursive backtracking algorithm to create unique mazes every time
- **Path Validation**: Ensures every generated maze has a valid path from start to finish using BFS
- **Print-Optimized**: Output is sized perfectly for standard 8.5"x11" paper at 300 DPI
- **Visual Markers**: Circle (○) symbol for start, square (■) symbol for finish, with legend
- **High Quality**: Crisp black walls on white background for excellent print quality

## Usage

### Running the Application

```bash
go run main.go
```

The application will:
1. Generate a 25x25 maze using recursive backtracking
2. Randomly place start and finish points
3. Validate that a path exists between start and finish
4. Render the maze to a PNG file with timestamp (e.g., `maze_20250628_093000.png`)
5. Display generation details and file information

### Output

- **File Format**: PNG image
- **Dimensions**: ~2100x2700 pixels (optimized for 8.5"x11" printing)
- **Colors**: 
  - Black walls
  - White paths
  - Legend showing circle (○) for start and square (■) for finish

## Project Structure

```
mazegenerator/
├── main.go              # Main application entry point
├── maze/
│   ├── types.go         # Core data structures and types
│   ├── generator.go     # Maze generation using recursive backtracking
│   ├── validator.go     # Path validation using BFS
│   └── renderer.go      # PNG rendering and image creation
├── go.mod              # Go module definition
└── README.md           # This file
```

## Algorithm Details

### Maze Generation
- **Algorithm**: Recursive Backtracking
- **Process**: 
  1. Start from random cell
  2. Mark current cell as visited
  3. Randomly select unvisited neighbor
  4. Remove wall between current and neighbor
  5. Recursively visit neighbor
  6. Backtrack when no unvisited neighbors remain

### Path Validation
- **Algorithm**: Breadth-First Search (BFS)
- **Purpose**: Ensures every maze is solvable
- **Process**:
  1. Start from the start position
  2. Explore all reachable cells level by level
  3. Return true if finish position is reached
  4. Retry start/finish placement if no path exists

### Rendering
- **Cell Size**: 84 pixels per cell
- **Wall Thickness**: 8 pixels
- **Markers**: Circle symbol for start position, square symbol for finish position, with legend header
- **Image Format**: RGBA PNG with high contrast colors

## Configuration

The default configuration is optimized for printing, but can be customized by modifying the `DefaultRenderConfig()` function in `maze/types.go`:

```go
func DefaultRenderConfig() RenderConfig {
    return RenderConfig{
        CellSize:      84,                             // Size of each cell in pixels
        WallThickness: 8,                              // Thickness of walls in pixels
        ImageWidth:    2100,                           // ~7" at 300 DPI
        ImageHeight:   2700,                           // ~9" at 300 DPI
        WallColor:     color.RGBA{0, 0, 0, 255},       // Black
        PathColor:     color.RGBA{255, 255, 255, 255}, // White
        TextColor:     color.RGBA{0, 0, 0, 255},       // Black text
        FontSize:      12,                             // Font size for text markers
    }
}
```

## Requirements

- Go 1.16 or later
- No external dependencies (uses only Go standard library)

## Building

To build a standalone executable:

```bash
go build -o mazegenerator main.go
```

Then run:
```bash
./mazegenerator
```

## Example Output

```
Maze Generator
==============
Generating 25x25 maze...
Placing start and finish points...
✓ Path verified from start to finish!
Rendering maze to PNG (maze_20250628_093332.png)...
Image dimensions: 2108x2708 pixels
✓ Maze saved as 'maze_20250628_093332.png'
Start: (0, 24) - marked with circle (○) symbol
Finish: (24, 0) - marked with square (■) symbol

The maze is optimized for printing on 8.5"x11" paper.
Ready to print and solve!
```

## License

This project is open source and available under the MIT License.
