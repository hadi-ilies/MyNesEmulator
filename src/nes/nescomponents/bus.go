package nescomponents

import (
	// "os"
	"log"
)

// Mirroring Modes

const (
	MirrorHorizontal = 0
	MirrorVertical   = 1
	MirrorSingle0    = 2
	MirrorSingle1    = 3
	MirrorFour       = 4
)

type BUS struct {
	cpu          *CPU
	cpuRam       [2048]byte //fake ram
	ppu          *PPU
	cartridge    *Cartridge
	mapper       *Mapper
	clockCounter uint //nb clock
}

//create bus
func NewBus(cartridge *Cartridge) *BUS {
	var bus BUS

	bus.cartridge = cartridge
	bus.mapper = &cartridge.Mapper
	bus.cpu = NewCpu(&bus)
	bus.ppu = NewPpu(&bus)
	//bus.clockCounter = 0
	return &bus
}

//BUS READ/WRITE
func (bus *BUS) CpuWrite(address uint16, data byte) {
	//if bus.cartridge.CpuWrite(address, data) { //check cartrige addr
	//	println("LOL i have made an optimization")
	// if address >= 0x0000 && address <= 0x1FFF { //8KB range
	// 	bus.cpuRam[address&0x0800] = data
	// } else if address >= 0x2000 && address <= 0x3FFF {
	// 	bus.ppu.CpuWrite(address&0x0007, data)
	// }
	switch {
	case address < 0x2000:
		bus.cpuRam[address%0x0800] = data
	case address < 0x4000:
		bus.ppu.CpuWrite(0x2000+address%8, data)
	case address < 0x4014:
		//mem.console.APU.writeRegister(address, data)
	case address == 0x4014:
		bus.ppu.CpuWrite(address, data)
	case address == 0x4015:
		//mem.console.APU.writeRegister(address, data)
	case address == 0x4016:
		//mem.console.Controller1.Write data)
		//mem.console.Controller2.Write data)
	case address == 0x4017:
		//mem.console.APU.writeRegister(address, data)
	case address < 0x6000:
		// TODO: I/O registers
	case address >= 0x6000:
		bus.cartridge.Mapper.Write(address, data)
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

func (bus *BUS) CpuRead(address uint16) byte {
	//var data byte = 0x00

	// if bus.cartridge.CpuRead(address, &data) { //check cartrige addr
	// 	println("LOL i have made another optimization")
	// } else if address >= 0x0000 && address <= 0x1FFF { //8KB range
	// 	data = bus.cpuRam[address&0x07FF]
	// } else if address >= 0x2000 && address <= 0x3FFF {
	// 	bus.ppu.CpuRead(address & 0x0007)
	// }
	switch {
	case address < 0x2000:
		return bus.cpuRam[address%0x0800]
	case address < 0x4000:
		return bus.ppu.CpuRead(0x2000 + address%8)
	case address == 0x4014:
		return bus.ppu.CpuRead(address)
	case address == 0x4015:
		//return mem.console.APU.readRegister(address)
	case address == 0x4016:
		//return mem.console.Controller1.Read()
	case address == 0x4017:
		//return mem.console.Controller2.Read()
	case address < 0x6000:
		// TODO: I/O registers
	case address >= 0x6000:
		return bus.cartridge.Mapper.Read(address) //todo check mapper
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0
}

//System interface
func (bus *BUS) Reset() {
	bus.cpu.reset() //reset cpu flags and clocks
	//bus.clockCounter = 0 // nb clock useless
}

func (bus *BUS) Clock() {

}

func (bus *BUS) InsertCartridge(cartridge *Cartridge) {
	bus.cartridge = cartridge
	bus.ppu.ConnectCartridge(cartridge)
}

/*getter*/

func (bus *BUS) GetCpu() *CPU {
	return bus.cpu
}

func (bus *BUS) GetPpu() *PPU {
	return bus.ppu
}

func (bus *BUS) GetCartridge() *Cartridge {
	return bus.cartridge
}
