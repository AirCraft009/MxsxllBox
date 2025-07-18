package cpu

import (
	"fmt"
	"sync"
)

type CPU struct {
	Registers  [NumRegisters]uint16
	PC         uint16
	SP         uint16
	Flags      Flags
	Mem        *Memory
	Halted     bool
	Handlers   map[byte]func(cpu *CPU, instruction *HandlerInstructions)
	Tasks      []*Task
	ActiveTask uint16
	NewId      int
	Mutex      sync.Mutex
}

func NewCPU(mem *Memory) *CPU {
	cpu := &CPU{
		Mem:      mem,
		SP:       StackInit, // stack grows downward
		Handlers: make(map[byte]func(cpu *CPU, instruction *HandlerInstructions)),
		Tasks:    make([]*Task, 0),
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
	cpu.Handlers[ADDI] = handleAddi
	cpu.Handlers[DIVI] = handleDivi
	cpu.Handlers[SUBI] = handleSubi
	cpu.Handlers[MULI] = handleMuli
	cpu.Handlers[STORE] = handleStore
	cpu.Handlers[LOAD] = handleLoad
	cpu.Handlers[PUSH] = handlePush
	cpu.Handlers[POP] = handlePop
	cpu.Handlers[CALL] = handleCall
	cpu.Handlers[RET] = handleRet
	cpu.Handlers[PRINTSTR] = handlePrintstr
	cpu.Handlers[JNZ] = handleJnz
	cpu.Handlers[JNC] = handleJnc
	cpu.Handlers[CMP] = handleCmp
	cpu.Handlers[CMPI] = handleCmpi
	cpu.Handlers[TEST] = handleTest
	cpu.Handlers[TSTI] = handleTsti
	cpu.Handlers[JL] = handleJL
	cpu.Handlers[JLE] = handleJLE
	cpu.Handlers[JG] = handleJG
	cpu.Handlers[JGE] = handleJGE
	cpu.Handlers[STZ] = handleSTZ
	cpu.Handlers[STC] = handleSTC
	cpu.Handlers[CLZ] = handleCLZ
	cpu.Handlers[CLC] = handleCLC
	cpu.Handlers[MOV] = handleMov
	cpu.Handlers[MOD] = handleMod
	cpu.Handlers[MODI] = handleModi
	cpu.Handlers[RS] = handleRs
	cpu.Handlers[LS] = handleLs
	cpu.Handlers[AND] = handleAnd
	cpu.Handlers[OR] = handleOr

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
