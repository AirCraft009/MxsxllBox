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
	Parsers    map[string]func(parameters []string, currPC uint16) (pc uint16, code []byte, syntax error)
	Labels     map[string]uint16
	LabelCodes map[uint16]bool
}

func newParser() *Parser {
	parser := &Parser{
		Parsers:    make(map[string]func(parameters []string, currPC uint16) (pc uint16, code []byte, syntax error)),
		Labels:     make(map[string]uint16),
		LabelCodes: make(map[uint16]bool),
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

func parseFormatOP(parameters []string, currPC uint16) (pc uint16, code []byte, syntax error) {
	return currPC + 1, []byte{opCodes[parameters[OpCLoc]]}, nil
}

func parseFormatOPRegAddr(parameters []string, currPC uint16) (pc uint16, code []byte, syntax error) {
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
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatOPRegReg(parameters []string, currPC uint16) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 2)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	ry = regMap[parameters[RegsLoc2]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry)
	currPC += uint16(len(code))
	return currPC, code, nil
}

func parseFormatOPAddr(parameters []string, currPC uint16) (pc uint16, code []byte, syntax error) {
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
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatOPReg(parameters []string, currPC uint16) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 2)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry)
	currPC += uint16(len(code))
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

func FirstPass(data [][]string, parser *Parser) *Parser {
	var PC uint16
	for _, line := range data {
		if len(line) == 1 && strings.Contains(line[0], ":") {
			parser.Labels[line[0]] = PC
			parser.LabelCodes[PC] = true
		} else {
			parser.LabelCodes[PC] = false
		}
		PC += uint16(len(line))
	}
	return parser
}

func Assemble(data string) (code []byte) {
	formattedData := ParseLines(data)
	parser := newParser()
	parser = FirstPass(formattedData, parser)
	return SecondPass(formattedData, parser)
}

func SecondPass(data [][]string, parser *Parser) (code []byte) {
	code = make([]byte, 4)
	PC := uint16(0)
	for _, line := range data {
		if parser.LabelCodes[PC] {
			PC++
			continue
		}
		if parsfunc, ok := parser.Parsers[line[0]]; ok {
			codeSnippet := make([]byte, 2)
			var err error
			PC, codeSnippet, err = parsfunc(line, PC)
			if err != nil {
				panic(err)
			}
			code = append(code, codeSnippet...)
		}
	}
	return code
}
