package dashboard

import (
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

var (
	c                                                     *cpusimple.CPU
	CPUStatus                                             string
	br0, br1, br2, br3, br4, br5, br6, br7, br8           binding.ExternalInt
	br9, br10, br11, br12, br13, br14, br15, br16         binding.ExternalInt
	br0s, br1s, br2s, br3s, br4s, br5s, br6s, br7s, br8s  *widget.Label
	br9s, br10s, br11s, br12s, br13s, br14s, br15s, br16s *widget.Label
	w                                                     fyne.Window
	status                                                string
	//statusBarBound                                        binding.ExternalString
	statusScroll *widget.TextGrid
	//statusScrollBar *fyne.Container
	// statusScroller   = container.NewVScroll
	stackDisplay     string
	stackLabelWidget *widget.Label
	memoryDisplay    string
	memoryGridLabel  *widget.Label
	memoryLabel      *widget.Label
	inputCPUClock    *widget.Entry

	//console         *fyne.Container
	//consoleScroller *fyne.Container
)

var Console = container.NewVBox()
var ConsoleScroller = container.NewVScroll(Console)

func New(cpu *cpusimple.CPU, load func(), run func(), step func(), halt func(), reset func(), pause func()) fyne.Window {

	c = cpu
	a := app.NewWithID("simpleCPU")
	w = a.NewWindow("Simple CPU Simulator")

	//statusBarBound = binding.BindString(&status)
	status = "CPU status is displayed here."
	// statusScroll = widget.NewTextGrid()
	// statusScroll.ShowLineNumbers = true
	// statusScroll.SetText(status)
	// statusScrolling := container.NewScroll(statusScroll)
	//statusScrollBar := container.NewVBox()
	//statusScroller := container.NewVScoll(statusScrollBar)

	cpu.SetRunning(true)

	inputCPUClock = widget.NewEntry()
	inputCPUClock.SetText("0")

	loadButton := widget.NewButton("Load", load)

	runButton := widget.NewButton("Run", run)

	haltButton := widget.NewButton("Halt", halt)

	stepButton := widget.NewButton("Step", step)

	resetButton := widget.NewButton("Reset", reset)

	pauseButton := widget.NewButton("Pause", pause)

	registerHeader := container.New(layout.NewHBoxLayout(), canvas.NewText("Registers", color.Black))
	registerHeader2 := container.New(layout.NewHBoxLayout(), canvas.NewText(" ", color.Black))

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
	br1s = widget.NewLabelWithData(binding.IntToStringWithFormat(br1, "R1: x%04x"))
	br2s = widget.NewLabelWithData(binding.IntToStringWithFormat(br2, "R2: x%04x"))
	br3s = widget.NewLabelWithData(binding.IntToStringWithFormat(br3, "R3: x%04x"))
	br4s = widget.NewLabelWithData(binding.IntToStringWithFormat(br4, "R4: x%04x"))
	br5s = widget.NewLabelWithData(binding.IntToStringWithFormat(br5, "R5: x%04x"))
	br6s = widget.NewLabelWithData(binding.IntToStringWithFormat(br6, "R6: x%04x"))
	br7s = widget.NewLabelWithData(binding.IntToStringWithFormat(br7, "R7: x%04x"))
	br8s = widget.NewLabelWithData(binding.IntToStringWithFormat(br8, "R8: x%04x"))
	br9s = widget.NewLabelWithData(binding.IntToStringWithFormat(br9, "R9: x%04x"))
	br10s = widget.NewLabelWithData(binding.IntToStringWithFormat(br10, "R10: x%04x"))
	br11s = widget.NewLabelWithData(binding.IntToStringWithFormat(br11, "R11: x%04x"))
	br12s = widget.NewLabelWithData(binding.IntToStringWithFormat(br12, "R12: x%04x"))
	br13s = widget.NewLabelWithData(binding.IntToStringWithFormat(br13, "R13: x%04x"))
	br14s = widget.NewLabelWithData(binding.IntToStringWithFormat(br14, "R14: x%04x"))
	br15s = widget.NewLabelWithData(binding.IntToStringWithFormat(br15, "R15: x%04x"))
	br16s = widget.NewLabelWithData(binding.IntToStringWithFormat(br16, "R16: x%04x"))

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
	stackHeader := widget.NewLabel("Stack\nContent")
	stackHeader.TextStyle.Monospace = true
	stackDisplay = cpu.GetStack()
	stackLabelWidget = widget.NewLabel(stackDisplay)
	stackContainer := container.New(layout.NewVBoxLayout(), stackHeader, stackLabelWidget)

	// Memory
	memoryDisplay = cpu.GetAllMemory()
	memoryLabel = widget.NewLabel("Contents of Memory:\n")
	memoryGridLabel = widget.NewLabel(memoryDisplay)
	memoryContainer := container.New(
		layout.NewVBoxLayout(),
		memoryLabel,
		memoryGridLabel,
	)

	// Speed entry
	speedContainer := container.New(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		container.NewHBox(
			canvas.NewText("Clock = ", color.Black),
			inputCPUClock,
			canvas.NewText("sec  ", color.Black),
			widget.NewButton("Save", func() {
				if s, err := strconv.ParseInt(inputCPUClock.Text, 0, 32); err == nil {
					cpu.Clock = s
				}
				log.Println("Clock speed input value:", cpu.Clock, " seconds")
				stringValue := strconv.FormatInt(cpu.Clock, 10)
				SetStatus("Clock set to " + stringValue + " seconds")
			})),
		canvas.NewText("Set clock speed in seconds. Zero sets clock to full speed.  ", color.Black),
		layout.NewSpacer(),
	)

	buttonsContainer := container.New(layout.NewHBoxLayout(), loadButton, runButton, haltButton, stepButton, resetButton, pauseButton)

	settingsContainer := container.New(layout.NewVBoxLayout(), buttonsContainer, speedContainer)

	//statusContainer := container.NewVBox(widget.NewLabelWithData(statusBarBound), statusScrolling)
	statusContainer := container.NewVBox(ConsoleScroller)

	registerContainer := container.NewHBox(registerContainerCol1, registerContainerCol2)

	w.SetContent(container.NewBorder(settingsContainer, statusContainer, registerContainer, stackContainer, memoryContainer))

	return w
}

func UpdateAll() {
	// log.Println("UpdateAll():\n" + memoryDisplay)
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
	stackDisplay = c.GetStack()
	stackLabelWidget.Text = stackDisplay
	memoryDisplay = c.GetAllMemory()
	memoryGridLabel.SetText(memoryDisplay)
	stackLabelWidget.Refresh()
	//statusBarBound.Reload()
}

func SetStatus(s string) {
	status = s
	// statusBarBound.Reload()
	//statusScroll.SetText(strings.TrimPrefix(statusScroll.Text()+"\n"+status, "\n"))
	//statusScroll.Refresh()
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
