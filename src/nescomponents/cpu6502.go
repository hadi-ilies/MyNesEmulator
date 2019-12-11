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

//instruction of addr modes
type addrModes func()

func abs() {

}

func absX() {

}

func absY() {

}

func accumulator() {

}

func immediate() {

}

func implied() {

}

func indexedIndirect() {

}

func indirect() {

}

func indirectIndexed() {

}

func relative() {

}

func zeroPage() {

}

func zeroPageX() {

}

func zeroPageY() {

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
}

//this followed matrix shows the 210 op code (illegal/NOP are not counted) associated with the R65C00 family CPU devices.
var opCodeMatrix = [256]opCode{
	opCode{instructionName: "BRK", instructionMode: modeImplied, instructionSize: 2, nbCycle: 7}, opCode{instructionName: "ORA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "SLO", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "PHP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "ASL", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "ANC", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 0}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 3, nbCycle: 6},
	opCode{instructionName: "BPL", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "ORA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "SLO", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "CLC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "NOP", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "SLO", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 7}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 7}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 7},
	opCode{instructionName: "JSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "AND", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "RLA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "BIT", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "AND", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "ROL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "RLA", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "PLP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4}, opCode{instructionName: "AND", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "ANC", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "BIT", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "AND", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ROL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "RLA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6},
	opCode{instructionName: "BMI", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "AND", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "RLA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "AND", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "ROL", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "RLA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "SEC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "AND", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "AND", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ROL", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7},
	opCode{instructionName: "RTI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6}, opCode{instructionName: "EOR", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "SRE", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "PHA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "LSR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "ALR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "JMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 3}, opCode{instructionName: "EOR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "SRE", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6},
	opCode{instructionName: "BVC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "EOR", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "SRE", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "CLI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "SRE", instructionMode: modeAccumulator, instructionSize: 0, nbCycle: 7}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "EOR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LSR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 7},
	opCode{instructionName: "RTS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6}, opCode{instructionName: "ADC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "RRA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "ADC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "ROR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "RRA", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "PLA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4}, opCode{instructionName: "ADC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "ROR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "ARR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "JMP", instructionMode: modeIndirect, instructionSize: 3, nbCycle: 5}, opCode{instructionName: "ADC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ROR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "RRA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6},
	opCode{instructionName: "BVS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "ADC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "RRA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "ADC", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "ROR", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "RRA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "SEI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "ROR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7},
	opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "STA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "SAX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "STX", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "SAX", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 3}, opCode{instructionName: "DEY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "TXA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "XAA", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "STY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "STA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "STX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "SAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4},
	opCode{instructionName: "BCC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "STA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "AHX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "STX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "SAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4}, opCode{instructionName: "TYA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "STA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 5}, opCode{instructionName: "TXS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "TAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "SHY", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "STA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 5}, opCode{instructionName: "SHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "AHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5},
	opCode{instructionName: "LDY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "LDA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "LDX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "LDY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "LDA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "LDX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "LAX", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 3}, opCode{instructionName: "TAY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "LDA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "TAX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "LAX", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "LDY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LDA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LDX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4},
	opCode{instructionName: "BCS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "LDA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "LDY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "LDA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "LDX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "LAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4}, opCode{instructionName: "CLV", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "TSX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "LAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4}, opCode{instructionName: "LDY", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LDX", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "LAX", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4},
	opCode{instructionName: "CPY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "CMP", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "DCP", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "CPY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "CMP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "DEC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "DCP", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "INY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "CMP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "DEX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "AXS", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "CPY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "CMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "DEC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "DCP", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6},
	opCode{instructionName: "BNE", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "CMP", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "DCP", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "CMP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "DEC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "DCP", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "CLD", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7},
	opCode{instructionName: "CPX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "SBC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "ISC", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "CPX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "SBC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3}, opCode{instructionName: "INC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "ISC", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5}, opCode{instructionName: "INX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "CPX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "SBC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "INC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6}, opCode{instructionName: "ISC", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6},
	opCode{instructionName: "BEQ", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2}, opCode{instructionName: "SBC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2}, opCode{instructionName: "ISC", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "SBC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4}, opCode{instructionName: "INC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6}, opCode{instructionName: "ISC", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6}, opCode{instructionName: "SED", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7},
}

//THE CPU

type CPU struct {
	//Memory                        // memory interface
	Cycles    uint64              // number of cycles
	PC        uint16              // program counter
	SP        byte                // stack pointer
	A         byte                // accumulator
	X         byte                // x register
	Y         byte                // y register
	C         byte                // carry flag
	Z         byte                // zero flag
	I         byte                // interrupt disable flag
	D         byte                // decimal mode flag
	B         byte                // break command flag
	U         byte                // unused flag
	V         byte                // overflow flag
	N         byte                // negative flag
	interrupt byte                // interrupt type to perform
	stall     int                 // number of cycles to stall
	modes     map[int32]addrModes // map of instruction for each addr modes
	//table     [256]func(*stepInfo)
}

func createModesTables() map[int32]addrModes {
	modes := map[int32]addrModes{
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

	cpu.modes = createModesTables()

	return &cpu
}
