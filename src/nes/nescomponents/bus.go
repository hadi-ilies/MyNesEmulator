package nescomponents

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
	bus.mapper = &cartridge.mapper
	bus.cpu = NewCpu(&bus)
	bus.ppu = NewPpu(bus.cpu, bus.cartridge)
	bus.clockCounter = 0
	return &bus
}

//BUS READ/WRITE
func (bus *BUS) CpuWrite(address uint16, data byte) {
	if bus.cartridge.CpuWrite(address, data) { //check cartrige addr
		println("LOL i have made an optimization")
	} else if address >= 0x0000 && address <= 0x1FFF { //8KB range
		bus.cpuRam[address&0x07FF] = data
	} else if address >= 0x2000 && address <= 0x3FFF {
		bus.ppu.CpuWrite(address&0x0007, data)
	}
}

func (bus *BUS) CpuRead(address uint16) byte {
	var data byte = 0x00

	if bus.cartridge.CpuRead(address, &data) { //check cartrige addr
		println("LOL i have made another optimization")
	} else if address >= 0x0000 && address <= 0x1FFF { //8KB range
		data = bus.cpuRam[address&0x07FF]
	} else if address >= 0x2000 && address <= 0x3FFF {
		bus.ppu.CpuRead(address & 0x0007)
	}

	return data
}

//System interface
func (bus *BUS) Reset() {
	bus.cpu.reset()      //reset cpu flags and clocks
	bus.clockCounter = 0 // nb clock
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
