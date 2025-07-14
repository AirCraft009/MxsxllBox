package helper

func EncodeRegs(rx, ry byte) byte {
	return ((rx & 0x07) << 5) | ((ry & 0x07) << 2)
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
