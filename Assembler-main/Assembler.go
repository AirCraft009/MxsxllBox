package main

import (
	"MxsxllBox/Assembly-process/linker"

	flag "github.com/spf13/pflag"
)
import (
	"fmt"
	"os"
)

func main() {
	// define flags
	outFile := flag.String("o", "out.bin", "output file")
	nums := flag.IntSlice("n", []int{}, "integer values (will be aligned to files)")

	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Usage: assembler file1.asm file2.asm ... -n 0,1 -o out.bin")
		os.Exit(1)
	}

	// if fewer numbers than files, extend with last number
	if len(*nums) > 0 && len(*nums) < len(files) {
		last := (*nums)[len(*nums)-1]
		for len(*nums) < len(files) {
			*nums = append(*nums, last)
		}
	}

	// if no numbers at all, default all to 0
	if len(*nums) == 0 {
		for range files {
			*nums = append(*nums, 0)
		}
	}

	paths := make(map[string]uint16, len(files))
	for i, f := range files {
		paths[f] = uint16((*nums)[i])
	}
	code, _, err := linker.LinkModules(paths)
	if err != nil {
		return
	}
	err = os.WriteFile(*outFile, code, 0644)
}
