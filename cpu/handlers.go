package cpu

import (
	"fmt"
)

const (
	//for any operation that doesn't use the addr
	instructionSizeShort = 2
	//for any operation that does use the addr
	instructionSizeLong = 4
)

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
	opcode = cpu.Mem.ReadByte(cpu.PC)
	regs := cpu.Mem.ReadByte(cpu.PC + 1)
	rx, ry := decodeReg(regs)
	/**
	addr is twice as long so 16 bits we calculate it by reading two times,
	then upshifting the first by 8 and fusing them with the second read

	read PC+3
	adrr = 0000000010111001
	addr << 8
	adrr = 1011100100000000
	read PC + 3
	bitwise or
	adrr = 1011100101101011
	*/
	addr := uint16(cpu.Mem.ReadByte(cpu.PC+instructionSizeShort))<<8 | (uint16(cpu.Mem.ReadByte(cpu.PC + 3)))
	instructions = newHandlerInstructions(rx, ry, addr)
	return opcode, instructions
}

func decodeReg(reg byte) (rx byte, ry byte) {
	/**
	reg contains both rx and ry
	rx = bits 7-5
	ry = bits 4-2
	flags, etc. = bits 0 - 1

	>> rightshifts all bits by the following number
	& bitwise and looks at each number does and
	Example for decoding
	reg = 11010101
	reg >> 5 = 00000110
	reg & 0x07 = 00000110 & 00000111
	rx = 00000110
	rx = 6

	reg = 11010101
	reg >> 3 = 00110101
	reg & 0x07 = 00110101 & 00000111
	ry = 00000101
	ry = 5

	*/
	rx = (reg >> 5) & 0x07
	ry = (reg >> 3) & 0x07
	return rx, ry
}

func handlePush(cpu *CPU, instruction *HandlerInstructions) {
	val := cpu.Registers[instruction.Rx]
	cpu.SP -= instructionSizeShort
	cpu.Mem.WriteByte(cpu.SP, byte(val>>8))
	cpu.Mem.WriteByte(cpu.SP+1, byte(val&0xff))
	cpu.PC += instructionSizeShort
}

func handlePop(cpu *CPU, instruction *HandlerInstructions) {
	hi := cpu.Mem.ReadByte(cpu.SP)
	lo := cpu.Mem.ReadByte(cpu.SP + 1)
	cpu.Registers[instruction.Rx] = uint16(hi)<<8 | uint16(lo)
	cpu.PC += instructionSizeLong
	cpu.SP += instructionSizeShort
}

func handleCall(cpu *CPU, instruction *HandlerInstructions) {
	cpu.SP -= instructionSizeShort
	cpu.Mem.WriteWord(cpu.SP, cpu.PC)
	handleJmp(cpu, instruction)
}

func handleRet(cpu *CPU, instruction *HandlerInstructions) {
	instruction.Addr = cpu.Mem.ReadWord(cpu.SP) + instructionSizeLong
	cpu.PC += instructionSizeLong
	cpu.SP += instructionSizeShort
	handleJmp(cpu, instruction)
}

func handleReadWriteSize(addr uint16) bool {
	//true means 1 byte access false means full word access
	return addr >= VideoStart && addr <= VideoEnd
}

func handleNop(cpu *CPU, instructions *HandlerInstructions) {
	cpu.PC++
	return
}

func handleLoad(cpu *CPU, instructions *HandlerInstructions) {
	if handleReadWriteSize(instructions.Addr) {
		handleLoadB(cpu, instructions)
		return
	}
	handleLoadW(cpu, instructions)
}

func handleStore(cpu *CPU, instructions *HandlerInstructions) {
	if handleReadWriteSize(instructions.Addr) {
		handleStoreB(cpu, instructions)
		return
	}
	handleStoreW(cpu, instructions)
}

func handleLoadB(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = uint16(cpu.Mem.ReadByte(instructions.Addr))
	cpu.PC += instructionSizeLong
}

func handleLoadW(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = uint16(cpu.Mem.ReadByte(instructions.Addr)) << 8
	cpu.Registers[instructions.Rx] |= uint16(cpu.Mem.ReadByte(instructions.Addr + 1))
	cpu.PC += instructionSizeLong
}

func handleStoreB(cpu *CPU, instructions *HandlerInstructions) {
	val := byte(cpu.Registers[instructions.Rx] & 0xFF)
	cpu.Mem.WriteByte(instructions.Addr, val)
	cpu.PC += instructionSizeLong
}

func handleStoreW(cpu *CPU, instructions *HandlerInstructions) {
	val := cpu.Registers[instructions.Rx]
	cpu.Mem.WriteByte(instructions.Addr, byte((val>>8)&0xFF))
	cpu.Mem.WriteByte(instructions.Addr+1, byte(val&0xff))
	cpu.PC += instructionSizeLong
}

func handleAdd(cpu *CPU, instructions *HandlerInstructions) {
	rx := instructions.Rx
	a := cpu.Registers[rx]
	b := cpu.Registers[instructions.Ry]
	result := uint32(a) + uint32(b)

	cpu.Registers[rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[rx] == 0x00
	cpu.Flags.Carry = result > 0xffff

	cpu.PC += instructionSizeShort
}

func handleSub(cpu *CPU, instructions *HandlerInstructions) {
	rx := instructions.Rx
	a := cpu.Registers[rx]
	b := cpu.Registers[instructions.Ry]
	result := uint32(a) - uint32(b)

	cpu.Registers[instructions.Rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[rx] == 0x00
	cpu.Flags.Carry = result < 0x00

	cpu.PC += instructionSizeShort
}

func handleMul(cpu *CPU, instructions *HandlerInstructions) {
	rx := instructions.Rx
	a := cpu.Registers[rx]
	b := cpu.Registers[instructions.Ry]
	result := uint32(a) * uint32(b)

	cpu.Registers[instructions.Rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[rx] == 0x00
	cpu.Flags.Carry = result > 0xffff

	cpu.PC += instructionSizeShort
}

func handleDiv(cpu *CPU, instructions *HandlerInstructions) {
	rx := instructions.Rx
	a := cpu.Registers[rx]
	b := cpu.Registers[instructions.Ry]
	result := uint32(a) / uint32(b)
	cpu.Registers[instructions.Rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[rx] == 0x00
	cpu.Flags.Carry = false
	cpu.PC += instructionSizeShort
}

func handleJmp(cpu *CPU, instructions *HandlerInstructions) {
	if instructions.Addr <= ProgramEnd {
		cpu.PC = instructions.Addr
	}
}

func handleJc(cpu *CPU, instructions *HandlerInstructions) {
	if cpu.Flags.Carry && instructions.Addr <= ProgramEnd {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handleJz(cpu *CPU, instructions *HandlerInstructions) {
	if cpu.Flags.Zero && instructions.Addr <= ProgramEnd {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handlePrint(cpu *CPU, instructions *HandlerInstructions) {
	cpu.PC += instructionSizeShort
	fmt.Println(cpu.Registers[instructions.Rx])
}

func handleMovi(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = instructions.Addr
	cpu.PC += instructionSizeLong
}

func handleAddi(cpu *CPU, instructions *HandlerInstructions) {
	rx := cpu.Registers[instructions.Rx]
	result := uint32(rx) + uint32(instructions.Addr)
	cpu.Registers[instructions.Rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[instructions.Rx] == 0x00
	cpu.Flags.Carry = result > 0xffff
	cpu.PC += instructionSizeLong
}

func handleSubi(cpu *CPU, instructions *HandlerInstructions) {
	rx := cpu.Registers[instructions.Rx]
	result := uint32(rx) + -uint32(instructions.Addr)
	cpu.Registers[instructions.Rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[instructions.Rx] == 0x00
	cpu.Flags.Carry = result < 0x00
	cpu.PC += instructionSizeLong
}

func handleMuli(cpu *CPU, instructions *HandlerInstructions) {
	rx := cpu.Registers[instructions.Rx]
	result := uint32(rx) * uint32(instructions.Addr)
	cpu.Registers[instructions.Rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[instructions.Rx] == 0x00
	cpu.Flags.Carry = result > 0xffff
	cpu.PC += instructionSizeLong
}

func handleDivi(cpu *CPU, instructions *HandlerInstructions) {
	rx := cpu.Registers[instructions.Rx]
	result := uint32(rx) / uint32(instructions.Addr)
	cpu.Registers[instructions.Rx] = uint16(result)
	cpu.Flags.Zero = cpu.Registers[instructions.Rx] == 0x00
	cpu.Flags.Carry = false
	cpu.PC += instructionSizeLong
}

func handleHalt(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Halted = true
}
