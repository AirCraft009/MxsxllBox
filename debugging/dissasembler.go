package debugging

import (
	"MxsxllBox/Assembly-process/assembler"
	"MxsxllBox/helper"
	"io"
	"log"
	"strconv"
	"strings"
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

func DissasembleForDebugging(code []byte, lblocations map[uint16]string) (file string, PcToLine map[uint16]int) {
	PcToLine = make(map[uint16]int)

	revOpCodes := ReverseMaps(assembler.OpCodes)
	revRegMap := ReverseMaps(assembler.RegMap)
	var line string
	var nopCount int
	for i := 0; i < len(code); i += 0 {
		PcToLine[uint16(i)] = len(strings.Split(file, "\n")) - 1
		var args []string
		ByteInstruction := code[i]
		instruction := revOpCodes[ByteInstruction]
		offset := assembler.OffsetMap[instruction]
		if lbl, ok := lblocations[uint16(i)]; ok {
			line = "\n" + lbl + "\n"
			PcToLine[uint16(i)] += len(strings.Split(line, "\n")) - 1
			line += instruction
		} else {
			line = instruction
		}

		if ByteInstruction == 0 {
			i, nopCount = condenseNop(i, code)
			line = addToString(line, []string{strconv.Itoa(nopCount)})
			line += "\n\n"
			file += line
			continue
		}

		switch offset {
		case OP:
			break
		case OPREG:
			reg1Encoded, reg2Encoded := code[i+1], code[i+2]
			reg1Decoded, reg2Decoded, _ := helper.DecodeRegs(reg1Encoded, reg2Encoded)
			reg1, reg2 := revRegMap[reg1Decoded], revRegMap[reg2Decoded]
			args = append(args, reg1, reg2)
			line = addToString(line, args)
			break
		case OPADDR:
			reg1Encoded, reg2Encoded, addrBit1, addrBit2 := code[i+1], code[i+2], code[i+3], code[i+4]
			reg1Decoded, reg2Decoded, _ := helper.DecodeRegs(reg1Encoded, reg2Encoded)
			reg1, reg2 := revRegMap[reg1Decoded], revRegMap[reg2Decoded]
			addr := helper.DecodeAddr(addrBit1, addrBit2)

			stringAddr := strconv.Itoa(int(addr))
			if lbl, ok := lblocations[addr]; ok {
				stringAddr = lbl
			}
			args = append(args, reg1, reg2, stringAddr)
			line = addToString(line, args)
		default:
			panic("Unknown offsetLen " + strconv.Itoa(int(offset)))
		}

		line += "\n"
		file += line
		i += int(offset)
	}
	return file, PcToLine
}

func init() {
	// Suppress all standard logs (including Fyne logs using `log.Print`)
	log.SetOutput(io.Discard)
}
