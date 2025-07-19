package cpu

type Register uint16

const (
	NumRegisters = 32
)

type Flags struct {
	Zero  bool
	Carry bool
}
