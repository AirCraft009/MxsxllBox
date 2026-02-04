package main

import (
	"fmt"
	"os"

	"github.com/AirCraft009/MxsxllBox/internal/IO"
	"github.com/AirCraft009/MxsxllBox/internal/MxsxllBox/cpu"
	"github.com/AirCraft009/MxsxllBox/internal/debugging"
	"github.com/AirCraft009/mcc/pkg"
	flag "github.com/spf13/pflag"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./MxsxllBox path to program [flags]")
		return
	}

	fs := flag.NewFlagSet("MxsxllBox", flag.ExitOnError)
	debug := fs.Bool("debug", false, "launch with a debug window and no Normal Window")

	fs.Parse(os.Args[1:])

	args := fs.Args()
	inputF := args[0]
	if *debug {
		debugging.DebugStart(inputF)
	}
	mem := cpu.NewMemory()
	data, _, _, err := pkg.ReadMxBinary(inputF)
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
