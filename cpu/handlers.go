package cpu

import "fmt"

type HandlerInstructions struct {
	Rx   byte
	Ry   byte
	Addr uint16
}

func newHandlerInstructions(rx byte, ry byte, addr uint16) *HandlerInstructions {
	return &HandlerInstructions{
		Rx:   rx,
		Ry:   ry,
		Addr: addr,
	}
}

func getInstruction(cpu *CPU) (opcode byte, instructions *HandlerInstructions) {
	opcode = cpu.Mem.Read(cpu.PC)
	regs := cpu.Mem.Read(cpu.PC + 1)
	rx, ry := decodeReg(regs)
	/**
	addr is twice as long so 16 bits we calculate it by reading two times,
	then upshifting the first by 8 and fusing them with the second read

	read PC+2
	adrr = 0000000010111001
	addr << 8
	adrr = 1011100100000000
	read PC + 3
	bitwise or
	adrr = 1011100101101011
	*/
	addr := uint16(cpu.Mem.Read(cpu.PC+2))<<8 | (uint16(cpu.Mem.Read(cpu.PC + 3)))
	instructions = newHandlerInstructions(rx, ry, addr)
	return opcode, instructions
}

func decodeReg(reg byte) (rx byte, ry byte) {
	/**
	reg contains both rx and ry
	rx = bits 7-5
	ry = bits 4-2
	flags, etc. = bits 0 - 2

	>> rightshifts all bits by the following number
	& bitwise and looks at each number does and
	Example for decoding
	reg = 11010101
	reg >> 5 = 00000110
	reg & 0x07 = 00000110 & 00000111
	rx = 00000110
	rx = 6

	reg = 11010101
	reg >> 2 = 00110101
	reg & 0x07 = 00110101 & 00000111
	ry = 00000101
	ry = 5

	*/
	rx = (reg >> 5) & 0x07
	ry = (reg >> 2) & 0x07
	return rx, ry
}

func handleNop(cpu *CPU, instructions *HandlerInstructions) {
	return
}

func handleLoadB(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = instructions.Addr
	cpu.PC += 4
}

func handleLoadW(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = uint16(cpu.Mem.Read(instructions.Addr)) << 8
	cpu.Registers[instructions.Rx] |= uint16(cpu.Mem.Read(instructions.Addr + 1))
	cpu.PC += 4
}

func handleStoreB(cpu *CPU, instructions *HandlerInstructions) {
	val := byte(cpu.Registers[instructions.Rx] & 0xFF)
	cpu.Mem.Write(instructions.Addr, val)
	cpu.PC += 4
}

func handleStoreW(cpu *CPU, instructions *HandlerInstructions) {
	val := cpu.Registers[instructions.Rx]
	cpu.Mem.Write(instructions.Addr, byte(val&0xFF))
	cpu.Mem.Write(instructions.Addr+1, byte(val>>8))
	cpu.PC += 4
}

func handleAdd(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] += cpu.Registers[instructions.Ry]
	cpu.PC += 2
}

func handleSub(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] -= cpu.Registers[instructions.Ry]
	cpu.PC += 2
}

func handleMul(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] *= cpu.Registers[instructions.Ry]
	cpu.PC += 2
}

func handleDiv(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] /= cpu.Registers[instructions.Ry]
	cpu.PC += 2
}

func handleJmp(cpu *CPU, instructions *HandlerInstructions) {
	cpu.PC = instructions.Addr
}

func handleJc(cpu *CPU, instructions *HandlerInstructions) {
	//implement later
	cpu.PC += 1
}

func handleJz(cpu *CPU, instructions *HandlerInstructions) {
	//implement later
	cpu.PC += 1
}

func handlePrint(cpu *CPU, instructions *HandlerInstructions) {
	cpu.PC += 1
	fmt.Println(cpu.Registers[instructions.Rx])
}

func handleHalt(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Halted = true
}
