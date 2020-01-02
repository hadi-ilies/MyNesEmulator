package nescomponents

import "image"

//cpu can read only 8 addrs form the ppu
const (
	control = 0x0000
	Mask    = 0x0001
	Status  = 0x0002
	//The OAM (Object Attribute Memory) is internal memory inside the PPU that contains a display
	//list of up to 64 sprites, where each sprite's information occupies 4 bytes.
	oamAddr = 0x0003
	oamData = 0x0004
	Scroll  = 0x0005
	ppuAddr = 0x0006
	ppuData = 0x0007
)

//Comunication with main BUS
func (ppu *PPU) CpuWrite(address uint16, data byte) {
	switch address {
	case control:
		break
	case Mask:
		break
	case Status:
		break
	case oamAddr:
		break
	case oamData:
		break
	case Scroll:
		break
	case ppuAddr:
		break
	case ppuData:
		break
	}
}

func (ppu *PPU) CpuRead(address uint16) byte {
	var data byte = 0x00

	switch address {
	case control:
		break
	case Mask:
		break
	case Status:
		break
	case oamAddr:
		break
	case oamData:
		break
	case Scroll:
		break
	case ppuAddr:
		break
	case ppuData:
		break
	}
	return data
}

//Comunication  with the second "PPU" BUS
func (ppu *PPU) Read(address uint16) byte {
	var data byte = 0x00
	address &= 0x3FFF

	if ppu.cartridge.PpuRead(address, &data) {

	}
	return data
}

func (ppu *PPU) Write(address uint16, data byte) {
	address &= 0x3FFF

	if ppu.cartridge.PpuWrite(address, data) {

	}
}

//picture processing units
func (ppu *PPU) ConnectCartridge(cartridge *Cartridge) {
	ppu.cartridge = cartridge
}

func (ppu *PPU) clock() {

}

//indexs ppumask
const (
	flagGrayscale          = iota // 0: color; 1: grayscale
	flagShowLeftBackground        // 0: hide; 1: show
	flagShowLeftSprites           // 0: hide; 1: show
	flagShowBackground            // 0: hide; 1: show
	flagShowSprites               // 0: hide; 1: show
	flagRedTint                   // 0: normal; 1: emphasized
	flagGreenTint                 // 0: normal; 1: emphasized
	flagBlueTint                  // 0: normal; 1: emphasized
)

type PPU struct {
	cpu       *CPU       //pointer on nes's Cpu
	cartridge *Cartridge // the gamePAk
	// storage variables
	nameTable    [2][1024]byte
	paletteTable [32]byte    //ram connected to ppu that strored the palace info there are 32 entries
	oam          [256]byte   // (Object Attribute Memory)
	front        *image.RGBA // front ground that generate sprites
	back         *image.RGBA // back ground

	//circuit variable
	Cycle    int    // 0-340 nb cycles
	ScanLine int    // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame    uint64 // frame counter

	// NMI flags/vars
	nmiOccurred bool
	nmiOutput   bool
	nmiPrevious bool
	nmiDelay    byte

	// $2002 PPUSTATUS
	flagSpriteZeroHit  bool
	flagSpriteOverflow bool

	// $2001 PPUMASK
	ppuMask [8]byte
}

func (ppu *PPU) GetFront() *image.RGBA {
	return ppu.front
}

func (ppu *PPU) Reset() {
	ppu.Cycle = 340
	ppu.ScanLine = 240
	ppu.Frame = 0
	//ppu.writeControl(0)
	//ppu.writeMask(0)
	//ppu.writeOAMAddress(0)
}

func NewPpu(cpu *CPU, cartridge *Cartridge) *PPU {
	var ppu PPU

	ppu.cpu = cpu
	ppu.cartridge = cartridge
	ppu.front = image.NewRGBA(image.Rect(0, 0, 256, 240))
	return &ppu
}

func (ppu *PPU) nmiChange() {
	nmi := ppu.nmiOutput && ppu.nmiOccurred
	if nmi && !ppu.nmiPrevious {
		// TODO: this fixes some games but the delay shouldn't have to be so
		// long, so the timings are off somewhere
		ppu.nmiDelay = 15
	}
	ppu.nmiPrevious = nmi

}

// Start of vertical blanking: Set NMI_occurred in PPU to true.
// End of vertical blanking, sometime in pre-render scanline: Set NMI_occurred to false.
// Read PPUSTATUS: Return old status of NMI_occurred in bit 7, then set NMI_occurred to false.
// Write to PPUCTRL: Set NMI_output to bit 7.
func (ppu *PPU) setVerticalBlank() {
	ppu.front, ppu.back = ppu.back, ppu.front
	ppu.nmiOccurred = true
	ppu.nmiChange()
}

func (ppu *PPU) clearVerticalBlank() {
	ppu.nmiOccurred = false
	ppu.nmiChange()
}

func (ppu *PPU) Step() {
	// vblank logic
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.setVerticalBlank()
	}
	if ppu.ScanLine == 261 && ppu.Cycle == 1 {
		ppu.clearVerticalBlank()
		ppu.flagSpriteZeroHit, ppu.flagSpriteOverflow = false, false
	}
}
