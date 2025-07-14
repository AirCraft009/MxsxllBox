package main

import "MxsxllBox/cpu"

func main() {
	mem := &cpu.Memory{}

	vm := cpu.NewCPU(mem)
	vm.Run()
}
