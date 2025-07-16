package main

import (
	"MxsxllBox/cpu"
	"MxsxllBox/linker"
)

func main() {
	mem := &cpu.Memory{}
	filesWLocation := make(map[string]uint16)
	filesWLocation["C:/Users/cocon/Documents/Projects/Musa-Allmer/MxsxllBox/program.bin"] = 0
	filesWLocation["C:/Users/cocon/Documents/Projects/Musa-Allmer/MxsxllBox/stdlib/strings.bin"] = 300

	data, _ := linker.LinkModuels(filesWLocation)

	copy(mem.Data[:], data)
	/**
	program := []byte{}

	lo, hi := helper.EncodeAddr(1000)
	h, n := helper.EncodeAddr(10)
	j, i := helper.EncodeAddr(29)
	x, y := helper.EncodeAddr(10)
	program = append(program, cpu.MOVI, helper.EncodeRegs(2, 0), lo, hi)
	program = append(program, cpu.PRINT, helper.EncodeRegs(2, 0))
	program = append(program, cpu.CALL, 0x00, x, y)
	program = append(program, cpu.SUBI, helper.EncodeRegs(2, 0), h, n)
	program = append(program, cpu.JZ, helper.EncodeRegs(0, 0), j, i)
	program = append(program, cpu.HALT)
	program = append(program, cpu.RET)

	copy(mem.Data[:cpu.MemorySize], program)
	*/
	vm := cpu.NewCPU(mem)
	vm.Run()

}
