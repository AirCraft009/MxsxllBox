package cpu

import "fmt"

type CPU struct {
	Registers [NumRegisters]uint16
	PC        uint16
	SP        uint16
	Flags     Flags
	Mem       *Memory
	Halted    bool
	Handlers  map[byte]func(cpu *CPU, instruction *HandlerInstructions)
}

func NewCPU(mem *Memory) *CPU {
	cpu := &CPU{
		Mem:      mem,
		SP:       MemorySize - 1, // stack grows downward
		Handlers: make(map[byte]func(cpu *CPU, instruction *HandlerInstructions)),
	}

	cpu.Handlers[NOP] = handleNop
	cpu.Handlers[LOADB] = handleLoadB
	cpu.Handlers[LOADW] = handleLoadW
	cpu.Handlers[STOREB] = handleStoreB
	cpu.Handlers[STOREW] = handleStoreW
	cpu.Handlers[ADD] = handleAdd
	cpu.Handlers[SUB] = handleSub
	cpu.Handlers[MUL] = handleMul
	cpu.Handlers[DIV] = handleDiv
	cpu.Handlers[JMP] = handleJmp
	cpu.Handlers[JZ] = handleJz
	cpu.Handlers[JC] = handleJc
	cpu.Handlers[PRINT] = handlePrint
	cpu.Handlers[HALT] = handleHalt
	cpu.Handlers[MOVI] = handleMovi

	return cpu
}

func (cpu *CPU) Step() {
	opCode, instructions := getInstruction(cpu)
	if handler, ok := cpu.Handlers[opCode]; ok {
		handler(cpu, instructions)
	} else {
		panic(fmt.Sprintf("unknown opcode: %d", opCode))
	}
}

func (cpu *CPU) Run() {
	for !cpu.Halted {
		cpu.Step()
	}
}
