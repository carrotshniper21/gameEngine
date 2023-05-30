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

func drawParticle(win *pixelgl.Window, p Particle, size float64) {
	imd := imdraw.New(nil)
	imd.Color = color.RGBA{R: p.Color.R, G: p.Color.G, B: p.Color.B, A: p.Color.A}
	imd.Push(p.Position.Sub(pixel.V(size/2, size/2)), p.Position.Add(pixel.V(size/2, size/2)))
	imd.Rectangle(0)
	imd.Draw(win)
}
func drawRedTrail(win *pixelgl.Window, particles []Particle) {
	for _, p := range particles {
		drawParticle(win, p, 2.0)
	}
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

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Coom.io",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)

	pic, err := loadPicture("arrow.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	var (
		camPos   = pixel.ZV
		camSpeed = 500.0
		trail    []Particle
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

		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			camPos.Y += camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), 0)
			newParticle := Particle{Position: pixel.V(camPos.X, camPos.Y-1), Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}}
			trail = append(trail, newParticle)
		}
		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			camPos.Y -= camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), math.Pi)
			newParticle := Particle{Position: pixel.V(camPos.X, camPos.Y+1), Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}}
			trail = append(trail, newParticle)
		}
		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			camPos.X -= camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), math.Pi/2)
			newParticle := Particle{Position: pixel.V(camPos.X-1, camPos.Y), Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}}
			trail = append(trail, newParticle)
		}
		if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			camPos.X += camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), (math.Pi*3)/2)
			newParticle := Particle{Position: pixel.V(camPos.X+1, camPos.Y), Color: color.RGBA{R: 255, G: 0, B: 0, A: 255}}
			trail = append(trail, newParticle)
		}

		sprite.Draw(win, mat)
		drawRedTrail(win, trail)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
