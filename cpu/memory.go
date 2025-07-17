package cpu

const (
	MemorySize = 64 * 1024 // 64 KB total memory

	// ProgramStart ───── Code Region (8 KB) ─────
	ProgramStart       = 0x0000
	ProgramUserEnd     = 0x1FFF // 8 KB (User + StdLib)
	ProgramStdLibStart = 0x1800 // Last 2 KB for stdlib
	ProgramEnd         = 0x1FFF

	// HeapStart ───── Heap (16 KB) ─────
	HeapStart      = 0x2000
	HeapEnd        = 0x5FFF
	HeapSize       = HeapEnd - HeapStart
	BlockSize      = 0x10
	BitMapSize     = HeapSize / BlockSize
	HeapInfoOffset = 0x10

	// StackStart ───── Stack (8 KB) ─────
	StackStart = 0x6000
	StackEnd   = 0x7FFF
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
	Data   [MemorySize]byte
	Bitmap [BitMapSize]byte
}

func (mem *Memory) AllocBlocks(blockAmmount uint16) (addr uint16) {
	/**
	This Method Allocates Blocks of 16 or What the constant BlockSize is set to
	It doesn't allow allocating anything smaller it takes the ammount of blocks to free
	*/
	var freeBlocks uint16
	size := blockAmmount * BlockSize
	if !(size <= HeapSize && blockAmmount > 0) {
		panic("cannot allocate block ammount")
	}
	start := -1
	var activeAddr uint16
	for index, open := range mem.Bitmap {
		if open == 0 {
			freeBlocks++
			if start == -1 {
				start = index
			}
		} else {
			start = -1
			freeBlocks = 0
		}
		if freeBlocks == blockAmmount {
			MetaData := uint16(HeapStart + start*BlockSize)
			activeAddr = MetaData + HeapInfoOffset
			for i := uint16(start); i < uint16(index+1); i++ {
				mem.Bitmap[i] = 1
			}
			mem.WriteWord(MetaData, blockAmmount)
			return activeAddr
		}
	}
	return 0
}

func (mem *Memory) Free(addr uint16) {
	if addr == 0 {
		return
	}
	MetaData := addr - HeapInfoOffset
	ammount := mem.ReadWord(MetaData)
	startBlock := (MetaData - HeapStart) / BlockSize
	for i := uint16(0); i < ammount; i++ {
		mem.Bitmap[startBlock+i] = 0
	}
}

func (mem *Memory) ReadByte(addr uint16) byte {
	return mem.Data[addr]
}

func (mem *Memory) ReadReg(addr uint16) (byte, byte) {
	return mem.Data[addr], mem.Data[addr+1]
}

func (mem *Memory) WriteByte(addr uint16, value byte) {
	mem.Data[addr] = value
}

func (mem *Memory) WriteWord(addr uint16, val uint16) {
	mem.Data[addr] = byte(val >> 8)
	mem.Data[addr+1] = byte(val & 0xFF)
}

func (mem *Memory) ReadWord(addr uint16) uint16 {
	hi := uint16(mem.Data[addr])
	lo := uint16(mem.Data[addr+1])
	return (hi << 8) | lo
}

func (mem *Memory) LoadProgram(program []uint16) {
	for i, word := range program {
		mem.Data[i*2] = byte(word & 0xFF) // low byte first
		mem.Data[i*2+1] = byte(word >> 8)
	}
}
