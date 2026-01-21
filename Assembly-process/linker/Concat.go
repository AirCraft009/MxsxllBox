package linker

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const includeSignifier string = "#include"

// FindIncludes
// finds all include statements inside the file
// only accepts includes at the start
// after a line that isn't an include statement it returns\
//
// it also creates temp.asm file with the concatenated file
func FindIncludes(filePath string) (concatedData []string, err error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(filePath)

	stringData := string(file)
	for _, line := range strings.Split(stringData, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, includeSignifier) {
			break
		}
		// line should contain the relative path from the line file location to the include file
		line = strings.TrimSpace(strings.TrimPrefix(line, includeSignifier))

	}

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
