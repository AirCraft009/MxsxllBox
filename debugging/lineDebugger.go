func dissasemble(code []byte, Ops map[byte]string) {
pc := 0
for pc < len(code) {
opcode := code[pc]
mnemonic, exists := Ops[opcode]
if !exists {
fmt.Printf("0x%04X: UNKNOWN OPCODE 0x%02X\n", pc, opcode)
pc++
continue
}

offset := int(assembler.OffsetMap[mnemonic])
bytes := code[pc : pc+offset]

fmt.Printf("0x%04X: %-8s", pc, mnemonic)

switch offset {
case 1:
// NOP, HALT, RET, etc.
fmt.Println()

case 3:
// One register (like PUSH Rx, POP Rx, PRINT Rx, GPC Rx, etc.)
rx, _ := assembler.DecodeRegs(bytes[1], bytes[2])
fmt.Printf("R%d\n", rx)

case 5:
if assembler.IsRegAddrFormat(mnemonic) {
// OP Rx, addr
rx, isImm := assembler.DecodeSingleReg(bytes[1], bytes[2])
addr := assembler.DecodeAddr(bytes[3], bytes[4])
if isImm {
fmt.Printf("R%d, %d\n", rx, addr)
} else {
fmt.Printf("R%d, R%d\n", rx, addr) // treating addr as reg in that case
}
} else if assembler.IsRegRegFormat(mnemonic) {
// OP Rx, Ry
rx, ry := assembler.DecodeRegs(bytes[1], bytes[2])
fmt.Printf("R%d, R%d\n", rx, ry)
} else if assembler.IsJmpFormat(mnemonic) {
// JMP label (currently displays raw address)
addr := assembler.DecodeAddr(bytes[3], bytes[4])
fmt.Printf("0x%04X\n", addr)
} else {
fmt.Println("???") // fallback
}

default:
fmt.Println("??")
}

pc += offset
}
}
