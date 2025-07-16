package linker

import (
	"MxsxllBox/assembler"
	"MxsxllBox/helper"
)

const (
	strncpy     = 0x0C01
	ProgramSize = 4096
)

func LinkModuels(filePaths map[string]uint16) ([]byte, error) {
	finalCode := make([]byte, ProgramSize)
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
