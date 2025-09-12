package main

import (
	"MxsxllBox/Assembly-process/linker"
	"MxsxllBox/IO"
	cpu2 "MxsxllBox/VM/cpu"
	"fmt"
)

func main() {
	mem := cpu2.NewMemory()

	copy(mem.Data[:], linker.CompileForOs("program.asm", "MxsxllOS"))
	vm := cpu2.NewCPU(mem)

	inputs := IO.NewScreenKeyboardBuffer()
	fmt.Println("Program started")
	go cpu2.InitTicker(vm)
	go vm.Run()
	inputs.Run(vm)
}
