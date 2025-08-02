package debugging

import (
	"MxsxllBox/Assembly-process/assembler"
	"fmt"
)

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

}
