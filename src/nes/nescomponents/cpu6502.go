package nescomponents

// interrupt types
const (
	_ = iota
	interruptNone
	interruptNMI
	interruptIRQ
)

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

// pagesDiffer returns true if the two addresses reference different pages
func pagesDiffer(a uint16, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}

// pull pops a byte from the stack
func (cpu *CPU) pull() byte {
	cpu.SP++
	return cpu.bus.CpuRead(0x100 | uint16(cpu.SP))
}

// pull16 pops two bytes from the stack
func (cpu *CPU) pull16() uint16 {
	var low uint16 = uint16(cpu.pull())
	var high uint16 = uint16(cpu.pull())

	return high<<8 | low
}

// Read16 reads two bytes using Read to return a double-word value
func (cpu *CPU) Read16(address uint16) uint16 {
	var low uint16 = uint16(cpu.bus.CpuRead(address))
	var high uint16 = uint16(cpu.bus.CpuRead(address + 1))

	return high<<8 | low
}

// read16bug emulates a 6502 bug that caused the low byte to wrap without
// incrementing the high byte
func (cpu *CPU) read16bug(address uint16) uint16 {
	a := address
	b := (a & 0xFF00) | uint16(byte(a)+1)
	lo := cpu.bus.CpuRead(a)
	hi := cpu.bus.CpuRead(b)
	return uint16(hi)<<8 | uint16(lo)
}

//function that corespond to the execution of an instruction
type execinstructions func(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool)

//_____________________________________________________________________________________________________________________

//addr modes of each instruction
type addrModes func(cpu *CPU) uint16

func abs(cpu *CPU) uint16 {
	var address uint16 = cpu.Read16(cpu.PC + 1)

	return address
}

func absX(cpu *CPU) uint16 {
	var address uint16 = cpu.Read16(cpu.PC+1) + uint16(cpu.X)

	return address
}

func absY(cpu *CPU) uint16 {
	var address uint16 = cpu.Read16(cpu.PC+1) + uint16(cpu.Y)

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
	return 0
}

func indexedIndirect(cpu *CPU) uint16 {
	var address uint16 = cpu.read16bug(uint16(cpu.bus.CpuRead(cpu.PC+1) + cpu.X))

	return address
}

func indirect(cpu *CPU) uint16 {
	var address uint16 = cpu.read16bug(cpu.Read16(cpu.PC + 1))

	return address
}

func indirectIndexed(cpu *CPU) uint16 {
	var address uint16 = cpu.read16bug(uint16(cpu.bus.CpuRead(cpu.PC+1))) + uint16(cpu.Y)

	return address
}

func relative(cpu *CPU) uint16 {
	offset := uint16(cpu.bus.CpuRead(cpu.PC + 1))
	var address uint16 = cpu.PC + 2 + offset - 0x100

	if offset < 0x80 {
		address = cpu.PC + 2 + offset
	}
	return address
}

func zeroPage(cpu *CPU) uint16 {
	var address uint16 = uint16(cpu.bus.CpuRead(cpu.PC + 1))

	return address
}

func zeroPageX(cpu *CPU) uint16 {
	var address uint16 = uint16(cpu.bus.CpuRead(cpu.PC+1)+cpu.X) & 0xff

	return address
}

func zeroPageY(cpu *CPU) uint16 {
	var address uint16 = uint16(cpu.bus.CpuRead(cpu.PC+1)+cpu.Y) & 0xff
	return address
}

//___________________________________________________ instructions functions__________________________________________________________________

// break instruction which means force interrupt
func brk(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	high := byte(cpu.PC >> 8)
	low := byte(cpu.PC & 0xFF)

	cpu.bus.CpuWrite(0x100|uint16(cpu.SP), high)
	cpu.SP--
	cpu.bus.CpuWrite(0x100|uint16(cpu.SP), low)
	cpu.SP--
	php(cpu, address, pc, isAnAccumulator)
	sei(cpu, address, pc, isAnAccumulator)
	cpu.PC = cpu.Read16(0xFFFE)
}

// ORA - Logical Inclusive OR on the accumulator
func ora(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A |= cpu.bus.CpuRead(address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// PHP - Push Processor Status
// Instruction: Push Status Register to Stack
// Function:    status -> stack
// Note:        Break flag is set to 1 before push
func php(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var currentFlags byte

	currentFlags |= cpu.C << 0
	currentFlags |= cpu.Z << 1
	currentFlags |= cpu.I << 2
	currentFlags |= cpu.D << 3
	currentFlags |= cpu.B << 4
	currentFlags |= cpu.U << 5
	currentFlags |= cpu.V << 6
	currentFlags |= cpu.N << 7

	// push all current cpu flags bytes into the stack
	// 0x100 + uint16(cpu.SP) works too
	cpu.bus.CpuWrite(0x100|uint16(cpu.SP), currentFlags|0x10)
	cpu.SP--
}

//branch if position
// Instruction: Branch if Positive
// Function:    if(N == 0) pc = address
func bpl(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.N == 0 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???

	}
}

//carry clear
// Instruction: Clear Carry Flag
// Function:    C = 0
func clc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.C = 0
}

// JSR - Jump to Subroutine
func jsr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	value := cpu.PC - 1
	// push two bytes onto the stack
	//extract the first byte
	high := byte(value >> 8)
	//extract the second
	low := byte(value & 0xFF)
	//write bytes on bus
	cpu.bus.CpuWrite(0x100|uint16(cpu.SP), high)
	cpu.SP--
	cpu.bus.CpuWrite(0x100|uint16(cpu.SP), low)
	cpu.SP--
	cpu.PC = address
}

//this instruction is simply an 'and' logic gate
func and(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A &= cpu.bus.CpuRead(address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// ROR - Rotate Right
func rol(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var c byte = cpu.C

	if isAnAccumulator {
		cpu.C = (cpu.A >> 7) & 1
		cpu.A = (cpu.A << 1) | c
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		var value byte = cpu.bus.CpuRead(address)

		cpu.C = (value >> 7) & 1
		value = (value << 1) | c
		cpu.bus.CpuWrite(address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}
}

// PLP - Pull Processor Status
func plp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var flags byte = cpu.pull()&0xEF | 0x20

	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1
	cpu.B = (flags >> 4) & 1
	cpu.U = (flags >> 5) & 1
	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1

}

//branch if minus
//Instruction: Branch if Negative
// Function:    if(N == 1) pc = address
func bmi(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.N == 1 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???

	}
}

// set carry
// SEC - Set Carry Flag
func sec(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.C = 1
}

// RTI - Return from Interrupt
func rti(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var flags byte = cpu.pull()&0xEF | 0x20

	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1
	cpu.B = (flags >> 4) & 1
	cpu.U = (flags >> 5) & 1
	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1
	cpu.PC = cpu.pull16()
}

// EOR - Exclusive OR
func eor(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.A ^ cpu.bus.CpuRead(address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// LSR - Logical Shift Right
func lsr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if isAnAccumulator {
		cpu.C = cpu.A & 1
		cpu.A >>= 1
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		var value byte = cpu.bus.CpuRead(address)

		cpu.C = value & 1
		value >>= 1
		cpu.bus.CpuWrite(address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}
}

// PHA - Push Accumulator
// Instruction: Push Accumulator to Stack
// Function:    A -> stack write function allow me to access to the BUS
func pha(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	// push A byte into the stack
	// 0x100 + uint16(cpu.SP) works too
	cpu.bus.CpuWrite(0x100|uint16(cpu.SP), cpu.A)
	cpu.SP--
}

//branch if overflow clear
// Instruction: Branch if Overflow Clear
// Function:    if(V == 0) pc = address
func bvc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.V == 0 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???

	}
}

// Instruction: Disable Interrupts / Clear Interrupt Flag
// Function:    I = 0
func cli(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.I = 0
}

// RTS - Return from Subroutine
func rts(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.PC = cpu.pull16() + 1

}

// Instruction: Add with Carry In
// Function:    A = A + M + C
// Flags Out:   C, V, N, Z
//
// Explanation:
// The purpose of this function is to add a value to the accumulator and a carry bit. If
// the result is > 255 there is an overflow setting the carry bit. Ths allows you to
// chain together ADC instructions to add numbers larger than 8-bits. This in itself is
// simple, however the 6502 supports the concepts of Negativity/Positivity and Signed Overflow.
//
// 10000100 = 128 + 4 = 132 in normal circumstances, we know this as unsigned and it allows
// us to represent numbers between 0 and 255 (given 8 bits). The 6502 can also interpret
// this word as something else if we assume those 8 bits represent the range -128 to +127,
// i.e. it has become signed.
//
// Since 132 > 127, it effectively wraps around, through -128, to -124. This wraparound is
// called overflow, and this is a useful to know as it indicates that the calculation has
// gone outside the permissable range, and therefore no longer makes numeric sense.
//
// Note the implementation of ADD is the same in binary, this is just about how the numbers
// are represented, so the word 10000100 can be both -124 and 132 depending upon the
// context the programming is using it in. We can prove this!
//
//  10000100 =  132  or  -124
// +00010001 = + 17      + 17
//  ========    ===       ===     See, both are valid additions, but our interpretation of
//  10010101 =  149  or  -107     the context changes the value, not the hardware!
//
// In principle under the -128 to 127 range:
// 10000000 = -128, 11111111 = -1, 00000000 = 0, 00000000 = +1, 01111111 = +127
// therefore negative numbers have the most significant set, positive numbers do not
//
// To assist us, the 6502 can set the overflow flag, if the result of the addition has
// wrapped around. V <- ~(A^M) & A^(A+M+C) :D lol, let's work out why!
//
// Let's suppose we have A = 30, M = 10 and C = 0
//          A = 30 = 00011110
//          M = 10 = 00001010+
//     RESULT = 40 = 00101000
//
// Here we have not gone out of range. The resulting significant bit has not changed.
// So let's make a truth table to understand when overflow has occurred. Here I take
// the MSB of each component, where R is RESULT.
//
// A  M  R | V | A^R | A^M |~(A^M) |
// 0  0  0 | 0 |  0  |  0  |   1   |
// 0  0  1 | 1 |  1  |  0  |   1   |
// 0  1  0 | 0 |  0  |  1  |   0   |
// 0  1  1 | 0 |  1  |  1  |   0   |  so V = ~(A^M) & (A^R)
// 1  0  0 | 0 |  1  |  1  |   0   |
// 1  0  1 | 0 |  0  |  1  |   0   |
// 1  1  0 | 1 |  1  |  0  |   1   |
// 1  1  1 | 0 |  0  |  0  |   1   |
//
// We can see how the above equation calculates V, based on A, M and R. V was chosen
// based on the following hypothesis:
//       Positive Number + Positive Number = Negative Result -> Overflow
//       Negative Number + Negative Number = Positive Result -> Overflow
//       Positive Number + Negative Number = Either Result -> Cannot Overflow
//       Positive Number + Positive Number = Positive Result -> OK! No Overflow
//       Negative Number + Negative Number = Negative Result -> OK! NO Overflow
func adc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var a byte = cpu.A
	var m byte = cpu.bus.CpuRead(address)
	var c byte = cpu.C

	cpu.A = a + m + c
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
	if cpu.A > 0xFF {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^m)&0x80 == 0 && (a^cpu.A)&0x80 == 1 { //!= 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}

// ROR - Rotate Right
func ror(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var c byte = cpu.C

	if isAnAccumulator {
		cpu.C = cpu.A & 1
		cpu.A = (cpu.A >> 1) | (c << 7)
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		var value byte = cpu.bus.CpuRead(address)

		cpu.C = value & 1
		value = (value >> 1) | (c << 7)
		cpu.bus.CpuWrite(address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}
}

// PLA - Pull Accumulator
func pla(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.pull()
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// Instruction: Jump To Location
// Function:    pc = address
func jmp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.PC = address
}

// Instruction: Branch if Overflow Set
// Function:    if(V == 1) pc = address
func bvs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.V == 1 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???
	}
}

// SEI - Set Interrupt Disable
func sei(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.I = 1
}

// STA - Store Accumulator
func sta(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.bus.CpuWrite(address, cpu.A)
}

//DEY - Decrement Y Register

func dey(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y--
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// TXA - Transfer X register to Accumulator
func txa(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.X
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// STY - Store Y Register
func sty(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.bus.CpuWrite(address, cpu.Y)
}

// STX - Store X Register
func stx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.bus.CpuWrite(address, cpu.X)
}

//this function branch if the carry is clear
func bcc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.C == 0 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???

	}
}

//transfer y register to Accumulator
func tya(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.Y
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)

}

//TXS transfer X register to stack pointer
func txs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.SP = cpu.X
}

//TAS transfer accumulator to stack pointer
func tas(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.SP = cpu.A
}

// LDY - Load Y Register
func ldy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y = cpu.bus.CpuRead(address)
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// LDA - Load Accumulator
func lda(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.bus.CpuRead(address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

// LDX - Load X Register
func ldx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X = cpu.bus.CpuRead(address)
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)

}

// TAY - Transfer Accumulator to Y register
func tay(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y = cpu.A
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// TAX - Transfer Accumulator to X register
func tax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X = cpu.A
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// Instruction: Branch if Carry Set
// Function:    if(C == 1) pc = address
func bcs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.C == 1 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???
	}

}

// Instruction: Clear Overflow Flag
// Function:    V = 0
func clv(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.V = 0
}

//TSX transfer stackpointer to X register
func tsx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X = cpu.SP
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// Instruction: Compare Y Register
// Function:    C <- Y >= M      Z <- (Y - M) == 0
// Flags Out:   N, C, Z
func cpy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var value byte = cpu.bus.CpuRead(address)

	cpu.setZ(cpu.Y - value)
	cpu.setN(cpu.Y - value)
	if cpu.Y >= value {
		cpu.C = 1
	} else {
		cpu.C = 0
	}

}

// Instruction: Compare Accumulator
// Function:    C <- A >= M      Z <- (A - M) == 0
// Flags Out:   N, C, Z
func cmp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var value byte = cpu.bus.CpuRead(address)

	cpu.setZ(cpu.A - value)
	cpu.setN(cpu.A - value)
	if cpu.A >= value {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

// DEC - Decrement Memory
func dec(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var decAddr byte = cpu.bus.CpuRead(address) - 1

	cpu.bus.CpuWrite(address, decAddr)
	cpu.setZ(decAddr)
	cpu.setN(decAddr)
}

// INY - Increment Y Register
func iny(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y++
	cpu.setZ(cpu.Y)
	cpu.setN(cpu.Y)
}

// DEX - Decrement X Register
func dex(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X--
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

//branch if not equal
// Instruction: Branch if Not Equal
// Function:    if(Z == 0) pc = address
func bne(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if Z == 0 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???
	}
}

// Instruction: Clear decimal Flag
// Function:    D = 0
func cld(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.D = 0
}

// Instruction: Compare X Register
// Function:    C <- X >= M      Z <- (X - M) == 0
// Flags Out:   N, C, Z
func cpx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var value byte = cpu.bus.CpuRead(address)

	cpu.setZ(cpu.X - value)
	cpu.setN(cpu.X - value)
	if cpu.X >= value {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

// INC - Increment Memory
func inc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var incAddr byte = cpu.bus.CpuRead(address) + 1

	cpu.bus.CpuWrite(address, incAddr)

	cpu.setZ(incAddr)
	cpu.setN(incAddr)
}

// INX - Increment X Register
func inx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X++
	cpu.setZ(cpu.X)
	cpu.setN(cpu.X)
}

// Instruction: Subtraction with Borrow In
// Function:    A = A - M - (1 - C)
// Flags Out:   C, V, N, Z
//
// Explanation:
// Given the explanation for ADC above, we can reorganise our data
// to use the same computation for addition, for subtraction by multiplying
// the data by -1, i.e. make it negative
//
// A = A - M - (1 - C)  ->  A = A + -1 * (M - (1 - C))  ->  A = A + (-M + 1 + C)
//
// To make a signed positive number negative, we can invert the bits and add 1
// (OK, I lied, a little bit of 1 and 2s complement :P)
//
//  5 = 00000101
// -5 = 11111010 + 00000001 = 11111011 (or 251 in our 0 to 255 range)
//
// The range is actually unimportant, because if I take the value 15, and add 251
// to it, given we wrap around at 256, the result is 10, so it has effectively
// subtracted 5, which was the original intention. (15 + 251) % 256 = 10
//
// Note that the equation above used (1-C), but this got converted to + 1 + C.
// This means we already have the +1, so all we need to do is invert the bits
// of M, the data(!) therfore we can simply add, exactly the same way we did
// before.
func sbc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var a byte = cpu.A
	var m byte = cpu.bus.CpuRead(address)
	var c byte = cpu.C

	cpu.A += (-m + 1 + c)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
	if cpu.A > 0xFF {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^m)&0x80 == 0 && (a^cpu.A)&0x80 == 1 { //!= 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}

}

//branch if equal
// Instruction: Branch if Equal
// Function:    if(Z == 1) pc = address
func beq(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.Z == 1 {
		cpu.PC = address
		// adds a cycle for taking a branch and adds another cycle
		// if the branch jumps to a new page
		cpu.Cycles++
		addrAbs := pc + address

		//if the two addresses reference different pages
		if (addrAbs & 0xFF00) != (pc & 0xFF00) {
			cpu.Cycles++
		}
		//pc = addrAbs //Todo should i save the new addr ???
	}
}

// SED - Set Decimal Flag

func sed(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.D = 1
}

// BIT - Bit Test
func bit(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	value := cpu.bus.CpuRead(address)
	cpu.V = (value >> 6) & 1
	cpu.setZ(value & cpu.A)
	cpu.setN(value)
}

// ASL - Arithmetic Shift Left
func asl(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if isAnAccumulator {
		cpu.C = (cpu.A >> 7) & 1
		cpu.A <<= 1
		cpu.setZ(cpu.A)
		cpu.setN(cpu.A)
	} else {
		var value byte = cpu.bus.CpuRead(address)

		cpu.C = (value >> 7) & 1
		value <<= 1
		cpu.bus.CpuWrite(address, value)
		cpu.setZ(value)
		cpu.setN(value)
	}

}

//________________________________________________________illegal opcodes below_____________________________________________________________________________________

func isc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func kil(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func slo(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func nop(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func dcp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func axs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func las(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func lax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func shy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func shx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func ahx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func xaa(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func sax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func rra(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}
func arr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}
func sre(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func alr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}
func anc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}
func rla(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

//_________________________________________________________________________________________________________________________________

// setZ sets the zero flag if the argument is zero
func (cpu *CPU) setZ(value byte) {
	if value == 0x00 {
		cpu.Z = 1
	} else {
		cpu.Z = 0
	}
}

// setN sets the negative flag if the argument is negative (high bit is set)
func (cpu *CPU) setN(value byte) {
	if value&0x80 != 0x00 {
		cpu.N = 1
	} else {
		cpu.N = 0
	}
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
	nbPageCycles    uint16 //// instructionPageCycles indicates the number of cycles used by each instruction when a page is crossed
	instructionExec execinstructions
}

//this followed matrix shows the 210 op code (illegal/NOP are not counted) associated with the R65C00 family CPU devices.
//map of instruction
var opCodeMatrix = [256]opCode{
	opCode{instructionName: "BRK", instructionMode: modeImplied, instructionSize: 2, nbCycle: 7, nbPageCycles: 0, instructionExec: brk}, opCode{instructionName: "ORA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SLO", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "PHP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, nbPageCycles: 0, instructionExec: php}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "ANC", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 0, nbPageCycles: 0, instructionExec: anc}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: slo},
	opCode{instructionName: "BPL", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bpl}, opCode{instructionName: "ORA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: ora}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SLO", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "CLC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: clc}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ora}, opCode{instructionName: "NOP", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SLO", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: slo},
	opCode{instructionName: "JSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: jsr}, opCode{instructionName: "AND", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RLA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "BIT", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: bit}, opCode{instructionName: "AND", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "PLP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, nbPageCycles: 0, instructionExec: plp}, opCode{instructionName: "AND", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ANC", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: anc}, opCode{instructionName: "BIT", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: bit}, opCode{instructionName: "AND", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ROL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rla},
	opCode{instructionName: "BMI", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bmi}, opCode{instructionName: "AND", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: and}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RLA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "AND", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "SEC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: sec}, opCode{instructionName: "AND", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: and}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "AND", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rla},
	opCode{instructionName: "RTI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, nbPageCycles: 0, instructionExec: rti}, opCode{instructionName: "EOR", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SRE", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "PHA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, nbPageCycles: 0, instructionExec: pha}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "ALR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: alr}, opCode{instructionName: "JMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 3, nbPageCycles: 0, instructionExec: jmp}, opCode{instructionName: "EOR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: sre},
	opCode{instructionName: "BVC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bvc}, opCode{instructionName: "EOR", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: eor}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SRE", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "CLI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: cli}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: eor}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SRE", instructionMode: modeAccumulator, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: sre},
	opCode{instructionName: "RTS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, nbPageCycles: 0, instructionExec: rts}, opCode{instructionName: "ADC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RRA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "PLA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, nbPageCycles: 0, instructionExec: pla}, opCode{instructionName: "ADC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "ARR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: arr}, opCode{instructionName: "JMP", instructionMode: modeIndirect, instructionSize: 3, nbCycle: 5, nbPageCycles: 0, instructionExec: jmp}, opCode{instructionName: "ADC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rra},
	opCode{instructionName: "BVS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bvs}, opCode{instructionName: "ADC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: adc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RRA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "SEI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: sei}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: adc}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rra},
	opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "STA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SAX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: sax}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 3, nbPageCycles: 0, instructionExec: sax}, opCode{instructionName: "DEY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: dey}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "TXA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: txa}, opCode{instructionName: "XAA", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: xaa}, opCode{instructionName: "STY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: sax},
	opCode{instructionName: "BCC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bcc}, opCode{instructionName: "STA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 6, nbPageCycles: 1, instructionExec: sta}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "AHX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: ahx}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: sax}, opCode{instructionName: "TYA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tya}, opCode{instructionName: "STA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 5, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "TXS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: txs}, opCode{instructionName: "TAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: tas}, opCode{instructionName: "SHY", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: shy}, opCode{instructionName: "STA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 5, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "SHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: shx}, opCode{instructionName: "AHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: ahx},
	opCode{instructionName: "LDY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 3, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "TAY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tay}, opCode{instructionName: "LDA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "TAX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tax}, opCode{instructionName: "LAX", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: lax},
	opCode{instructionName: "BCS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bcs}, opCode{instructionName: "LDA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: lda}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 5, nbPageCycles: 1, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "CLV", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: clv}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: lda}, opCode{instructionName: "TSX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tsx}, opCode{instructionName: "LAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, nbPageCycles: 1, instructionExec: las}, opCode{instructionName: "LDY", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, nbPageCycles: 1, instructionExec: lax},
	opCode{instructionName: "CPY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "DCP", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "CPY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "INY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: iny}, opCode{instructionName: "CMP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: dex}, opCode{instructionName: "AXS", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: axs}, opCode{instructionName: "CPY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: dcp},
	opCode{instructionName: "BNE", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bne}, opCode{instructionName: "CMP", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: cmp}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "DCP", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "CMP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "CLD", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: cld}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: cmp}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: dcp},
	opCode{instructionName: "CPX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ISC", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "CPX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "INX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: inx}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "CPX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: isc},
	opCode{instructionName: "BEQ", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: beq}, opCode{instructionName: "SBC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: sbc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "ISC", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "SED", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: sed}, opCode{instructionName: "SBC", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ISC", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: cmp}, opCode{instructionName: "INC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "ISC", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: dcp},
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
	bus        *BUS               // Linkage to the communications bus
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

// NMI Non-Maskable Interrupt
// IRQ IRQ interrupt
//interruptMode == interrupt nmi -> NMI
//interruptMode == interrupt irq -> IRQ
func (cpu *CPU) cpuInterruptions(interruptMode byte) {
	if cpu.I == 0 && interruptMode != interruptNone {
		value := cpu.PC
		// push two bytes onto the stack
		//extract the first byte
		high := byte(value >> 8)
		//extract the second
		low := byte(value & 0xFF)
		//write bytes on bus
		cpu.bus.CpuWrite(0x100|uint16(cpu.SP), high)
		cpu.SP--
		cpu.bus.CpuWrite(0x100|uint16(cpu.SP), low)
		cpu.SP--

		php(cpu, 0, 0, false) // get current flags
		//if NMI
		if interruptMode == interruptNMI {
			cpu.PC = cpu.Read16(0xFFFA)
		} else { // else IRq
			cpu.PC = cpu.Read16(0xFFFE)
		}
		cpu.I = 1
		cpu.Cycles += 7
	}
}

// triggerNMI causes a non-maskable interrupt to occur on the next cycle
func (cpu *CPU) triggerNmi() {
	cpu.interrupt = interruptNMI
}

// Reset resets the CPU to its initial powerup state
func (cpu *CPU) reset() {
	var flags byte = 0x24

	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1
	cpu.B = (flags >> 4) & 1
	cpu.U = (flags >> 5) & 1
	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1
	cpu.PC = cpu.Read16(0xFFFC)
	cpu.SP = 0xFD
}

//init and create nes CPU
func NewCpu(bus *BUS) *CPU {
	var cpu CPU = CPU{bus: bus}

	cpu.modesTable = createModesTables()
	cpu.reset()
	return &cpu
}

//check if page crossed
func (cpu *CPU) isPageCrossed(op opCode, address uint16) bool {
	var isPageCrossed bool = false

	if op.instructionMode == modeAbsoluteX {
		isPageCrossed = pagesDiffer(address-uint16(cpu.X), address)
	}
	if op.instructionMode == modeAbsoluteY {
		isPageCrossed = pagesDiffer(address-uint16(cpu.Y), address)
	}
	if op.instructionMode == modeIndirectIndexed {
		isPageCrossed = pagesDiffer(address-uint16(cpu.Y), address)
	}
	return isPageCrossed
}

// Step executes a single CPU instruction
func (cpu *CPU) Step() uint64 {
	if cpu.stall > 0 {
		cpu.stall--
		return 1
	}
	var startNbCycles uint64 = cpu.Cycles
	var opCodeIndex byte = cpu.bus.CpuRead(cpu.PC)
	var op opCode = opCodeMatrix[opCodeIndex]
	var isAnAccumulator bool = false
	var address uint16

	cpu.cpuInterruptions(cpu.interrupt)
	cpu.interrupt = interruptNone // ??
	cpu.PC += uint16(op.instructionSize)
	cpu.Cycles += uint64(op.nbCycle)

	if op.instructionMode == modeAccumulator {
		isAnAccumulator = true
	}
	//get address from bus
	address = cpu.modesTable[op.instructionMode](cpu) //return addr mode
	//check wheter all the data are on the same page
	if cpu.isPageCrossed(op, address) {
		cpu.Cycles += uint64(op.nbPageCycles)
	}
	//exec instruction
	op.instructionExec(cpu, address, op.instructionSize, isAnAccumulator)
	return cpu.Cycles - startNbCycles
}
