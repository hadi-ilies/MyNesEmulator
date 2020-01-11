package nescomponents

import (
	"image"
	"log"
)

//addr cpu instruc
// $2000: PPUCTRL
func (ppu *PPU) writeControl(value byte) {
	ppu.ppuCtrl[flagNameTable] = (value >> 0) & 3
	ppu.ppuCtrl[flagIncrement] = (value >> 2) & 1
	ppu.ppuCtrl[flagSpriteTable] = (value >> 3) & 1
	ppu.ppuCtrl[flagBackgroundTable] = (value >> 4) & 1
	ppu.ppuCtrl[flagSpriteSize] = (value >> 5) & 1
	ppu.ppuCtrl[flagMasterSlave] = (value >> 6) & 1
	ppu.nmiOutput = (value>>7)&1 == 1
	ppu.nmiChange()
	// t: ....BA.. ........ = d: ......BA
	ppu.t = (ppu.t & 0xF3FF) | ((uint16(value) & 0x03) << 10)
}

// $2001: PPUMASK
func (ppu *PPU) writeMask(value byte) {
	ppu.ppuMask[flagGrayscale] = (value >> 0) & 1
	ppu.ppuMask[flagShowLeftBackground] = (value >> 1) & 1
	ppu.ppuMask[flagShowLeftSprites] = (value >> 2) & 1
	ppu.ppuMask[flagShowBackground] = (value >> 3) & 1
	ppu.ppuMask[flagShowSprites] = (value >> 4) & 1
	ppu.ppuMask[flagRedTint] = (value >> 5) & 1
	ppu.ppuMask[flagGreenTint] = (value >> 6) & 1
	ppu.ppuMask[flagBlueTint] = (value >> 7) & 1
}

func (ppu *PPU) boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// $2002: PPUSTATUS
func (ppu *PPU) readStatus() byte {
	result := ppu.register & 0x1F
	result |= ppu.boolToByte(ppu.flagSpriteOverflow) << 5
	result |= ppu.boolToByte(ppu.flagSpriteZeroHit) << 6
	if ppu.nmiOccurred {
		result |= 1 << 7
	}
	ppu.nmiOccurred = false
	ppu.nmiChange()
	// w:                   = 0
	ppu.w = 0
	return result
}

// $2003: OAMADDR
func (ppu *PPU) writeOAMAddress(value byte) {
	ppu.oamAddress = value
}

// $2004: OAMDATA (read)
func (ppu *PPU) readOamData() byte {
	return ppu.oam[ppu.oamAddress]
}

// $2004: OAMDATA (write)
func (ppu *PPU) writeOAMData(value byte) {
	ppu.oam[ppu.oamAddress] = value
	ppu.oamAddress++
}

// $2005: PPUSCROLL
func (ppu *PPU) writeScroll(value byte) {
	if ppu.w == 0 {
		// t: ........ ...HGFED = d: HGFED...
		// x:               CBA = d: .....CBA
		// w:                   = 1
		ppu.t = (ppu.t & 0xFFE0) | (uint16(value) >> 3)
		ppu.x = value & 0x07
		ppu.w = 1
	} else {
		// t: .CBA..HG FED..... = d: HGFEDCBA
		// w:                   = 0
		ppu.t = (ppu.t & 0x8FFF) | ((uint16(value) & 0x07) << 12)
		ppu.t = (ppu.t & 0xFC1F) | ((uint16(value) & 0xF8) << 2)
		ppu.w = 0
	}
}

// $2006: PPUADDR
func (ppu *PPU) writeAddress(value byte) {
	if ppu.w == 0 {
		// t: ..FEDCBA ........ = d: ..FEDCBA
		// t: .X...... ........ = 0
		// w:                   = 1
		ppu.t = (ppu.t & 0x80FF) | ((uint16(value) & 0x3F) << 8)
		ppu.w = 1
	} else {
		// t: ........ HGFEDCBA = d: HGFEDCBA
		// v                    = t
		// w:                   = 0
		ppu.t = (ppu.t & 0xFF00) | uint16(value)
		ppu.v = ppu.t
		ppu.w = 0
	}
}

// $2007: PPUDATA (read)
func (ppu *PPU) readData() byte {
	value := ppu.Read(ppu.v)
	// emulate buffered reads
	if ppu.v%0x4000 < 0x3F00 {
		buffered := ppu.bufferedData
		ppu.bufferedData = value
		value = buffered
	} else {
		ppu.bufferedData = ppu.Read(ppu.v - 0x1000)
	}
	// increment address
	if ppu.ppuCtrl[flagIncrement] == 0 {
		ppu.v += 1
	} else {
		ppu.v += 32
	}
	return value
}

// $2007: PPUDATA (write)
func (ppu *PPU) writeData(value byte) {
	ppu.Write(ppu.v, value)
	if ppu.ppuCtrl[flagIncrement] == 0 {
		ppu.v += 1
	} else {
		ppu.v += 32
	}
}

// $4014: OAMDMA
func (ppu *PPU) writeDMA(value byte) {
	address := uint16(value) << 8
	for i := 0; i < 256; i++ {
		ppu.oam[ppu.oamAddress] = ppu.bus.CpuRead(address)
		ppu.oamAddress++
		address++
	}
	ppu.bus.cpu.stall += 513
	if ppu.bus.cpu.Cycles%2 == 1 {
		ppu.bus.cpu.stall++
	}
}

//cpu can read only 8 addrs form the ppu
const (
	control = 0x2000
	Mask    = 0x2001
	Status  = 0x2002
	//The OAM (Object Attribute Memory) is internal memory inside the PPU that contains a display
	//list of up to 64 sprites, where each sprite's information occupies 4 bytes.
	oamAddr = 0x2003
	oamData = 0x2004
	Scroll  = 0x2005
	ppuAddr = 0x2006
	ppuData = 0x2007
)

//Comunication with main BUS
func (ppu *PPU) CpuWrite(address uint16, data byte) {
	ppu.register = data
	switch address {
	case control:
		ppu.writeControl(data)
	case Mask:
		ppu.writeMask(data)
	case oamAddr:
		ppu.writeOAMAddress(data)
	case oamData:
		ppu.writeOAMData(data)
	case Scroll:
		ppu.writeScroll(data)
	case ppuAddr:
		ppu.writeAddress(data)
	case ppuData:
		ppu.writeData(data)
	case 0x4014:
		ppu.writeDMA(data)
	}
}

func (ppu *PPU) CpuRead(address uint16) byte {
	var data byte = 0x00

	switch address {
	case Status:
		data = ppu.readStatus()
	case oamData:
		data = ppu.readOamData()
	case ppuData:
		data = ppu.readData()
	}
	return data
}

//MIRROR ADDR
var MirrorLookup = [...][4]uint16{
	{0, 0, 1, 1},
	{0, 1, 0, 1},
	{0, 0, 0, 0},
	{1, 1, 1, 1},
	{0, 1, 2, 3},
}

func (ppu *PPU) mirrorAddress(mode byte, address uint16) uint16 {
	address = (address - 0x2000) % 0x1000
	table := address / 0x0400
	offset := address % 0x0400
	return 0x2000 + MirrorLookup[mode][table]*0x0400 + offset
}

//Comunication  with the second "PPU" BUS
func (ppu *PPU) Read(address uint16) byte {
	// var data byte = 0x00
	// address &= 0x3FFF

	// if ppu.cartridge.PpuRead(address, &data) {

	// }
	// return data
	var data byte = 0
	address = address % 0x4000
	switch {
	case address < 0x2000:
		return ppu.cartridge.mapper.Read(address)
	case address < 0x3F00:
		mode := ppu.cartridge.mirror
		data = ppu.nameTable[ppu.mirrorAddress(mode, address)%2048]
	case address < 0x4000:
		data = ppu.readPalette(address % 32)
	default:
		log.Fatalf("unhandled ppu memory read at address: 0x%04X", address)
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
	bus       *BUS       //pointer on bus to get nes's Cpu
	cartridge *Cartridge // the gamePAk
	// storage variables
	nameTable    [2048]byte  //[2][1024]byte
	paletteTable [32]byte    //ram connected to ppu that strored the palace info there are 32 entries
	oam          [256]byte   // (Object Attribute Memory)
	front        *image.RGBA // front ground that generate sprites
	back         *image.RGBA // back ground

	// PPU registers
	v        uint16 // current vram address (15 bit)
	t        uint16 // temporary vram address (15 bit)
	x        byte   // fine x scroll (3 bit)
	w        byte   // write toggle (1 bit)
	f        byte   // even/odd frame flag (1 bit)
	register byte
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

	// $2003 OAMADDR
	oamAddress byte

	// $2007 PPUDATA
	bufferedData byte // for buffered reads
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

func NewPpu(bus *BUS) *PPU {
	var ppu PPU

	ppu.bus = bus
	ppu.cartridge = bus.cartridge
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

// NTSC Timing Helper Functions

func (ppu *PPU) incrementX() {
	// increment hori(v)
	// if coarse X == 31
	if ppu.v&0x001F == 31 {
		// coarse X = 0
		ppu.v &= 0xFFE0
		// switch horizontal nametable
		ppu.v ^= 0x0400
	} else {
		// increment coarse X
		ppu.v++
	}
}

func (ppu *PPU) incrementY() {
	// increment vert(v)
	// if fine Y < 7
	if ppu.v&0x7000 != 0x7000 {
		// increment fine Y
		ppu.v += 0x1000
	} else {
		// fine Y = 0
		ppu.v &= 0x8FFF
		// let y = coarse Y
		y := (ppu.v & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			ppu.v ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		ppu.v = (ppu.v & 0xFC1F) | (y << 5)
	}
}

func (ppu *PPU) copyX() {
	// hori(v) = hori(t)
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
}

func (ppu *PPU) copyY() {
	// vert(v) = vert(t)
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
}

func (ppu *PPU) fetchTileData() uint32 {
	return uint32(ppu.tileData >> 32)
}

func (ppu *PPU) backgroundPixel() byte {
	if ppu.ppuMask[flagShowBackground] == 0 {
		return 0
	}
	data := ppu.fetchTileData() >> ((7 - ppu.x) * 4)
	return byte(data & 0x0F)
}

func (ppu *PPU) spritePixel() (byte, byte) {
	if ppu.ppuMask[flagShowSprites] == 0 {
		return 0, 0
	}
	for i := 0; i < ppu.spriteCount; i++ {
		offset := (ppu.Cycle - 1) - int(ppu.spritePositions[i])
		if offset < 0 || offset > 7 {
			continue
		}
		offset = 7 - offset
		color := byte((ppu.spritePatterns[i] >> byte(offset*4)) & 0x0F)
		if color%4 == 0 {
			continue
		}
		return byte(i), color
	}
	return 0, 0
}

func (ppu *PPU) readPalette(address uint16) byte {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	return ppu.paletteTable[address]
}

func (ppu *PPU) writePalette(address uint16, value byte) {
	if address >= 16 && address%4 == 0 {
		address -= 16
	}
	ppu.paletteTable[address] = value
}

func (ppu *PPU) renderPixel() {
	x := ppu.Cycle - 1
	y := ppu.ScanLine
	var background byte = ppu.backgroundPixel()
	i, sprite := ppu.spritePixel()
	if x < 8 && ppu.ppuMask[flagShowLeftBackground] == 0 {
		background = 0
	}
	if x < 8 && ppu.ppuMask[flagShowLeftSprites] == 0 {
		sprite = 0
	}
	b := background%4 != 0
	s := sprite%4 != 0
	var color byte
	if !b && !s {
		color = 0
	} else if !b && s {
		color = sprite | 0x10
	} else if b && !s {
		color = background
	} else {
		if ppu.spriteIndexes[i] == 0 && x < 255 {
			ppu.flagSpriteZeroHit = true
		}
		if ppu.spritePriorities[i] == 0 {
			color = sprite | 0x10
		} else {
			color = background
		}
	}
	c := Palette[ppu.readPalette(uint16(color))%64]
	ppu.back.SetRGBA(x, y, c)
}

// update updates Cycle, ScanLine and Frame counters
func (ppu *PPU) update() {
	if ppu.nmiDelay > 0 {
		ppu.nmiDelay--
		if ppu.nmiDelay == 0 && ppu.nmiOutput && ppu.nmiOccurred {
			ppu.bus.cpu.triggerNmi()
		}
	}

	if ppu.ppuMask[flagShowBackground] != 0 || ppu.ppuMask[flagShowSprites] != 0 {
		if ppu.f == 1 && ppu.ScanLine == 261 && ppu.Cycle == 339 {
			ppu.Cycle = 0
			ppu.ScanLine = 0
			ppu.Frame++
			ppu.f ^= 1
			return
		}
	}
	ppu.Cycle++
	if ppu.Cycle > 340 {
		ppu.Cycle = 0
		ppu.ScanLine++
		if ppu.ScanLine > 261 {
			ppu.ScanLine = 0
			ppu.Frame++
			ppu.f ^= 1
		}
	}
}

func (ppu *PPU) Step() {
	ppu.update()
	visibleLine := ppu.ScanLine < 240
	preLine := ppu.ScanLine == 261
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle
	renderLine := preLine || visibleLine

	//background
	if ppu.isRenderingEnabled() {
		if visibleLine && visibleCycle {
			ppu.renderPixel()
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
		if preLine && ppu.Cycle >= 280 && ppu.Cycle <= 304 {
			ppu.copyY()
		}
		if renderLine {
			if fetchCycle && ppu.Cycle%8 == 0 {
				ppu.incrementX()
			}
			if ppu.Cycle == 256 {
				ppu.incrementY()
			}
			if ppu.Cycle == 257 {
				ppu.copyX()
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
