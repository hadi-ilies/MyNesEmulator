package ui

import (
	//	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/hadi-ilies/MyNesEmulator/src/nes"
)

//Ui my ui struct
type Ui struct {
	window     *glfw.Window
	actualView View
	timestamp  float64
}

//NewUI is the constructor of my ui
func NewUI(width int, height int, uiTitle string) *Ui {
	var ui Ui

	//create and Init window
	window, err := glfw.CreateWindow(width, height, uiTitle, nil, nil)

	if err != nil {
		println("print usage and error")
	}
	window.MakeContextCurrent()

	ui.window = window
	ui.timestamp = 0
	return &ui
}

//GetWindow return a pointer on the ui's window
func (ui *Ui) GetWindow() *glfw.Window { // i have created this func because i can't access elem of struct in another packet golang allow method only. have to check this on internet.
	return ui.window
}

func (ui *Ui) getInView(view View) {
	if ui.actualView != nil {
		ui.actualView.End()
	}
	ui.actualView = view
	if ui.actualView != nil {
		ui.actualView.Start()
	}
	ui.timestamp = glfw.GetTime()
}

//playGame
func (ui *Ui) loadGame(gamePath string) {
	var nes nes.Nes = nes.NewNes(gamePath)

	ui.getInView(NewGameView(ui, &nes))
}

func (ui *Ui) displayView() {
	timestamp := glfw.GetTime()
	difftime := timestamp - ui.timestamp
	ui.timestamp = timestamp
	if ui.actualView != nil {
		ui.actualView.Update(difftime)
	}
}

//start UI it is the main loop
func (ui *Ui) Run(gamePath string) {
	//load the emulator and views
	ui.loadGame(gamePath)
	//main loop
	for !ui.window.ShouldClose() {
		// clear screen at each loop's turn.
		gl.Clear(gl.COLOR_BUFFER_BIT)
		// display ui screen
		ui.displayView()
		// SwapBuffers swaps the front and back buffers of the window
		ui.window.SwapBuffers()
		//well be useful when i will code the controllers
		glfw.PollEvents()
	}
	ui.getInView(nil) //tmp maybe useless
}
