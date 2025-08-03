package main

import (
	"MxsxllBox/Assembly-process/linker"
	"MxsxllBox/IO/KeyboardBuffer"
	cpu2 "MxsxllBox/VM/cpu"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	code, lblLocations := linker.CompileForDebug("program.asm", "EchoKeys")
	fmt.Println(lblLocations)
	mem := &cpu2.Memory{}
	copy(mem.Data[:], code)

	DebugVm := cpu2.NewCPU(mem)
	go KeyboardBuffer.WriteKeyboardToBuffer(DebugVm)

	myApp := app.New()
	myWindow := myApp.NewWindow("Editor with Mode Toggle")

	// Text editor
	editor := widget.NewMultiLineEntry()
	editor.SetText("You can type and highlight text here.")

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

	// Layout: button on top, editor below
	content := container.NewVBox(
		modeButton,
		editor,
	)

	myWindow.SetContent(content)
	myWindow.Resize(content.MinSize())
	myWindow.ShowAndRun()
}
