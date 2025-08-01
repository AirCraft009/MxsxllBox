package helper

func EncodeRegs(rx byte, ry byte, addrnecs bool) (byte, byte) {
	var addrFlag byte
	if addrnecs {
		addrFlag = 0x1
	} else {
		addrFlag = 0x0
	}
	return rx, ry<<1 | addrFlag
}

func EncodeAddr(addr uint16) (byte, byte) {
	if addr <= 255 {
		return byte(0), byte(addr)
	}
	return byte(addr>>8) & 0xff, byte(addr & 0xff)
}

func InsertMatrixAtIndex(dest, insert [][]string, index int) [][]string {
	if index < 0 {
		index = 0
	} else if index > len(dest) {
		index = len(dest)
	}

	result := make([][]string, 0, len(dest)+len(insert))
	result = append(result, dest[:index]...)
	result = append(result, insert...)
	result = append(result, dest[index:]...)

	return result
}

func DeleteMatrixRow(matrix [][]string, index int) [][]string {
	if index < 0 || index >= len(matrix) {
		return matrix
	}

	return append(matrix[:index], matrix[index+1:]...)
}

func DecodeAddr(lo byte, hi byte) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}

func ConcactSliceAtIndex(dest, input []byte, index int) []byte {
	if len(dest)-int(index) < len(input) {
		return dest
	}
	for i := 0; i < len(input); i++ {
		dest[index+i] = input[i]
	}
	return dest
}

func DecodeRegs(reg1, reg2Wflag byte) (rx byte, ry byte, addresNec bool) {
	/**
	old: theory still applies{
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
