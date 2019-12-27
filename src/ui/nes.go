package ui

import (
	"github.com/hadi-ilies/MyNesEmulator/src/nescomponents"
)

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
