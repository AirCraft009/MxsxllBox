package cpu

type debugCpu struct {
	Cpu         *CPU
	debugTicker int
}

func NewDebugCpu(mem *Memory) *debugCpu {
	return &debugCpu{
		Cpu:         NewCPU(mem),
		debugTicker: 0,
	}
}

func (debugcpu *debugCpu) StepDebug() {
	debugcpu.Cpu.Step()
	debugcpu.debugTicker++
	if debugcpu.debugTicker == 10000 {
		debugcpu.debugTicker = 0
		debugcpu.Cpu.InterruptPending = true
		debugcpu.Cpu.InterruptId = TimerInterrupt
	}
}
