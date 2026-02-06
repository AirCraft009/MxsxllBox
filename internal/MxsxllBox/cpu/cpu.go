package cpu

import (
	"fmt"
	"os"
	"time"

	"github.com/AirCraft009/MxsxllBox/internal/helper"
)

type id uint16

const (
	NumRegisters             = 64
	JmpOffset                = 5
	InterruptHandlerLocation = 23965
	I1                       = 24
)

const (
	KeyboardInterrupt id = (1 + iota) * JmpOffset
	TimerInterrupt
)

func (cpu *CPU) Step() {
	if cpu.InterruptPending && helper.IsInterruptActivated(int(cpu.InterruptId), cpu.InterruptMask) {
		cpu.Interrupt = true
	}

	opCode, instructions := getInstruction(cpu)
	if handler, ok := cpu.Handlers[opCode]; ok {
		handler(cpu, instructions)
	} else {
		panic(fmt.Sprintf("unknown opcode: %d", opCode))
	}

	if cpu.Interrupt {
		cpu.Registers[I1] = uint16(cpu.InterruptId)
		*cpu.SP -= 2
		cpu.Mem.WriteWordKeyboardSafe(*cpu.SP, *cpu.PC)
		*cpu.PC = InterruptHandlerLocation
		cpu.Mutex.Lock()
		cpu.InterruptPending = false
		cpu.Interrupt = false
		cpu.Mutex.Unlock()
	}
}

func (cpu *CPU) Run() {
	tStart := time.Now()
	for !cpu.Halted {
		cpu.Step()
	}
	tEnd := time.Now()
	dur := tEnd.Sub(tStart)
	fmt.Printf("The Proccess took %f seconds\n", dur.Seconds())
	cpu.HardwareTimer.Stop()
	os.Exit(0)
}
