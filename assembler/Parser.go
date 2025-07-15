package assembler

import (
	"MxsxllBox/helper"
	"strconv"
	"strings"
)

const (
	OpCLoc       = 0
	RegsLoc1     = 1
	RegsLoc2     = 2
	RegsLocOut   = 1
	AddrLoc1     = 1
	AddrLoc2     = 2
	AddrOutLocHi = 2
	AddrOutLocLo = 3
)

type Parser struct {
	Parsers    map[string]func(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error)
	Labels     map[string]uint16
	LabelCodes map[uint16]bool
}

func newParser() *Parser {
	parser := &Parser{
		Parsers:    make(map[string]func(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error)),
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
	parser.Parsers["JMP"] = parseFormatOPLbl
	parser.Parsers["JZ"] = parseFormatOPLbl
	parser.Parsers["JC"] = parseFormatOPLbl
	parser.Parsers["CALL"] = parseFormatOPLbl
	parser.Parsers["PUSH"] = parseFormatOPReg
	parser.Parsers["POP"] = parseFormatOPReg
	parser.Parsers["PRINT"] = parseFormatOPReg
	parser.Parsers["RET"] = parseFormatOP
	parser.Parsers["HALT"] = parseFormatOP
	parser.Parsers["ALLOC"] = parseFormatOPRegReg
	parser.Parsers["FREE"] = parseFormatOPReg

	return parser
}

func parseFormatOP(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	return currPC + 1, []byte{opCodes[parameters[OpCLoc]]}, nil
}

func parseFormatOPRegAddr(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 2)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	ry, ok := regMap[parameters[RegsLoc2]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry, !ok)
	if !ok {
		addr, syntax := strconv.Atoi(parameters[AddrLoc2])
		if syntax != nil {
			panic("syntax error: " + syntax.Error())
		}
		hi, lo := helper.EncodeAddr(uint16(addr))
		code = append(code, hi, lo)
	}
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatString(parameters []string, currPC uint16, parser *Parser) (pc, code []byte, syntax error) {
	var _, _ byte
	_ = regMap[parameters[RegsLoc1]]
	_ = regMap[parameters[RegsLoc2]]
	inputString := parameters[AddrLoc1]
	inputString = strings.ReplaceAll(inputString, "\"", "")
	/**
	length := len(inputString)
	MOVI reg1 part
	ADDI reg2 0
	STOREB reg1 reg2
	code = make([]byte)
	*/
	return nil, nil, nil
}

func parseFormatOPRegReg(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 2)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	ry = regMap[parameters[RegsLoc2]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry, false)
	currPC += uint16(len(code))
	return currPC, code, nil
}

func parseFormatOPAddr(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 4)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	addr, syntax := strconv.Atoi(parameters[AddrLoc1])
	if syntax != nil {
		panic("syntax error: " + syntax.Error())
	}
	hi, lo := helper.EncodeAddr(uint16(addr))
	code[RegsLocOut] = helper.EncodeRegs(rx, ry, true)
	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatOPLbl(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 4)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	addr, ok := parser.Labels[parameters[AddrLoc1]]
	if !ok {
		panic("label not found: " + parameters[AddrLoc1])
	}
	hi, lo := helper.EncodeAddr(uint16(addr))
	code[RegsLocOut] = helper.EncodeRegs(rx, ry, true)
	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatOPReg(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 2)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	code[RegsLocOut] = helper.EncodeRegs(rx, ry, false)
	currPC += uint16(len(code))
	return currPC, code, nil
}

func ParseLines(data string) [][]string {
	//turns into array removes comments
	stringLines := strings.Split(data, "\n")
	stringParts := make([][]string, len(stringLines))
	stringPartIndex := 0
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
	return stringParts[:stringPartIndex]
}

func FirstPass(data [][]string, parser *Parser) *Parser {
	var PC uint16
	for _, line := range data {
		if len(line) == 1 && strings.Contains(line[0], ":") {
			parser.Labels[line[0][:len(line[0])-1]] = PC
			parser.LabelCodes[PC] = true
			continue
		}
		PC += uint16(getOffset(line[0]))
	}
	return parser
}

func getOffset(OP string) byte {
	offset := offsetMap[OP]
	if offset == 0 {
		return 2
	}
	return offset
}

func Assemble(data string) (code []byte) {
	formattedData := ParseLines(data)
	parser := newParser()
	parser = FirstPass(formattedData, parser)
	return SecondPass(formattedData, parser)
}

func SecondPass(data [][]string, parser *Parser) (code []byte) {
	code = make([]byte, 0)
	PC := uint16(0)
	for _, line := range data {
		if parsfunc, ok := parser.Parsers[line[0]]; ok {
			codeSnippet := make([]byte, 2)
			var err error
			PC, codeSnippet, err = parsfunc(line, PC, parser)
			if err != nil {
				panic(err)
			}
			code = append(code, codeSnippet...)
		}
	}
	return code
}
