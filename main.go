package main

import (
	"fmt"
	"log"
	"time"

	"mazegenerator/maze"
)

const (
	// Default maze dimensions
	DefaultWidth  = 25
	DefaultHeight = 25

	// Maximum retries for maze generation
	MaxRetries = 5
)

func main() {
	fmt.Println("Maze Generator")
	fmt.Println("==============")

	// Create generator and renderer
	generator := maze.NewGenerator()
	renderer := maze.NewDefaultRenderer()

	fmt.Printf("Generating %dx%d maze...\n", DefaultWidth, DefaultHeight)

	// Generate maze with validation
	mazeObj := generator.GenerateWithValidation(DefaultWidth, DefaultHeight, MaxRetries)

	fmt.Println("Placing start and finish points...")

	// Validate the final maze
	validator := maze.NewValidator()
	if !validator.HasPath(mazeObj) {
		fmt.Println("Warning: Generated maze may not have a valid path from start to finish")
	} else {
		fmt.Println("✓ Path verified from start to finish!")
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("maze_%s.png", timestamp)

	fmt.Printf("Rendering maze to PNG (%s)...\n", filename)

	// Get image dimensions for user info
	width, height := renderer.GetImageDimensions(mazeObj)
	fmt.Printf("Image dimensions: %dx%d pixels\n", width, height)

	// Render to PNG
	err := renderer.RenderToPNG(mazeObj, filename)
	if err != nil {
		log.Fatalf("Error rendering maze: %v", err)
	}

	// Success message
	fmt.Printf("✓ Maze saved as '%s'\n", filename)
	fmt.Printf("Start: (%d, %d) - marked with circle (○)\n", mazeObj.Start.X, mazeObj.Start.Y)
	fmt.Printf("Finish: (%d, %d) - marked with square (■)\n", mazeObj.Finish.X, mazeObj.Finish.Y)
	fmt.Println("\nThe maze is optimized for printing on 8.5\"x11\" paper.")
	fmt.Println("Legend is shown at the top of the maze.")
	fmt.Println("Ready to print and solve!")
}
