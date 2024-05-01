package dashboard

import (
	"image/color"
	"log"
	"os"
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

var (
	logger                                                *log.Logger
	c                                                     *cpusimple.CPU
	CPUStatus                                             string
	br0, br1, br2, br3, br4, br5, br6, br7, br8           binding.ExternalInt
	br9, br10, br11, br12, br13, br14, br15, br16         binding.ExternalInt
	sp, pc                                                binding.ExternalInt
	br0s, br1s, br2s, br3s, br4s, br5s, br6s, br7s, br8s  *widget.Label
	br9s, br10s, br11s, br12s, br13s, br14s, br15s, br16s *widget.Label
	w                                                     fyne.Window
	status                                                string

	stackDisplay          string
	stackLabelWidget      *widget.Label
	stackHeader           *widget.Label
	memoryDisplay         string
	memoryGridLabel       *widget.Label
	memoryLabel           *widget.Label
	inputCPUClock         *widget.Entry
	loadButton            *widget.Button
	runButton             *widget.Button
	haltButton            *widget.Button
	stepButton            *widget.Button
	resetButton           *widget.Button
	pauseButton           *widget.Button
	buttonsContainer      *fyne.Container
	settingsContainer     *fyne.Container
	statusContainer       *fyne.Container
	registerContainer     *fyne.Container
	memoryContainer       *fyne.Container
	stackContainer        *fyne.Container
	cpuInternalsContainer *fyne.Container
	speedContainer        *fyne.Container
	centerContainer       *fyne.Container
)

var Console = container.NewVBox()
var ConsoleScroller = container.NewVScroll(Console)

func New(cpu *cpusimple.CPU, load func(), run func(), step func(), halt func(), reset func(), pause func()) fyne.Window {

	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	c = cpu
	a := app.NewWithID("simpleCPU")
	w = a.NewWindow("Simple CPU Simulator")
	cpu.SetRunning(true)

	inputCPUClock = widget.NewEntry()
	inputCPUClock.SetText("0")

	status = "CPU status is displayed here."

	loadButton = widget.NewButton("Load", load)
	runButton = widget.NewButton("Run", run)
	haltButton = widget.NewButton("Halt", halt)
	stepButton = widget.NewButton("Step", step)
	resetButton = widget.NewButton("Reset", reset)
	pauseButton = widget.NewButton("Pause", pause)

	registerHeader := widget.NewLabel("Registers")
	registerHeader.TextStyle.Monospace = true
	registerHeader.TextStyle.Bold = true
	registerHeader2 := widget.NewLabel("Content")
	registerHeader2.TextStyle.Monospace = true
	registerHeader2.TextStyle.Bold = true

	br0 = binding.BindInt(&cpu.Registers[0])
	br1 = binding.BindInt(&cpu.Registers[1])
	br2 = binding.BindInt(&cpu.Registers[2])
	br3 = binding.BindInt(&cpu.Registers[3])
	br4 = binding.BindInt(&cpu.Registers[4])
	br5 = binding.BindInt(&cpu.Registers[5])
	br6 = binding.BindInt(&cpu.Registers[6])
	br7 = binding.BindInt(&cpu.Registers[7])
	br8 = binding.BindInt(&cpu.Registers[8])
	br9 = binding.BindInt(&cpu.Registers[9])
	br10 = binding.BindInt(&cpu.Registers[10])
	br11 = binding.BindInt(&cpu.Registers[11])
	br12 = binding.BindInt(&cpu.Registers[12])
	br13 = binding.BindInt(&cpu.Registers[13])
	br14 = binding.BindInt(&cpu.Registers[14])
	br15 = binding.BindInt(&cpu.Registers[15])
	br16 = binding.BindInt(&cpu.Registers[16])

	br0s = widget.NewLabelWithData(binding.IntToStringWithFormat(br0, "R0: x%04x"))
	br0s.TextStyle.Monospace = true
	br1s = widget.NewLabelWithData(binding.IntToStringWithFormat(br1, "R1: x%04x"))
	br1s.TextStyle.Monospace = true
	br2s = widget.NewLabelWithData(binding.IntToStringWithFormat(br2, "R2: x%04x"))
	br2s.TextStyle.Monospace = true
	br3s = widget.NewLabelWithData(binding.IntToStringWithFormat(br3, "R3: x%04x"))
	br3s.TextStyle.Monospace = true
	br4s = widget.NewLabelWithData(binding.IntToStringWithFormat(br4, "R4: x%04x"))
	br4s.TextStyle.Monospace = true
	br5s = widget.NewLabelWithData(binding.IntToStringWithFormat(br5, "R5: x%04x"))
	br5s.TextStyle.Monospace = true
	br6s = widget.NewLabelWithData(binding.IntToStringWithFormat(br6, "R6: x%04x"))
	br6s.TextStyle.Monospace = true
	br7s = widget.NewLabelWithData(binding.IntToStringWithFormat(br7, "R7: x%04x"))
	br7s.TextStyle.Monospace = true
	br8s = widget.NewLabelWithData(binding.IntToStringWithFormat(br8, "R8: x%04x"))
	br8s.TextStyle.Monospace = true
	br9s = widget.NewLabelWithData(binding.IntToStringWithFormat(br9, "R9: x%04x"))
	br9s.TextStyle.Monospace = true
	br10s = widget.NewLabelWithData(binding.IntToStringWithFormat(br10, "R10: x%04x"))
	br10s.TextStyle.Monospace = true
	br11s = widget.NewLabelWithData(binding.IntToStringWithFormat(br11, "R11: x%04x"))
	br11s.TextStyle.Monospace = true
	br12s = widget.NewLabelWithData(binding.IntToStringWithFormat(br12, "R12: x%04x"))
	br12s.TextStyle.Monospace = true
	br13s = widget.NewLabelWithData(binding.IntToStringWithFormat(br13, "R13: x%04x"))
	br13s.TextStyle.Monospace = true
	br14s = widget.NewLabelWithData(binding.IntToStringWithFormat(br14, "R14: x%04x"))
	br14s.TextStyle.Monospace = true
	br15s = widget.NewLabelWithData(binding.IntToStringWithFormat(br15, "R15: x%04x"))
	br15s.TextStyle.Monospace = true
	br16s = widget.NewLabelWithData(binding.IntToStringWithFormat(br16, "R16: x%04x"))
	br16s.TextStyle.Monospace = true

	r0 := container.New(layout.NewHBoxLayout(), br0s)
	r1 := container.New(layout.NewHBoxLayout(), br1s)
	r2 := container.New(layout.NewHBoxLayout(), br2s)
	r3 := container.New(layout.NewHBoxLayout(), br3s)
	r4 := container.New(layout.NewHBoxLayout(), br4s)
	r5 := container.New(layout.NewHBoxLayout(), br5s)
	r6 := container.New(layout.NewHBoxLayout(), br6s)
	r7 := container.New(layout.NewHBoxLayout(), br7s)
	r8 := container.New(layout.NewHBoxLayout(), br8s)
	r9 := container.New(layout.NewHBoxLayout(), br9s)
	r10 := container.New(layout.NewHBoxLayout(), br10s)
	r11 := container.New(layout.NewHBoxLayout(), br11s)
	r12 := container.New(layout.NewHBoxLayout(), br12s)
	r13 := container.New(layout.NewHBoxLayout(), br13s)
	r14 := container.New(layout.NewHBoxLayout(), br14s)
	r15 := container.New(layout.NewHBoxLayout(), br15s)
	r16 := container.New(layout.NewHBoxLayout(), br16s)

	registerContainerCol1 := container.New(layout.NewVBoxLayout(), registerHeader, r0, r1, r2, r3, r4, r5, r6, r7, r8)
	registerContainerCol2 := container.New(layout.NewVBoxLayout(), registerHeader2, r9, r10, r11, r12, r13, r14, r15, r16)

	// Stack
	stackHeader = widget.NewLabel("Stack\nContent")
	stackHeader.TextStyle.Monospace = true
	stackHeader.TextStyle.Bold = true
	stackDisplay = cpu.GetStack()
	stackLabelWidget = widget.NewLabel(stackDisplay)
	stackLabelWidget.TextStyle.Monospace = true
	stackLabelWidget.TextStyle.Bold = true
	stackContainer = container.New(
		layout.NewVBoxLayout(),
		stackHeader,
		stackLabelWidget,
	)

	// Memory
	memoryDisplay = cpu.GetAllMemory()
	memoryLabel = widget.NewLabel("Contents of Memory:\n")
	memoryLabel.TextStyle.Monospace = true
	memoryLabel.TextStyle.Bold = true
	memoryGridLabel = widget.NewLabel(memoryDisplay)
	memoryGridLabel.TextStyle.Monospace = true
	memoryContainer = container.New(
		layout.NewVBoxLayout(),
		memoryLabel,
		memoryGridLabel,
	)

	// Speed entry
	/* speedContainer = container.New(
		layout.NewHBoxLayout(),
		//layout.NewSpacer(),
		container.NewHBox(
			canvas.NewText("Clock = ", color.Black),
			inputCPUClock,
			canvas.NewText("sec  ", color.Black),
			widget.NewButton("Save", func() {
				if s, err := strconv.ParseInt(inputCPUClock.Text, 0, 32); err == nil {
					cpu.Clock = s
				}
				//logger.Println("Clock speed input value:", cpu.Clock, " seconds")
				stringValue := strconv.FormatInt(cpu.Clock, 10)
				SetStatus("Clock set to " + stringValue + " seconds")
			})),
		canvas.NewText("Set clock speed in seconds. Zero sets clock to full speed.  ", color.Black),
		layout.NewSpacer(),
	) */

	speedContainer = container.New(
		layout.NewHBoxLayout(),
		//layout.NewSpacer(),
		canvas.NewText("Clock = ", color.Black),
		inputCPUClock,
		canvas.NewText("sec  ", color.Black),
		widget.NewButton("Save", func() {
			if s, err := strconv.ParseInt(inputCPUClock.Text, 0, 32); err == nil {
				cpu.Clock = s
			}
			stringValue := strconv.FormatInt(cpu.Clock, 10)
			SetStatus("Clock set to " + stringValue + " seconds")
		}),
		canvas.NewText("Set clock speed in seconds. Zero sets clock to full speed.  ", color.Black),
		layout.NewSpacer(),
	)

	// CPU Internals: PC, SP
	pc = binding.BindInt(&cpu.PC)
	sp = binding.BindInt(&cpu.SP)
	cpuInternalsContainer = container.New(
		layout.NewHBoxLayout(),
		//layout.NewSpacer(),
		container.NewHBox(
			widget.NewLabelWithData(binding.IntToStringWithFormat(pc, "PC: x%04x")),
			//layout.NewSpacer(),
			widget.NewLabelWithData(binding.IntToStringWithFormat(sp, "SP: x%04x")),
		),
		//layout.NewSpacer(),
	)

	buttonsContainer = container.New(layout.NewHBoxLayout(), loadButton, runButton, haltButton, stepButton, resetButton, pauseButton)
	settingsContainer = container.New(layout.NewVBoxLayout(), buttonsContainer, speedContainer, cpuInternalsContainer)
	statusContainer = container.NewVBox(ConsoleScroller)
	registerContainer = container.NewHBox(registerContainerCol1, registerContainerCol2)
	centerContainer = container.NewHBox(memoryContainer, stackContainer)

	w.SetContent(container.NewBorder(settingsContainer, statusContainer, registerContainer, nil, centerContainer))

	return w
}

func UpdateAll() {

	// Reload
	br0.Reload()
	br1.Reload()
	br2.Reload()
	br3.Reload()
	br4.Reload()
	br5.Reload()
	br6.Reload()
	br7.Reload()
	br8.Reload()
	br9.Reload()
	br10.Reload()
	br11.Reload()
	br12.Reload()
	br13.Reload()
	br14.Reload()
	br15.Reload()
	br16.Reload()
	pc.Reload()
	sp.Reload()
	stackDisplay = c.GetStack()
	stackLabelWidget.Text = stackDisplay
	memoryDisplay = c.GetAllMemory()
	memoryGridLabel.SetText(memoryDisplay)

	// Refresh
	stackLabelWidget.Refresh()
	memoryGridLabel.Refresh()
	memoryContainer.Refresh()
	stackContainer.Refresh()
	settingsContainer.Refresh()
	buttonsContainer.Refresh()
	cpuInternalsContainer.Refresh()
}

func SetStatus(s string) {
	status = s
	ConsoleWrite(status)
}

func ConsoleWrite(text string) {
	Console.Add(&canvas.Text{
		Text:      text,
		Color:     color.Black,
		TextSize:  12,
		TextStyle: fyne.TextStyle{Monospace: true},
	})

	if len(Console.Objects) > 100 {
		Console.Remove(Console.Objects[0])
	}
	delta := (Console.Size().Height - ConsoleScroller.Size().Height) - ConsoleScroller.Offset.Y

	if delta < 50 {
		ConsoleScroller.ScrollToBottom()
	}
	Console.Refresh()
}
