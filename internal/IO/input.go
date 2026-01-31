package IO

import (
	"MxsxllBox/internal/IO/KeyboardBuffer"
	"MxsxllBox/internal/IO/Screen"
	"MxsxllBox/internal/VM/cpu"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type ScreenKeyboardBuffer struct {
	Screen   *Screen.Screen
	keyboard *KeyboardBuffer.RingBuffer
}

func NewScreenKeyboardBuffer() *ScreenKeyboardBuffer {
	return &ScreenKeyboardBuffer{
		Screen:   Screen.NewScreen(),
		keyboard: KeyboardBuffer.NewRingBuffer(),
	}
}

func (s *ScreenKeyboardBuffer) Run(Cpu *cpu.CPU) {
	defer s.Screen.Window.Destroy()
	defer s.Screen.Renderer.Destroy()
	defer s.Screen.Texture.Destroy()

	// enable text input (so SDL2 sends TextInputEvent)
	sdl.StartTextInput()
	defer sdl.StopTextInput()

	for !Cpu.Halted {
		s.Screen.Refresh(Cpu)
		s.Screen.LastDraw = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				Cpu.Halted = true

			case *sdl.TextInputEvent:
				b := e.Text[0]
				if b < 128 { // keep ASCII only
					s.keyboard.Write(byte(b), Cpu)
				}

			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					if val, ok := s.Screen.Keymap[e.Keysym.Sym]; ok {
						s.keyboard.Write(val, Cpu)
					}
				}
			}
			sdl.Delay(15)
		}
	}
}
