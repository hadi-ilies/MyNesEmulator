package ui

import (
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/hadi-ilies/MyNesEmulator/src/constant"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// initialize opengl
func initOpengl() bool {
	// initialize opengl
	if err := gl.Init(); err != nil {
		return false
	}
	gl.Enable(gl.TEXTURE_2D)
	return true
}

//init whole emulator and start it
func Start(gamePath string) bool {

	err := glfw.Init()
	if err != nil {
		return false
	}
	defer glfw.Terminate() //destroy all opengl stuff when func is terminated
	//create the ui
	ui := NewUI(constant.WindowWidth*constant.Scale, constant.WindowHeight*constant.Scale, constant.UITitle)

	if !initOpengl() {
		return false
	}

	ui.Run(gamePath)
	return true
}
