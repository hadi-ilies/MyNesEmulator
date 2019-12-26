package ui

import (
	"github.com/hadi-ilies/MyNesEmulator/src/nescomponents"
)

type Nes struct {
	bus *nescomponents.BUS
}

func NewNes(gamePath string) *Nes {
	var nes Nes = Nes{nescomponents.NewBus(nescomponents.NewCartridge(gamePath))}

	return &nes
}
