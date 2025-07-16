package linker

import "fmt"

const (
	strncpy     = 0x0C01
	ProgramSize = 4096
)

func linkModuels(modules [][]byte, programCode []byte) ([]byte, error) {
	linked := make([]byte, ProgramSize)
	offset := 0

	for i, mod := range modules {
		if offset+len(mod) > ProgramSize {
			return nil, fmt.Errorf("module %d too large to fit in program memory", i)
		}
		copy(linked[offset:], mod)
		offset += len(mod)
	}

	return linked, nil
}
