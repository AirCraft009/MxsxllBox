package cpu

import (
	"sync"
	"time"

	"github.com/AirCraft009/mcc/pkg"
)

type Register uint16
type Flags struct {
	Zero  bool
	Carry bool
}
type CPU struct {
	Registers         [NumRegisters]uint16
	PC                *uint16
	SP                *uint16
	Flags             Flags
	Mem               *Memory
	Halted            bool
	Handlers          map[byte]func(cpu *CPU, instruction *HandlerInstructions)
	Mutex             sync.Mutex
	InterruptPending  bool
	InterruptId       id
	Interrupt         bool //makes sure that the interrupt is executed after the next step
	HardwareTimer     *time.Ticker
	PrevInterruptMask byte
	InterruptMask     byte
}

func InitTicker(cpu *CPU) {
	cpu.HardwareTimer = time.NewTicker(15 * time.Millisecond)
	for {
		select {
		case _ = <-cpu.HardwareTimer.C:

			cpu.InterruptPending = true
			cpu.InterruptId = TimerInterrupt
		}
	}
}

func NewCPU(mem *Memory) *CPU {
	cpu := &CPU{
		Mem:      mem, // stack grows downward
		Handlers: make(map[byte]func(cpu *CPU, instruction *HandlerInstructions)),
	}
	cpu.Registers[pkg.SPRegister] = pkg.StackInit
	cpu.SP = &cpu.Registers[pkg.SPRegister]
	cpu.PC = &cpu.Registers[pkg.PCRegister]

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
	cpu.Handlers[MOVA] = handleMova
	cpu.Handlers[GPC] = handleGPc
	cpu.Handlers[SPC] = handleSPc
	cpu.Handlers[GSP] = handleGSp
	cpu.Handlers[SSP] = handleSSp
	cpu.Handlers[GRFN] = handleGrfn
	cpu.Handlers[GF] = handleGf
	cpu.Handlers[SF] = handleSf
	cpu.Handlers[SRFN] = handleSrfn
	cpu.Handlers[YIELD] = handleYield
	cpu.Handlers[UNYIELD] = handleUnyield
	cpu.Handlers[STINTI] = handleSTINTI
	cpu.Handlers[STINT] = handleSTINT
	cpu.Handlers[XOR] = handleXor
	cpu.Handlers[DRAWPX] = handleStorePixel
	cpu.Handlers[STOREBLOCK] = handleStoreSection

	return cpu
}
