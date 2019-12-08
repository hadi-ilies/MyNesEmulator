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
