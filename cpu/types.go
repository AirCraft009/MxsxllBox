package cpu

type Register byte

const (
	R0 Register = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	NumRegisters
)

type Flags struct {
	Zero  bool
	Carry bool
}
