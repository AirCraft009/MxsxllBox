package assembler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func SaveObjectFile(obj *ObjectFile, w io.Writer) error {
	w.Write([]byte("MXBI")) //MxsxllBinary header
	binary.Write(w, binary.LittleEndian, uint16(len(obj.Code)))
	binary.Write(w, binary.LittleEndian, uint16(len(obj.Symbols)))
	binary.Write(w, binary.LittleEndian, uint16(len(obj.Relocs)))
	if obj.Entry {
		w.Write([]byte{1})
	} else {
		w.Write([]byte{0})
	}
	w.Write(obj.Code)

	for name, addr := range obj.Symbols {
		w.Write([]byte{byte(len(name))})
		w.Write([]byte(name))
		binary.Write(w, binary.LittleEndian, addr)
		if obj.Globals[addr] {
			w.Write([]byte{1})
		} else {
			w.Write([]byte{0})
		}
	}

	for _, rel := range obj.Relocs {
		binary.Write(w, binary.LittleEndian, rel.Offset)
		w.Write([]byte{byte(len(rel.Lbl))})
		w.Write([]byte(rel.Lbl))
	}

	return nil
}

func ReadObjectFile(path string) (*ObjectFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)

	header := make([]byte, 4)
	if _, err := buf.Read(header); err != nil {
		return nil, err
	}
	if string(header) != "MXBI" {
		return nil, fmt.Errorf("invalid object file format: missing MXBI header")
	}

	var symbolCount uint16
	var codeLen uint16
	var relocCount uint16
	var entryfile byte

	//.Read reads from r in this case buf into arg2 so &codelen
	binary.Read(buf, binary.LittleEndian, &codeLen)
	binary.Read(buf, binary.LittleEndian, &symbolCount)
	binary.Read(buf, binary.LittleEndian, &relocCount)
	binary.Read(buf, binary.LittleEndian, &entryfile)

	code := make([]byte, codeLen)
	buf.Read(code)

	symbols := make(map[string]uint16)
	globals := make(map[uint16]bool)
	for i := 0; i < int(symbolCount); i++ {
		var nameLen uint8
		binary.Read(buf, binary.LittleEndian, &nameLen)
		nameBytes := make([]byte, nameLen)
		buf.Read(nameBytes)

		var addr uint16
		binary.Read(buf, binary.LittleEndian, &addr)

		var global byte
		binary.Read(buf, binary.LittleEndian, &global)
		if global == 1 {
			globals[addr] = true
		} else if global == 0 {
			globals[addr] = false
		} else {
			panic(fmt.Errorf("invalid global flag: %d", global))
		}

		symbols[string(nameBytes)] = addr
	}

	relocs := make([]RelocationEntry, relocCount)
	for i := 0; i < int(relocCount); i++ {
		binary.Read(buf, binary.LittleEndian, &relocs[i].Offset)
		var labelLen byte
		binary.Read(buf, binary.LittleEndian, &labelLen)
		name := make([]byte, labelLen)
		buf.Read(name)
		relocs[i].Lbl = string(name)
	}

	return &ObjectFile{Code: code, Symbols: symbols, Relocs: relocs, Globals: globals, Entry: entryfile == 1}, nil
}
