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

//___________________________________________________ instructions functions__________________________________________________________________

// break instruction
func brk(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func ora(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func kil(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func slo(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func nop(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func php(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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

func jsr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

//this instruction is simply an 'and' logic gate
func and(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A &= cpu.bus.Read(address)
	cpu.setZ(cpu.A)
	cpu.setN(cpu.A)
}

func rla(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func rol(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func plp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func anc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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

func sec(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func rti(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func eor(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func lsr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func pha(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func alr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

//branch if overflowe clear
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

func sre(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

// Instruction: Disable Interrupts / Clear Interrupt Flag
// Function:    I = 0
func cli(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.I = 0
}

func rts(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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
	var m byte = cpu.bus.Read(address)
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

func ror(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func pla(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func arr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func jmp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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

func rra(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func sei(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func sta(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func sax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func dey(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func txa(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func xaa(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func sty(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func stx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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

func ahx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func tya(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func txs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func tas(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func shy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func shx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func ldy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func lda(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func ldx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func lax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func tay(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func tax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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

func tsx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func las(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func cpy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func cmp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func dec(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func iny(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func dex(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func axs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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

func dcp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

// Instruction: Clear decimal Flag
// Function:    D = 0
func cld(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.D = 0
}

func cpx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func isc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func inc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func inx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

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
	var m byte = cpu.bus.Read(address)
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

func sed(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func asl(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

func bit(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {

}

//_________________________________________________________________________________________________________________________________

// setZ sets the zero flag if the argument is zero
func (cpu *CPU) setZ(value byte) {
	if value == 0 {
		cpu.Z = 1
	} else {
		cpu.Z = 0
	}
}

// setN sets the negative flag if the argument is negative (high bit is set)
func (cpu *CPU) setN(value byte) {
	if value&0x80 != 0 {
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
	instructionExec execinstructions
}

//this followed matrix shows the 210 op code (illegal/NOP are not counted) associated with the R65C00 family CPU devices.
//map of instruction
var opCodeMatrix = [256]opCode{
	opCode{instructionName: "BRK", instructionMode: modeImplied, instructionSize: 2, nbCycle: 7, instructionExec: brk}, opCode{instructionName: "ORA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: ora}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "SLO", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: slo}, opCode{instructionName: "PHP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, instructionExec: php}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: asl}, opCode{instructionName: "ANC", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 0, instructionExec: anc}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 3, nbCycle: 6, instructionExec: slo},
	opCode{instructionName: "BPL", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: bpl}, opCode{instructionName: "ORA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: ora}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "SLO", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 4, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 6, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 6, instructionExec: slo}, opCode{instructionName: "CLC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: clc}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4, instructionExec: ora}, opCode{instructionName: "NOP", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "SLO", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 7, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 7, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 7, instructionExec: slo},
	opCode{instructionName: "JSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: jsr}, opCode{instructionName: "AND", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: and}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "RLA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: rla}, opCode{instructionName: "BIT", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: bit}, opCode{instructionName: "AND", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, instructionExec: rla}, opCode{instructionName: "PLP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, instructionExec: plp}, opCode{instructionName: "AND", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: and}, opCode{instructionName: "ANC", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: anc}, opCode{instructionName: "BIT", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: bit}, opCode{instructionName: "AND", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: and}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: ora}, opCode{instructionName: "ROL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: rla},
	opCode{instructionName: "BMI", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: bmi}, opCode{instructionName: "AND", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: and}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "RLA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: rla}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "AND", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6, instructionExec: rla}, opCode{instructionName: "SEC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: sec}, opCode{instructionName: "AND", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: and}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: rla}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "AND", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: rla},
	opCode{instructionName: "RTI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, instructionExec: rti}, opCode{instructionName: "EOR", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: eor}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "SRE", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5, instructionExec: sre}, opCode{instructionName: "PHA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, instructionExec: pha}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: lsr}, opCode{instructionName: "ALR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: alr}, opCode{instructionName: "JMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 3, instructionExec: jmp}, opCode{instructionName: "EOR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: sre},
	opCode{instructionName: "BVC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: bvc}, opCode{instructionName: "EOR", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: eor}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "SRE", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, instructionExec: sre}, opCode{instructionName: "CLI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: cli}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 3, nbCycle: 4, instructionExec: eor}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "SRE", instructionMode: modeAccumulator, instructionSize: 0, nbCycle: 7, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 7, instructionExec: sre},
	opCode{instructionName: "RTS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, instructionExec: rts}, opCode{instructionName: "ADC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: adc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "RRA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 5, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 5, instructionExec: rra}, opCode{instructionName: "PLA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, instructionExec: pla}, opCode{instructionName: "ADC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, instructionExec: ror}, opCode{instructionName: "ARR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: arr}, opCode{instructionName: "JMP", instructionMode: modeIndirect, instructionSize: 3, nbCycle: 5, instructionExec: jmp}, opCode{instructionName: "ADC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: rra},
	opCode{instructionName: "BVS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: bvs}, opCode{instructionName: "ADC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: adc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "RRA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 6, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 6, instructionExec: rra}, opCode{instructionName: "SEI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: sei}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: adc}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: rra},
	opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "STA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: sta}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "SAX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6, instructionExec: sax}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 3, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 3, instructionExec: sax}, opCode{instructionName: "DEY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: dey}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "TXA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: txa}, opCode{instructionName: "XAA", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: xaa}, opCode{instructionName: "STY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, instructionExec: sax},
	opCode{instructionName: "BCC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: bcc}, opCode{instructionName: "STA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 6, instructionExec: sta}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "AHX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6, instructionExec: ahx}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, instructionExec: sax}, opCode{instructionName: "TYA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: tya}, opCode{instructionName: "STA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 5, instructionExec: sta}, opCode{instructionName: "TXS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: txs}, opCode{instructionName: "TAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 5, instructionExec: tas}, opCode{instructionName: "SHY", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, instructionExec: shy}, opCode{instructionName: "STA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 5, instructionExec: sta}, opCode{instructionName: "SHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, instructionExec: shx}, opCode{instructionName: "AHX", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, instructionExec: ahx},
	opCode{instructionName: "LDY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 6, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 3, instructionExec: lax}, opCode{instructionName: "TAY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: tay}, opCode{instructionName: "LDA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: lda}, opCode{instructionName: "TAX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: tax}, opCode{instructionName: "LAX", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, instructionExec: lax},
	opCode{instructionName: "BCS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: bcs}, opCode{instructionName: "LDA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: lda}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 5, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, instructionExec: lax}, opCode{instructionName: "CLV", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: clv}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: lda}, opCode{instructionName: "TSX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: tsx}, opCode{instructionName: "LAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, instructionExec: las}, opCode{instructionName: "LDY", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, instructionExec: lax},
	opCode{instructionName: "CPY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: cmp}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "DCP", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: dcp}, opCode{instructionName: "CPY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, instructionExec: dcp}, opCode{instructionName: "INY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: iny}, opCode{instructionName: "CMP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: cmp}, opCode{instructionName: "DEX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: dex}, opCode{instructionName: "AXS", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: axs}, opCode{instructionName: "CPY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: dcp},
	opCode{instructionName: "BNE", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: bne}, opCode{instructionName: "CMP", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: cmp}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "DCP", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "CMP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, instructionExec: dcp}, opCode{instructionName: "CLD", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: cld}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: cmp}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: dcp},
	opCode{instructionName: "CPX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, instructionExec: sbc}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "ISC", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, instructionExec: isc}, opCode{instructionName: "CPX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, instructionExec: isc}, opCode{instructionName: "INX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: inx}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, instructionExec: sbc}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, instructionExec: sbc}, opCode{instructionName: "CPX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, instructionExec: isc},
	opCode{instructionName: "BEQ", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, instructionExec: beq}, opCode{instructionName: "SBC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, instructionExec: sbc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, instructionExec: kil}, opCode{instructionName: "ISC", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, instructionExec: isc}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, instructionExec: isc}, opCode{instructionName: "SED", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: sed}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, instructionExec: nop}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: nop}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, instructionExec: dcp},
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

//init and create nes CPU
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
	var isAnAccumulator bool = false
	var address uint16

	cpu.PC += uint16(op.instructionSize)
	cpu.Cycles += uint64(op.nbCycle)

	// if pageCrossed {
	// 	cpu.Cycles += uint64(instructionPageCycles[opcode])
	// }

	if op.instructionMode == modeAccumulator {
		isAnAccumulator = true
	}
	address = cpu.modesTable[op.instructionMode](cpu) //return addr mode
	op.instructionExec(cpu, address, op.instructionSize, isAnAccumulator)
	return cpu.Cycles - startNbCycles
}
