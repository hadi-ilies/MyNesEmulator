package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func createWindow(width uint, height uint, bitPerPixel uint, title string) {
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "NES-EMULATOR",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	for !win.Closed() {
		win.Update()
	}
}

func Run(path string) bool {
	pixelgl.Run(run)
	return true
}
