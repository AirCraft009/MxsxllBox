package debugging

func reverseOpCodes(opcodes map[string]uint8) (reverseOpCodes map[uint8]string) {
	for k, v := range opcodes {
		reverseOpCodes[v] = k
	}
	return reverseOpCodes
}
