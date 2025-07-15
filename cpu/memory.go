package cpu

const (
	MemorySize = 16 * 1024 // 16 KB total memory

	ProgramStart       = 0x0000
	programUserEnd     = 0x0C00
	ProgramStdLibStart = 0x0C01
	ProgramEnd         = 0x0FFF // 4 KB for program (0x0000–0x0FFF)

	HeapStart      = 0x1000
	HeapEnd        = 0x1FFF // 4 KB heap (0x1000–0x1FFF)
	HeapSize       = HeapEnd - HeapStart
	BlockSize      = 0x10
	BitMapSize     = HeapSize / BlockSize
	HeapInfoOffset = 16

	StackStart = 0x2000
	StackEnd   = 0x27FF // 2 KB stack (0x2000–0x27FF)
	StackInit  = 0x2800 // stack pointer start address (just past stack end)

	VideoStart = 0x2800
	VideoEnd   = 0x3FFF // 5 KB video RAM (0x2800–0x3FFF)

	ReservedStart = 0x4000
	ReservedEnd   = 0x43FF // 1 KB reserved for IO, keyboard buffers, etc.

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
