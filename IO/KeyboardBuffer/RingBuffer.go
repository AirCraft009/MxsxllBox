package KeyboardBuffer

import "C"
import (
	"MxsxllBox/VM/cpu"
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
		lenght: cpu.RingBufferSize,
		mutex:  sync.Mutex{},
	}
}

// WriteKeyboardToBuffer : old method used with stdin
func WriteKeyboardToBuffer(Cpu *cpu.CPU) {
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

func (ringBuffer *RingBuffer) Write(char byte, Cpu *cpu.CPU) bool {
	ringBuffer.mutex.Lock()

	Cpu.Mem.WriteByte(cpu.RingBufferStart+ringBuffer.writePtr, char)
	ringBuffer.writePtr = (ringBuffer.writePtr + 1) % ringBuffer.lenght
	Cpu.Mem.WriteByte(cpu.WritePtr, byte(ringBuffer.writePtr))
	defer ringBuffer.mutex.Unlock()

	Cpu.InterruptPending = true
	Cpu.InterruptId = cpu.KeyboardInterrupt
	return true
}
