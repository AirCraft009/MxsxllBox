package cpu

import "unsafe"

const (
	MemorySize = 16 * 1024 // 16 KB total memory

	ProgramStart = 0x0000
	ProgramEnd   = 0x0FFF // 4 KB for program (0x0000–0x0FFF)

	HeapStart = 0x1000
	HeapEnd   = 0x1FFF // 4 KB heap (0x1000–0x1FFF)

	StackStart = 0x2000
	StackEnd   = 0x27FF // 2 KB stack (0x2000–0x27FF)
	StackInit  = 0x2800 // stack pointer start address (just past stack end)

	VideoStart = 0x2800
	VideoEnd   = 0x3FFF // 5 KB video RAM (0x2800–0x3FFF)

	ReservedStart = 0x4000
	ReservedEnd   = 0x43FF // 1 KB reserved for IO, keyboard buffers, etc.
)

type Memory struct {
	Data  [MemorySize]byte
	HeapP uint16
}

func (mem *Memory) Alloc(size uint16) uint16 {
	addr := mem.HeapP
	if mem.HeapP+size+instructionSizeShort > HeapEnd {
		panic("Heap overflow")
	}
	return addr
}

func (mem *Memory) Free(ptr unsafe.Pointer) bool {
	return false
}

func NewMemory() *Memory {
	return &Memory{
		HeapP: HeapStart,
	}
}

func (mem *Memory) ReadByte(addr uint16) byte {
	return mem.Data[addr]
}

func (mem *Memory) WriteByte(addr uint16, value byte) {
	mem.Data[addr] = value
}

func (m *Memory) WriteWord(addr uint16, val uint16) {
	m.Data[addr] = byte(val >> 8)
	m.Data[addr+1] = byte(val & 0xFF)
}

func (m *Memory) ReadWord(addr uint16) uint16 {
	hi := uint16(m.Data[addr])
	lo := uint16(m.Data[addr+1])
	return (hi << 8) | lo
}

func (mem *Memory) LoadProgram(program []uint16) {
	for i, word := range program {
		mem.Data[i*2] = byte(word & 0xFF) // low byte first
		mem.Data[i*2+1] = byte(word >> 8)
	}
}
