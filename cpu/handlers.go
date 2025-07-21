package cpu

import (
	"fmt"
)

const (
	//for any operation that doesn't use the addr
	instructionSizeShort = 3
	//for any operation that does use the addr
	instructionSizeLong = 5
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
	if cpu.PC > ProgramEnd {
		fmt.Println(cpu.PC)
		fmt.Println(cpu.Mem.Data[8257:8795])
		panic("program out of memory")
	} else if cpu.SP < StackStart {
		fmt.Println(cpu.SP)
		fmt.Println(cpu.Mem.Data[8379:8795])
		panic("stack out of memory")
	}
	opcode = cpu.Mem.ReadByte(cpu.PC)
	regs1, flagbyte := cpu.Mem.ReadReg(cpu.PC + 1)
	rx, ry, addresnec := decodeReg(regs1, flagbyte)
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
	var addr uint16
	if addresnec {
		addr = cpu.Mem.ReadWord(cpu.PC + instructionSizeShort)
	}
	instructions = newHandlerInstructions(rx, ry, addr)
	return opcode, instructions
}

func decodeReg(reg1, reg2Wflag byte) (rx byte, ry byte, addresNec bool) {
	/**
	old: theory still applies{
		new config:
		rx = bits 8-4
		ry = bits 3-0
		addrFlag = byte 2 bit 0
	}
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
	reg >> 2 = 00110101
	reg & 0x07 = 00110101 & 00000111
	ry = 00000101
	ry = 5

	*/
	rx = reg1
	ry = reg2Wflag >> 1
	addrnec := (reg2Wflag) & 0x01
	return rx, ry, addrnec != 0x0
}

func handleOr(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] |= cpu.Registers[instructions.Ry]
	cpu.PC += instructionSizeShort
}

func handleGf(cpu *CPU, instructions *HandlerInstructions) {
	var z, c uint16
	if cpu.Flags.Zero {
		z = 1
	}
	if cpu.Flags.Carry {
		c = 1
	}
	flag := (c << 1) | z
	cpu.Registers[instructions.Rx] = flag
	cpu.PC += instructionSizeShort
}

func handleSf(cpu *CPU, instructions *HandlerInstructions) {
	flag := cpu.Registers[instructions.Rx]
	c := flag >> 1
	z := flag & 0x01
	cpu.Flags.Zero, cpu.Flags.Carry = z == 1, c == 1
	cpu.PC += instructionSizeShort
}

func handleSrfn(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[cpu.Registers[instructions.Rx]] = cpu.Registers[instructions.Ry]
	cpu.PC += instructionSizeShort
}

func handleGrfn(cpu *CPU, instructions *HandlerInstructions) { //num dst
	cpu.Registers[instructions.Ry] = cpu.Registers[cpu.Registers[instructions.Rx]]
	cpu.PC += instructionSizeShort
}

func handleGPc(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = cpu.PC
	cpu.PC += instructionSizeShort
}

func handleGSp(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = cpu.SP
	cpu.PC += instructionSizeShort
}

func handleSPc(cpu *CPU, instructions *HandlerInstructions) {
	cpu.PC = cpu.Registers[instructions.Rx]
}

func handleSSp(cpu *CPU, instructions *HandlerInstructions) {
	cpu.SP = cpu.Registers[instructions.Rx]
	cpu.PC += instructionSizeShort
}

func handleMova(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = instructions.Addr
	cpu.PC += instructionSizeLong
}

func handleAnd(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] &= cpu.Registers[instructions.Ry]
	cpu.PC += instructionSizeShort
}

func handleRs(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] <<= cpu.Registers[instructions.Ry]
	cpu.PC += instructionSizeShort
}

func handleLs(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] >>= cpu.Registers[instructions.Ry]
	cpu.PC += instructionSizeShort
}

func handleMov(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] = cpu.Registers[instructions.Ry]
	cpu.PC += instructionSizeShort
}

func handleModi(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] %= instructions.Addr
	cpu.Flags.Zero = false
	if cpu.Registers[instructions.Rx] == 0 {
		cpu.Flags.Zero = true
	}
	cpu.PC += instructionSizeLong
}

func handleMod(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Registers[instructions.Rx] %= cpu.Registers[instructions.Ry]
	cpu.Flags.Zero = false
	if cpu.Registers[instructions.Rx] == 0 {
		cpu.Flags.Zero = true
	}
	cpu.PC += instructionSizeShort
}

func handleSTZ(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Zero = true
	cpu.PC++
}

func handleSTC(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Carry = true
	cpu.PC++
}

func handleCLZ(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Zero = false
	cpu.PC++
}

func handleCLC(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Carry = false
	cpu.PC++
}

func handleJG(cpu *CPU, instructions *HandlerInstructions) {
	if !cpu.Flags.Zero && cpu.Flags.Carry {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handleJGE(cpu *CPU, instructions *HandlerInstructions) {
	if cpu.Flags.Zero && cpu.Flags.Carry {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handleJLE(cpu *CPU, instructions *HandlerInstructions) {
	if cpu.Flags.Zero && !cpu.Flags.Carry {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handleJL(cpu *CPU, instructions *HandlerInstructions) {
	if !cpu.Flags.Zero && !cpu.Flags.Carry {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handleTsti(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Zero = false
	result := cpu.Registers[instructions.Rx] & instructions.Addr
	if result != 0 {
		cpu.Flags.Zero = true
	}
	cpu.PC += instructionSizeLong
}

func handleTest(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Zero = false
	result := cpu.Registers[instructions.Rx] & cpu.Registers[instructions.Ry]
	if result == 0 {
		cpu.Flags.Zero = true
	}
	cpu.PC += instructionSizeShort
}

func handleCmpi(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Zero = false
	cpu.Flags.Carry = false
	if cpu.Registers[instructions.Rx] == instructions.Addr {
		cpu.Flags.Zero = true
		cpu.Flags.Carry = false
	} else if cpu.Registers[instructions.Rx] > instructions.Addr {
		cpu.Flags.Carry = true
		cpu.Flags.Zero = false
	}
	cpu.PC += instructionSizeLong
}

func handleCmp(cpu *CPU, instructions *HandlerInstructions) {
	cpu.Flags.Zero = false
	cpu.Flags.Carry = false
	if cpu.Registers[instructions.Rx] == cpu.Registers[instructions.Ry] {
		cpu.Flags.Zero = true
		cpu.Flags.Carry = false
	} else if cpu.Registers[instructions.Rx] > cpu.Registers[instructions.Ry] {
		cpu.Flags.Carry = true
		cpu.Flags.Zero = false
	}
	cpu.PC += instructionSizeShort
}

func handleJnz(cpu *CPU, instructions *HandlerInstructions) {
	if !cpu.Flags.Zero {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handleJnc(cpu *CPU, instructions *HandlerInstructions) {
	if !cpu.Flags.Carry {
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handlePrintstr(cpu *CPU, instructions *HandlerInstructions) {
	lenght := cpu.Mem.ReadByte(cpu.Registers[instructions.Rx])
	outPutStr := ""
	for i := uint16(1); i <= uint16(lenght); i++ {
		outPutStr += string(cpu.Mem.ReadByte(cpu.Registers[instructions.Rx] + i))
	}
	fmt.Println(outPutStr)
	cpu.PC += instructionSizeShort
}

func handlePush(cpu *CPU, instruction *HandlerInstructions) {
	val := cpu.Registers[instruction.Rx]
	cpu.SP -= 2
	cpu.Mem.WriteByte(cpu.SP, byte(val>>8))
	cpu.Mem.WriteByte(cpu.SP+1, byte(val&0xff))
	cpu.PC += instructionSizeShort
}

func handlePop(cpu *CPU, instruction *HandlerInstructions) {
	hi := cpu.Mem.ReadByte(cpu.SP)
	lo := cpu.Mem.ReadByte(cpu.SP + 1)
	cpu.Registers[instruction.Rx] = uint16(hi)<<8 | uint16(lo)
	cpu.PC += instructionSizeShort
	cpu.SP += 2
}

func handleCall(cpu *CPU, instruction *HandlerInstructions) {
	cpu.SP -= 2
	cpu.Mem.WriteWord(cpu.SP, cpu.PC)
	handleJmp(cpu, instruction)
}

func handleRet(cpu *CPU, instruction *HandlerInstructions) {
	instruction.Addr = cpu.Mem.ReadWord(cpu.SP) + instructionSizeLong
	cpu.PC += instructionSizeLong
	cpu.SP += 2
	handleJmp(cpu, instruction)
}

func handleReadWriteSize(addr uint16, regAddr uint16) bool {
	//true means 1 byte access false means full word access
	return addr >= VideoStart && addr <= VideoEnd || regAddr >= VideoStart && regAddr <= VideoEnd
}

func handleNop(cpu *CPU, instructions *HandlerInstructions) {
	cpu.PC++
	return
}

func handleLoad(cpu *CPU, instructions *HandlerInstructions) {
	if handleReadWriteSize(instructions.Addr, cpu.Registers[instructions.Ry]) {
		handleLoadB(cpu, instructions)
		return
	}
	handleLoadW(cpu, instructions)
}

func handleStore(cpu *CPU, instructions *HandlerInstructions) {
	if handleReadWriteSize(instructions.Addr, cpu.Registers[instructions.Ry]) {
		handleStoreB(cpu, instructions)
		return
	}
	handleStoreW(cpu, instructions)
}

func handleLoadB(cpu *CPU, instructions *HandlerInstructions) {
	if instructions.Addr == 0 && cpu.Registers[instructions.Ry] != 0 {
		addr := cpu.Registers[instructions.Ry]
		cpu.Registers[instructions.Rx] = uint16(cpu.Mem.ReadByte(addr))
		cpu.PC += instructionSizeLong
		return
	}
	cpu.Registers[instructions.Rx] = uint16(cpu.Mem.ReadByte(instructions.Addr))
	cpu.PC += instructionSizeLong
}

func handleLoadW(cpu *CPU, instructions *HandlerInstructions) {
	if instructions.Addr == 0 && cpu.Registers[instructions.Ry] != 0 {
		addr := cpu.Registers[instructions.Ry]
		cpu.Registers[instructions.Rx] = cpu.Mem.ReadWord(addr)
		cpu.PC += instructionSizeLong
		return
	}
	cpu.Registers[instructions.Rx] = cpu.Mem.ReadWord(instructions.Addr)
	cpu.PC += instructionSizeLong
}

func handleStoreB(cpu *CPU, instructions *HandlerInstructions) {
	val := byte(cpu.Registers[instructions.Rx] & 0xFF)
	if instructions.Addr == 0 {

		cpu.Mem.WriteByte(cpu.Registers[instructions.Ry], val)
		cpu.PC += instructionSizeLong
		return
	}
	cpu.Mem.WriteByte(instructions.Addr, val)
	cpu.PC += instructionSizeLong
}

func handleStoreW(cpu *CPU, instructions *HandlerInstructions) {
	val := cpu.Registers[instructions.Rx]
	if instructions.Addr == 0 && cpu.Registers[instructions.Ry] != 0 {
		cpu.Mem.WriteWord(cpu.Registers[instructions.Ry], val)
		cpu.PC += instructionSizeLong
		return
	}
	cpu.Mem.WriteWord(instructions.Addr, val)
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
		cpu.Flags.Carry = false
		cpu.PC = instructions.Addr
		return
	}
	cpu.PC += instructionSizeLong
}

func handleJz(cpu *CPU, instructions *HandlerInstructions) {
	if cpu.Flags.Zero && instructions.Addr <= ProgramEnd {
		cpu.Flags.Zero = false
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
	//fmt.Printf("adding: %d + %d = %d\n", rx, instructions.Addr, result)
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
