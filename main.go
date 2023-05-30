package main

import (
	"image"
	"math"
	"os"
	"time"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

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

		// TODO Handle logic for key presses
		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			camPos.Y -= camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), 0)
		}
		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			camPos.Y += camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), math.Pi)
		}
		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyD) {
			camPos.X += camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), math.Pi/2)
		}
		if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyA) {
			camPos.X -= camSpeed * dt
			mat = mat.Rotated(win.Bounds().Center(), (math.Pi*3)/2)
		}

		sprite.Draw(win, mat)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}