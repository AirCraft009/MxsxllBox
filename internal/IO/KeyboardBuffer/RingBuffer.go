package KeyboardBuffer

import "C"
import (
	cpu2 "MxsxllBox/internal/VM/cpu"
	"os"
	"sync"

	"golang.org/x/term"
)

type RingBuffer struct {
	writePtr uint16
	readPtr  uint16
	lenght   uint16
	mutex    sync.Mutex
}

func NewRingBuffer() *RingBuffer {
	return &RingBuffer{
		lenght: cpu2.RingBufferSize,
		mutex:  sync.Mutex{},
	}
}

// WriteKeyboardToBuffer : old method used with stdin
func WriteKeyboardToBuffer(Cpu *cpu2.CPU) {
	ringBuffer := NewRingBuffer()

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
		if buf[0] > 127 {
			buf = make([]byte, 1)
			continue
		}
		ringBuffer.Write(buf[0], Cpu)
	}
}

func (ringBuffer *RingBuffer) Write(char byte, Cpu *cpu2.CPU) bool {
	ringBuffer.mutex.Lock()

	Cpu.Mem.WriteByte(cpu2.RingBufferStart+ringBuffer.writePtr, char)
	ringBuffer.writePtr = (ringBuffer.writePtr + 1) % ringBuffer.lenght
	Cpu.Mem.WriteByte(cpu2.WritePtr, byte(ringBuffer.writePtr))
	defer ringBuffer.mutex.Unlock()

	Cpu.InterruptPending = true
	Cpu.InterruptId = cpu2.KeyboardInterrupt
	return true
}
