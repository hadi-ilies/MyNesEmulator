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

var opCodeMatrix = [256]opCode{
	opCode{
		instructionName: "BRK",
		instructionMode: modeImplied,
		instructionSize: 2,
		nbCycle:         7,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeIndexedIndirect,
		instructionSize: 2,
		nbCycle:         6,
	}, opCode{
		instructionName: "LOL",
		instructionMode: modeImplied,
		instructionSize: 0,
		nbCycle:         2,
	}, opCode{
		instructionName: "LOL",
		instructionMode: modeIndexedIndirect,
		instructionSize: 0,
		nbCycle:         8,
	}, opCode{
		instructionName: "NOP", //new op code
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         3,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         3,
	}, opCode{
		instructionName: "ASL",
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         5,
	}, opCode{
		instructionName: "SLO",
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         5,
	}, opCode{
		instructionName: "PHP",
		instructionMode: modeImplied,
		instructionSize: 1,
		nbCycle:         3,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeImmediate,
		instructionSize: 2,
		nbCycle:         2,
	}, opCode{ // i have to finish this tomorrow
		instructionName: "ASL",
		instructionMode: modeAccumulator,
		instructionSize: 1,
		nbCycle:         2,
	}, opCode{
		instructionName: "ANC",
		instructionMode: modeImmediate,
		instructionSize: 3,
		nbCycle:         0,
	}, opCode{
		instructionName: "NOP",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         6,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "ASL",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         6,
	}, opCode{
		instructionName: "SLO",
		instructionMode: modeZeroPage,
		instructionSize: 3,
		nbCycle:         6,
	},

	//new line
	opCode{
		instructionName: "BPL",
		instructionMode: modeRelative,
		instructionSize: 2,
		nbCycle:         2,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeIndirectIndexed,
		instructionSize: 2,
		nbCycle:         5,
	}, opCode{
		instructionName: "KIL",
		instructionMode: modeImplied,
		instructionSize: 0,
		nbCycle:         2,
	}, opCode{
		instructionName: "SLO",
		instructionMode: modeIndirectIndexed,
		instructionSize: 0,
		nbCycle:         8,
	}, opCode{
		instructionName: "NOP", //new op code
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         4,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         4,
	}, opCode{
		instructionName: "ASL",
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         6,
	}, opCode{
		instructionName: "SLO",
		instructionMode: modeZeroPage,
		instructionSize: 0,
		nbCycle:         6,
	}, opCode{
		instructionName: "CLC",
		instructionMode: modeImplied,
		instructionSize: 1,
		nbCycle:         2,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeImmediate,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "NOP",
		instructionMode: modeAccumulator,
		instructionSize: 1,
		nbCycle:         2,
	}, opCode{
		instructionName: "SLO",
		instructionMode: modeImmediate,
		instructionSize: 0,
		nbCycle:         7,
	}, opCode{
		instructionName: "NOP",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "ASL",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         7,
	}, opCode{
		instructionName: "SLO",
		instructionMode: modeZeroPage,
		instructionSize: 0,
		nbCycle:         7,
	},

	//new line
	opCode{
		instructionName: "JSR",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         6,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeIndexedIndirect,
		instructionSize: 2,
		nbCycle:         6,
	}, opCode{
		instructionName: "KIL",
		instructionMode: modeImplied,
		instructionSize: 0,
		nbCycle:         2,
	}, opCode{
		instructionName: "RLA",
		instructionMode: modeIndexedIndirect,
		instructionSize: 0,
		nbCycle:         8,
	}, opCode{
		instructionName: "BIT", //new op code
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         3,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         3,
	}, opCode{
		instructionName: "ROL",
		instructionMode: modeZeroPage,
		instructionSize: 2,
		nbCycle:         5,
	}, opCode{
		instructionName: "RLA",
		instructionMode: modeZeroPage,
		instructionSize: 0,
		nbCycle:         5,
	}, opCode{
		instructionName: "PLP",
		instructionMode: modeImplied,
		instructionSize: 1,
		nbCycle:         4,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeImmediate,
		instructionSize: 2,
		nbCycle:         2,
	}, opCode{
		instructionName: "ANC",
		instructionMode: modeAccumulator,
		instructionSize: 1,
		nbCycle:         2,
	}, opCode{
		instructionName: "BIT",
		instructionMode: modeImmediate,
		instructionSize: 0,
		nbCycle:         2,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "ORA",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "ROL",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         6,
	}, opCode{
		instructionName: "RLA",
		instructionMode: modeAbsolute,
		instructionSize: 0,
		nbCycle:         6,
	},

	//new line
	opCode{
		instructionName: "BMI",
		instructionMode: modeRelative,
		instructionSize: 2,
		nbCycle:         2,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeIndirectIndexed,
		instructionSize: 2,
		nbCycle:         5,
	}, opCode{
		instructionName: "KIL",
		instructionMode: modeImplied,
		instructionSize: 0,
		nbCycle:         2,
	}, opCode{
		instructionName: "RLA",
		instructionMode: modeIndirectIndexed,
		instructionSize: 0,
		nbCycle:         8,
	}, opCode{
		instructionName: "NOP", //new op code
		instructionMode: modeZeroPageY,
		instructionSize: 2,
		nbCycle:         4,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeZeroPageY,
		instructionSize: 2,
		nbCycle:         4,
	}, opCode{
		instructionName: "ROL",
		instructionMode: modeZeroPageY,
		instructionSize: 2,
		nbCycle:         6,
	}, opCode{
		instructionName: "RLA",
		instructionMode: modeZeroPageY,
		instructionSize: 0,
		nbCycle:         6,
	}, opCode{
		instructionName: "SEC",
		instructionMode: modeImplied,
		instructionSize: 1,
		nbCycle:         2,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeAbsoluteY,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "NOP",
		instructionMode: modeImplied,
		instructionSize: 1,
		nbCycle:         2,
	}, opCode{
		instructionName: "RLA",
		instructionMode: modeAbsoluteY,
		instructionSize: 0,
		nbCycle:         7,
	}, opCode{
		instructionName: "NOP",
		instructionMode: modeAbsoluteX,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "AND",
		instructionMode: modeAbsoluteX,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "ROL",
		instructionMode: modeAbsoluteX,
		instructionSize: 3,
		nbCycle:         7,
	}, opCode{
		instructionName: "RLA",
		instructionMode: modeAbsoluteX,
		instructionSize: 0,
		nbCycle:         7,
	},

	//new line
	opCode{
		instructionName: "RTI",
		instructionMode: modeImplied,
		instructionSize: 1,
		nbCycle:         6,
	}, opCode{
		instructionName: "EOR",
		instructionMode: modeIndexedIndirect,
		instructionSize: 2,
		nbCycle:         6,
	}, opCode{
		instructionName: "KIL",
		instructionMode: modeImplied,
		instructionSize: 0,
		nbCycle:         2,
	}, opCode{
		instructionName: "SRE",
		instructionMode: modeIndexedIndirect,
		instructionSize: 0,
		nbCycle:         8,
	}, opCode{
		instructionName: "NOP", //new op code
		instructionMode: modeZeroPageX,
		instructionSize: 2,
		nbCycle:         3,
	}, opCode{
		instructionName: "EOR",
		instructionMode: modeZeroPageX,
		instructionSize: 2,
		nbCycle:         3,
	}, opCode{
		instructionName: "LSR",
		instructionMode: modeZeroPageX,
		instructionSize: 2,
		nbCycle:         5,
	}, opCode{
		instructionName: "SRE",
		instructionMode: modeZeroPageX,
		instructionSize: 0,
		nbCycle:         5,
	}, opCode{
		instructionName: "PHA",
		instructionMode: modeImplied,
		instructionSize: 1,
		nbCycle:         3,
	}, opCode{
		instructionName: "EOR",
		instructionMode: modeImmediate,
		instructionSize: 2,
		nbCycle:         2,
	}, opCode{ // i have to finish this tomorrow
		instructionName: "LSR",
		instructionMode: modeAccumulator,
		instructionSize: 1,
		nbCycle:         2,
	}, opCode{
		instructionName: "ALR",
		instructionMode: modeImmediate,
		instructionSize: 0,
		nbCycle:         2,
	}, opCode{
		instructionName: "JMP",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         3,
	}, opCode{
		instructionName: "EOR",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         4,
	}, opCode{
		instructionName: "LSR",
		instructionMode: modeAbsolute,
		instructionSize: 3,
		nbCycle:         6,
	}, opCode{
		instructionName: "SRE",
		instructionMode: modeAbsolute,
		instructionSize: 0,
		nbCycle:         6,
	},

	//new line
	//name
	//mode
	//size
	//cycles

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

func CreateCpu() *CPU {
	var cpu CPU = CPU{}

	return &cpu
}
