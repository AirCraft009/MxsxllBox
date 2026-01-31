package main

import (
	"MxsxllBox/internal/IO"
	"MxsxllBox/internal/VM/cpu"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./VM-main.exe path to program")
		return
	}
	fmt.Println(os.Args[1])
	mem := cpu.NewMemory()
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
		return
	}
	copy(mem.Data[:], data)
	vm := cpu.NewCPU(mem)

	inputs := IO.NewScreenKeyboardBuffer()
	fmt.Println("Program started")
	go cpu.InitTicker(vm)
	go vm.Run()
	inputs.Run(vm)
}
