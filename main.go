package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	config := pixelgl.WindowConfig{
		Title:  "Game of Life",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	window, err := pixelgl.NewWindow(config)
	if err != nil {
		panic(err)
	}
	for !window.Closed() {
		window.Update()
	}
}
