package cpu

type DebugCpu struct {
	Cpu         *CPU
	debugTicker int
}

func NewDebugCpu(mem *Memory) *DebugCpu {
	return &DebugCpu{
		Cpu:         NewCPU(mem),
		debugTicker: 0,
	}
}

func (debugcpu *DebugCpu) StepDebug() {
	debugcpu.Cpu.Step()
	debugcpu.debugTicker++
	if debugcpu.debugTicker == 10000 {
		debugcpu.debugTicker = 0
		debugcpu.Cpu.InterruptPending = true
		debugcpu.Cpu.InterruptId = TimerInterrupt
	}
}
