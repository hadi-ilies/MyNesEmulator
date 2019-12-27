package views

import (
	"image"

	"github.com/hadi-ilies/MyNesEmulator/src/ui"
)

type GameView struct {
	nes     ui.Nes
	ui      ui.Ui // lol there is no inerittance in golang, I am a noob ':(
	texture uint32
	frames  []image.Image
	// director *Director
	// console  *nes.Console
	// title    string
	// hash     string
	// texture  uint32
	// record   bool
	// frames   []image.Image
}
