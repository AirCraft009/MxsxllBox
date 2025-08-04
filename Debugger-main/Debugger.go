package main

import (
	"MxsxllBox/Assembly-process/assembler"
	"MxsxllBox/Assembly-process/linker"
	"MxsxllBox/VM/cpu"
	"MxsxllBox/debugging"
	"MxsxllBox/helper"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
	_ "time"
)

const (
	OP = 2*iota + 1
	OPREG
	OPADDR
)

func addToString(input string, args []string) string {
	for arg := range args {
		input += " " + args[arg]
	}
	return input
}

func condenseNop(index int, code []byte) (newIndex, nopCount int) {
	for i := index; i < len(code); i++ {
		if code[i] == 0 {
			nopCount++
			continue
		}
		return i, nopCount
	}
	return len(code), len(code) - index
}

func dissasembleForDebugging(code []byte, lblocations map[uint16]string) (file string, PcToLine map[uint16]int) {
	PcToLine = make(map[uint16]int)

	revOpCodes := debugging.ReverseMaps(assembler.OpCodes)
	revRegMap := debugging.ReverseMaps(assembler.RegMap)
	var line string
	var nopCount int
	for i := 0; i < len(code); i += 0 {
		PcToLine[uint16(i)] = len(strings.Split(file, "\n")) - 1
		var args []string
		ByteInstruction := code[i]
		instruction := revOpCodes[ByteInstruction]
		offset := assembler.OffsetMap[instruction]
		if lbl, ok := lblocations[uint16(i)]; ok {
			line = "\n" + lbl + "\n"
			PcToLine[uint16(i)] += len(strings.Split(line, "\n")) - 1
			line += instruction
		} else {
			line = instruction
		}

		if ByteInstruction == 0 {
			i, nopCount = condenseNop(i, code)
			line = addToString(line, []string{strconv.Itoa(nopCount)})
			line += "\n\n"
			file += line
			continue
		}

		switch offset {
		case OP:
			break
		case OPREG:
			reg1Encoded, reg2Encoded := code[i+1], code[i+2]
			reg1Decoded, reg2Decoded, _ := helper.DecodeRegs(reg1Encoded, reg2Encoded)
			reg1, reg2 := revRegMap[reg1Decoded], revRegMap[reg2Decoded]
			args = append(args, reg1, reg2)
			line = addToString(line, args)
			break
		case OPADDR:
			reg1Encoded, reg2Encoded, addrBit1, addrBit2 := code[i+1], code[i+2], code[i+3], code[i+4]
			reg1Decoded, reg2Decoded, _ := helper.DecodeRegs(reg1Encoded, reg2Encoded)
			reg1, reg2 := revRegMap[reg1Decoded], revRegMap[reg2Decoded]
			addr := helper.DecodeAddr(addrBit1, addrBit2)

			stringAddr := strconv.Itoa(int(addr))
			if lbl, ok := lblocations[addr]; ok {
				stringAddr = lbl
			}
			args = append(args, reg1, reg2, stringAddr)
			line = addToString(line, args)
		default:
			panic("Unknown offsetLen " + strconv.Itoa(int(offset)))
		}

		line += "\n"
		file += line
		i += int(offset)
	}
	return file, PcToLine
}

func init() {
	// Suppress all standard logs (including Fyne logs using `log.Print`)
	log.SetOutput(io.Discard)
}

func main() {
	code, lblLocations := linker.CompileForDebug("program.asm", "EchoKeys")
	mem := &cpu.Memory{}
	copy(mem.Data[:], code)
	debugVm := cpu.NewCPU(mem)

	file, PcToLine := dissasembleForDebugging(code, lblLocations)
	lines := strings.Split(file, "\n")

	currentLine := 0
	myApp := app.New()

	var lineBoxes []fyne.CanvasObject
	var lineBackgrounds []*canvas.Rectangle

	for i, text := range lines {
		bg := canvas.NewRectangle(color.Black)
		fg := color.White
		if i == currentLine {
			bg.FillColor = color.RGBA{0, 0, 255, 255}
		}

		label := canvas.NewText(text, fg)
		label.TextSize = 16

		lineBackgrounds = append(lineBackgrounds, bg)
		lineBoxes = append(lineBoxes, container.NewMax(bg, label))
	}

	textList := container.NewVBox(lineBoxes...)
	scroll := container.NewScroll(textList)

	currentMode := "Step"
	var modeButton *widget.Button
	modeButton = widget.NewButton("Mode: Step", func() {
		if currentMode == "Step" {
			currentMode = "Run"
		} else {
			currentMode = "Step"
		}
		modeButton.SetText("Mode: " + currentMode)
	})

	topBar := container.NewHBox(layout.NewSpacer(), modeButton)

	myWindow := myApp.NewWindow("Debugger UI")
	myWindow.Resize(container.NewVBox(topBar, scroll).MinSize())
	myWindow.SetContent(container.NewBorder(topBar, nil, nil, nil, scroll))

	// Line highlighter
	highlightLine := func(newIndex int) {
		if newIndex < 0 || newIndex >= len(lineBackgrounds) {
			return
		}
		lineBackgrounds[currentLine].FillColor = color.Black
		lineBackgrounds[currentLine].Refresh()

		lineBackgrounds[newIndex].FillColor = color.RGBA{B: 255, A: 255}
		lineBackgrounds[newIndex].Refresh()

		currentLine = newIndex
	}

	// Simulate stepping
	go func() {

		fmt.Println(PcToLine)
		for {
			time.Sleep(1 * time.Second)
			fmt.Println(PcToLine[debugVm.PC])
			newIndex := PcToLine[debugVm.PC]
			debugVm.Step()

			highlightLine(newIndex)
		}
	}()

	myWindow.ShowAndRun()
}
