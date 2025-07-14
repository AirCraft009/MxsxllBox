package cpu

const (
	NOP    byte = 0x00
	LOADB  byte = 0x01 // LOAD Rx, [addr]
	LOADW  byte = 0x02
	STOREB byte = 0x03 // STORE Rx, [addr]
	STOREW byte = 0x04
	ADD    byte = 0x05 // ADD Rx, Ry
	SUB    byte = 0x06 // SUB Rx, Ry
	MUL    byte = 0x07 // MUL Rx, Ry
	DIV    byte = 0x08 // DIV Rx, Ry
	JMP    byte = 0x09 // JMP addr
	JZ     byte = 0x0A // Jump if zero flag
	JC     byte = 0x0B //Jump if carry flag
	PRINT  byte = 0x0C // Print register as char
	MOVI   byte = 0x0D
	HALT   byte = 0xFF //stop the program
)
