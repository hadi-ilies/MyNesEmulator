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

//indexs ppuctrl
const (
	flagNameTable       = iota // 0: $2000; 1: $2400; 2: $2800; 3: $2C00
	flagIncrement              // 0: add 1; 1: add 32
	flagSpriteTable            // 0: $0000; 1: $1000; ignored in 8x16 mode
	flagBackgroundTable        // 0: $0000; 1: $1000
	flagSpriteSize             // 0: 8x8; 1: 8x16
	flagMasterSlave            // 0: read EXT; 1: write EXT
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

	// PPU registers
	v uint16 // current vram address (15 bit)
	t uint16 // temporary vram address (15 bit)
	x byte   // fine x scroll (3 bit)
	w byte   // write toggle (1 bit)
	f byte   // even/odd frame flag (1 bit)

	//circuit variable
	Cycle    int    // 0-340 nb cycles
	ScanLine int    // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame    uint64 // frame counter

	// NMI flags/vars
	nmiOccurred bool
	nmiOutput   bool
	nmiPrevious bool
	nmiDelay    byte

	// background temporary variables
	nameTableByte      byte
	attributeTableByte byte
	lowTileByte        byte
	highTileByte       byte
	tileData           uint64

	// sprite temporary variables
	spriteCount      int
	spritePatterns   [8]uint32
	spritePositions  [8]byte
	spritePriorities [8]byte
	spriteIndexes    [8]byte

	// $2002 PPUSTATUS
	flagSpriteZeroHit  bool
	flagSpriteOverflow bool

	// $2000 PPUCTRL
	ppuCtrl [6]byte

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

func (ppu *PPU) isRenderingEnabled() bool {
	return ppu.ppuMask[flagShowBackground] != 0 || ppu.ppuMask[flagShowSprites] != 0
}

func (ppu *PPU) fetchSpritePattern(i, row int) uint32 {
	tile := ppu.oam[i*4+1]
	attributes := ppu.oam[i*4+2]
	var address uint16
	if ppu.ppuCtrl[flagSpriteSize] == 0 {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		table := ppu.ppuCtrl[flagSpriteTable]
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	} else {
		if attributes&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE
		if row > 7 {
			tile++
			row -= 8
		}
		address = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	}
	a := (attributes & 3) << 2
	lowTileByte := ppu.Read(address)
	highTileByte := ppu.Read(address + 8)
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (lowTileByte & 1) << 0
			p2 = (highTileByte & 1) << 1
			lowTileByte >>= 1
			highTileByte >>= 1
		} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (ppu *PPU) evaluateSprites() {
	var h int

	if ppu.ppuCtrl[flagSpriteSize] == 0 {
		h = 8
	} else {
		h = 16
	}
	count := 0
	for i := 0; i < 64; i++ {
		y := ppu.oam[i*4]
		a := ppu.oam[i*4+2]
		x := ppu.oam[i*4+3]
		row := ppu.ScanLine - int(y)
		if row < 0 || row >= h {
			continue
		}
		if count < 8 {
			ppu.spritePatterns[count] = ppu.fetchSpritePattern(i, row)
			ppu.spritePositions[count] = x
			ppu.spritePriorities[count] = (a >> 5) & 1
			ppu.spriteIndexes[count] = byte(i)
		}
		count++
	}
	if count > 8 {
		count = 8
		ppu.flagSpriteOverflow = true
	}
	ppu.spriteCount = count
}

func (ppu *PPU) fetchNameTableByte() {
	v := ppu.v
	address := 0x2000 | (v & 0x0FFF)
	ppu.nameTableByte = ppu.Read(address)
}

func (ppu *PPU) fetchAttributeTableByte() {
	v := ppu.v
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	ppu.attributeTableByte = ((ppu.Read(address) >> shift) & 3) << 2
}

func (ppu *PPU) fetchLowTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := ppu.ppuCtrl[flagBackgroundTable]
	tile := ppu.ppuCtrl[ppu.nameTableByte]
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.lowTileByte = ppu.Read(address)
}

func (ppu *PPU) fetchHighTileByte() {
	fineY := (ppu.v >> 12) & 7
	table := ppu.ppuCtrl[flagBackgroundTable]
	tile := ppu.nameTableByte
	address := 0x1000*uint16(table) + uint16(tile)*16 + fineY
	ppu.highTileByte = ppu.Read(address + 8)
}

func (ppu *PPU) storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		a := ppu.attributeTableByte
		p1 := (ppu.lowTileByte & 0x80) >> 7
		p2 := (ppu.highTileByte & 0x80) >> 6
		ppu.lowTileByte <<= 1
		ppu.highTileByte <<= 1
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	ppu.tileData |= uint64(data)
}

func (ppu *PPU) Step() {
	visibleLine := ppu.ScanLine < 240
	preLine := ppu.ScanLine == 261
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle
	renderLine := preLine || visibleLine

	//background
	if ppu.isRenderingEnabled() {
		if visibleLine && visibleCycle {
			//ppu.renderPixel()
		}
		if renderLine && fetchCycle {
			ppu.tileData <<= 4
			switch ppu.Cycle % 8 {
			case 1:
				ppu.fetchNameTableByte()
			case 3:
				ppu.fetchAttributeTableByte()
			case 5:
				ppu.fetchLowTileByte()
			case 7:
				ppu.fetchHighTileByte()
			case 0:
				ppu.storeTileData()
			}
		}
	}
	//sprite logic, forback logic
	if ppu.isRenderingEnabled() && ppu.Cycle == 257 {
		if visibleLine { //if scanline visible
			ppu.evaluateSprites()
		} else {
			ppu.spriteCount = 0
		}
	}
	// vblank logic
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.setVerticalBlank()
	}
	if ppu.ScanLine == 261 && ppu.Cycle == 1 {
		ppu.clearVerticalBlank()
		ppu.flagSpriteZeroHit, ppu.flagSpriteOverflow = false, false
	}
}
