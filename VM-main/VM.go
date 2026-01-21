package main

import (
	"MxsxllBox/IO"
	cpu2 "MxsxllBox/VM/cpu"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./VM-main.exe path to program")
		return
	}
	fmt.Println(os.Args[1])
	mem := cpu2.NewMemory()
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
		return
	}
	copy(mem.Data[:], data)
	vm := cpu2.NewCPU(mem)

	inputs := IO.NewScreenKeyboardBuffer()
	fmt.Println("Program started")
	go cpu2.InitTicker(vm)
	go vm.Run()
	inputs.Run(vm)
}
