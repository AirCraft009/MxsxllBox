package debugging

func ReverseMaps(opcodes map[string]uint8) (reverseMap map[uint8]string) {
	reverseMap = make(map[uint8]string)
	for k, v := range opcodes {
		reverseMap[v] = k
	}
	return reverseMap
}
