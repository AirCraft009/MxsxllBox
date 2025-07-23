package KeyboardBuffer

import "C"
import (
	"MxsxllBox/cpu"
	"golang.org/x/term"
	"os"
	"sync"
)

type RingBuffer struct {
	writePtr uint16
	readPtr  uint16
	lenght   uint16
	mutex    sync.Mutex
}

func newRingBuffer() *RingBuffer {
	return &RingBuffer{
		lenght: cpu.RingBufferSize,
		mutex:  sync.Mutex{},
	}
}

func WriteKeyboardToBuffer(Cpu *cpu.CPU) {
	ringBuffer := newRingBuffer()
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 1)

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}
		if buf[0] < 32 || buf[0] > 126 {
			buf = make([]byte, 1)
			continue
		}
		ringBuffer.write(buf[0], Cpu)
		Cpu.InterruptPending = true
		Cpu.InterruptId = cpu.KeyboardInterrupt
	}
}

func (ringBuffer *RingBuffer) write(char byte, Cpu *cpu.CPU) bool {
	ringBuffer.mutex.Lock()
	if byte((ringBuffer.writePtr+1)%ringBuffer.lenght) == Cpu.Mem.ReadByte(cpu.ReadPtr) {
		return false
	}

	Cpu.Mem.WriteByte(cpu.RingBufferStart+ringBuffer.writePtr, char)
	ringBuffer.writePtr = (ringBuffer.writePtr + 1) % ringBuffer.lenght
	Cpu.Mem.WriteByte(cpu.WritePtr, byte(ringBuffer.writePtr))
	defer ringBuffer.mutex.Unlock()
	return true
}
