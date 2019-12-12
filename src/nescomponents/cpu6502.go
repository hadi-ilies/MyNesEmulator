package nescomponents

//CPU_6502 Flags
const (
	C = 1 << 0 // Carry Bit, Decimal value 1
	Z = 1 << 1 // Zero Decimal, Decimal value 2
	I = 1 << 2 // Disable Interrupts , Decimal value 4
	D = 1 << 3 // Decimal Mode , Decimal value 8
	B = 1 << 4 // Break , Decimal value 16
	U = 1 << 5 // Unused, Decimal value 32
	V = 1 << 6 // Overflow , Decimal value 64
	N = 1 << 7 // Negative , Decimal value 128
)

// Read16 reads two bytes using Read to return a double-word value
func (cpu *CPU) Read16(address uint16) uint16 {
	var lo uint16 = uint16(cpu.bus.Read(address))
	var hi uint16 = uint16(cpu.bus.Read(address + 1))

	return hi<<8 | lo
}

//function that corespond to the execution of an instruction
type execinstructions func(address uint16, pc uint16, isAnAccumulator bool)

//addr modes of each instruction
type addrModes func(cpu *CPU) uint16

func abs(cpu *CPU) uint16 {
	var address uint16 = cpu.Read16(cpu.PC + 1)

	return address
}

func absX(cpu *CPU) uint16 {
	var address uint16 = cpu.Read16(cpu.PC+1) + uint16(cpu.X)
	//pageCrossed = pagesDiffer(address-uint16(cpu.X), address)

	return address
}

func absY(cpu *CPU) uint16 {
	var address uint16 = cpu.Read16(cpu.PC+1) + uint16(cpu.Y)
	// pageCrossed = pagesDiffer(address-uint16(cpu.Y), address)

	return address
}

func accumulator(cpu *CPU) uint16 {
	var address uint16 = 0
	return address
}

func immediate(cpu *CPU) uint16 {
	var address uint16 = cpu.PC + 1
	return address
}

func implied(cpu *CPU) uint16 {
	//address = 0
	return 0
}

func indexedIndirect(cpu *CPU) uint16 {
	//address = cpu.read16bug(uint16(cpu.Read(cpu.PC+1) + cpu.X))
	return 0
}

func indirect(cpu *CPU) uint16 {
	//address = cpu.read16bug(cpu.Read16(cpu.PC + 1))
	return 0
}

func indirectIndexed(cpu *CPU) uint16 {
	// address = cpu.read16bug(uint16(cpu.Read(cpu.PC+1))) + uint16(cpu.Y)
	// pageCrossed = pagesDiffer(address-uint16(cpu.Y), address)
	return 0
}

func relative(cpu *CPU) uint16 {
	offset := uint16(cpu.bus.Read(cpu.PC + 1))
	var address uint16 = cpu.PC + 2 + offset - 0x100

	if offset < 0x80 {
		address = cpu.PC + 2 + offset
	}
	return address
}

func zeroPage(cpu *CPU) uint16 {
	var address uint16 = uint16(cpu.bus.Read(cpu.PC + 1))

	return address
}

func zeroPageX(cpu *CPU) uint16 {
	var address uint16 = uint16(cpu.bus.Read(cpu.PC+1)+cpu.X) & 0xff

	return address
}

func zeroPageY(cpu *CPU) uint16 {
	var address uint16 = uint16(cpu.bus.Read(cpu.PC+1)+cpu.Y) & 0xff
	return address
}

// instructions functions
func brk(address uint16, pc uint16, isAnAccumulator bool) {

}

func lol(address uint16, pc uint16, isAnAccumulator bool) {

}

// addressing modes
const (
	modeAbsolute = iota + 1 //iota allow variable to work like enum in c "auto incrementation"
	modeAbsoluteX
	modeAbsoluteY
	modeAccumulator
	modeImmediate
	modeImplied
	modeIndexedIndirect
	modeIndirect
	modeIndirectIndexed
	modeRelative
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
)

//OPCODE MATRIX look doc page 11
type opCode struct {
	instructionName string
	instructionMode byte
	instructionSize uint16
	nbCycle         uint16
	instructionExec execinstructions
}

//this followed matrix shows the 210 op code (illegal/NOP are not counted) associated with the R65C00 family CPU devices.
//map of instruction
var opCodeMatrix = [256]opCode{
	opCode{instructionName: "BRK", instructionMode: modeImplied, instructionSize: 2, nbCycle: 7, instructionExec: brk}, opCode{instructionName: "ORA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SLO", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "PHP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ASL", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ANC", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 0, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 3, nbCycle: 6, instructionExec: lol},
	opCode{instructionName: "BPL", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SLO", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "CLC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SLO", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 7, instructionExec: lol},
	opCode{instructionName: "JSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "RLA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "BIT", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "ROL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "RLA", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "PLP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ANC", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "BIT", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ROL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "RLA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: lol},
	opCode{instructionName: "BMI", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "RLA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ROL", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "RLA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "SEC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "AND", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ROL", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: lol},
	opCode{instructionName: "RTI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SRE", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "PHA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LSR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ALR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "JMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "SRE", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: lol},
	opCode{instructionName: "BVC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SRE", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "CLI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SRE", instructionMode: modeAccumulator, instructionSize: 0, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "EOR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LSR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 7, instructionExec: lol},
	opCode{instructionName: "RTS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "RRA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "ROR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "RRA", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "PLA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ROR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ARR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "JMP", instructionMode: modeIndirect, instructionSize: 3, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ROR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "RRA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: lol},
	opCode{instructionName: "BVS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "RRA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ROR", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "RRA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "SEI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "ROR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: lol},
	opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "STA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SAX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "STX", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "SAX", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "DEY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "TXA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "XAA", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "STY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "STA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "STX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "SAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, instructionExec: lol},
	opCode{instructionName: "BCC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "STA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "AHX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "STX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "SAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "TYA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "STA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "TXS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "TAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "SHY", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "STA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "SHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "AHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, instructionExec: lol},
	opCode{instructionName: "LDY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "LDX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "LDY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "LDX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "LAX", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "TAY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "TAX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LAX", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LDY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LDX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, instructionExec: lol},
	opCode{instructionName: "BCS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "LDY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LDX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "CLV", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "TSX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "LAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LDY", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LDX", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "LAX", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, instructionExec: lol},
	opCode{instructionName: "CPY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "CPY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "DEC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "INY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "DEX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "AXS", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "CPY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "DEC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: lol},
	opCode{instructionName: "BNE", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "DEC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "CLD", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: lol},
	opCode{instructionName: "CPX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SBC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ISC", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "CPX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "SBC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lol}, opCode{instructionName: "INC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "ISC", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "INX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "CPX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "SBC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "INC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "ISC", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: lol},
	opCode{instructionName: "BEQ", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "SBC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lol}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "ISC", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "SBC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "INC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "ISC", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, instructionExec: lol}, opCode{instructionName: "SED", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lol}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: lol}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: lol},
}

//THE CPU

type CPU struct {
	//Memory                      // memory interface
	Cycles     uint64             // number of cycles
	PC         uint16             // program counter
	SP         byte               // stack pointer
	A          byte               // accumulator
	X          byte               // x register
	Y          byte               // y register
	C          byte               // carry flag
	Z          byte               // zero flag
	I          byte               // interrupt disable flag
	D          byte               // decimal mode flag
	B          byte               // break command flag
	U          byte               // unused flag
	V          byte               // overflow flag
	N          byte               // negative flag
	interrupt  byte               // interrupt type to perform
	stall      int                // number of cycles to stall
	modesTable map[byte]addrModes // address for each modes
	bus        BUS                // Linkage to the communications bus
}

func createModesTables() map[byte]addrModes {
	modes := map[byte]addrModes{
		modeAbsolute:        abs,
		modeAbsoluteX:       absX,
		modeAbsoluteY:       absY,
		modeAccumulator:     accumulator,
		modeImmediate:       immediate,
		modeImplied:         implied,
		modeIndexedIndirect: indexedIndirect,
		modeIndirect:        indirect,
		modeIndirectIndexed: indirectIndexed,
		modeRelative:        relative,
		modeZeroPage:        zeroPage,
		modeZeroPageX:       zeroPageX,
		modeZeroPageY:       zeroPageY,
	}
	return modes
}

func CreateCpu() *CPU {
	cpu := CPU{}

	cpu.modesTable = createModesTables()

	return &cpu
}

// Step executes a single CPU instruction
func (cpu *CPU) Step() uint64 {
	var startNbCycles uint64 = cpu.Cycles
	var opCodeIndex byte = cpu.bus.Read(cpu.PC)
	var op opCode = opCodeMatrix[opCodeIndex]

	cpu.PC += uint16(op.instructionSize)
	cpu.Cycles += uint64(op.nbCycle)

	// if pageCrossed {
	// 	cpu.Cycles += uint64(instructionPageCycles[opcode])
	// }

	//var address uint16 = cpu.modesTable[op.instructionMode](cpu) //return addr mode
	return cpu.Cycles - startNbCycles
}
