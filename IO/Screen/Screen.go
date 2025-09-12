package Screen

import (
	"MxsxllBox/VM/cpu"
	"fmt"
	"time"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	VmBlack = byte(iota)
	VmWhite
	VmRed
	VmBlue
)

const (
	Black = 0xFF000000
	White = 0xFFFFFFFF
	Red   = 0xFF0000FF
	Blue  = 0xFFFF0000
)

var (
	title         = "MxsxllBox-VM"
	Width         = int32(256)
	height        = int32(256)
	upscale       = int32(4)
	Bpp           = 2
	PpB           = 4
	transitionMap = make(map[byte]uint32)
	keyMap        = map[sdl.Keycode]byte{
		sdl.K_RETURN:    13,
		sdl.K_TAB:       9,
		sdl.K_BACKSPACE: 8,
		sdl.K_ESCAPE:    27,
		sdl.K_UP:        1,
		sdl.K_DOWN:      3,
		sdl.K_LEFT:      2,
		sdl.K_RIGHT:     4,
	}
)

type Screen struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Texture  *sdl.Texture
	LastDraw time.Time
	Keymap   map[sdl.Keycode]byte
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewScreen() *Screen {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, Width*upscale, height*upscale, sdl.WINDOW_SHOWN)
	checkError(err)
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	checkError(err)
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, Width, height)
	checkError(err)
	window.SetAlwaysOnTop(true)
	window.Raise()

	return &Screen{Window: window, Renderer: renderer, Texture: tex, LastDraw: time.Now(), Keymap: keyMap}
}

func init() {
	transitionMap[VmBlack] = Black
	transitionMap[VmWhite] = White
	transitionMap[VmRed] = Red
	transitionMap[VmBlue] = Blue
}

func ExtractPixels(VideoBuffer []byte) []uint32 {
	pixels := make([]uint32, len(VideoBuffer)*PpB)
	for i, b := range VideoBuffer {
		for p := 0; p < 4; p++ {
			isoPx := (b >> (p * 2)) & 3
			pixels[i*PpB+p] = transitionMap[isoPx]
		}
	}
	return pixels
}

func (s *Screen) Refresh(Cpu *cpu.CPU) {
	pixels := ExtractPixels(Cpu.Mem.Data[cpu.VideoStart : cpu.VideoEnd+1])
	if len(pixels) != int(Width*height) {
		fmt.Println("Pixels length mismatch", len(pixels))
		return
	}

	err := s.Texture.Update(nil, unsafe.Pointer(&pixels[0]), int(Width*4))
	checkError(err)
	err = s.Renderer.Clear()
	//copy(Cpu.Mem.Data[cpu.VideoStart:cpu.VideoEnd+1], make([]byte, cpu.VideoEnd+1-cpu.VideoStart)) delete the framebuffer for speed testing
	checkError(err)
	dst := sdl.Rect{X: 0, Y: 0, W: Width * upscale, H: height * upscale}
	err = s.Renderer.Copy(s.Texture, nil, &dst)
	checkError(err)
	s.Renderer.Present()
}
