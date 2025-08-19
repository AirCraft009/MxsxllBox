package Screen

import (
	"MxsxllBox/VM/cpu"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// for now 2bpp
const (
	VM_BLACK = byte(iota)
	VM_WHITE
	VM_RED
	VM_BLUE
)

const (
	SCREEN_BLACK = uint32(0x000000FF)
	SCREEN_WHITE = uint32(0xFFFFFFFF)
	SCREEN_RED   = uint32(0xFF0000FF)
	SCREEN_BLUE  = uint32(0x0000FFFF)
)

var (
	title         = "MxsxllBox-VM"
	width         = int32(256)
	height        = int32(256)
	upscale       = 4
	Bpp           = 2
	PpB           = 4
	transitionMap = make(map[byte]uint32)
	renderer      *sdl.Renderer
	tex           *sdl.Texture
)

type Screen struct {
	Window   *sdl.Window
	LastDraw time.Time
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewScreen() *Screen {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	checkError(err)
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkError(err)
	tex = renderer.CreateTexture()
	return &Screen{Window: window, LastDraw: time.Now()}
}

func init() {
	transitionMap[VM_BLACK] = SCREEN_BLACK
	transitionMap[VM_WHITE] = SCREEN_WHITE
	transitionMap[VM_RED] = SCREEN_RED
	transitionMap[VM_BLUE] = SCREEN_BLUE
}

func ExtractPixels(VideoBuffer []byte) (pixels []uint32) {
	for i := 0; i < len(VideoBuffer)*PpB; i++ {
		pixels = append(pixels, transitionMap[VideoBuffer[i]])
	}
	return pixels
}

// Refresh doesn't really do anything only refreshes so the image changes redrawing is handled by the Vm
func (s *Screen) Refresh(Cpu *cpu.CPU, tex *sdl.Texture) {
	pixels := ExtractPixels(Cpu.Mem.Data[cpu.VideoStart:cpu.VideoEnd])
	err := tex.Update(nil, unsafe.Pointer(&pixels), upscale)
	checkError(err)
}

func (s *Screen) Run(Cpu *cpu.CPU) {

}
