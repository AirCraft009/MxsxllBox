package main

import (
	"MxsxllBox/cpu"
)

func EncodeRegs(rx, ry byte) byte {
	return ((rx & 0x07) << 5) | ((ry & 0x07) << 2)
}

func EncodeAddr(addr uint16) (byte, byte) {
	if addr <= 255 {
		return byte(0), byte(addr)
	}
	return byte(addr>>8) & 0xff, byte(addr & 0xff)
}

func main() {
	mem := &cpu.Memory{}

	program := []byte{
		cpu.MOVI,
	}
	hi, low := EncodeAddr('w')
	program = append(program, EncodeRegs(2, 0), hi, low)
	program = append(program, cpu.MOVI, EncodeRegs(4, 0), hi, low)
	program = append(program, cpu.ADD, EncodeRegs(2, 4))
	program = append(program, cpu.PRINT, EncodeRegs(2, 0), cpu.HALT)

	copy(mem.Data[:cpu.MemorySize], program)

	vm := cpu.NewCPU(mem)
	vm.Run()
}
