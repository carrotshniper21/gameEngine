package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width  = 800
	height = 600
)

type point struct {
	x, y float64
}

type snake struct {
	body        []point
	dir         point
	boostSpeed  float64
	boostActive bool
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Snake",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	snake := snake{
		body:       []point{{width / 2, height / 2}},
		dir:        point{10, 0},
		boostSpeed: 1.0,
	}
	food := point{rand.Float64() * width, rand.Float64() * height}

	for !win.Closed() {
		win.Clear(colornames.Black)
		drawBorder(win, width, height)

		// Handle input
		if win.Pressed(pixelgl.KeyW) || win.Pressed(pixelgl.KeyUp) {
			snake.dir = point{0, 10}
		}
		if win.Pressed(pixelgl.KeyS) || win.Pressed(pixelgl.KeyDown) {
			snake.dir = point{0, -10}
		}
		if win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyLeft) {
			snake.dir = point{-10, 0}
		}
		if win.Pressed(pixelgl.KeyD) || win.Pressed(pixelgl.KeyRight) {
			snake.dir = point{10, 0}
		}
		if win.Pressed(pixelgl.KeySpace) {
			if !snake.boostActive {
				snake.boostActive = true
				snake.boostSpeed = 2.0
			} else {
				snake.boostActive = false
				snake.boostSpeed = 1.0
			}
		}

		// Move snake
		head := snake.body[len(snake.body)-1]
		newHead := point{head.x + snake.dir.x*snake.boostSpeed, head.y + snake.dir.y*snake.boostSpeed}
		if newHead.x < 0 || newHead.x >= width || newHead.y < 0 || newHead.y >= height {
			break
		}
		for _, p := range snake.body {
			if p == newHead {
				break
			}
		}

		// Draw snake
		for _, p := range snake.body {
			imd := imdraw.New(nil)
			imd.Color = colornames.Green
			imd.Push(pixel.V(p.x, p.y))
			imd.Circle(10, 0)
			imd.Draw(win)
		}

		// Draw food
		imd := imdraw.New(nil)
		imd.Color = colornames.Red
		imd.Push(pixel.V(food.x, food.y))

		imd.Circle(10, 0)
		imd.Draw(win)

		// Check if the snake eats the food
		if pointsCloseEnough(newHead, food, 10) {
			food = point{rand.Float64() * width, rand.Float64() * height}
			snake.body = append(snake.body, newHead)
		} else {
			snake.body = append(snake.body[1:], newHead)
		}

		win.Update()

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	pixelgl.Run(run)
}

func pointsCloseEnough(a, b point, tolerance float64) bool {
	return math.Abs(a.x-b.x) < tolerance && math.Abs(a.y-b.y) < tolerance
}

func drawBorder(win *pixelgl.Window, width, height float64) {
	imd := imdraw.New(nil)
	imd.Color = colornames.White

	// Draw the top border
	imd.Push(pixel.V(0, height), pixel.V(width, height))

	// Draw the bottom border
	imd.Push(pixel.V(0, 0), pixel.V(width, 0))

	// Draw the left border
	imd.Push(pixel.V(0, 0), pixel.V(0, height))

	// Draw the right border
	imd.Push(pixel.V(width, 0), pixel.V(width, height))

	imd.Line(1)
	imd.Draw(win)
}