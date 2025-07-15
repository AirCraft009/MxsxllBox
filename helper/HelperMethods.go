package helper

func EncodeRegs(rx byte, ry byte, addrnecs bool) byte {
	var addrFlag byte
	if addrnecs {
		addrFlag |= 0x1
	} else {
		addrFlag |= 0x0
	}
	return ((rx & 0x07) << 5) | ((ry & 0x07) << 2) | (addrFlag & 0x03)
}

func EncodeAddr(addr uint16) (byte, byte) {
	if addr <= 255 {
		return byte(0), byte(addr)
	}
	return byte(addr>>8) & 0xff, byte(addr & 0xff)
}

func DecodeAddr(lo byte, hi byte) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}
