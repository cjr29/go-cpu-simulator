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
	MEMSIZE   = uint16(256)
	STACKHEAD = MEMSIZE - 3
)

var (
	cpu         = *cpusimple.NewCPU()
	logger      *log.Logger
	stepChan    = make(chan bool)
	runChan     = make(chan bool)
	pauseChan   = make(chan bool)
	ClockChange = make(chan bool) // Used by dashboard to notify CPU the clock speed has changed
	cpuclock    *time.Ticker

	/* program = []byte{
		0x05, 0x81, 0x06, 0xa0, 0x20, // SET R0=5, PUSH, SET R0=6, POP R1, R0=R0+R1
		0x11, // Halt
	} */

	program = []byte{
		0x15, 0x00, 0x10, // CALL x0010
		0x11,                               // Halt
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // Empty buffer
		// Subroutine
		0x05, 0x81, 0x06, 0xa0, 0x20, // SET R0=5, PUSH, SET R0=6, POP R1, R0=R0+R1
		0x01, 0x81, 0x02, 0x81, // SET R0=1, PUSH, SET R0=2, PUSH
		0x03, 0x81, 0x04, 0x81, // SET R0=3, PUSH, SET R0=4, PUSH
		0x05, 0x81, 0x06, 0x81, // SET R0=5, PUSH, SET R0=6, PUSH
		0x18, 0x00, 0x63, // XSET R0 <-- 0x63
		0x12, 0x00, 0x80, // STORE R0 at M[0x80]
		0x13, 0x00, 0x40, // LOAD R5 <-- word at M[0x40]
		0x10, 0x17, // NOOP, CMP R0 <-> R1, Set Flag
		0x14, 0x15, // SWAP R1 with R5
		0xa6, 0xa6, 0xa6, 0xa6, 0xa6, 0xa6, // Pop the six values off the stack before returning
		0x16,                         // RETurn
		0x11,                         // HALT
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // Filler
		0x00, 0x05, // 16-bit word Value at address x0040 to load in R0
	}

	/* program = []byte{
		0x05, 0x81, 0x06, 0xa0, 0x20, // SET R0=5, PUSH, SET R0=6, POP R1, R0=R0+R1
		0x01, 0x81, 0x02, 0x81, // SET R0=1, PUSH, SET R0=2, PUSH
		0x03, 0x81, 0x04, 0x81, // SET R0=3, PUSH, SET R0=4, PUSH
		0x05, 0x81, 0x06, 0x81, // SET R0=5, PUSH, SET R0=6, PUSH
		0x18, 0x00, 0x63, // XSET R0 <-- 0x63
		0x12, 0x00, 0x80, // STORE R0 at M[0x80]
		0x13, 0x00, 0x30, // LOAD R0 <-- word at M[0x30]
		0x14, 0x16, // SWAP R1 with R6
		0x01, 0x17, // CMP R0 <-> R1, Set Flag
		0x11,                                              // HALT
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // Filler
		0x00, 0xff, // 16-bit word Value to load in R0
	} */

	/* program = []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1, 0x10, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x11,
	} */

/* program = []byte{
	0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
	0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
	0xc1, 0x80, 0xa1, 0x00, 0x05, 0x80, 0x01, 0x00, 0x05, 0x00, 0x05,
} */

/*
	 program = []byte{
		0x00, 0x81, 0xa0, // INIT R1
		0x0a, 0x81, 0xa2, // INIT R2
		0x01, 0x81, 0xa4, // INIT R3
		0xe0, // LABEL 0
		0x80, 0xa1, 0x22, 0x81, 0xa0, 0x82, 0xa1, 0x44, 0x81, 0xa2,
		0xc1, // GOTO 0
		0x80, 0xa1,
	}
*/
)

func main() {

	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	os.Setenv("FYNE_THEME", "light")

	cpu.CPUStatus = make(chan string, 10)
	cpu.InitMemory(MEMSIZE)
	cpu.InitStack(STACKHEAD)
	cpu.SetClock(1)         // Default to no delay
	go g_monitorCPUStatus() // Set up background CPU monitor

	cpuclock = time.NewTicker(time.Duration(cpu.Clock) * time.Millisecond)
	// Set up Fyne window before trying to write to Status line!!!
	var w fyne.Window = dashboard.New(&cpu, reset, load, step, run, pause, exit)

	go clock()

	// Activate dashboard process
	w.ShowAndRun()

}

func load() {
	// Loads code in []program into CPU memory at index 0
	cpu.Reset()
	cpu.Load(program, len(program))
	cpu.Preprocess(program, uint16(len(program)))
	dashboard.SetStatus("Program loaded.")
	dashboard.UpdateAll()
	cpu.RunFlag = false
}

func run() {
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		return
	}
	cpu.RunFlag = true
	go g_monitorCPUStatus() // Set up background CPU monitor
	go g_Run(runChan)
	go clock()
}

func step() {
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		return
	}
	go g_monitorCPUStatus() // Set up background CPU monitor
	go g_Step(stepChan)     // Must do as goroutine because dashboard blocks
	go clock()
}

func reset() {
	cpu.Reset()
	dashboard.SetStatus("CPU and memory reset.")
	dashboard.UpdateAll()
	cpuclock.Stop()
	cpu.RunFlag = false
}

func pause() {
	dashboard.SetStatus("CPU paused. Press Run or Step to continue current program.")
	dashboard.UpdateAll()
	cpu.RunFlag = true
	cpuclock.Stop()
	go g_Pause(pauseChan)
}

func exit() {
	os.Exit(0)
}

// This function is the CPU clock
// Instructions are fetched and executed so long as the ticker is running
// Control is by channel messages from buttons in GUI
func clock() {
	for {
		select {
		case <-pauseChan:
			cpu.RunFlag = false
			dashboard.SetStatus("Clock paused.")
			dashboard.UpdateAll()
			cpuclock.Stop()
		case <-runChan:
			cpu.RunFlag = true
			dashboard.SetStatus("Clock started.")
			dashboard.UpdateAll()
			cpuclock = time.NewTicker(time.Duration(cpu.Clock) * time.Millisecond)
		case <-stepChan:
			// Fetch and execute next instruction only
			cpu.RunFlag = false
			cpu.FetchInstruction(cpu.Memory)
			dashboard.SetStatus(fmt.Sprintf("Single step. PC = x%04x, SP = x%04x, Flag = %t", cpu.PC, cpu.SP, cpu.Flag))
			dashboard.UpdateAll()
			cpuclock.Stop()
		case <-cpuclock.C: // Execute next instruction of running and not paused
			if cpu.RunFlag {
				// Fetch and execute next instruction loop
				cpu.FetchInstruction(cpu.Memory)
				dashboard.UpdateAll()
			} else {
				cpuclock.Stop()
			}
		}
	}
}

func g_monitorCPUStatus() {
	// Respond when channel message is received from CPU
	s := <-cpu.CPUStatus
	dashboard.SetStatus("From CPU channel monitor: " + s)
}

// Monitor dashboard Step button status
func g_Step(c chan bool) {
	c <- true
}

// Monitor dashboard Run button status
func g_Run(c chan bool) {
	// Return dashboard status message
	dashboard.SetStatus("From Run channel monitor: Running loaded program ...")
	c <- true
}

// Monitor dashboard Pause button status
func g_Pause(c chan bool) {
	// Return dashboard status message
	dashboard.SetStatus("From Pause channel monitor: Program paused.")
	c <- true
}

// Monitor dashboard Clock speed change channel
func g_Clock(c chan bool) {
	// Return dashboard status message
	dashboard.SetStatus("From Clock channel monitor: Clock speed changed.")
	c <- true
}
