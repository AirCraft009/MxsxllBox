package assembler

import (
	"MxsxllBox/helper"
	"os"
	"strconv"
	"strings"
)

const (
	OpCLoc         = 0
	RegsLoc1       = 1
	RegsLoc2       = 2
	RegsLocOut     = 1
	AddrLoc1       = 1
	AddrLoc2       = 2
	AddrOutLocHi   = 3
	AddrOutLocLo   = 4
	StrLoc         = 3
	RegWidthOffset = 1
)

type Parser struct {
	Parsers   map[string]func(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error)
	Formatter map[string]func(parameters []string) (formatted [][]string)
	Labels    map[string]uint16
	ObjFile   *ObjectFile
}

type ObjectFile struct {
	Code    []byte
	Symbols map[string]uint16
	Relocs  []RelocationEntry
	Globals map[uint16]bool
}

type RelocationEntry struct {
	Offset uint16 // Where in Code the label is called/JMP'd to
	Lbl    string
}

func newObjectFile() *ObjectFile {
	return &ObjectFile{
		Code:    nil,
		Symbols: make(map[string]uint16),
		Relocs:  make([]RelocationEntry, 0),
		Globals: make(map[uint16]bool),
	}
}

func newParser() *Parser {
	parser := &Parser{
		Parsers:   make(map[string]func(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error)),
		Formatter: make(map[string]func(parameters []string) (formatted [][]string)),
		Labels:    make(map[string]uint16),
		ObjFile:   newObjectFile(),
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
	parser.Parsers["PRINTSTR"] = parseFormatOPReg
	parser.Parsers["JNZ"] = parseFormatOPLbl
	parser.Parsers["JNC"] = parseFormatOPLbl
	parser.Parsers["CMP"] = parseFormatOPRegReg
	parser.Parsers["CMPI"] = parseFormatOPRegAddr
	parser.Parsers["TEST"] = parseFormatOPRegReg
	parser.Parsers["TSTI"] = parseFormatOPRegAddr
	parser.Parsers["JL"] = parseFormatOPLbl
	parser.Parsers["JLE"] = parseFormatOPLbl
	parser.Parsers["JG"] = parseFormatOPLbl
	parser.Parsers["JGE"] = parseFormatOPLbl
	parser.Formatter["STRING"] = formatString
	parser.Parsers["STZ"] = parseFormatOP
	parser.Parsers["STC"] = parseFormatOP
	parser.Parsers["CLZ"] = parseFormatOP
	parser.Parsers["CLC"] = parseFormatOP
	parser.Parsers["MOD"] = parseFormatOPRegReg
	parser.Parsers["MODI"] = parseFormatOPRegAddr
	parser.Parsers["MOV"] = parseFormatOPRegReg

	return parser
}

func parseFormatOP(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	return currPC + 1, []byte{opCodes[parameters[OpCLoc]]}, nil
}

func parseFormatOPRegAddr(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 3)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	ry, ok := regMap[parameters[RegsLoc2]]
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, !ok)
	if !ok {
		addr, syntax := strconv.Atoi(parameters[AddrLoc2])
		if syntax != nil {
			panic("syntax error: " + syntax.Error())
		}
		hi, lo := helper.EncodeAddr(uint16(addr))
		code = append(code, hi, lo)
	} else {
		hi, lo := helper.EncodeAddr(uint16(0))
		code = append(code, hi, lo)
	}
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func formatString(parameters []string) (formatted [][]string) {
	//STRING RegUse RegAddr "STRING"
	var rx, ry string
	//rx = part
	//ry = addr
	rx = parameters[RegsLoc1]
	ry = parameters[RegsLoc2]
	inputStringParts := parameters[StrLoc:len(parameters)]
	inputString := ""
	for _, part := range inputStringParts {
		inputString += part + " "
	}
	inputString = strings.ReplaceAll(inputString, "\"", "")
	inputString = inputString[:len(inputString)-1]
	length := len(inputString)
	/**
	Jeder einzelne char wird mit diesen drei Op's dargestellt
	es ist lenght prefix based bedeutet das ersteByte, welches gelesen wird ist die länge des Strings
	MOVI reg1 part
	ADDI reg2 0
	STOREB reg1 reg2
	*/
	ascii := make([]byte, length+2)
	formatted = make([][]string, 0)
	inputString = inputString + "/"
	ascii[0] = byte(length)
	var line []string
	line = []string{"SUBI", ry, "1"}
	formatted = append(formatted, line)
	for i, part := range inputString {
		line = []string{}
		ascii[i+1] = byte(part)
		line = append(line, "MOVI", rx, strconv.Itoa(int(ascii[i]))) //4
		formatted = append(formatted, line)
		line = []string{}
		line = append(line, "ADDI", ry, "1") //4
		formatted = append(formatted, line)
		line = []string{}
		line = append(line, "STOREB", rx, ry) //4
		formatted = append(formatted, line)
	}
	line = []string{}
	line = append(line, "SUBI", ry, strconv.Itoa(length))
	formatted = append(formatted, line)
	return formatted
}

func parseFormatOPRegReg(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 3)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	ry = regMap[parameters[RegsLoc2]]
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, false)
	currPC += uint16(len(code))
	return currPC, code, nil
}

func parseFormatOPAddr(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 5)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	addr, syntax := strconv.Atoi(parameters[AddrLoc1])
	if syntax != nil {
		panic("syntax error: " + syntax.Error())
	}
	hi, lo := helper.EncodeAddr(uint16(addr))
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, true)
	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	currPC += uint16(len(code))
	return currPC, code, syntax
}

func parseFormatOPLbl(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 5)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	/**
	Lbladdr, ok := parser.Labels[parameters[AddrLoc1]]
	if !ok {
		panic("label not found: " + parameters[AddrLoc1])
	}

	hi, lo := helper.EncodeAddr(uint16(addr))

	code[AddrOutLocHi] = hi
	code[AddrOutLocLo] = lo
	*/
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, true)
	currPC += AddrOutLocHi
	parser.ObjFile.Relocs = append(parser.ObjFile.Relocs, RelocationEntry{
		Offset: currPC,
		Lbl:    parameters[AddrLoc1],
	})
	code[AddrOutLocHi], code[AddrOutLocLo] = 0x00, 0x00
	currPC += uint16(len(code)) - AddrOutLocHi
	return currPC, code, syntax
}

func parseFormatOPReg(parameters []string, currPC uint16, parser *Parser) (pc uint16, code []byte, syntax error) {
	var rx, ry byte
	code = make([]byte, 3)
	code[OpCLoc] = opCodes[parameters[OpCLoc]]
	rx = regMap[parameters[RegsLoc1]]
	code[RegsLocOut], code[RegsLocOut+RegWidthOffset] = helper.EncodeRegs(rx, ry, false)
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

func FirstPass(data [][]string, parser *Parser) (*Parser, [][]string) {
	var PC uint16
	formattedExtraCode := make(map[int][][]string)

	for i, line := range data {
		if len(line) == 1 && strings.Contains(line[0], ":") {
			parser.Labels[line[0][:len(line[0])-1]] = PC
			if strings.HasPrefix(line[0], "_") {
				parser.ObjFile.Globals[PC] = true
			}
			continue
		} else if formatter, ok := parser.Formatter[line[0]]; ok {
			formatted := formatter(data[i])
			formattedExtraCode[i] = formatted
			for _, formatLine := range formatted {
				PC += uint16(getOffset(formatLine[0]))
			}
			continue
		}
		PC += uint16(getOffset(line[0]))
	}
	for i, extraData := range formattedExtraCode {
		data = helper.DeleteMatrixRow(data, i)
		data = helper.InsertMatrixAtIndex(data, extraData, i)
	}
	parser.ObjFile.Symbols = parser.Labels
	return parser, data
}

func getOffset(OP string) byte {
	offset := offsetMap[OP]
	if offset == 0 {
		return 2
	}
	return offset
}

func Assemble(data string) {
	parsedData := ParseLines(data)
	parser := newParser()
	var formattedData [][]string
	parser, formattedData = FirstPass(parsedData, parser)
	/**
	for lbl, adr := range parser.Labels {
		fmt.Printf("Label: %s :addr: %d\n", lbl, adr)
	}

	*/
	ObjFile := SecondPass(formattedData, parser)
	f, err := os.OpenFile("C:\\Users\\cocon\\Documents\\Projects\\Musa-Allmer\\MxsxllBox/program.bin", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	SaveObjectFile(ObjFile, f)
}

func SecondPass(data [][]string, parser *Parser) (ObjFile *ObjectFile) {
	code := make([]byte, 0)

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
	parser.ObjFile.Code = code
	return parser.ObjFile
}
