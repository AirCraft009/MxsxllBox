package cpu

import (
	"io"
	"os"
	"path/filepath"
)

func WritetoRom(filename string, offset int) {
	wd, _ := os.Getwd()
	path := "VM/cpu/ROM.mem"
	fileToWrite, err := os.OpenFile(filepath.Join(wd, path), os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fileToWrite.Close()
	fileToRead, err := os.OpenFile(wd+"/"+filename, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fileToRead.Close()
	_, err = fileToWrite.Seek(int64(offset), io.SeekStart)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(fileToWrite, fileToRead)
	if err != nil {
		panic(err)
	}

}
