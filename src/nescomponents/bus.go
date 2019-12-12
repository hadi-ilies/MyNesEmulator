package nescomponents

type BUS struct {
	cpu *CPU
	ram [2048]byte
}

func (bus *BUS) Write(address uint16, data byte) {
	if address >= 0x0000 && address <= 0xFFFF {
		bus.ram[address] = data
	}
}

func (bus *BUS) Read(address uint16) byte {
	if address >= 0x0000 && address <= 0xFFFF {
		return bus.ram[address]
	}
	return 0x00
}
