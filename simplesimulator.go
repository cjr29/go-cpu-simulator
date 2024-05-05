package main

import (
	"fmt"

	"log"
	"os"
	"time"

	"chrisriddick.net/cpusimple"
	"chrisriddick.net/dashboard"
	"fyne.io/fyne/v2"
)

const (
	memSize   = 256
	stackSize = 16
)

var (
	cpu       = *cpusimple.NewCPU()
	logger    *log.Logger
	loadChan  = make(chan string)
	resetChan = make(chan string)
	stepChan  = make(chan string)
	runChan   = make(chan string)
	pauseChan = make(chan string)
	haltChan  = make(chan string)

	program = []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1, 0x11, 0x10,
	}
	/* program = []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1, 0x00, 0x05, 0x80, 0x01, 0x00, 0x05, 0x00, 0x05,
	} */

	/* program = []byte{
		0x00, 0x81, 0xa0, // INIT R1
		0x0a, 0x81, 0xa2, // INIT R2
		0x01, 0x81, 0xa4, // INIT R3
		0xe0, // LABEL 0
		0x80, 0xa1, 0x22, 0x81, 0xa0, 0x82, 0xa1, 0x44, 0x81, 0xa2,
		0xc1, // GOTO 0
		0x80, 0xa1,
	} */
)

func main() {

	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	//cpu.InitChan() // Initialize CPU status channel for goroutines
	cpu.CPUStatus = make(chan string)
	cpu.SetMemSize(memSize)
	cpu.SetStackSize(stackSize)
	cpu.SetClock(0)

	// Set up Fyne window before trying to write to Status line!!!
	var w fyne.Window = dashboard.New(&cpu, reset, load, step, run, pause, halt, exit)

	// Activate dashboard process
	w.ShowAndRun()

}

func load() {
	// Loads code in []program into CPU memory at index 0
	// logger.Println("Entered load().")
	cpu.Reset()
	cpu.Load(program, len(program))
	cpu.Preprocess(program, len(program))
	//logger.Println("Program loaded")
	dashboard.SetStatus("Program loaded.")
	dashboard.UpdateAll()
	//go g_monitorHalt() // Start CPU monitor in background
	//go g_monitorCPUStatus()
	go g_Load(loadChan)
}

func run() {
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		cpu.CPUStatus <- "No program is loaded."
		//logger.Println("ERROR: No program loaded.")
		return
	}
	if !cpu.GetRunning() {
		// CPU isn't running, so retstart monitoring and set running flag
		go monitorCPUStatus()
		//go g_monitorHalt()
		go g_monitorCPUStatus()
		cpu.SetRunning(true)
	}
	//logger.Println("Running loaded program, standby...")
	dashboard.SetStatus("Running loaded program ...")
	go g_Run(runChan)
	//go monitorCPUStatus()
	//go g_monitorCPUStatus()
	go cpu.RunFromPC(len(program))
}

func step() {
	logger.Println("Single-step in program.")
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		//logger.Println("ERROR: No program loaded.")
		return
	}
	//if cpu.GetHalt() {
	// If CPU is stopped, restart monitoring and set running flag
	//logger.Println("Starting goroutines for CPUStatus and Halt")
	//go monitorCPUStatus()
	//go g_monitorCPUStatus()
	//go g_monitorHalt()
	//cpu.SetRunning(true)
	//}
	if cpu.PC < len(program) {
		cpu.SetRunning(true)
		go g_Step(stepChan)
		cpu.FetchInstruction(cpu.Memory)
		//log.Printf("PC = x%04x, SP = %d", cpu.PC, cpu.SP)
		dashboard.SetStatus(fmt.Sprintf("Step: PC = %d, SP = %d, S[0] = %d", cpu.PC, cpu.SP, cpu.Stack[0]))
		dashboard.UpdateAll()
		cpu.SetRunning(false)
		// //logger.Println("Sleep ", cpu.Clock, " seconds")
	} else {
		cpu.SetRunning(false) // Stop CPU
		//logger.Println("End of memory reached, reset and load new program, or press halt to quit application.")
		//logger.Printf("End of memory, R0 = %d; PC = %d, SP = %d, S[0] = %d\n", cpu.Registers[0], cpu.PC, cpu.SP, cpu.Stack[0])
		dashboard.SetStatus(fmt.Sprintf("End of memory reached, reset and load new program, or press halt. PC = %d, SP = %d, S[0] = %d\n", cpu.PC, cpu.SP, cpu.Stack[0]))
		dashboard.UpdateAll()
		dashboard.UpdateAll()
		return
	}
}

func halt() {
	//logger.Println("Halt button pressed.")
	cpu.SetHalt(true)
	cpu.SetRunning(false)
	dashboard.SetStatus("CPU is halted. User stopped program.")
	dashboard.UpdateAll()
	//cpu.CPUStatus <- "User stopped program."
	go g_Halt(haltChan)
}

func reset() {
	logger.Println("Reset button pressed.")
	cpu.Reset()
	dashboard.SetStatus("CPU and memory reset.")
	dashboard.UpdateAll()
	go g_Reset(resetChan)
}

func pause() {
	logger.Println("Pause button pressed.")
	cpu.SetRunning(false)
	dashboard.SetStatus("CPU paused. Press Run or Step to continue current program.")
	dashboard.UpdateAll()
	//cpu.CPUStatus <- "CPU is paused."
	go g_Pause(pauseChan)
}

func exit() {
	logger.Println("User exited simulator.")
	os.Exit(0)
}

// This function is a goroutine that watches the Running Flag in the CPU
// and updates status display periodically
func monitorCPUStatus() {
	for {
		t := time.Now()
		disptime := t.Format(time.TimeOnly)
		if cpu.GetRunning() {
			dashboard.SetStatus("CPU is running, " + disptime)
			dashboard.UpdateAll()
		} else {
			dashboard.SetStatus("CPU is not running, " + disptime)
			dashboard.UpdateAll()
			return
		}
		time.Sleep(time.Duration(3) * time.Second)
	}
}

func g_monitorCPUStatus() {
	// Respond when channel message is received from CPU
	s := <-cpu.CPUStatus
	logger.Println("From channel monitor: " + s)
	dashboard.SetStatus("From channel monitor: " + s)
}

// Monitor dashboard Reset button status
func g_Reset(c chan string) {
	// Return dashboard status message
	logger.Println("From channel monitor: CPU and memory reset.")
	dashboard.SetStatus("From channel monitor: CPU and memory reset.")
	c <- "From channel monitor: CPU and memory reset."
}

// Monitor dashboard Load button status
func g_Load(c chan string) {
	// Return dashboard status message
	logger.Println("From channel monitor: Program loaded.")
	dashboard.SetStatus("From channel monitor: Program loaded.")
	c <- "From channel monitor: Program loaded."
}

// Monitor dashboard Step button status
func g_Step(c chan string) {
	// Return dashboard status message
	logger.Println("From channel monitor: Single step.")
	dashboard.SetStatus("From channel monitor: Single step.")
	c <- "From channel monitor: Single step."
}

// Monitor dashboard Run button status
func g_Run(c chan string) {
	// Return dashboard status message
	logger.Println("From channel monitor: Running loaded program ...")
	dashboard.SetStatus("From channel monitor: Running loaded program ...")
	c <- "From channel monitor: Running loaded program ..."
}

// Monitor dashboard Pause button status
func g_Pause(c chan string) {
	// Return dashboard status message
	logger.Println("From channel monitor: Program paused.")
	dashboard.SetStatus("From channel monitor: Program paused.")
	c <- "From channel monitor: Program paused."
}

// Monitor dashboard Halt button status
func g_Halt(c chan string) {
	// Return dashboard status message
	logger.Println("From channel monitor: Program halted.")
	dashboard.SetStatus("From channel monitor: Program halted.")
	c <- "From channel monitor: Program halted."
}

/* func g_monitorHalt() {
	// Respond when channel message is received from CPU
	s := <-cpu.CPUHalt
	logger.Println("From channel monitor: " + s)
	dashboard.SetStatus("From channel monitor: " + s)
} */
