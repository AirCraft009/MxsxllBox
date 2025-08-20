package main

import (
	"MxsxllBox/Assembly-process/linker"
	"MxsxllBox/IO/KeyboardBuffer"
	"MxsxllBox/IO/Screen"
	cpu2 "MxsxllBox/VM/cpu"
	"fmt"
)

func main() {
	mem := &cpu2.Memory{}

	copy(mem.Data[:], linker.CompileForOs("program.asm", "EchoKeys"))
	vm := cpu2.NewCPU(mem)

	screen := Screen.NewScreen()
	go KeyboardBuffer.WriteKeyboardToBuffer(vm)
	go screen.Run(vm)
	fmt.Println("Program started")
	go cpu2.InitTicker(vm)
	go vm.Run()

	select {}
}
