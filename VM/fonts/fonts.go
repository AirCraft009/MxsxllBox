package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	path := wd + "/VM/fonts/"
	data, err := os.ReadFile(path + "font8x8_basic.bin")
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(data)
	realchar := 8 * 95
	waste := make([]byte, realchar)
	_, err = buf.Read(waste)
	if err != nil {
		panic(err)
	}
	characterBuf := make([]byte, 8)
	_, err = buf.Read(characterBuf)
	if err != nil {
		panic(err)
	}

	render(characterBuf)
}

func render(bitmap []byte) {
	var x, y int
	var set byte
	for x = 0; x < 8; x++ {
		for y = 0; y < 8; y++ {
			set = bitmap[x] & (1 << y)
			if set == 0 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%c", 'x')
			}
		}
		fmt.Printf("\n")
	}
}
