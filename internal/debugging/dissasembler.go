package debugging

import (
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/AirCraft009/MxsxllBox/internal/MxsxllBox/cpu"
	"github.com/AirCraft009/MxsxllBox/internal/helper"

	"github.com/AirCraft009/mcc/pkg"
)

const (
	OP = 2*iota + 1
	OPREG
	OPADDR
)

func addToString(input string, args []string) string {
	for arg := range args {
		input += " " + args[arg]
	}
	return input
}

func condenseNop(index int, code []byte) (newIndex, nopCount int) {
	for i := index; i < len(code); i++ {
		if code[i] == 0 {
			nopCount++
			continue
		}
		return i, nopCount
	}
	return len(code), len(code) - index
}

func DisassembleForDebugging(code []byte, lblocations map[uint16]string) (string, map[uint16]int) {
	pcToLine := make(map[uint16]int)

	revOpCodes := InvertMaps(pkg.OpCodes)
	revRegMap := InvertMaps(pkg.RegMap)

	var out strings.Builder
	lineNo := 0 // visual line number in output
	pc := 0     // program counter into code[]

	for pc < len(code) {
		op := code[pc]

		instr, ok := revOpCodes[op]
		if !ok {
			panic("Debugger: unknown opcode: " + strconv.Itoa(int(op)))
		}

		offset := pkg.OffsetMap[instr]

		// ------------------------------------------------------------------
		// Labels are *metadata*, not executable instructions.
		// They must NEVER be mapped as PC -> line targets.
		// Therefore:
		//   1. Emit labels first
		//   2. Advance lineNo
		//   3. Map PC -> *instruction* line AFTER labels
		// ------------------------------------------------------------------
		if lbl, ok := lblocations[uint16(pc)]; ok {
			out.WriteString(lbl)
			out.WriteByte('\n')
			lineNo++
		}

		// This mapping MUST point to the instruction mnemonic line.
		// Never move this above label emission.
		pcToLine[uint16(pc)] = lineNo

		out.WriteString(instr)

		// ------------------------------------------------------------------
		// Special case: NOP compression
		// NOP is still an instruction, so mapping above is correct.
		// ------------------------------------------------------------------
		if op == cpu.NOP {
			var count int
			pc, count = condenseNop(pc, code)

			out.WriteByte(' ')
			out.WriteString(strconv.Itoa(count))
			out.WriteString("\n\n")

			// Two newlines = instruction line + spacer
			lineNo += 2
			continue
		}

		// ------------------------------------------------------------------
		// Operand decoding depends on instruction format (offset)
		// ------------------------------------------------------------------
		switch offset {
		case OP:
			// no operands

		case OPREG:
			if pc+2 >= len(code) {
				panic("truncated OPREG")
			}

			r1e, r2e := code[pc+1], code[pc+2]
			r1d, r2d, _ := helper.DecodeRegs(r1e, r2e)

			out.WriteByte(' ')
			out.WriteString(revRegMap[r1d])
			out.WriteByte(' ')
			out.WriteString(revRegMap[r2d])

		case OPADDR:
			if pc+4 >= len(code) {
				panic("truncated OPADDR")
			}

			r1e, r2e := code[pc+1], code[pc+2]
			a1, a2 := code[pc+3], code[pc+4]

			r1d, r2d, _ := helper.DecodeRegs(r1e, r2e)
			addr := helper.DecodeAddr(a1, a2)

			out.WriteByte(' ')
			out.WriteString(revRegMap[r1d])
			out.WriteByte(' ')
			out.WriteString(revRegMap[r2d])
			out.WriteByte(' ')

			// Prefer symbolic labels for addresses when available
			if lbl, ok := lblocations[addr]; ok {
				out.WriteString(lbl)
			} else {
				out.WriteString(strconv.Itoa(int(addr)))
			}

		default:
			panic("unknown instruction format")
		}

		out.WriteByte('\n')
		lineNo++
		pc += int(offset)
	}

	return out.String(), pcToLine
}

// InvertMaps
// returns the map with keys as values and values as keys
func InvertMaps(opcodes map[string]uint8) (reverseMap map[uint8]string) {
	reverseMap = make(map[uint8]string)
	for k, v := range opcodes {
		reverseMap[v] = k
	}
	return reverseMap
}

func init() {
	// Suppress all standard logs (including Fyne logs using `log.Print`)
	log.SetOutput(io.Discard)
}
