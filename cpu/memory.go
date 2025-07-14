package cpu

const (
	MemorySize   = 4096
	ProgramStart = 0x0000
	ProgramEnd   = 0x03FF

	HeapStart = 0x0400
	HeapEnd   = 0x05FF

	VideoStart = 0x0600
	VideoEnd   = 0x0BFF

	StackStart = 0x0C00
	StackEnd   = 0x0FFF
	StackInit  = 0x1000
)

type Memory struct {
	Data    [MemorySize]byte
	Program uint16
	Video   uint16
	Stack   uint16
}

func (mem *Memory) Read(addr uint16) byte {
	return mem.Data[addr]
}

func (mem *Memory) Write(addr uint16, value byte) {
	mem.Data[addr] = value
}
