package Assembler

import (
	"MxsxllBox/helper"
	"fmt"
	"strconv"
	"strings"
)

const (
	OpCLoc   = 0
	RegsLoc1 = 1
	RegsLoc2 = 2
	RegsLocOut
	AddrLoc1     = 2
	AddrLoc2     = 3
	AddrOutLocHi = 2
	AddrOutLocLo = 3
)

type parser struct {
	Parsers map[byte]func(parameters []string, currPC int) (pc int, code []byte, syntax error)
}

func parseNop(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	return currPC + 1, []byte{0x00}, nil
}

func parseLoadB(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	return parseLoad(parameters, currPC)
}

func parseLoadW(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	return parseLoad(parameters, currPC)
}

func parseLoad(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	var rx, ry byte
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry)
	addr, syntax := strconv.Atoi(parameters[AddrLoc1])
	hi, lo := helper.EncodeAddr(uint16(addr))
	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	currPC += len(code)
	return currPC, code, syntax
}

func parseStore(parameters []string, currPC int) (pc int, code []byte, syntax error) {

}

func ParseLines(data string) [][]string {
	//turns into array removes comments
	stringLines := strings.Split(data, "\n")
	stringParts := make([][]string, len(stringLines))
	fmt.Println(stringLines)
	for index, line := range stringLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		commentIndex := strings.Index(line, "#")
		if commentIndex != -1 {
			line = line[:commentIndex]
			stringLines[index] = line
		}
		lineParts := strings.Split(line, " ")
		stringParts[index] = lineParts
	}
	return stringParts
}
