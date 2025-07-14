package main

import (
	"MxsxllBox/cpu"
)

func main() {
	mem := &cpu.Memory{}

	program := []byte{
		cpu.MOVI, 0x00, 0x00, 0x2A,
		cpu.PRINT, 0x00,
		cpu.HALT,
	}

	copy(mem.Data[:cpu.MemorySize], program)

	vm := cpu.NewCPU(mem)
	vm.Run()
}
