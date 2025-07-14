package Assembler

import (
	"fmt"
	"strings"
)

type parser struct {
	Parsers map[byte]func(string) (int, []byte, error)
}

func parse()

func ParseLines(data string) []string {
	//turns into array removes comments
	stringLines := strings.Split(data, "\n")
	fmt.Println(stringLines)
	for index, line := range stringLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		commentIndex := strings.Index(line, "#")
		if commentIndex != -1 {
			line = line[:commentIndex]
			stringLines[index] = line
		}
	}
	return stringLines
}

func firstPass(data []string) (out []string, pc int) {

}
