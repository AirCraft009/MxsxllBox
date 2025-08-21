package cpu

import (
	"io"
	"os"
)

func WritetoRom(filename string, offset int) {
	wd, _ := os.Getwd()
	path := "/VM/cpu/ROM.mem"
	fileToWrite, _ := os.OpenFile(wd+path, os.O_WRONLY, 0666)
	defer fileToWrite.Close()
	fileToRead, _ := os.OpenFile(wd+"/"+filename, os.O_RDONLY, 0666)
	defer fileToRead.Close()
	_, _ = fileToWrite.Seek(int64(offset), 0)
	_, _ = io.Copy(fileToWrite, fileToRead)
}
