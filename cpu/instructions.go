package cpu

const (
	NOP      byte = 0x00
	LOADB    byte = 0x01 // LOAD Rx, [addr]
	LOADW    byte = 0x02
	STOREB   byte = 0x03 // STORE Rx, [addr]
	STOREW   byte = 0x04
	ADD      byte = 0x05 // ADD Rx, Ry
	SUB      byte = 0x06 // SUB Rx, Ry
	MUL      byte = 0x07 // MUL Rx, Ry
	DIV      byte = 0x08 // DIV Rx, Ry
	JMP      byte = 0x09 // JMP addr
	JZ       byte = 0x0A // Jump if zero flag
	JC       byte = 0x0B //Jump if carry flag
	PRINT    byte = 0x0C // Print register as char
	MOVI     byte = 0x0D // MOVI Regs, immideate
	ADDI     byte = 0x0E
	SUBI     byte = 0x0F
	MULI     byte = 0x10
	DIVI     byte = 0x11
	LOAD     byte = 0x12
	STORE    byte = 0x13
	PUSH     byte = 0x14
	POP      byte = 0x15
	CALL     byte = 0x16
	RET      byte = 0x17
	ALLOC    byte = 0x18 // now implemented in stdlib kept to avoid refactoring
	FREE     byte = 0x19 //  -||-
	PRINTSTR byte = 0x1A
	JNZ      byte = 0x1B
	JNC      byte = 0x1C
	CMP      byte = 0x1D
	CMPI     byte = 0x1E
	TEST     byte = 0x1F
	TSTI     byte = 0x20
	JL       byte = 0x21
	JLE      byte = 0x22
	JG       byte = 0x23
	JGE      byte = 0x24
	STZ      byte = 0x25
	STC      byte = 0x26
	CLZ      byte = 0x27
	CLC      byte = 0x28
	MOD      byte = 0x29
	MOV      byte = 0x2A
	MODI     byte = 0x02B
	RS       byte = 0x2C
	LS       byte = 0x2D
	OR       byte = 0x2E
	AND      byte = 0x2F
	MOVA     byte = 0x30 //MOVA Reg Lbl-name
	GPC      byte = 0x31 // Get Program Counter
	SPC      byte = 0x32 // Set Program Counter
	GSP      byte = 0x33 // Get Stack Pointer
	SSP      byte = 0x34 // Set Stack Pointer
	GRFN     byte = 0x35 // Get Register from number
	GF       byte = 0x36 // Get Flags bit 1 = carry  bit 0 = zero
	SF       byte = 0x37 // Set Flags bit 1 = catty	bit 0 = zero
	SRFN     byte = 0x38
	HALT     byte = 0xFF //stop the program
)
