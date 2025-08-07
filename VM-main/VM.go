package main

import (
	"MxsxllBox/Assembly-process/linker"
	"MxsxllBox/IO/KeyboardBuffer"
	cpu2 "MxsxllBox/VM/cpu"
	"fmt"
	"runtime/debug"
)

func main() {
	mem := &cpu2.Memory{}

	copy(mem.Data[:], linker.CompileForOs("program.asm", "EchoKeys"))
	vm := cpu2.NewCPU(mem)
	go func() {
		if r := recover(); r != nil {
			fmt.Println("Program crashed with panic:", r)
			fmt.Printf("PC, OpCode: %d, %d\n", vm.PC, vm.Mem.Data[vm.PC])
			fmt.Printf("stack pointer: %d\n")
			fmt.Printf("stack trace: %s\n", string(debug.Stack()))
		}
	}()
	go KeyboardBuffer.WriteKeyboardToBuffer(vm)
	fmt.Println("Program started")
	go vm.Run()

	select {}
}
