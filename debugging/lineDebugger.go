package debugging

import "fmt"
import "MxsxllBox/assembler"

func reverseOpCodes(opcodes map[string]uint8) (reverseOpCodes map[uint8]string) {
	for k, v := range opcodes {
		reverseOpCodes[v] = k
	}
	return reverseOpCodes
}

func debugWPC(code []byte) {
	fmt.Println("DEBUG With PC")
	Ops := reverseOpCodes(assembler.OpCodes)
	dissasemble(code, Ops)
}

func dissasemble(code []byte, Ops map[byte]string) {
	stringCode := ""
	var offset byte
	for Pc, v := range code {
		// at first, it takes the opCode so that it knows the offset
		Op := Ops[v]
		offset = assembler.OffsetMap[Op]

	}
}
