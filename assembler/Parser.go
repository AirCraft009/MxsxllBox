package assembler

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
	AddrLoc1     = 1
	AddrLoc2     = 2
	AddrOutLocHi = 2
	AddrOutLocLo = 3
)

type Parser struct {
	Parsers map[string]func(parameters []string, currPC int) (pc int, code []byte, syntax error)
	Labels  map[string]uint16
}

func newParser() *Parser {
	parser := &Parser{
		Parsers: make(map[string]func(parameters []string, currPC int) (pc int, code []byte, syntax error)),
	}

	parser.Parsers["NOP"] = parseFormatOP
	parser.Parsers["LOADB"] = parseFormatOPRegAddr
	parser.Parsers["LOADW"] = parseFormatOPRegAddr
	parser.Parsers["LOAD"] = parseFormatOPRegAddr
	parser.Parsers["STOREB"] = parseFormatOPRegAddr
	parser.Parsers["STOREW"] = parseFormatOPRegAddr
	parser.Parsers["STORE"] = parseFormatOPRegAddr
	parser.Parsers["MOVI"] = parseFormatOPRegAddr
	parser.Parsers["DIVI"] = parseFormatOPRegAddr
	parser.Parsers["MULI"] = parseFormatOPRegAddr
	parser.Parsers["SUBI"] = parseFormatOPRegAddr
	parser.Parsers["ADDI"] = parseFormatOPRegAddr
	parser.Parsers["ADD"] = parseFormatOPRegReg
	parser.Parsers["SUB"] = parseFormatOPRegReg
	parser.Parsers["DIV"] = parseFormatOPRegReg
	parser.Parsers["MUL"] = parseFormatOPRegReg
	parser.Parsers["JMP"] = parseFormatOPAddr
	parser.Parsers["JZ"] = parseFormatOPAddr
	parser.Parsers["JC"] = parseFormatOPAddr
	parser.Parsers["CALL"] = parseFormatOPAddr
	parser.Parsers["PUSH"] = parseFormatOPReg
	parser.Parsers["POP"] = parseFormatOPReg
	parser.Parsers["PRINT"] = parseFormatOPReg
	parser.Parsers["RET"] = parseFormatOP
	parser.Parsers["HALT"] = parseFormatOP

	return parser
}

func parseFormatOP(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	return currPC + 1, []byte{opCodes[parameters[OpCLoc]]}, nil
}

func parseFormatOPRegAddr(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 4)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry)
	addr, syntax := strconv.Atoi(parameters[AddrLoc2])
	if syntax != nil {
		panic("syntax error: " + syntax.Error())
	}
	hi, lo := helper.EncodeAddr(uint16(addr))
	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	currPC += len(code)
	return currPC, code, syntax
}

func parseFormatOPRegReg(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 2)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	ry = regMap[parameters[RegsLoc2]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry)
	currPC += len(code)
	return currPC, code, nil
}

func parseFormatOPAddr(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 4)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	addr, syntax := strconv.Atoi(parameters[AddrLoc1])
	if syntax != nil {
		panic("syntax error: " + syntax.Error())
	}
	hi, lo := helper.EncodeAddr(uint16(addr))
	code[RegsLocOut] = helper.EncodeRegs(rx, ry)
	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	currPC += len(code)
	return currPC, code, syntax
}

func parseFormatOPReg(parameters []string, currPC int) (pc int, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 2)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry)
	currPC += len(code)
	return currPC, code, nil
}

func ParseLines(data string) [][]string {
	//turns into array removes comments
	stringLines := strings.Split(data, "\n")
	stringParts := make([][]string, len(stringLines))
	stringPartIndex := 0
	fmt.Println(stringLines)
	for index, line := range stringLines {
		line = strings.TrimSpace(line)
		commentIndex := strings.Index(line, "#")
		if commentIndex != -1 {
			line = line[:commentIndex]
			stringLines[index] = line
		}
		if line != "" {
			lineParts := strings.Fields(line)
			stringParts[stringPartIndex] = lineParts
			stringPartIndex++
		}
	}
	outPut := make([][]string, stringPartIndex)
	copy(outPut, stringParts)
	return outPut
}

func FirstPass(data [][]string) {
	for index, line := range data {
		if len(line) == 1 && strings.Contains(line[0], ":") {

		}
	}
}
