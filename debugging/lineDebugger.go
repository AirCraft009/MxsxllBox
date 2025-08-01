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

}

func dissasemble(code []byte) {
	stringCode := ""
	for Pc, v := range code {
		// at first
	}
}
