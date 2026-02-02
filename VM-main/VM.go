package main

import (
	"fmt"
	"os"

	"github.com/AirCraft009/MxsxllBox/internal/IO"
	"github.com/AirCraft009/MxsxllBox/internal/MxsxllBox/cpu"
	"github.com/AirCraft009/mcc/pkg"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./MxsxllBox-main.exe path to program")
		return
	}
	fmt.Println(os.Args[1])
	mem := cpu.NewMemory()
	data, _, _, err := pkg.ReadMxBinary(os.Args[1])
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
