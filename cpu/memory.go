package cpu

import "unsafe"

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
