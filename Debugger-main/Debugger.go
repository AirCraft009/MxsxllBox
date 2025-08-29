package main

import (
	"MxsxllBox/Assembly-process/assembler"
	"MxsxllBox/Assembly-process/linker"
	"MxsxllBox/IO/KeyboardBuffer"
	"MxsxllBox/VM/cpu"
	"MxsxllBox/debugging"
	"fmt"
	"image/color"
	"os"
	_ "strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildRegisterPanel(cpuState *cpu.CPU, revRegMap map[uint8]string) (*fyne.Container, [18][2]*widget.Label) {
	var regLabels [18][2]*widget.Label
	var rows []fyne.CanvasObject

	addRow := func(left, right *widget.Label) {
		rows = append(rows, container.NewHBox(layout.NewSpacer(), left, layout.NewSpacer(), right, layout.NewSpacer()))
	}

	for i := 0; i < 18; i += 2 {
		left := widget.NewLabel(fmt.Sprintf("%s: %d", revRegMap[uint8(i)], cpuState.Registers[i]))
		right := widget.NewLabel(fmt.Sprintf("%s: %d", revRegMap[uint8(i+1)], cpuState.Registers[i+1]))
		regLabels[i/2][0] = left
		regLabels[i/2][1] = right
		addRow(left, right)
	}

	stackTop := widget.NewLabel(fmt.Sprintf("Stack-top: %d", cpuState.Mem.ReadWordStack(cpuState.SP)))
	pc := widget.NewLabel(fmt.Sprintf("PC: %d", cpuState.PC))
	regLabels[len(cpuState.Registers)/2][0] = stackTop
	regLabels[len(cpuState.Registers)/2][1] = pc
	addRow(stackTop, pc)

	cFlag := widget.NewLabel(fmt.Sprintf("C-flag: %t", cpuState.Flags.Carry))
	zFlag := widget.NewLabel(fmt.Sprintf("Z-flag: %t", cpuState.Flags.Zero))
	regLabels[len(cpuState.Registers)/2+1][0] = cFlag
	regLabels[len(cpuState.Registers)/2+1][1] = zFlag
	addRow(cFlag, zFlag)

	return container.NewVBox(rows...), regLabels
}

func updateRegisterPanel(regLabels [18][2]*widget.Label, cpuState *cpu.CPU, revRegMap map[uint8]string) {
	for i := 0; i < len(cpuState.Registers); i++ {
		row, col := i/2, i%2
		regLabels[row][col].SetText(fmt.Sprintf("%s: %d", revRegMap[uint8(i)], cpuState.Registers[i]))
	}

	lastRow := len(regLabels) - 2
	regLabels[lastRow][0].SetText(fmt.Sprintf("Stack-top: %d", cpuState.Mem.ReadWordStack(cpuState.SP)))
	regLabels[lastRow][1].SetText(fmt.Sprintf("PC: %d", cpuState.PC))

	flagRow := lastRow + 1
	regLabels[flagRow][0].SetText(fmt.Sprintf("C-flag: %t", cpuState.Flags.Carry))
	regLabels[flagRow][1].SetText(fmt.Sprintf("Z-flag: %t", cpuState.Flags.Zero))
}

func buildCodeView(lines []string, breakpoints map[int]bool) ([]fyne.CanvasObject, []*canvas.Rectangle) {
	var boxes []fyne.CanvasObject
	var backgrounds []*canvas.Rectangle

	for i, text := range lines {
		bg := canvas.NewRectangle(color.Black)
		label := canvas.NewText(text, color.White)
		label.TextSize = 16

		index := i
		button := widget.NewButton("", func() {
			breakpoints[index] = !breakpoints[index]
			if breakpoints[index] {
				backgrounds[index].FillColor = color.RGBA{R: 180, G: 0, B: 0, A: 255}
			} else {
				backgrounds[index].FillColor = color.Black
			}
			backgrounds[index].Refresh()
		})
		button.Importance = widget.LowImportance

		if i == 0 {
			bg.FillColor = color.RGBA{B: 255, A: 255}
		}

		backgrounds = append(backgrounds, bg)
		boxes = append(boxes, container.NewStack(bg, label, button))
	}

	return boxes, backgrounds
}

func main() {
	reverseRegMap := debugging.ReverseMaps(assembler.RegMap)
	breakpoints := make(map[int]bool)

	code, labels := linker.CompileForDebug("program.asm", "MxsxllOS")
	mem := cpu.NewMemory()
	copy(mem.Data[:], code)
	vm := cpu.NewDebugCpu(mem)
	go KeyboardBuffer.WriteKeyboardToBuffer(vm.Cpu)

	disasm, pcMap := debugging.DissasembleForDebugging(code, labels)
	lines := strings.Split(disasm, "\n")
	fmt.Println(lines[pcMap[616]])

	currentLine := 0
	stepChan := make(chan struct{}, 1)
	resumeChan := make(chan struct{}, 1)

	myApp := app.New()
	codeBoxes, codeBackgrounds := buildCodeView(lines, breakpoints)
	codeScroll := container.NewScroll(container.NewVBox(codeBoxes...))

	highlight := func(newIndex int, jump bool) {
		if newIndex < 0 || newIndex >= len(codeBackgrounds) {
			return
		}
		codeBackgrounds[currentLine].FillColor = color.Black
		codeBackgrounds[currentLine].Refresh()

		codeBackgrounds[newIndex].FillColor = color.RGBA{B: 255, A: 255}
		codeBackgrounds[newIndex].Refresh()

		currentLine = newIndex
		if jump {
			go func() {
				time.Sleep(10 * time.Millisecond)
				pos := codeBoxes[newIndex].Position()
				size := codeBoxes[newIndex].Size()
				viewSize := codeScroll.Size()

				y := pos.Y - viewSize.Height/2 + size.Height/2
				if y < 0 {
					y = 0
				}
				codeScroll.Offset = fyne.NewPos(0, y)
				codeScroll.Refresh()
			}()
		}
	}

	regPanel, regLabels := buildRegisterPanel(vm.Cpu, reverseRegMap)
	mode := "Step"
	var modeBtn *widget.Button
	modeBtn = widget.NewButton("Mode: Step", func() {
		if mode == "Step" {
			mode = "Run"
			modeBtn.SetText("Mode: Run")
			go func() { resumeChan <- struct{}{} }()
		} else {
			highlight(pcMap[vm.Cpu.PC], true)
			updateRegisterPanel(regLabels, vm.Cpu, reverseRegMap)
			regPanel.Refresh()
			mode = "Step"
			modeBtn.SetText("Mode: Step")
		}
	})

	topBar := container.NewHBox(layout.NewSpacer(), modeBtn)
	mainSplit := container.NewHSplit(codeScroll, regPanel)
	mainWin := myApp.NewWindow("Debugger UI")
	mainWin.Resize(container.NewVBox(topBar, mainSplit).MinSize())
	mainWin.SetContent(container.NewBorder(topBar, nil, nil, nil, mainSplit))

	mainWin.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if mode == "Step" && ev.Name == fyne.KeyRight {
			select {
			case stepChan <- struct{}{}:
			default:
			}
		}
	})

	go func() {
		for {
			if mode == "Step" {
				select {
				case <-stepChan:
					vm.StepDebug()
					highlight(pcMap[vm.Cpu.PC], true)
					updateRegisterPanel(regLabels, vm.Cpu, reverseRegMap)
					regPanel.Refresh()
				case <-time.After(50 * time.Millisecond):
				}
			} else if mode == "Run" {
				if breakpoints[currentLine] {
					highlight(currentLine, true)
					mode = "Step"
					modeBtn.SetText("Mode: Step")
					continue
				}
				vm.StepDebug()
				highlight(pcMap[vm.Cpu.PC], false)
			}
		}
	}()

	mainWin.ShowAndRun()
	os.Exit(0)
}
