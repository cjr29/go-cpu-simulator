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
	cpu       = *cpusimple.NewCPU()
	logger    *log.Logger
	stepChan  = make(chan bool)
	runChan   = make(chan bool)
	pauseChan = make(chan bool)
	ticker    *time.Ticker

	/* program = []byte{
		0x05, 0x81, 0x06, 0xa0, 0x20, // SET R0=5, PUSH, SET R0=6, POP R1, R0=R0+R1
		0x11, // Halt
	} */

	/* program = []byte{
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
		0x13, 0x00, 0x30, // LOAD R0 <-- word at M[0x40]
		0x14, 0x16, // SWAP R1 with R5
		0x01, 0x17, // CMP R0 <-> R1, Set Flag
		0xa0, 0xa0, 0xa0, 0xa0, 0xa0, 0xa0, // Pop the six values off the stack before returning
		0x16,                         // RETurn
		0x11,                         // HALT
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // Filler
		0x00, 0xff, // 16-bit word Value to load in R0
	} */

	program = []byte{
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
	}

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
	cpu.SetClock(1)

	ticker = time.NewTicker(time.Duration(cpu.Clock) * time.Millisecond)

	// Set up Fyne window before trying to write to Status line!!!
	var w fyne.Window = dashboard.New(&cpu, reset, load, step, run, pause, exit)

	// Activate dashboard process
	w.ShowAndRun()

}

func load() {
	// Loads code in []program into CPU memory at index 0
	cpu.SetRunning(false)
	cpu.Reset()
	cpu.Load(program, len(program))
	cpu.Preprocess(program, uint16(len(program)))
	dashboard.SetStatus("Program loaded.")
	dashboard.UpdateAll()
}

func run() {
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		cpu.SetRunning(false)
		return
	}
	// If CPU already running, do nothing
	if cpu.GetRunning() {
		dashboard.SetStatus("Program already running.")
		dashboard.UpdateAll()
		return
	}
	cpu.SetRunning(true)
	go g_monitorCPUStatus()
	go g_Run(runChan)
	go startClock()
}

func step() {
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		cpu.SetRunning(false)
		return
	}
	cpu.SetRunning(true)
	go g_monitorCPUStatus()
	go g_Step(stepChan)
	go startClock()
}

func reset() {
	cpu.SetRunning(false)
	cpu.Reset()
	dashboard.SetStatus("CPU and memory reset.")
	dashboard.UpdateAll()
}

func pause() {
	cpu.SetRunning(false)
	dashboard.SetStatus("CPU paused. Press Run or Step to continue current program.")
	dashboard.UpdateAll()
	go g_Pause(pauseChan)
}

func exit() {
	os.Exit(0)
}

// This function is the CPU clock
// Instructions are fetched and executed so long as the ticker is running
func startClock() {
	for {
		go g_monitorCPUStatus()
		select {
		case <-pauseChan:
			logger.Println("Clock paused.")
			cpu.CPUStatus <- "Clock paused."
			dashboard.SetStatus("Clock paused.")
			dashboard.UpdateAll()
			ticker.Stop()
		case <-runChan:
			logger.Println("Clock started.")
			cpu.CPUStatus <- "Clock started."
			dashboard.SetStatus("Clock started.")
			dashboard.UpdateAll()
			ticker = time.NewTicker(time.Duration(cpu.Clock) * time.Millisecond)
		case <-stepChan:
			logger.Println("Single step.")
			cpu.CPUStatus <- "Single step."
			// Fetch and execute next instruction
			cpu.FetchInstruction(cpu.Memory)
			cpu.SetRunning(false)
			dashboard.SetStatus(fmt.Sprintf("Single step. PC = x%04x, SP = x%04x, Flag = %t", cpu.PC, cpu.SP, cpu.Flag))
			dashboard.UpdateAll()
		case <-ticker.C:
			// logger.Println("Fetch and Run instruction at", t)
			if cpu.GetRunning() {
				// Fetch and execute next instruction
				cpu.FetchInstruction(cpu.Memory)
				dashboard.UpdateAll()
			} else {
				ticker.Stop()
			}
		}
	}
}

// This function is a goroutine that watches the Running Flag in the CPU
// and updates status display periodically
// Generates excessive status scrolling. Use only for debugging
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
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func g_monitorCPUStatus() {
	// Respond when channel message is received from CPU
	s := <-cpu.CPUStatus
	//logger.Println("From channel monitor: " + s)
	dashboard.SetStatus("From channel monitor: " + s)
}

// Monitor dashboard Step button status
func g_Step(c chan bool) {
	// Return dashboard status message
	//logger.Println("From channel monitor: Single step.")
	dashboard.SetStatus("From channel monitor: Single step.")
	c <- true
}

// Monitor dashboard Run button status
func g_Run(c chan bool) {
	// Return dashboard status message
	//logger.Println("From channel monitor: Running loaded program ...")
	dashboard.SetStatus("From channel monitor: Running loaded program ...")
	c <- true
}

// Monitor dashboard Pause button status
func g_Pause(c chan bool) {
	// Return dashboard status message
	//logger.Println("From channel monitor: Program paused.")
	dashboard.SetStatus("From channel monitor: Program paused.")
	c <- true
}
