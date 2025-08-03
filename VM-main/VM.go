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
