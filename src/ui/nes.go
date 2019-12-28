package ui

import (
	"github.com/hadi-ilies/MyNesEmulator/src/nescomponents"
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
	nes.Reset()
}

func (nes *Nes) Display() {

}

//get the circuit that is linked with all nes components
func (nes *Nes) GetComponents() *nescomponents.BUS {
	return nes.bus
}
