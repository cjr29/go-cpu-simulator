package dashboard

import (
	"fmt"
	"image/color"
	"strconv"

	"chrisriddick.net/cpusimple"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	c                     *cpusimple.CPU
	CPUStatus             string
	sps, pcs              *widget.Label
	w                     fyne.Window
	status                string = "CPU status is displayed here."
	stackDisplay          string
	stackLabelWidget      *widget.Label
	stackHeader           *widget.Label
	registerHeader        *widget.Label
	registerDisplay       string
	registerDisplayWidget *widget.Label
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
	exitButton            *widget.Button
	mainContainer         *fyne.Container
	buttonsContainer      *fyne.Container
	settingsContainer     *fyne.Container
	statusContainer       *fyne.Container
	registerContainer     *fyne.Container
	memoryContainer       *fyne.Container
	stackContainer        *fyne.Container
	cpuInternalsContainer *fyne.Container
	speedContainer        *fyne.Container
	centerContainer       *fyne.Container
	middleContainer       *fyne.Container
)

var Console = container.NewVBox()
var ConsoleScroller = container.NewVScroll(Console)

func New(cpu *cpusimple.CPU, reset func(), load func(), step func(), run func(), pause func(), halt func(), exit func()) fyne.Window {

	c = cpu // All data comes from the CPU structure object
	a := app.NewWithID("simpleCPU")
	w = a.NewWindow("Simple CPU Simulator")
	cpu.SetRunning(true)

	// Color backgrounds to be used in container stacks
	registerBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	stackBackground := canvas.NewRectangle(color.RGBA{R: 173, G: 219, B: 156, A: 200})
	memoryBackground := canvas.NewRectangle(color.RGBA{R: 223, G: 159, B: 173, A: 200})

	// Control buttons
	loadButton = widget.NewButton("Load", load)
	runButton = widget.NewButton("Run", run)
	haltButton = widget.NewButton("Halt", halt)
	stepButton = widget.NewButton("Step", step)
	resetButton = widget.NewButton("Reset", reset)
	pauseButton = widget.NewButton("Pause", pause)
	exitButton = widget.NewButton("Exit", exit)

	// Clock settings line
	inputCPUClock = widget.NewEntry()
	inputCPUClock.SetText("0")
	speedContainer = container.NewHBox(
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
	pcs = widget.NewLabel(fmt.Sprintf("PC: x%04x", cpu.PC))
	pcs.TextStyle.Monospace = true
	sps = widget.NewLabel(fmt.Sprintf("SP: x%04x", cpu.SP))
	sps.TextStyle.Monospace = true
	cpuInternalsContainer = container.NewHBox(
		pcs,
		sps,
	)

	// Stack
	stackHeader = widget.NewLabel("Stack\n")
	stackHeader.TextStyle.Monospace = true
	stackHeader.TextStyle.Bold = true
	stackDisplay = cpu.GetStack()
	stackLabelWidget = widget.NewLabel(stackDisplay)
	stackLabelWidget.TextStyle.Monospace = true
	stackLabelWidget.TextStyle.Bold = true
	stackContainer = container.NewStack(
		stackBackground,
		container.NewVBox(
			stackHeader,
			stackLabelWidget,
		))

	// Registers
	registerHeader = widget.NewLabel("Registers\n")
	registerHeader.TextStyle.Monospace = true
	registerHeader.TextStyle.Bold = true
	registerDisplay = cpu.GetRegisters()
	registerDisplayWidget = widget.NewLabel(registerDisplay)
	registerDisplayWidget.TextStyle.Monospace = true
	registerDisplayWidget.TextStyle.Bold = true
	registerContainer = container.NewStack(
		registerBackground,
		container.NewVBox(
			registerHeader,
			registerDisplayWidget,
		))

	// Memory
	memoryDisplay = cpu.GetAllMemory()
	memoryLabel = widget.NewLabel("Memory\n")
	memoryLabel.TextStyle.Monospace = true
	memoryLabel.TextStyle.Bold = true
	memoryGridLabel = widget.NewLabel(memoryDisplay)
	memoryGridLabel.TextStyle.Monospace = true
	memoryContainer = container.NewStack(
		memoryBackground,
		container.NewVBox(
			memoryLabel,
			memoryGridLabel,
		))

	buttonsContainer = container.NewHBox(
		resetButton,
		loadButton,
		runButton,
		stepButton,
		pauseButton,
		haltButton,
		exitButton,
	)

	settingsContainer = container.NewVBox(
		buttonsContainer,
		speedContainer,
		cpuInternalsContainer,
	)

	middleContainer = container.NewHBox(
		registerContainer,
		memoryContainer,
		stackContainer,
	)

	statusContainer = container.NewVBox(ConsoleScroller)
	registerContainer = container.NewHBox(registerContainer)
	centerContainer = container.NewHBox(memoryContainer, stackContainer)

	mainContainer = container.NewVBox(
		settingsContainer,
		middleContainer,
		statusContainer,
	)

	w.SetContent(mainContainer)

	return w
}

func UpdateAll() {

	// Reload
	pcs.SetText(fmt.Sprintf("PC: x%04x", c.PC))
	sps.SetText(fmt.Sprintf("SP: x%04x", c.SP))
	stackDisplay = c.GetStack()
	stackLabelWidget.Text = stackDisplay
	memoryDisplay = c.GetAllMemory()
	memoryGridLabel.SetText(memoryDisplay)
	registerDisplay = c.GetRegisters()
	registerDisplayWidget.Text = registerDisplay

	// Refresh
	buttonsContainer.Refresh()
	speedContainer.Refresh()
	cpuInternalsContainer.Refresh()
	settingsContainer.Refresh()
	stackLabelWidget.Refresh()
	stackContainer.Refresh()
	memoryGridLabel.Refresh()
	memoryContainer.Refresh()
	registerContainer.Refresh()
	middleContainer.Refresh()
	statusContainer.Refresh()
	centerContainer.Refresh()
	mainContainer.Refresh()
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
