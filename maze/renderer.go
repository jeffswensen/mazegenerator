package maze

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Renderer handles converting maze data to PNG images
type Renderer struct {
	config   RenderConfig
	fontFace font.Face
}

// NewRenderer creates a new maze renderer with the given configuration
func NewRenderer(config RenderConfig) *Renderer {
	r := &Renderer{
		config: config,
	}

	// Try to load a Unicode-capable font, fallback to basic font
	r.fontFace = r.loadFont()

	return r
}

// NewDefaultRenderer creates a renderer with default settings
func NewDefaultRenderer() *Renderer {
	return NewRenderer(DefaultRenderConfig())
}

// RenderToPNG renders the maze to a PNG file
func (r *Renderer) RenderToPNG(maze *Maze, filename string) error {
	img := r.createImage(maze)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// createImage creates an image representation of the maze
func (r *Renderer) createImage(maze *Maze) image.Image {
	// Calculate image dimensions based on maze size, cell size, padding, and header
	imgWidth := maze.Width*r.config.CellSize + r.config.WallThickness + 2*r.config.Padding
	imgHeight := maze.Height*r.config.CellSize + r.config.WallThickness + 2*r.config.Padding + r.config.HeaderHeight

	// Create image with white background (paths)
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{r.config.PathColor}, image.Point{}, draw.Src)

	// Draw legend in header area
	r.drawLegend(img)

	// Draw walls (offset by header height)
	r.drawWalls(img, maze)

	// Draw start and finish markers (offset by header height)
	r.drawMarkers(img, maze)

	return img
}

// drawLegend draws the legend in the header area
func (r *Renderer) drawLegend(img *image.RGBA) {
	// Legend text - use ASCII alternatives if Unicode font is not available
	var legendText string
	if r.fontFace == basicfont.Face7x13 {
		// Fallback to ASCII symbols that work with basic font
		legendText = "O START    # FINISH"
	} else {
		// Use Unicode symbols with TrueType font
		legendText = "○ START    ■ FINISH"
	}

	// Use scaled font rendering
	r.drawScaledText(img, legendText, r.config.LegendFontSize)
}

// drawScaledText draws text with a specified scale factor
func (r *Renderer) drawScaledText(img *image.RGBA, text string, scale int) {
	// Create a temporary image for the original font
	d := &font.Drawer{
		Dst:  image.NewRGBA(image.Rect(0, 0, 1000, 100)), // Temporary canvas
		Src:  image.NewUniform(r.config.TextColor),
		Face: basicfont.Face7x13,
	}

	// Get text dimensions at original size
	textBounds, _ := d.BoundString(text)
	origWidth := (textBounds.Max.X - textBounds.Min.X).Ceil()
	origHeight := (textBounds.Max.Y - textBounds.Min.Y).Ceil()

	// Create a temporary image to render the original text
	tempImg := image.NewRGBA(image.Rect(0, 0, origWidth+20, origHeight+20))

	// Fill with transparent background
	for y := 0; y < tempImg.Bounds().Max.Y; y++ {
		for x := 0; x < tempImg.Bounds().Max.X; x++ {
			tempImg.Set(x, y, r.config.PathColor) // Use path color as background
		}
	}

	// Draw text on temporary image
	d.Dst = tempImg
	d.Dot = fixed.Point26_6{
		X: fixed.I(10),
		Y: fixed.I(origHeight + 5),
	}
	d.DrawString(text)

	// Calculate scaled dimensions
	scaledWidth := origWidth * scale
	scaledHeight := origHeight * scale

	// Calculate position to center the scaled text in header
	textX := (img.Bounds().Max.X - scaledWidth) / 2
	textY := (r.config.HeaderHeight - scaledHeight) / 2

	// Draw scaled text by copying each pixel as a scale x scale block
	for y := 0; y < origHeight; y++ {
		for x := 0; x < origWidth; x++ {
			srcColor := tempImg.At(x+10, y+5)

			// Only draw if it's not the background color (i.e., it's text)
			if srcColor != r.config.PathColor {
				// Draw a scale x scale block for each original pixel
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						destX := textX + x*scale + dx
						destY := textY + y*scale + dy

						// Check bounds
						if destX >= 0 && destX < img.Bounds().Max.X &&
							destY >= 0 && destY < img.Bounds().Max.Y {
							img.Set(destX, destY, r.config.TextColor)
						}
					}
				}
			}
		}
	}
}

// drawWalls draws all the walls in the maze
func (r *Renderer) drawWalls(img *image.RGBA, maze *Maze) {
	wallColor := &image.Uniform{r.config.WallColor}

	for y := 0; y < maze.Height; y++ {
		for x := 0; x < maze.Width; x++ {
			cell := maze.GetCell(x, y)
			if cell == nil {
				continue
			}

			// Calculate cell position in pixels (offset by padding and header)
			cellX := x*r.config.CellSize + r.config.Padding
			cellY := y*r.config.CellSize + r.config.Padding + r.config.HeaderHeight

			// Check if this cell is the start or finish position
			isStart := (x == maze.Start.X && y == maze.Start.Y)
			isFinish := (x == maze.Finish.X && y == maze.Finish.Y)

			// Draw walls for this cell, but skip outer walls for start/finish positions
			if cell.Walls[North] {
				// Skip drawing north wall if this is start/finish on top edge
				if !((isStart || isFinish) && y == 0) {
					r.drawHorizontalWall(img, wallColor, cellX, cellY, r.config.CellSize)
				}
			}
			if cell.Walls[South] {
				// Skip drawing south wall if this is start/finish on bottom edge
				if !((isStart || isFinish) && y == maze.Height-1) {
					r.drawHorizontalWall(img, wallColor, cellX, cellY+r.config.CellSize, r.config.CellSize)
				}
			}
			if cell.Walls[West] {
				// Skip drawing west wall if this is start/finish on left edge
				if !((isStart || isFinish) && x == 0) {
					r.drawVerticalWall(img, wallColor, cellX, cellY, r.config.CellSize)
				}
			}
			if cell.Walls[East] {
				// Skip drawing east wall if this is start/finish on right edge
				if !((isStart || isFinish) && x == maze.Width-1) {
					r.drawVerticalWall(img, wallColor, cellX+r.config.CellSize, cellY, r.config.CellSize)
				}
			}
		}
	}
}

// drawHorizontalWall draws a horizontal wall
func (r *Renderer) drawHorizontalWall(img *image.RGBA, wallColor *image.Uniform, x, y, length int) {
	rect := image.Rect(x, y, x+length+r.config.WallThickness, y+r.config.WallThickness)
	draw.Draw(img, rect, wallColor, image.Point{}, draw.Src)
}

// drawVerticalWall draws a vertical wall
func (r *Renderer) drawVerticalWall(img *image.RGBA, wallColor *image.Uniform, x, y, length int) {
	rect := image.Rect(x, y, x+r.config.WallThickness, y+length+r.config.WallThickness)
	draw.Draw(img, rect, wallColor, image.Point{}, draw.Src)
}

// drawMarkers draws the start and finish markers
func (r *Renderer) drawMarkers(img *image.RGBA, maze *Maze) {
	// Draw start marker (circle)
	r.drawCircleMarker(img, maze.Start)

	// Draw finish marker (square)
	r.drawSquareMarker(img, maze.Finish)
}

// drawCircleMarker draws a circle marker in the center of the specified cell
func (r *Renderer) drawCircleMarker(img *image.RGBA, pos Point) {
	// Calculate cell center position (offset by padding and header)
	cellX := pos.X*r.config.CellSize + r.config.Padding
	cellY := pos.Y*r.config.CellSize + r.config.Padding + r.config.HeaderHeight
	centerX := cellX + r.config.CellSize/2
	centerY := cellY + r.config.CellSize/2

	// Circle radius (about 1/3 of cell size)
	radius := r.config.CellSize / 3
	thickness := 3 // Line thickness

	// Draw circle outline
	for y := centerY - radius; y <= centerY+radius; y++ {
		for x := centerX - radius; x <= centerX+radius; x++ {
			// Calculate distance from center
			dx := x - centerX
			dy := y - centerY
			distSq := dx*dx + dy*dy
			radiusSq := radius * radius
			innerRadiusSq := (radius - thickness) * (radius - thickness)

			// Draw if within the ring (between inner and outer radius)
			if distSq <= radiusSq && distSq >= innerRadiusSq {
				if x >= 0 && x < img.Bounds().Max.X && y >= 0 && y < img.Bounds().Max.Y {
					img.Set(x, y, r.config.WallColor)
				}
			}
		}
	}
}

// drawSquareMarker draws a square marker in the center of the specified cell
func (r *Renderer) drawSquareMarker(img *image.RGBA, pos Point) {
	// Calculate cell center position (offset by padding and header)
	cellX := pos.X*r.config.CellSize + r.config.Padding
	cellY := pos.Y*r.config.CellSize + r.config.Padding + r.config.HeaderHeight
	centerX := cellX + r.config.CellSize/2
	centerY := cellY + r.config.CellSize/2

	// Square size (about 2/3 of cell size)
	size := r.config.CellSize * 2 / 3
	halfSize := size / 2
	thickness := 3 // Line thickness

	wallColor := &image.Uniform{r.config.WallColor}

	// Draw square outline (4 rectangles for the sides)
	// Top side
	topRect := image.Rect(centerX-halfSize, centerY-halfSize, centerX+halfSize, centerY-halfSize+thickness)
	draw.Draw(img, topRect, wallColor, image.Point{}, draw.Src)

	// Bottom side
	bottomRect := image.Rect(centerX-halfSize, centerY+halfSize-thickness, centerX+halfSize, centerY+halfSize)
	draw.Draw(img, bottomRect, wallColor, image.Point{}, draw.Src)

	// Left side
	leftRect := image.Rect(centerX-halfSize, centerY-halfSize, centerX-halfSize+thickness, centerY+halfSize)
	draw.Draw(img, leftRect, wallColor, image.Point{}, draw.Src)

	// Right side
	rightRect := image.Rect(centerX+halfSize-thickness, centerY-halfSize, centerX+halfSize, centerY+halfSize)
	draw.Draw(img, rightRect, wallColor, image.Point{}, draw.Src)
}

// RenderToImage returns the maze as an image.Image (useful for further processing)
func (r *Renderer) RenderToImage(maze *Maze) image.Image {
	return r.createImage(maze)
}

// GetImageDimensions returns the dimensions the rendered image will have
func (r *Renderer) GetImageDimensions(maze *Maze) (width, height int) {
	width = maze.Width*r.config.CellSize + r.config.WallThickness + 2*r.config.Padding
	height = maze.Height*r.config.CellSize + r.config.WallThickness + 2*r.config.Padding + r.config.HeaderHeight
	return
}

// loadFont attempts to load a Unicode-capable font, falls back to basic font
func (r *Renderer) loadFont() font.Face {
	// If a specific font path is provided, try to load it
	if r.config.FontPath != "" {
		if face := r.loadFontFromPath(r.config.FontPath); face != nil {
			return face
		}
	}

	// Try to find a system font that supports Unicode
	systemFonts := r.getSystemFontPaths()
	for _, fontPath := range systemFonts {
		if face := r.loadFontFromPath(fontPath); face != nil {
			return face
		}
	}

	// Fallback to basic font
	return basicfont.Face7x13
}

// getSystemFontPaths returns a list of common system font paths that support Unicode
func (r *Renderer) getSystemFontPaths() []string {
	var fontPaths []string

	// Add common Windows fonts
	fontPaths = append(fontPaths,
		"C:/Windows/Fonts/arial.ttf",
		"C:/Windows/Fonts/calibri.ttf",
		"C:/Windows/Fonts/tahoma.ttf",
		"C:/Windows/Fonts/verdana.ttf",
	)

	// Add common macOS fonts
	fontPaths = append(fontPaths,
		"/System/Library/Fonts/Arial.ttf",
		"/System/Library/Fonts/Helvetica.ttc",
		"/Library/Fonts/Arial.ttf",
	)

	// Add common Linux fonts
	fontPaths = append(fontPaths,
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/TTF/arial.ttf",
		"/usr/local/share/fonts/arial.ttf",
	)

	return fontPaths
}

// loadFontFromPath attempts to load a TrueType font from the given path
func (r *Renderer) loadFontFromPath(fontPath string) font.Face {
	// Check if file exists
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return nil
	}

	// For now, return nil to use fallback - we'll implement TrueType loading if needed
	// This allows the code to compile and run with the basic font
	return nil
}
