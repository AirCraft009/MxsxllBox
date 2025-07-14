package main

import (
	"MxsxllBox/cpu"
	"MxsxllBox/helper"
)

func main() {
	mem := cpu.NewMemory()

	program := []byte{}
	hi, low := helper.EncodeAddr('w')
	divhi, divlow := helper.EncodeAddr(20000)
	beginhi, beginlow := helper.EncodeAddr(0)
	program = append(program, cpu.MOVI, helper.EncodeRegs(2, 0), hi, low)
	program = append(program, cpu.PRINT, helper.EncodeRegs(2, 0))
	program = append(program, cpu.MULI, helper.EncodeRegs(2, 0), divhi, divlow)
	program = append(program, cpu.JZ, 0x00, 0x00, beginhi, beginlow)
	program = append(program, cpu.HALT)

	copy(mem.Data[:cpu.MemorySize], program)

	vm := cpu.NewCPU(mem)
	vm.Run()
}
