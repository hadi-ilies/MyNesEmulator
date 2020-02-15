package nescomponents

import (
	"log"
)

type Mapper1 struct {
	cartridge     *Cartridge
	shiftRegister byte
	control       byte
	prgMode       byte
	chrMode       byte
	prgBank       byte
	chrBank0      byte
	chrBank1      byte
	prgOffsets    [2]int
	chrOffsets    [2]int
}

func NewMapper1(cartridge *Cartridge) Mapper {
	mapper := Mapper1{}
	mapper.cartridge = cartridge
	mapper.shiftRegister = 0x10
	mapper.prgOffsets[1] = mapper.prgBankOffset(-1)
	return &mapper
}

func (mapper *Mapper1) Step() {
}

func (mapper *Mapper1) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		bank := address / 0x1000
		offset := address % 0x1000
		return mapper.cartridge.chr[mapper.chrOffsets[bank]+int(offset)]
	case address >= 0x8000:
		address = address - 0x8000
		bank := address / 0x4000
		offset := address % 0x4000

		return mapper.cartridge.prg[mapper.prgOffsets[bank]+int(offset)]
	case address >= 0x6000:
		return mapper.cartridge.sram[int(address)-0x6000]
	default:
		log.Fatalf("unhandled mapper1 read at address: 0x%04X", address)
	}
	return 0
}

func (mapper *Mapper1) Write(address uint16, value byte) bool {
	switch {
	case address < 0x2000:
		bank := address / 0x1000
		offset := address % 0x1000
		mapper.cartridge.chr[mapper.chrOffsets[bank]+int(offset)] = value
	case address >= 0x8000:
		mapper.loadRegister(address, value)
	case address >= 0x6000:
		mapper.cartridge.sram[int(address)-0x6000] = value
	default:
		log.Fatalf("unhandled mapper1 write at address: 0x%04X", address)
		return false
	}
	return true
}

func (mapper *Mapper1) loadRegister(address uint16, value byte) {
	if value&0x80 == 0x80 {
		mapper.shiftRegister = 0x10
		mapper.writeControl(mapper.control | 0x0C)
	} else {
		complete := mapper.shiftRegister&1 == 1
		mapper.shiftRegister >>= 1
		mapper.shiftRegister |= (value & 1) << 4
		if complete {
			mapper.writeRegister(address, mapper.shiftRegister)
			mapper.shiftRegister = 0x10
		}
	}
}

func (mapper *Mapper1) writeRegister(address uint16, value byte) {
	switch {
	case address <= 0x9FFF:
		mapper.writeControl(value)
	case address <= 0xBFFF:
		mapper.writeCHRBank0(value)
	case address <= 0xDFFF:
		mapper.writeCHRBank1(value)
	case address <= 0xFFFF:
		mapper.writePRGBank(value)
	}
}

// Control (internal, $8000-$9FFF)
func (mapper *Mapper1) writeControl(value byte) {
	mapper.control = value
	mapper.chrMode = (value >> 4) & 1
	mapper.prgMode = (value >> 2) & 3
	mirror := value & 3
	switch mirror {
	case 0:
		mapper.cartridge.mirror = MirrorSingle0
	case 1:
		mapper.cartridge.mirror = MirrorSingle1
	case 2:
		mapper.cartridge.mirror = MirrorVertical
	case 3:
		mapper.cartridge.mirror = MirrorHorizontal
	}
	mapper.updateOffsets()
}

// CHR bank 0 (internal, $A000-$BFFF)
func (mapper *Mapper1) writeCHRBank0(value byte) {
	mapper.chrBank0 = value
	mapper.updateOffsets()
}

// CHR bank 1 (internal, $C000-$DFFF)
func (mapper *Mapper1) writeCHRBank1(value byte) {
	mapper.chrBank1 = value
	mapper.updateOffsets()
}

// PRG bank (internal, $E000-$FFFF)
func (mapper *Mapper1) writePRGBank(value byte) {
	mapper.prgBank = value & 0x0F
	mapper.updateOffsets()
}

func (mapper *Mapper1) prgBankOffset(index int) int {
	if index >= 0x80 {
		index -= 0x100
	}
	index %= len(mapper.cartridge.prg) / 0x4000
	offset := index * 0x4000
	if offset < 0 {
		offset += len(mapper.cartridge.prg)
	}
	return offset
}

func (mapper *Mapper1) chrBankOffset(index int) int {
	if index >= 0x80 {
		index -= 0x100
	}
	index %= len(mapper.cartridge.chr) / 0x1000
	offset := index * 0x1000
	if offset < 0 {
		offset += len(mapper.cartridge.chr)
	}
	return offset
}

// PRG ROM bank mode (0, 1: switch 32 KB at $8000, ignoring low bit of bank number;
//                    2: fix first bank at $8000 and switch 16 KB bank at $C000;
//                    3: fix last bank at $C000 and switch 16 KB bank at $8000)
// CHR ROM bank mode (0: switch 8 KB at a time; 1: switch two separate 4 KB banks)
func (mapper *Mapper1) updateOffsets() {
	switch mapper.prgMode {
	case 0, 1:
		mapper.prgOffsets[0] = mapper.prgBankOffset(int(mapper.prgBank & 0xFE))
		mapper.prgOffsets[1] = mapper.prgBankOffset(int(mapper.prgBank | 0x01))
	case 2:
		mapper.prgOffsets[0] = 0
		mapper.prgOffsets[1] = mapper.prgBankOffset(int(mapper.prgBank))
	case 3:
		mapper.prgOffsets[0] = mapper.prgBankOffset(int(mapper.prgBank))
		mapper.prgOffsets[1] = mapper.prgBankOffset(-1)
	}
	switch mapper.chrMode {
	case 0:
		mapper.chrOffsets[0] = mapper.chrBankOffset(int(mapper.chrBank0 & 0xFE))
		mapper.chrOffsets[1] = mapper.chrBankOffset(int(mapper.chrBank0 | 0x01))
	case 1:
		mapper.chrOffsets[0] = mapper.chrBankOffset(int(mapper.chrBank0))
		mapper.chrOffsets[1] = mapper.chrBankOffset(int(mapper.chrBank1))
	}
}
