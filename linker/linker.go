package linker

import (
	"MxsxllBox/assembler"
	"MxsxllBox/cpu"
	"MxsxllBox/helper"
	"os"
	"path/filepath"
	"strings"
)

const (
	outName    = "Out"
	objOutName = outName + "/ObjOut"
)

func LinkModules(filePaths map[string]uint16) ([]byte, error) {
	finalCode := make([]byte, cpu.MemorySize)
	globalLookupTable := make(map[string]uint16)
	allObjFiles := make(map[*assembler.ObjectFile]uint16, 0)

	for filePath, location := range filePaths {
		objFile, _ := assembler.ReadObjectFile(filePath)
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
			}
			symbol += location
			hi, lo := helper.EncodeAddr(symbol)
			objFile.Code[relo.Offset] = hi
			objFile.Code[relo.Offset+1] = lo
		}
		finalCode = helper.ConcactSliceAtIndex(finalCode, objFile.Code, int(location))
	}
	return finalCode, nil
}

func CompileAndLinkFiles(files map[string]uint16, Name string) []byte {
	//for now this funcion will recomplile all files
	//It will take relative paths
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	objFiles := make(map[string]uint16)
	for filePath, location := range files {
		NewFilePath := filepath.Join(wd, filePath)
		filePath = strings.ReplaceAll(filePath, ".asm", ".obj")

		data, err := os.ReadFile(NewFilePath)
		if err != nil {
			panic(err)
		}
		OutFilePath := filepath.Join(wd, objOutName, Name, filePath)
		assembler.Assemble(string(data), OutFilePath)
		objFiles[OutFilePath] = location
	}

	LinkedCode, err := LinkModules(objFiles)
	if err != nil {
		panic(err)
	}

	if Name == "" {
		panic("Empty Name")
	}
	finalOutPath := filepath.Join(wd, outName, Name, "program.bin")
	os.WriteFile(finalOutPath, LinkedCode, 0644)
	return LinkedCode
}
