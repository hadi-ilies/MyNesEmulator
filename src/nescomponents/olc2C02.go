package nescomponents

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

type PPU struct {
	cpu       *CPU       //pointer on nes's Cpu
	cartridge *Cartridge // the gamePAk
	// storage variables
	nameTable    [2][1024]byte
	paletteTable [32]byte  //ram connected to ppu that strored the palace info there are 32 entries
	oam          [256]byte // (Object Attribute Memory)
	//front         *image.RGBA
	//back          *image.RGBA

}
