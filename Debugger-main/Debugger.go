package main

import (
	"MxsxllBox/Assembly-process/assembler"
	"MxsxllBox/Assembly-process/linker"
	"MxsxllBox/IO/KeyboardBuffer"
	"MxsxllBox/VM/cpu"
	"MxsxllBox/debugging"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
	"strings"
	"time"
)

func initDebugRegs(nameValueRows []fyne.CanvasObject, cpu *cpu.CPU, revRegMap map[uint8]string) (*fyne.Container, [17][2]*widget.Label) {
	//bad code I know but this is only for debugging
	var regLbls [len(cpu.Registers)/2 + 1][2]*widget.Label
	for i := 0; i < len(cpu.Registers); i += 2 {
		left := widget.NewLabel(fmt.Sprintf("%s: %d", revRegMap[uint8(i)], cpu.Registers[i]))
		right := widget.NewLabel(fmt.Sprintf("%s: %d", revRegMap[uint8(i+1)], cpu.Registers[i+1]))
		regLbls[i/2][0] = left
		regLbls[i/2][1] = right

		row := container.NewHBox(
			layout.NewSpacer(), left,
			layout.NewSpacer(), right,
			layout.NewSpacer(),
		)

		nameValueRows = append(nameValueRows, row)
	}
	right := widget.NewLabel(fmt.Sprintf("PC: %d", cpu.PC))
	left := widget.NewLabel(fmt.Sprintf("Stack-top: %d", cpu.Mem.ReadWord(cpu.SP)))
	regLbls[len(cpu.Registers)/2][0] = left
	regLbls[len(cpu.Registers)/2][1] = right

	row := container.NewHBox(
		layout.NewSpacer(), left,
		layout.NewSpacer(), right,
		layout.NewSpacer(),
	)

	nameValueRows = append(nameValueRows, row)
	nameValuePanel := container.NewVBox(nameValueRows...)
	return nameValuePanel, regLbls
}

func setRegDebug(regLbls [17][2]*widget.Label, cpu *cpu.CPU, revRegMap map[uint8]string) {
	for i := 0; i < len(cpu.Registers); i += 1 {
		name := fmt.Sprintf("%s:", revRegMap[uint8(i)])
		val := strconv.Itoa(int(cpu.Registers[i]))
		row := i / 2
		col := i % 2

		regLbls[row][col].Text = fmt.Sprintf("%s %s", name, val)
	}
	pcName := "PC:"
	pcVal := strconv.Itoa(int(cpu.PC))
	row := len(regLbls) - 1
	col := 0
	regLbls[row][col].Text = fmt.Sprintf("%s %s", pcName, pcVal)
	SpName := "Stack-Top:"
	SPVal := strconv.Itoa(int(cpu.Mem.ReadWord(cpu.SP)))
	col = 1
	regLbls[row][col].Text = fmt.Sprintf("%s %s", SpName, SPVal)
}

func main() {
	reverseRegMap := debugging.ReverseMaps(assembler.RegMap)
	var breakpoints = make(map[int]bool)
	code, lblLocations := linker.CompileForDebug("program.asm", "EchoKeys")
	mem := &cpu.Memory{}
	copy(mem.Data[:], code)
	debugVm := cpu.NewCPU(mem)
	go KeyboardBuffer.WriteKeyboardToBuffer(debugVm)

	file, PcToLine := debugging.DissasembleForDebugging(code, lblLocations)
	lines := strings.Split(file, "\n")

	currentLine := 0
	myApp := app.New()

	var lineBoxes []fyne.CanvasObject
	var lineBackgrounds []*canvas.Rectangle

	stepChan := make(chan struct{}, 1)
	resumeChan := make(chan struct{}, 1)

	scroll := container.NewScroll(nil)

	highlightLine := func(newIndex int, jmpWith bool) {
		if newIndex < 0 || newIndex >= len(lineBackgrounds) {
			return
		}

		lineBackgrounds[currentLine].FillColor = color.Black
		lineBackgrounds[currentLine].Refresh()

		lineBackgrounds[newIndex].FillColor = color.RGBA{B: 255, A: 255}
		lineBackgrounds[newIndex].Refresh()

		currentLine = newIndex

		if jmpWith {
			go func() {
				time.Sleep(10 * time.Millisecond)

				lineObj := lineBoxes[newIndex]
				linePos := lineObj.Position()
				lineSize := lineObj.Size()
				scrollSize := scroll.Size()

				targetY := linePos.Y - scrollSize.Height/2 + lineSize.Height/2
				if targetY < 0 {
					targetY = 0
				}

				scroll.Offset = fyne.NewPos(0, targetY)
				scroll.Refresh()
			}()
		}
	}

	for i, text := range lines {
		bg := canvas.NewRectangle(color.Black)
		label := canvas.NewText(text, color.White)
		label.TextSize = 16
		index := i

		button := widget.NewButton("", func() {
			breakpoints[index] = !breakpoints[index]
			if breakpoints[index] {
				lineBackgrounds[index].FillColor = color.RGBA{R: 180, G: 0, B: 0, A: 255}
			} else {
				lineBackgrounds[index].FillColor = color.Black
			}
			lineBackgrounds[index].Refresh()
		})
		button.Importance = widget.LowImportance

		if i == currentLine {
			bg.FillColor = color.RGBA{B: 255, A: 255}
		}

		lineBackgrounds = append(lineBackgrounds, bg)
		lineBoxes = append(lineBoxes, container.NewStack(bg, label, button))
	}

	textList := container.NewVBox(lineBoxes...)
	scroll.Content = textList
	scroll.Refresh()

	currentMode := "Step"
	var modeButton *widget.Button

	modeButton = widget.NewButton("Mode: Step", func() {
		if currentMode == "Step" {
			currentMode = "Run"
			modeButton.SetText("Mode: Run")
			go func() { resumeChan <- struct{}{} }()
		} else {
			currentMode = "Step"
			modeButton.SetText("Mode: Step")
		}
	})

	topBar := container.NewHBox(layout.NewSpacer(), modeButton)

	var nameValueRows []fyne.CanvasObject
	nameValuePanel, lbls := initDebugRegs(nameValueRows, debugVm, reverseRegMap)

	splitView := container.NewHSplit(scroll, nameValuePanel)

	myWindow := myApp.NewWindow("Debugger UI")
	myWindow.Resize(container.NewVBox(topBar, splitView).MinSize())
	myWindow.SetContent(container.NewBorder(topBar, nil, nil, nil, splitView))

	myWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if currentMode == "Step" && ev.Name == fyne.KeyRight {
			select {
			case stepChan <- struct{}{}:
			default:
			}
		}
	})

	go func() {
		for {
			if currentMode == "Step" {
				select {
				case <-stepChan:
					debugVm.Step()
					newIndex := PcToLine[debugVm.PC]
					highlightLine(newIndex, true)
				case <-time.After(50 * time.Millisecond):
				}
			} else if currentMode == "Run" {
				if breakpoints[currentLine] {
					fmt.Println("break")
					currentMode = "Step"
					modeButton.SetText("Mode: Step")
					continue
				}
				debugVm.Step()
				newIndex := PcToLine[debugVm.PC]
				highlightLine(newIndex, false)
			}
			setRegDebug(lbls, debugVm, reverseRegMap)
			nameValuePanel.Refresh()
		}
	}()

	myWindow.ShowAndRun()
}
