package cpu

import (
	"os"
	"sync"

	"github.com/AirCraft009/mcc/pkg"
)

type Memory struct {
	Data       [pkg.MemorySize]byte
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
		Data:       [pkg.MemorySize]byte{},
		EEPROM:     data,
		keyboardMu: sync.Mutex{},
	}
}

func isKeyboardRegion(addr uint16) bool {
	return addr >= 0xC000 && addr <= 0xC020
}
func isStackRegion(addr uint16) bool {
	return addr >= pkg.StackStart && addr <= pkg.StackEnd
}

func isCodeRegion(addr uint16) bool {
	return addr <= pkg.ProgramEnd
}

// the stack is memMapped to a rom section for strings

func (mem *Memory) ReadByte(addr uint16) byte {
	if isCodeRegion(addr) {
		return mem.EEPROM[addr]
	}
	if isStackRegion(addr) {
		return mem.EEPROM[pkg.ProgramStart+uint16(addr-pkg.StackStart)]
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
		hi := mem.EEPROM[pkg.ProgramStart+uint16(addr-pkg.StackStart)]
		lo := mem.EEPROM[pkg.ProgramStart+(addr+1)-uint16(pkg.StackStart)]
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
