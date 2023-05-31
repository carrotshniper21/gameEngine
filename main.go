package main

import (
	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Particle struct {
	Position pixel.Vec
	Color    color.RGBA
}

type Grid struct {
	Size      pixel.Vec
	CellSize  float64
	CellColor color.RGBA
}

func drawParticle(win *pixelgl.Window, p Particle, size float64) {
	imd := imdraw.New(nil)
	imd.Color = color.RGBA{R: p.Color.R, G: p.Color.G, B: p.Color.B, A: p.Color.A}
	imd.Push(p.Position.Sub(pixel.V(size/2, size/2)), p.Position.Add(pixel.V(size/2, size/2)))
	imd.Rectangle(0)
	imd.Draw(win)
}

func drawRedTrail(win *pixelgl.Window, particles []Particle) {
	for _, p := range particles {
		drawParticle(win, p, 32.0)
	}
}

func drawGrid(win *pixelgl.Window, grid Grid) {
	imd := imdraw.New(nil)
	imd.Color = grid.CellColor

	for x := 0.0; x < grid.Size.X; x += grid.CellSize {
		imd.Push(pixel.V(x, 0), pixel.V(x, grid.Size.Y))
		imd.Line(1)
	}

	for y := 0.0; y < grid.Size.Y; y += grid.CellSize {
		imd.Push(pixel.V(0, y), pixel.V(grid.Size.X, y))
		imd.Line(1)
	}

	imd.Draw(win)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func gameloop(win *pixelgl.Window, sprite pixel.Picture, cfg pixelgl.WindowConfig) {
	var (
		camPos       = pixel.ZV
		camSpeed     = 500.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
		trail        []Particle
		currentKey   = pixelgl.KeyUnknown
		isPaused     = false
		gridSize     = pixel.V(800, 600) // Adjust the grid size as needed
		cellSize     = 32.0              // Adjust the cell size as needed
		gridColor    = colornames.Lightgray
		grid         = Grid{
			Size:      gridSize,
			CellSize:  cellSize,
			CellColor: gridColor,
		}
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(colornames.Seashell)

		mat := pixel.IM
		cam := mat.Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		mat = mat.Moved(win.Bounds().Center())

		if win.JustPressed(pixelgl.KeyUp) {
			currentKey = pixelgl.KeyUp
			isPaused = false
		} else if win.JustPressed(pixelgl.KeyDown) {
			currentKey = pixelgl.KeyDown
			isPaused = false
		} else if win.JustPressed(pixelgl.KeyLeft) {
			currentKey = pixelgl.KeyLeft
			isPaused = false
		} else if win.JustPressed(pixelgl.KeyRight) {
			currentKey = pixelgl.KeyRight
			isPaused = false
		}

		if win.JustPressed(pixelgl.KeySpace) {
			isPaused = !isPaused
		}

		if !isPaused {
			switch currentKey {
			case pixelgl.KeyUp:
				camPos.Y += camSpeed * dt
				mat = mat.Rotated(win.Bounds().Center(), 0)
				newParticle := Particle{
					Position: pixel.V(camPos.X, camPos.Y-1),
					Color:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
				}
				trail = append(trail, newParticle)
			case pixelgl.KeyDown:
				camPos.Y -= camSpeed * dt
				mat = mat.Rotated(win.Bounds().Center(), math.Pi)
				newParticle := Particle{
					Position: pixel.V(camPos.X, camPos.Y+1),
					Color:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
				}
				trail = append(trail, newParticle)
			case pixelgl.KeyLeft:
				camPos.X -= camSpeed * dt
				mat = mat.Rotated(win.Bounds().Center(), math.Pi/2)
				newParticle := Particle{
					Position: pixel.V(camPos.X-1, camPos.Y),
					Color:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
				}
				trail = append(trail, newParticle)
			case pixelgl.KeyRight:
				camPos.X += camSpeed * dt
				mat = mat.Rotated(win.Bounds().Center(), (math.Pi*3)/2)
				newParticle := Particle{
					Position: pixel.V(camPos.X+1, camPos.Y),
					Color:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
				}
				trail = append(trail, newParticle)
			}
		}

		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		drawGrid(win, grid)
		drawRedTrail(win, trail)
		win.Update()
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Game",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	sprite, err := loadPicture("arrow.png")
	if err != nil {
		panic(err)
	}

	gameloop(win, sprite, cfg)
}

func main() {
	pixelgl.Run(run)
}
