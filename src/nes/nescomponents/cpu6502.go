package nescomponents

import (
	"fmt"
	"os"
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

//prototype function that corespond to the execution of an instruction
type execinstructions func(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool)

//________________________________________________addr modes of each instruction______________________________________________________________

//prototype addr modes of each instruction
type addrModes func(cpu *CPU) uint16

func abs(cpu *CPU) uint16 {
	return cpu.Read16(cpu.PC + 1)
}

func absX(cpu *CPU) uint16 {
	return cpu.Read16(cpu.PC+1) + uint16(cpu.X)
}

func absY(cpu *CPU) uint16 {
	return cpu.Read16(cpu.PC+1) + uint16(cpu.Y)
}

func accumulator(cpu *CPU) uint16 {
	return 0
}

func immediate(cpu *CPU) uint16 {
	return cpu.PC + 1
}

func implied(cpu *CPU) uint16 {
	return 0
}

func indexedIndirect(cpu *CPU) uint16 {
	return cpu.read16bug(uint16(cpu.bus.CpuRead(cpu.PC+1) + cpu.X))
}

func indirect(cpu *CPU) uint16 {
	return cpu.read16bug(cpu.Read16(cpu.PC + 1))
}

func indirectIndexed(cpu *CPU) uint16 {
	return cpu.read16bug(uint16(cpu.bus.CpuRead(cpu.PC+1))) + uint16(cpu.Y)
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
	return uint16(cpu.bus.CpuRead(cpu.PC + 1))
}

func zeroPageX(cpu *CPU) uint16 {
	return  uint16(cpu.bus.CpuRead(cpu.PC+1)+cpu.X) & 0xff
}

func zeroPageY(cpu *CPU) uint16 {
	return uint16(cpu.bus.CpuRead(cpu.PC+1)+cpu.Y) & 0xff
}

//___________________________________________________ instructions functions__________________________________________________________________

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

	cpu.A += m + c
	cpu.setZN(cpu.A)

	//we write "int(a)+int(b)+int(c)" instead of cpu.A because cpu.A is a byte and CANT be over 255
	//if accumulator overflow
	if int(a)+int(m)+int(c) > 0xFF {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^m)&0x80 == 0 && (a^cpu.A)&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}

// setZN sets the zero flag and the negative flag
func (cpu *CPU) setZN(value byte) {
	cpu.setZ(value)
	cpu.setN(value)
}

// addBranchCycles adds a cycle for taking a branch and adds another cycle
// if the branch jumps to a new page
func (cpu *CPU) addBranchCycles(address uint16, pc uint16) {
	cpu.Cycles++
	if pagesDiffer(pc, address) {
		cpu.Cycles++
	}
}

// push pushes a byte onto the stack
func (cpu *CPU) push(value byte) {
	cpu.bus.CpuWrite(0x100|uint16(cpu.SP), value)
	cpu.SP--
}

// push16 pushes two bytes onto the stack
func (cpu *CPU) push16(value uint16) {
	//extract first bytes
	high := byte(value >> 8)
	//extract the seconde one
	low := byte(value & 0xFF)
	//push on the bus
	cpu.push(high)
	cpu.push(low)
}

func (cpu *CPU) compare(a, b byte) {
	cpu.setZN(a - b)
	if a >= b {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
}

// Flags returns the processor status flags
func (cpu *CPU) Flags() byte {
	var flags byte

	flags |= cpu.C << 0
	flags |= cpu.Z << 1
	flags |= cpu.I << 2
	flags |= cpu.D << 3
	flags |= cpu.B << 4
	flags |= cpu.U << 5
	flags |= cpu.V << 6
	flags |= cpu.N << 7
	return flags
}

// SetFlags sets the processor status flags
func (cpu *CPU) SetFlags(flags byte) {
	cpu.C = (flags >> 0) & 1
	cpu.Z = (flags >> 1) & 1
	cpu.I = (flags >> 2) & 1
	cpu.D = (flags >> 3) & 1
	cpu.B = (flags >> 4) & 1
	cpu.U = (flags >> 5) & 1
	cpu.V = (flags >> 6) & 1
	cpu.N = (flags >> 7) & 1
}

// AND - Logical AND
//this instruction is simply an 'and' logic gate
func and(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A &= cpu.bus.CpuRead(address)
	cpu.setZN(cpu.A)
}

// ASL - Arithmetic Shift Left
func asl(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if isAnAccumulator {
		cpu.C = (cpu.A >> 7) & 1
		cpu.A <<= 1
		cpu.setZN(cpu.A)
	} else {
		var value byte = cpu.bus.CpuRead(address)

		cpu.C = (value >> 7) & 1
		value <<= 1
		cpu.bus.CpuWrite(address, value)
		cpu.setZN(value)
	}
}

// BCC - Branch if Carry Clear
func bcc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.C == 0 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// BCS - Branch if Carry Set
// Function:    if(C == 1) pc = address
func bcs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.C == 1 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// BEQ - Branch if Equal
func beq(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.Z == 1 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// BIT - Bit Test
func bit(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var value byte = cpu.bus.CpuRead(address)

	cpu.V = (value >> 6) & 1
	cpu.setZ(value & cpu.A)
	cpu.setN(value)
}

// BMI - Branch if Minus
// Instruction: Branch if Negative
// Function:    if(N == 1) pc = address
func bmi(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.N == 1 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// BNE - Branch if Not Equal
// Instruction: Branch if Not Equal
// Function:    if(Z == 0) pc = address
func bne(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.Z == 0 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// BPL - Branch if Positive
// Instruction: Branch if Positive
// Function:    if(N == 0) pc = address
func bpl(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.N == 0 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// BRK - Force Interrupt
// break instruction which means force interrupt
func brk(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.push16(cpu.PC)
	php(cpu, address, pc, isAnAccumulator)
	sei(cpu, address, pc, isAnAccumulator)
	cpu.PC = cpu.Read16(0xFFFE)
}

// BVC - Branch if Overflow Clear
// Instruction: Branch if Overflow Clear
// Function:    if(V == 0) pc = address
func bvc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.V == 0 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// BVS - Branch if Overflow Set
// Function:    if(V == 1) pc = address
func bvs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if cpu.V == 1 {
		cpu.PC = address
		cpu.addBranchCycles(address, pc)
	}
}

// CLC - Clear Carry Flag
// Instruction: Clear Carry Flag
// Function:    C = 0
func clc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.C = 0
}

// CLD - Clear Decimal Mode
func cld(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.D = 0
}

// CLI - Clear Interrupt Disable
// Instruction: Disable Interrupts / Clear Interrupt Flag
// Function:    I = 0
func cli(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.I = 0
}

// CLV - Clear Overflow Flag
// Function:    V = 0
func clv(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.V = 0
}

// CMP - Compare the accumulator A
// Function:    C <- A >= M      Z <- (A - M) == 0
// Flags Out:   N, C, Z
func cmp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.compare(cpu.A, cpu.bus.CpuRead(address))
}

// CPX - Compare X Register
// Function:    C <- X >= M      Z <- (X - M) == 0
// Flags Out:   N, C, Z
func cpx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.compare(cpu.X, cpu.bus.CpuRead(address))
}

// CPY - Compare Y Register
// Instruction: Compare Y Register
// Function:    C <- Y >= M      Z <- (Y - M) == 0
// Flags Out:   N, C, Z
func cpy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.compare(cpu.Y, cpu.bus.CpuRead(address))
}

// DEC - Decrement Memory
func dec(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var value byte = cpu.bus.CpuRead(address) - 1

	cpu.bus.CpuWrite(address, value)
	cpu.setZN(value)
}

// DEX - Decrement X Register
func dex(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X--
	cpu.setZN(cpu.X)
}

// DEY - Decrement Y Register
func dey(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y--
	cpu.setZN(cpu.Y)
}

// EOR - Exclusive OR
// instruction: xor gate on the accumulator
func eor(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A ^= cpu.bus.CpuRead(address)
	cpu.setZN(cpu.A)
}

// INC - Increment Memory
func inc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var value byte = cpu.bus.CpuRead(address) + 1

	cpu.bus.CpuWrite(address, value)
	cpu.setZN(value)
}

// INX - Increment X Register
func inx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X++
	cpu.setZN(cpu.X)
}

// INY - Increment Y Register
func iny(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y++
	cpu.setZN(cpu.Y)
	//os.Exit(9)
}

// JMP - Jump
// Instruction: Jump To Location
// Function:    pc = address
func jmp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.PC = address
}

// JSR - Jump to Subroutine
func jsr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.push16(cpu.PC - 1)
	cpu.PC = address
}

// LDA - Load Accumulator
func lda(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.bus.CpuRead(address)
	cpu.setZN(cpu.A)
}

// LDX - Load X Register
func ldx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X = cpu.bus.CpuRead(address)
	cpu.setZN(cpu.X)
}

// LDY - Load Y Register
func ldy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y = cpu.bus.CpuRead(address)
	cpu.setZN(cpu.Y)
}

// LSR - Logical Shift Right
func lsr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	if isAnAccumulator {
		cpu.C = cpu.A & 1
		cpu.A >>= 1
		cpu.setZN(cpu.A)
	} else {
		var value byte = cpu.bus.CpuRead(address)

		cpu.C = value & 1
		value >>= 1
		cpu.bus.CpuWrite(address, value)
		cpu.setZN(value)
	}
}

// NOP - No Operation
func nop(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	//the hardest instruction ;)
}

// ORA - Logical Inclusive OR
// ORA - Logical Inclusive OR on the accumulator
func ora(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A |= cpu.bus.CpuRead(address)
	cpu.setZN(cpu.A)
}

// PHA - Push Accumulator
// Instruction: Push Accumulator to Stack
// Function:    A -> stack write function allow me to access to the BUS
func pha(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.push(cpu.A)
}

// PHP - Push Processor Status
// Instruction: Push Status Register to Stack
// Function:    status -> stack
// Note:        Break flag is set to 1 before push
func php(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.push(cpu.Flags() | 0x10)
}

// PLA - Pull Accumulator
func pla(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.pull()
	cpu.setZN(cpu.A)
}

// PLP - Pull Processor Status
func plp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.SetFlags(cpu.pull()&0xEF | 0x20)
}

// ROL - Rotate Left
func rol(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var tmp byte = cpu.C

	if isAnAccumulator {
		cpu.C = (cpu.A >> 7) & 1
		cpu.A = (cpu.A << 1) | tmp
		cpu.setZN(cpu.A)
	} else {
		value := cpu.bus.CpuRead(address)
		cpu.C = (value >> 7) & 1
		value = (value << 1) | tmp
		cpu.bus.CpuWrite(address, value)
		cpu.setZN(value)
	}
}

// ROR - Rotate Right
func ror(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	var tmp byte = cpu.C

	if isAnAccumulator {
		cpu.C = cpu.A & 1
		cpu.A = (cpu.A >> 1) | (tmp << 7)
		cpu.setZN(cpu.A)
	} else {
		value := cpu.bus.CpuRead(address)
		cpu.C = value & 1
		value = (value >> 1) | (tmp << 7)
		cpu.bus.CpuWrite(address, value)
		cpu.setZN(value)
	}
}

// RTI - Return from Interrupt
func rti(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.SetFlags(cpu.pull()&0xEF | 0x20)
	cpu.PC = cpu.pull16()
}

// RTS - Return from Subroutine
func rts(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.PC = cpu.pull16() + 1
}

// SBC - Subtract with Carry
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

	//cpu.A = a - m - (1 - c) == a + (m ^ 255) + c
	cpu.A += (m ^ 255) + c
	cpu.setZN(cpu.A)
	if int(a)-int(m)-int(1-c) >= 0x0 {
		cpu.C = 1
	} else {
		cpu.C = 0
	}
	if (a^m)&0x80 != 0 && (a^cpu.A)&0x80 != 0 {
		cpu.V = 1
	} else {
		cpu.V = 0
	}
}

// SEC - Set Carry Flag
func sec(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.C = 1
}

// SED - Set Decimal Flag
func sed(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.D = 1
}

// SEI - Set Interrupt Disable
func sei(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.I = 1
}

// STA - Store Accumulator
func sta(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.bus.CpuWrite(address, cpu.A)
}

// STX - Store X Register
func stx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.bus.CpuWrite(address, cpu.X)
}

// STY - Store Y Register
func sty(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.bus.CpuWrite(address, cpu.Y)
}

// TAX - Transfer Accumulator to X
func tax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X = cpu.A
	cpu.setZN(cpu.X)
}

// TAY - Transfer Accumulator to Y
func tay(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.Y = cpu.A
	cpu.setZN(cpu.Y)
}

// TSX - Transfer Stack Pointer to X
func tsx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.X = cpu.SP
	cpu.setZN(cpu.X)
}

// TXA - Transfer X to Accumulator
func txa(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.X
	cpu.setZN(cpu.A)
}

// TXS - Transfer X to Stack Pointer
func txs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.SP = cpu.X
}

// TYA - Transfer Y to Accumulator
func tya(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
	cpu.A = cpu.Y
	cpu.setZN(cpu.A)
}

// illegal opcodes below

func ahx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func alr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func anc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func arr(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func axs(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func dcp(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func isc(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func kil(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func las(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func lax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func rla(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func rra(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func sax(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func shx(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func shy(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func slo(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func sre(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func tas(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
}

func xaa(cpu *CPU, address uint16, pc uint16, isAnAccumulator bool) {
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

// interrupt types
const (
	interruptNone = iota + 1
	interruptNMI
	interruptIRQ
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

//OPCODE MATRIX look doc page 11
//an opCode is composed of:
// - instruction name
// - instruction addr mode
// - instruction size
// - the number of cycles that the instruction take
// - the number of page that the instruction take
// - the function that exec the instruction
type opCode struct {
	instructionName string
	instructionMode byte
	instructionSize uint16
	nbCycle         uint16
	nbPageCycles    uint16 //// instructionPageCycles indicates the number of cycles used by each instruction when a page is crossed
	instructionExec execinstructions
}

//this followed matrix shows the 210 op code (illegal/NOP are not counted) associated with the R65C00 family CPU devices.
//enjoy ;)
//map of instruction
var opCodeMatrix = [256]opCode{
	opCode{instructionName: "BRK", instructionMode: modeImplied, instructionSize: 2, nbCycle: 7, nbPageCycles: 0, instructionExec: brk}, opCode{instructionName: "ORA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SLO", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "PHP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, nbPageCycles: 0, instructionExec: php}, opCode{instructionName: "ORA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "ANC", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: anc}, opCode{instructionName: "NOP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: slo},
	opCode{instructionName: "BPL", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bpl}, opCode{instructionName: "ORA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: ora}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SLO", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "CLC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: clc}, opCode{instructionName: "ORA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ora}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SLO", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: slo}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "ORA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ora}, opCode{instructionName: "ASL", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: asl}, opCode{instructionName: "SLO", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: slo},
	opCode{instructionName: "JSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: jsr}, opCode{instructionName: "AND", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RLA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "BIT", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: bit}, opCode{instructionName: "AND", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "PLP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, nbPageCycles: 0, instructionExec: plp}, opCode{instructionName: "AND", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "ANC", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: anc}, opCode{instructionName: "BIT", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: bit}, opCode{instructionName: "AND", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rla},
	opCode{instructionName: "BMI", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bmi}, opCode{instructionName: "AND", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: and}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RLA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "AND", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "SEC", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: sec}, opCode{instructionName: "AND", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: and}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rla}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "AND", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: and}, opCode{instructionName: "ROL", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: rol}, opCode{instructionName: "RLA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rla},
	opCode{instructionName: "RTI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, nbPageCycles: 0, instructionExec: rti}, opCode{instructionName: "EOR", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SRE", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "PHA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 3, nbPageCycles: 0, instructionExec: pha}, opCode{instructionName: "EOR", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "ALR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: alr}, opCode{instructionName: "JMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 3, nbPageCycles: 0, instructionExec: jmp}, opCode{instructionName: "EOR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: sre},
	opCode{instructionName: "BVC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bvc}, opCode{instructionName: "EOR", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: eor}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "SRE", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "CLI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: cli}, opCode{instructionName: "EOR", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: eor}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SRE", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: sre}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "EOR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: eor}, opCode{instructionName: "LSR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: lsr}, opCode{instructionName: "SRE", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: sre},
	opCode{instructionName: "RTS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 6, nbPageCycles: 0, instructionExec: rts}, opCode{instructionName: "ADC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RRA", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "PLA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 4, nbPageCycles: 0, instructionExec: pla}, opCode{instructionName: "ADC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAccumulator, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "ARR", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: arr}, opCode{instructionName: "JMP", instructionMode: modeIndirect, instructionSize: 3, nbCycle: 5, nbPageCycles: 0, instructionExec: jmp}, opCode{instructionName: "ADC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rra},
	opCode{instructionName: "BVS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bvs}, opCode{instructionName: "ADC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: adc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "RRA", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "SEI", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: sei}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: adc}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rra}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "ADC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: adc}, opCode{instructionName: "ROR", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: ror}, opCode{instructionName: "RRA", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: rra},
	opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "STA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: sax}, opCode{instructionName: "STY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 3, nbPageCycles: 0, instructionExec: sax}, opCode{instructionName: "DEY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: dey}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "TXA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: txa}, opCode{instructionName: "XAA", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: xaa}, opCode{instructionName: "STY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: sax},
	opCode{instructionName: "BCC", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bcc}, opCode{instructionName: "STA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "AHX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: ahx}, opCode{instructionName: "STY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: sty}, opCode{instructionName: "STA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "STX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: stx}, opCode{instructionName: "SAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: sax}, opCode{instructionName: "TYA", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tya}, opCode{instructionName: "STA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 5, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "TXS", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: txs}, opCode{instructionName: "TAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: tas}, opCode{instructionName: "SHY", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: shy}, opCode{instructionName: "STA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 5, nbPageCycles: 0, instructionExec: sta}, opCode{instructionName: "SHX", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: shx}, opCode{instructionName: "AHX", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: ahx},
	opCode{instructionName: "LDY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeImmediate, instructionSize: 2 /*0*/, nbCycle: 2, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 3, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "TAY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tay}, opCode{instructionName: "LDA", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "TAX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tax}, opCode{instructionName: "LAX", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: lax},
	opCode{instructionName: "BCS", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bcs}, opCode{instructionName: "LDA", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: lda}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "LAX", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 5, nbPageCycles: 1, instructionExec: lax}, opCode{instructionName: "LDY", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeZeroPageY, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeZeroPageY, instructionSize: 0, nbCycle: 4, nbPageCycles: 0, instructionExec: lax}, opCode{instructionName: "CLV", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: clv}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: lda}, opCode{instructionName: "TSX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: tsx}, opCode{instructionName: "LAS", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, nbPageCycles: 1, instructionExec: las}, opCode{instructionName: "LDY", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ldy}, opCode{instructionName: "LDA", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: lda}, opCode{instructionName: "LDX", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: ldx}, opCode{instructionName: "LAX", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 4, nbPageCycles: 1, instructionExec: lax},
	opCode{instructionName: "CPY", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "DCP", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "CPY", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "INY", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: iny}, opCode{instructionName: "CMP", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: dex}, opCode{instructionName: "AXS", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: axs}, opCode{instructionName: "CPY", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: cpy}, opCode{instructionName: "CMP", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: dcp},
	opCode{instructionName: "BNE", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: bne}, opCode{instructionName: "CMP", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: cmp}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "DCP", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "CMP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "CLD", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: cld}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: cmp}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: dcp}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "CMP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: cmp}, opCode{instructionName: "DEC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: dec}, opCode{instructionName: "DCP", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: dcp},
	opCode{instructionName: "CPX", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeIndexedIndirect, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "NOP", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ISC", instructionMode: modeIndexedIndirect, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "CPX", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 3, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeZeroPage, instructionSize: 2, nbCycle: 5, nbPageCycles: 0, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeZeroPage, instructionSize: 0, nbCycle: 5, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "INX", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: inx}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 2, nbCycle: 2, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeImmediate, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "CPX", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: cpx}, opCode{instructionName: "SBC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 4, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeAbsolute, instructionSize: 3, nbCycle: 6, nbPageCycles: 0, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeAbsolute, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: isc},
	opCode{instructionName: "BEQ", instructionMode: modeRelative, instructionSize: 2, nbCycle: 2, nbPageCycles: 1, instructionExec: beq}, opCode{instructionName: "SBC", instructionMode: modeIndirectIndexed, instructionSize: 2, nbCycle: 5, nbPageCycles: 1, instructionExec: sbc}, opCode{instructionName: "KIL", instructionMode: modeImplied, instructionSize: 0, nbCycle: 2, nbPageCycles: 0, instructionExec: kil}, opCode{instructionName: "ISC", instructionMode: modeIndirectIndexed, instructionSize: 0, nbCycle: 8, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "NOP", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 4, nbPageCycles: 0, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeZeroPageX, instructionSize: 2, nbCycle: 6, nbPageCycles: 0, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeZeroPageX, instructionSize: 0, nbCycle: 6, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "SED", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: sed}, opCode{instructionName: "SBC", instructionMode: modeAbsoluteY, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: sbc}, opCode{instructionName: "NOP", instructionMode: modeImplied, instructionSize: 1, nbCycle: 2, nbPageCycles: 0, instructionExec: nop}, opCode{instructionName: "ISC", instructionMode: modeAbsoluteY, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: isc}, opCode{instructionName: "NOP", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: nop}, opCode{instructionName: "SBC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 4, nbPageCycles: 1, instructionExec: sbc}, opCode{instructionName: "INC", instructionMode: modeAbsoluteX, instructionSize: 3, nbCycle: 7, nbPageCycles: 0, instructionExec: inc}, opCode{instructionName: "ISC", instructionMode: modeAbsoluteX, instructionSize: 0, nbCycle: 7, nbPageCycles: 0, instructionExec: isc},
}

//CPU nes
type CPU struct {
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

//map of addr mode constructor
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

// triggerIRQ causes an IRQ interrupt to occur on the next cycle // NOTE : useless for now !
func (cpu *CPU) triggerIRQ() {
	if cpu.I == 0 {
		cpu.interrupt = interruptIRQ
		os.Exit(5)
	}
}

// NMI Non-Maskable Interrupt
// IRQ IRQ interrupt
//interruptMode == interrupt nmi -> NMI
//interruptMode == interrupt irq -> IRQ
func (cpu *CPU) cpuInterruptions(interruptMode byte) {
	//check i == 0
	if interruptMode != interruptNone {
		cpu.push16(cpu.PC)
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
	cpu.PC = cpu.Read16(0xFFFC)
	cpu.SP = 0xFD
	// Reset resets the CPU to its initial powerup state
	cpu.SetFlags(0x24)
	cpu.interrupt = interruptNone //tmp
}

//NewCpu function is the constructor of my CPU
func NewCpu(bus *BUS) *CPU {
	var cpu CPU = CPU{bus: bus}

	cpu.modesTable = createModesTables()
	// init cpu values by reseting the component. smart ;)
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

// PrintInstruction prints the current CPU state
//VERY USEFUL for debugging trust me !!!
func (cpu *CPU) PrintInstruction() {
	opCodeIndex := cpu.bus.CpuRead(cpu.PC)
	opcode := opCodeMatrix[opCodeIndex]
	bytes := opcode.instructionSize
	name := opcode.instructionName

	println("-------------------------------------------------- ", cpu.Cycles, " --------------------------------------------------")
	println("opcode.instructionName = ", opcode.instructionName, " opcode.instructionSize = ", opcode.instructionSize, " opcode.nbCycle = ", opcode.nbCycle, " opcode.nbPageCycles = ", opcode.nbPageCycles, "  opcode.instructionMode = ", opcode.instructionMode)
	w0 := fmt.Sprintf("%02X", cpu.bus.CpuRead(cpu.PC+0))
	w1 := fmt.Sprintf("%02X", cpu.bus.CpuRead(cpu.PC+1))
	w2 := fmt.Sprintf("%02X", cpu.bus.CpuRead(cpu.PC+2))

	if bytes < 2 {
		w1 = "  "
	}
	if bytes < 3 {
		w2 = "  "
	}
	fmt.Printf(
		"%4X  %s %s %s  %s %28s"+
			"A:%02X X:%02X Y:%02X SP:%02X CYC:%3d\n",
		cpu.PC, w0, w1, w2, name, "",
		cpu.A, cpu.X, cpu.Y, cpu.SP, (cpu.Cycles*3)%341)
}

// Step executes a single CPU instruction
func (cpu *CPU) Step() uint64 {
	if cpu.stall > 0 {
		cpu.stall--
		return 1
	}
	var startNbCycles uint64 = cpu.Cycles

	//check whether an Iterruption occur or not
	cpu.cpuInterruptions(cpu.interrupt)
	cpu.interrupt = interruptNone
	//get the opcode
	var opCodeIndex byte = cpu.bus.CpuRead(cpu.PC)
	var op opCode = opCodeMatrix[opCodeIndex]
	var isAnAccumulator bool = false
	var address uint16

	//get address from bus
	address = cpu.modesTable[op.instructionMode](cpu) //return addr mode
	cpu.PC += uint16(op.instructionSize)
	cpu.Cycles += uint64(op.nbCycle)

	if op.instructionMode == modeAccumulator {
		isAnAccumulator = true
	}
	//check wheter all the data are on the same page
	if cpu.isPageCrossed(op, address) {
		cpu.Cycles += uint64(op.nbPageCycles)
	}
	//exec instruction
	op.instructionExec(cpu, address, cpu.PC, isAnAccumulator)
	return cpu.Cycles - startNbCycles
}
