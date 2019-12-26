package ui

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
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
		return false
	}
	defer glfw.Terminate()
	//create window
	//glfw.WindowHint(glfw.ContextVersionMajor, 2)
	//glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(800, 600, "NES-EMULATOR", nil, nil)

	if err != nil {
		return false
		//panic(err)
	}
	window.MakeContextCurrent()

	// initialize opengl
	if err := gl.Init(); err != nil {
		log.Fatalln(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	//main loop
	for !window.ShouldClose() {
		// Do OpenGL stuff.
		gl.Clear(gl.COLOR_BUFFER_BIT) //clear screen
		window.SwapBuffers()
		glfw.PollEvents()
	}
	return true
}
