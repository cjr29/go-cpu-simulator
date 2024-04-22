package dashboard

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"chrisriddick.net/cpusimple"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func New(cpu *cpusimple.CPU) fyne.Window {

	program := []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1,
	}
	log.Println("Length of program = ", len(program))

	a := app.NewWithID("simpleCPU")
	w := a.NewWindow("Simple CPU Simulator")

	statusBar := binding.NewString()
	statusBar.Set("CPU status is displayed here.")

	inputCPUClock := widget.NewEntry()
	inputCPUClock.SetPlaceHolder("0.0")

	runButton := widget.NewButton("Run", func() {
		statusBar.Set("Run program.")
		cpu.Step = false
		cpu.Reset()
		cpu.Load(program, len(program))
		cpu.Preprocess(cpu.Memory[0:], len(program))
		res := cpu.RunProgram(len(program))
		// Print contents of CPU Memory
		for i := 0; i < len(cpu.Memory); i = i + 16 {
			fmt.Println(cpu.GetMemory(i))
		}
		statusBar.Set("R0 = " + strconv.Itoa(res))
	})

	haltButton := widget.NewButton("Halt", func() {
		log.Println("Halt pressed")
		statusBar.Set("Halt button pressed.")
		cpu.Halt = true
		panic(0)
	})

	stepButton := widget.NewButton("Step", func() {
		log.Println("Step pressed")
		statusBar.Set("Step button pressed.")
		cpu.Step = true
	})

	resetButton := widget.NewButton("Reset", func() {
		log.Println("Reset pressed")
		statusBar.Set("Reset button pressed.")
		cpu.Reset()
	})

	registerHeader := container.New(layout.NewHBoxLayout(), canvas.NewText("Register          Value", color.Black))

	r0 := container.New(layout.NewHBoxLayout(), canvas.NewText("R0", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[0]&0xFFFF), color.Black))
	r1 := container.New(layout.NewHBoxLayout(), canvas.NewText("R1", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[1]&0xFFFF), color.Black))
	r2 := container.New(layout.NewHBoxLayout(), canvas.NewText("R2", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[2]&0xFFFF), color.Black))
	r3 := container.New(layout.NewHBoxLayout(), canvas.NewText("R3", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[3]&0xFFFF), color.Black))
	r4 := container.New(layout.NewHBoxLayout(), canvas.NewText("R4", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[4]&0xFFFF), color.Black))
	r5 := container.New(layout.NewHBoxLayout(), canvas.NewText("R5", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[5]&0xFFFF), color.Black))
	r6 := container.New(layout.NewHBoxLayout(), canvas.NewText("R6", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[6]&0xFFFF), color.Black))
	r7 := container.New(layout.NewHBoxLayout(), canvas.NewText("R7", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[7]&0xFFFF), color.Black))
	r8 := container.New(layout.NewHBoxLayout(), canvas.NewText("R8", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[8]&0xFFFF), color.Black))
	r9 := container.New(layout.NewHBoxLayout(), canvas.NewText("R9", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[9]&0xFFFF), color.Black))
	r10 := container.New(layout.NewHBoxLayout(), canvas.NewText("R10", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[10]&0xFFFF), color.Black))
	r11 := container.New(layout.NewHBoxLayout(), canvas.NewText("R11", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[11]&0xFFFF), color.Black))
	r12 := container.New(layout.NewHBoxLayout(), canvas.NewText("R12", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[12]&0xFFFF), color.Black))
	r13 := container.New(layout.NewHBoxLayout(), canvas.NewText("R13", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[13]&0xFFFF), color.Black))
	r14 := container.New(layout.NewHBoxLayout(), canvas.NewText("R14", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[14]&0xFFFF), color.Black))
	r15 := container.New(layout.NewHBoxLayout(), canvas.NewText("R15", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[15]&0xFFFF), color.Black))
	r16 := container.New(layout.NewHBoxLayout(), canvas.NewText("R16", color.Black), layout.NewSpacer(), canvas.NewText(fmt.Sprintf("%04x", cpu.Registers[16]&0xFFFF), color.Black))

	// Copy stack into string array and build list for display
	stackContent := widget.NewList(
		func() int {
			return len(cpu.Stack)
		},
		func() fyne.CanvasObject {
			stackHeader := widget.NewLabel("CPU Stack")
			return stackHeader
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			t := cpu.Stack[i]
			o.(*widget.Label).SetText(fmt.Sprintf("%04x", t&0xFFFF))
		})

	stackContainer := container.New(
		layout.NewMaxLayout(),
		stackContent,
	)

	speedContainer := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		container.NewHBox(
			inputCPUClock, widget.NewButton("Save", func() {
				log.Println("Input value:", inputCPUClock)
				if s, err := strconv.ParseFloat(inputCPUClock.Text, 64); err == nil {
					cpu.Clock = s
				}
			})),
		canvas.NewText("Set clock speed in seconds. Zero sets clock to full speed.  ", color.Black),
		canvas.NewText("    CPU Stack", color.Black),
		layout.NewSpacer(),
	)

	buttonsContainer := container.New(layout.NewHBoxLayout(), runButton, haltButton, stepButton, resetButton)

	settingsContainer := container.New(layout.NewVBoxLayout(), buttonsContainer, speedContainer)

	statusContainer := container.NewHBox(
		widget.NewLabelWithData(statusBar),
	)

	memoryGrid := widget.NewTextGridFromString("Display memory grid here")
	memoryContainer := container.New(
		layout.NewCenterLayout(),
		memoryGrid,
	)

	registerContainer := container.New(layout.NewVBoxLayout(), registerHeader, r0, r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12, r13, r14, r15, r16)

	w.SetContent(container.NewBorder(settingsContainer, statusContainer, registerContainer, stackContainer, memoryContainer))

	w.ShowAndRun()

	return w
}
