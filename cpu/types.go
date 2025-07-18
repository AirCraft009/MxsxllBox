package cpu

type Register uint16

const (
	outputRegister          = O1
	R0             Register = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	O1
	O2
	O3
	O4
	NumRegisters
)

type Flags struct {
	Zero  bool
	Carry bool
}
