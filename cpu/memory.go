package cpu

import "sync"

const (
	MemorySize = 64 * 1024 // 64 KB total memory

	// ProgramStart ───── Code Region (8 KB) ─────
	ProgramStart       = 0x0000
	ProgramUserEnd     = 0x1FFF // 8 KB (User + StdLib)
	ProgramStdLibStart = 0x1800 // Last 2 KB for stdlib
	ProgramEnd         = 0x1FFF

	// HeapStart ───── Heap (16 KB) ─────
	HeapStart          = 0x2000
	HeapEnd            = 0x6000
	writeableHeapStart = 9628
	HeapSize           = HeapEnd - writeableHeapStart
	BlockSize          = 0x10

	// StackStart ───── Stack (8 KB) ─────
	StackStart = 0x600
	StackEnd   = 0x7FFF
	St         = StackInit + 910
	StackInit  = StackEnd + 1 // 0x8000 (Stack grows down)

	// VideoStart ───── Video RAM / Framebuffer (16 KB) ─────
	VideoStart = 0x8000
	VideoEnd   = 0xBFFF

	// KeyboardStart ReservedStart ───── Reserved for IO / Buffers / MMIO (8 KB) ─────
	KeyboardStart   = 0xC000
	ReadPtr         = 0xC000
	WritePtr        = 0xC001
	RingBufferStart = 0xC002
	RingBufferEnd   = 0xC020 //N = 30
	RingBufferSize  = RingBufferEnd - RingBufferStart
	ReservedEnd     = 0xDFFF

	// ExtraStart ───── Unused / Future Expansion / Paging Tables / Filesystem etc (8 KB) ─────
	ExtraStart = 0xE000
	ExtraEnd   = 0xFFFF
)

type Memory struct {
	Data       [MemorySize]byte
	keyboardMu sync.Mutex
}

func isKeyboardRegion(addr uint16) bool {
	return addr >= 0xC000 && addr <= 0xC020
}

func (mem *Memory) ReadByte(addr uint16) byte {
	if isKeyboardRegion(addr) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}
	return mem.Data[addr]
}

func (mem *Memory) ReadWord(addr uint16) uint16 {
	if isKeyboardRegion(addr) || isKeyboardRegion(addr+1) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}
	hi := uint16(mem.Data[addr])
	lo := uint16(mem.Data[addr+1])
	return (hi << 8) | lo
}

func (mem *Memory) ReadReg(addr uint16) (byte, byte) {
	if isKeyboardRegion(addr) || isKeyboardRegion(addr+1) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}
	return mem.Data[addr], mem.Data[addr+1]
}

func (mem *Memory) WriteByte(addr uint16, value byte) {
	if isKeyboardRegion(addr) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}
	mem.Data[addr] = value
}

func (mem *Memory) WriteWord(addr uint16, val uint16) {
	if isKeyboardRegion(addr) || isKeyboardRegion(addr+1) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}
	mem.Data[addr] = byte(val >> 8)
	mem.Data[addr+1] = byte(val & 0xFF)
}

func (mem *Memory) LoadProgram(program []uint16) {
	for i, word := range program {
		mem.Data[i*2] = byte(word & 0xFF) // low byte first
		mem.Data[i*2+1] = byte(word >> 8)
	}
}
