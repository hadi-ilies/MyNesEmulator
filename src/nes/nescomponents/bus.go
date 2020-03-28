package nescomponents

import (
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
	Controller1  *Controller
	mapper       *Mapper
	clockCounter uint //nb clock
}

//NewBus bus
func NewBus(cartridge *Cartridge) *BUS {
	var bus BUS

	bus.cartridge = cartridge
	bus.mapper = &cartridge.Mapper
	bus.cpu = NewCpu(&bus)
	bus.ppu = NewPpu(&bus)
	bus.Controller1 = NewController()
	//bus.clockCounter = 0
	return &bus
}

//CpuWrite BUS handle the memory
func (bus *BUS) CpuWrite(address uint16, data byte) {
	if address >= 0 && address <= 0x1FFF { //8KB range
		bus.cpuRam[address%0x0800] = data
	} else if address > 0x1FFF && address < 0x4000 {
		bus.ppu.CpuWrite(0x2000+address%8, data)
	} else if address > 0x4000 && address < 0x4014 {
		//mem.console.APU.writeRegister(address, data)
	} else if address == 0x4014 {
		bus.ppu.CpuWrite(address, data)
	} else if address == 0x4015 {
		//mem.console.APU.writeRegister(address, data)
	} else if address == 0x4016 {
		bus.Controller1.Write(data)
		//mem.console.Controller2.Write data)
	} else if address == 0x4017 {
		//mem.console.APU.writeRegister(address, data)
	} else if address < 0x6000 {
		// TODO: I/O registers
	} else if address >= 0x6000 {
		bus.cartridge.Mapper.Write(address, data)
	} else {
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

//CpuRead BUS handle the memory
func (bus *BUS) CpuRead(address uint16) byte {
	var data byte = 0x00

	if address >= 0 && address <= 0x1FFF { //8KB range
		data = bus.cpuRam[address%0x0800]
	} else if address > 0x1FFF && address < 0x4000 {
		data = bus.ppu.CpuRead(0x2000 + address%8)
	} else if address == 0x4014 {
		data = bus.ppu.CpuRead(address)
	} else if address == 0x4015 {
		//data = mem.console.APU.readRegister(address)
	} else if address == 0x4016 {
		data = bus.Controller1.Read()
	} else if address == 0x4017 {
		//data = mem.console.Controller2.Read()
	} else if address < 0x6000 {
		// TODO: I/O registers
	} else if address >= 0x6000 {
		return bus.cartridge.Mapper.Read(address) //todo check mapper
	} else {
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return data
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
