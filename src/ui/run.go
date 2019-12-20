package ui

import (
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

type Vector2i struct {
	X int32
	Y int32
}

func Run(path string) bool {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()
	window, err := glfw.CreateWindow(800, 600, "NES-EMULATOR", nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	for !window.ShouldClose() {
		// Do OpenGL stuff.
		window.SwapBuffers()
		glfw.PollEvents()
	}
	return true
}
