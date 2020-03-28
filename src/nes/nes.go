package nes

import (
	"image"

	"github.com/hadi-ilies/MyNesEmulator/src/nes/nescomponents"
)

/**
MEMO : FIRST LETTER of struct elem DECIDE WETHER THE ELEM IS Private or public
MAj -> public
MIN -> private
**/
type Nes struct {
	bus *nescomponents.BUS
}

func NewNes(gamePath string) Nes {
	var nes Nes = Nes{nescomponents.NewBus(nescomponents.NewCartridge(gamePath))} //load the cartridge file and insert it into the nes

	return nes
}

//reset the console
func (nes *Nes) Reset() {
	nes.bus.Reset()
}

//get the circuit that is linked with all nes components
func (nes *Nes) GetComponents() *nescomponents.BUS {
	return nes.bus
}

func (nes *Nes) PixelBuffer() *image.RGBA {
	return nes.bus.GetPpu().GetFront()
}

func (nes *Nes) Step() uint64 {
	var cpuCycles uint64 = nes.GetComponents().GetCpu().Step()
	ppuCycles := cpuCycles * 3
	var i uint64 = 0
	for i = 0; i < ppuCycles; i++ {
		nes.GetComponents().GetPpu().Step()              //todo check ppu
		nes.GetComponents().GetCartridge().Mapper.Step() //todo it depend the mapper search a fix for that, i have to test that on the other rep
	}
	// for i := 0; i < cpuCycles; i++ {
	// 	nes.APU.Step()
	// }
	return cpuCycles
}

func (nes *Nes) Run(seconds float64) {
	CPUFrequency := float64(1789773)
	cycles := int(CPUFrequency * seconds)
	for cycles > 0 {
		cycles -= int(nes.Step())
	}
}

//https://wiki.nesdev.com/w/index.php/Controller_reading_code
//???
const (
	KeyA = iota
	KeyB
	KeySelect
	KeyStart
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
)

func (nes *Nes) SetButtonToController(buttons [8]byte) {
	nes.bus.Controller1.SetButtons(buttons)
}
