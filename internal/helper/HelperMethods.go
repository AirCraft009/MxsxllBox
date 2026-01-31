package helper

func DecodeAddr(hi byte, lo byte) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}

func DecodeRegs(reg1, reg2Wflag byte) (rx byte, ry byte, addresNec bool) {
	/**
	!!!old: theory still applies{
		new config:
		rx = bits 8-4
		ry = bits 3-0
		addrFlag = byte 2 bit 0
	}
	reg contains both rx and ry
	rx = bits 7-5
	ry = bits 4-2
	flags, etc. = bits 0 - 1

	>> rightshifts all bits by the following number
	& bitwise and looks at each number does and
	Example for decoding
	reg = 11010101
	reg >> 5 = 00000110
	reg & 0x07 = 00000110 & 00000111
	rx = 00000110
	rx = 6

	reg = 11010101
	reg >> 2 = 00110101
	reg & 0x07 = 00110101 & 00000111
	ry = 00000101
	ry = 5

	*/
	rx = reg1
	ry = reg2Wflag >> 1
	addrnec := (reg2Wflag) & 0x01
	return rx, ry, addrnec != 0x0
}

func IsInterruptActivated(id int, mask byte) bool {
	collapsedId := (id / 5) - 1
	shiftedMask := mask >> collapsedId
	//fmt.Printf("mask: %d, shiftedMask: %d, id: %d\n", mask, shiftedMask, id)
	return shiftedMask%2 == 1
}
