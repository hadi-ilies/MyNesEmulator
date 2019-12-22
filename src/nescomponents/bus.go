package nescomponents

type BUS struct {
	cpu    *CPU
	cpuRam [2048]byte
	PPU    *PPU
}

func (bus *BUS) CpuWrite(address uint16, data byte) {
	//8KB range
	if address >= 0x0000 && address <= 0x1FFF {
		bus.cpuRam[address&0x07FF] = data
	}
}

func (bus *BUS) CpuRead(address uint16) byte {
	var data byte = 0x00

	//8KB range
	if address >= 0x0000 && address <= 0x1FFF {
		data = bus.cpuRam[address&0x07FF]
	}
	return data
}
