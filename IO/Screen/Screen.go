package Screen

import (
	"MxsxllBox/VM/cpu"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// for now 2bpp
const (
	BLACK = iota
	WHITE
	RED
	BLUE
)

type Screen struct {
	Window   *sdl.Window
	LastDraw time.Time
}

func NewScreen(window *sdl.Window) *Screen {
	return &Screen{Window: window, LastDraw: time.Now()}
}

// Refresh doesn't really do anything only refreshes so the image changes redrawing is handled by the Vm
func (s *Screen) Refresh(Cpu *cpu.CPU) {

}
