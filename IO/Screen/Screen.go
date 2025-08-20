package Screen

import (
	"MxsxllBox/VM/cpu"
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
	width         = int32(256)
	height        = int32(256)
	upscale       = int32(4)
	Bpp           = 2
	PpB           = 4
	transitionMap = make(map[byte]uint32)
	renderer      *sdl.Renderer
	tex           *sdl.Texture
)

type Screen struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Texture  *sdl.Texture
	LastDraw time.Time
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewScreen() *Screen {
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width*upscale, height*upscale, sdl.WINDOW_SHOWN)
	checkError(err)
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkError(err)
	tex, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	checkError(err)
	return &Screen{Window: window, Renderer: renderer, Texture: tex, LastDraw: time.Now()}
}

func init() {
	transitionMap[VmBlack] = Black
	transitionMap[VmWhite] = White
	transitionMap[VmRed] = Red
	transitionMap[VmBlue] = Blue
}

func ExtractPixels(VideoBuffer []byte) (pixels []uint32) {
	for i, j, p := 0, 0, 0; i < len(VideoBuffer)*PpB; i++ {
		isoPx := (VideoBuffer[j] >> (p * 2)) & 3
		pixels = append(pixels, transitionMap[isoPx])
		p = (p + 1) % PpB
		if p == 0 {
			j++
		}
	}
	return pixels
}

func (s *Screen) Refresh(Cpu *cpu.CPU) {
	pixels := ExtractPixels(Cpu.Mem.Data[cpu.VideoStart:cpu.VideoEnd])
	err := s.Texture.Update(nil, unsafe.Pointer(&pixels[0]), int(width*4))
	checkError(err)
	err = s.Renderer.Clear()
	checkError(err)
	dst := sdl.Rect{X: 0, Y: 0, W: width * upscale, H: height * upscale}
	err = s.Renderer.Copy(s.Texture, nil, &dst)
	checkError(err)
	s.Renderer.Present()
}

func (s *Screen) Run(Cpu *cpu.CPU) {
	defer s.Window.Destroy()
	defer s.Renderer.Destroy()
	defer s.Texture.Destroy()
	for !Cpu.Halted {
		s.Refresh(Cpu)
		s.LastDraw = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				Cpu.Halted = true
				break
			}
		}
		sdl.Delay(33)
	}
}
