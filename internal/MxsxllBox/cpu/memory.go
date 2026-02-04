package cpu

import (
	"os"
	"sync"
)

const (
	HardDriveSize = 2 * 1073741824 //2  GB
	EEPROM        = 16 * 1024      // the rom is 8 KB of storage I'm mapping to
)

const (
	MemorySize = 1024 * 64 // 64 KB total memory

	// ProgramStart ───── Code Region (8 KB) ─────
	ProgramStart       = 0x0000
	ProgramUserEnd     = 0x17FF // 8 KB (User + StdLib)
	ProgramStdLibStart = 0x1800 // Last 2 KB for stdlib
	ProgramEnd         = 0x1FFF

	// HeapStart ───── Heap (16 KB) ─────
	HeapStart          = 0x2000
	HeapEnd            = 0x6000
	writeableHeapStart = 9628
	writeableHeapEnd   = 23964
	HeapSize           = writeableHeapEnd - writeableHeapStart
	Interrupttable     = 23965
	InterruptTableSIze = HeapEnd - Interrupttable
	BlockSize          = 0x10

	// StackStart ─────  (8 KB) ─────
	StackStart = 0x6000
	StackEnd   = 0x7FFF
	StackInit  = StackEnd

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

	// ExtraStart ───── Unused / Future Expansion / Paging Tables / Filesystem etc (≈16KB KB) ─────
	ExtraStart          = 0xC021
	VideoCharTableStart = 0xC021
	VideoCharTableEnd   = VideoCharTableStart + 1024*2
	InputStringLen      = VideoCharTableEnd + 1
	InputStringMain     = InputStringLen + 2
	InputStringMainEnd  = InputStringMain + 64
	ExtraEnd            = 0xFFFF
	ExtraSize           = ExtraEnd - ExtraStart
)

type Memory struct {
	Data       [MemorySize]byte
	EEPROM     []byte
	keyboardMu sync.Mutex
}

func NewMemory() *Memory {
	wd, err := os.Getwd()
	data, err := os.ReadFile(wd + "/internal/MxsxllBox/cpu/EEPROM.mem")
	if err != nil {
		panic(err)
	}
	return &Memory{
		Data:       [MemorySize]byte{},
		EEPROM:     data,
		keyboardMu: sync.Mutex{},
	}
}

func isKeyboardRegion(addr uint16) bool {
	return addr >= 0xC000 && addr <= 0xC020
}
func isStackRegion(addr uint16) bool {
	return addr >= StackStart && addr <= StackEnd
}

func isCodeRegion(addr uint16) bool {
	return addr <= ProgramEnd
}

// the stack is memMapped to a rom section for strings

func (mem *Memory) ReadByte(addr uint16) byte {
	if isCodeRegion(addr) {
		return mem.EEPROM[addr]
	}
	if isStackRegion(addr) {
		return mem.EEPROM[ProgramStart+uint16(addr-StackStart)]
	}
	return mem.ReadByteKeyboardSafe(addr)
}

func (mem *Memory) ReadWord(addr uint16) uint16 {
	if isCodeRegion(addr) {
		hi := mem.EEPROM[addr]
		lo := mem.EEPROM[addr+1]
		return uint16(hi)<<8 | uint16(lo)
	}
	if isStackRegion(addr) || isStackRegion(addr+1) {
		hi := mem.EEPROM[ProgramStart+uint16(addr-StackStart)]
		lo := mem.EEPROM[ProgramStart+(addr+1)-uint16(StackStart)]
		return uint16(hi)<<8 | uint16(lo)
	}
	return mem.ReadWordKeyboardSafe(addr)
}

func (mem *Memory) WriteByte(addr uint16, value byte) {
	if isCodeRegion(addr) {
		return
	}
	if isStackRegion(addr) {
		return
	}
	mem.WriteByteKeyboardSafe(addr, value)
}

func (mem *Memory) WriteWord(addr uint16, val uint16) {
	if isCodeRegion(addr) {
		return
	}
	if isStackRegion(addr) || isStackRegion(addr+1) {
		return
	}
	mem.WriteWordKeyboardSafe(addr, val)
}

// ReadByteKeyboardSafe
//
// reads a byte from addr with the constraints:
// that if addr is in the keyboard buffer region
// the keyboardMutex is locked
func (mem *Memory) ReadByteKeyboardSafe(addr uint16) byte {
	if isKeyboardRegion(addr) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}
	return mem.Data[addr]
}

func (mem *Memory) ReadWordKeyboardSafe(addr uint16) uint16 {
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

func (mem *Memory) WriteByteKeyboardSafe(addr uint16, value byte) {
	if isKeyboardRegion(addr) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}

	mem.Data[addr] = value
}

func (mem *Memory) WriteWordKeyboardSafe(addr uint16, val uint16) {
	if isKeyboardRegion(addr) || isKeyboardRegion(addr+1) {
		mem.keyboardMu.Lock()
		defer mem.keyboardMu.Unlock()
	}
	mem.Data[addr] = byte(val >> 8)
	mem.Data[addr+1] = byte(val & 0xFF)
}
