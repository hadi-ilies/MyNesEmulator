package ui

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/hadi-ilies/MyNesEmulator/src/nes"
)

type Ui struct {
	window     *glfw.Window
	actualView View
	timestamp  float64
}

//UI constructor
func NewUI(width int, height int, uiTitle string) *Ui {
	var ui Ui

	//create and Init window
	window, err := glfw.CreateWindow(width, height, uiTitle, nil, nil)

	if err != nil {
		println("print usage and error")
	}
	window.MakeContextCurrent() //make context of window current

	ui.window = window
	ui.timestamp = 0
	return &ui
}

//return a pointer on the window of the ui
func (ui *Ui) GetWindow() *glfw.Window { // i have created this func because i can't access elem of struct in another packet golang allow method only. have to check this on internet.
	return ui.window
}

func (ui *Ui) displayView(view View) {
	if ui.actualView != nil {
		ui.actualView.End()
	}
	ui.actualView = view
	if ui.actualView != nil {
		ui.actualView.Start()
	}
	ui.timestamp = glfw.GetTime()
}

//ToDO refacto code move nes into nescomponent and create dir in it
//playGame
func (ui *Ui) loadGame(gamePath string) {
	var nes nes.Nes = nes.NewNes(gamePath)

	ui.displayView(NewGameView(ui, &nes))
}

//start UI it is the main loop
func (ui *Ui) Run(gamePath string) {
	//load the emulator and views
	ui.loadGame(gamePath)
	//main loop
	for !ui.window.ShouldClose() {
		// Do OpenGL stuff.
		//todo start emulator
		gl.Clear(gl.COLOR_BUFFER_BIT) //clear screen
		ui.window.SwapBuffers()
		glfw.PollEvents()
	}
}
