package linker

import (
	assembler2 "MxsxllBox/Assembly-process/assembler"
	"MxsxllBox/VM/cpu"
	"MxsxllBox/helper"
	"os"
	"path/filepath"
	"strings"
)

const (
	outName    = "VM-bin"
	objOutName = outName + "/ObjOut"
)

func LinkModules(filePaths map[string]uint16) (code []byte, debugLocations map[uint16]string, err error) {
	//debug locations are only necesarry if the debugger is used can be discarded otherwise
	finalCode := make([]byte, cpu.MemorySize)
	globalLookupTable := make(map[string]uint16)
	allObjFiles := make(map[*assembler2.ObjectFile]uint16)
	debugLocations = make(map[uint16]string)

	for filePath, location := range filePaths {
		objFile, _ := assembler2.ReadObjectFile(filePath)
		allObjFiles[objFile] = location
		for symbol, relAddr := range objFile.Symbols {
			if objFile.Globals[relAddr] {
				if _, ok := globalLookupTable[symbol]; ok {
					panic("Duplicate Lbl names: " + symbol + " file: " + filePath)
				}
				globalLookupTable[symbol] = location + relAddr
			}
		}
	}
	for objFile, location := range allObjFiles {
		for _, relo := range objFile.Relocs {
			symbol, ok := objFile.Symbols[relo.Lbl]
			if !ok {
				globalSymbol, k := globalLookupTable[relo.Lbl]
				if !k {
					panic("Label not found: " + relo.Lbl)
				}
				symbol = globalSymbol
			} else {
				symbol += location
			}
			debugLocations[symbol] = relo.Lbl
			hi, lo := helper.EncodeAddr(symbol)
			objFile.Code[relo.Offset] = hi
			objFile.Code[relo.Offset+1] = lo
		}
		finalCode = helper.ConcactSliceAtIndex(finalCode, objFile.Code, int(location))
	}
	return finalCode, debugLocations, nil
}

func CompileAndLinkFiles(files map[string]uint16, Name string) (code []byte, debugLocations map[uint16]string) {
	//for now this funcion will recomplile all files
	//It will take relative paths
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	genericOutPath := filepath.Join(wd, objOutName, Name)

	objFiles := make(map[string]uint16)
	locations := make(map[uint16]uint16)
	for filePath, location := range files {
		NewFilePath := filepath.Join(wd, filePath)
		filePath = strings.ReplaceAll(filePath, ".asm", ".obj")

		data, err := os.ReadFile(NewFilePath)
		if err != nil {
			panic(err)
		}
		OutFilePath := filepath.Join(genericOutPath, filePath)
		code := assembler2.Assemble(string(data), OutFilePath).Code
		if value, ok := locations[location]; ok {
			objFiles[OutFilePath] = location + value
			locations[location] = uint16(len(code)) + value
		} else {
			objFiles[OutFilePath] = location
			locations[location] = uint16(len(code))
		}
	}

	LinkedCode, debugLocations, err := LinkModules(objFiles)
	if err != nil {
		panic(err)
	}

	if Name == "" {
		panic("Empty Name")
	}
	finalOutPath := filepath.Join(wd, outName, Name, "program.bin")
	os.WriteFile(finalOutPath, LinkedCode, 0644)
	return LinkedCode, debugLocations
}

func setBasePaths(fileName string) map[string]uint16 {
	paths := make(map[string]uint16, 6)
	paths[fileName] = 0x00
	paths["\\VM\\stdlib\\io.asm"] = cpu.ProgramStdLibStart
	paths["\\VM\\stdlib\\math.asm"] = cpu.ProgramStdLibStart
	paths["\\VM\\stdlib\\string.asm"] = cpu.ProgramStdLibStart
	paths["\\VM\\stdlib\\sys.asm"] = cpu.ProgramStdLibStart
	paths["\\VM\\stdlib\\utils.asm"] = cpu.ProgramStdLibStart
	paths["\\VM\\interrupts\\interruptTable.asm"] = cpu.Interrupttable
	paths["\\VM\\scheduler\\scheduler.asm"] = 300
	return paths
}

func CompileForOs(fileName, Name string) []byte {
	paths := setBasePaths(fileName)
	code, _ := CompileAndLinkFiles(paths, Name)
	return code
}

func CompileForDebug(fileName, Name string) ([]byte, map[uint16]string) {
	paths := setBasePaths(fileName)
	return CompileAndLinkFiles(paths, Name)
}
