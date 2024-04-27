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
	memSize   = 128
	stackSize = 10
)

var (
	cpu = *cpusimple.NewCPU()

	program = []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1,
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

	cpu.SetMemSize(memSize)
	cpu.SetStackSize(stackSize)
	cpu.SetClock(0)

	// Set up Fyne window before trying to write to Status line!!!
	var w fyne.Window = dashboard.New(&cpu, load, run, step, halt, reset, pause)

	// Activate dashboard process
	w.ShowAndRun()

}

func load() {
	// Loads code in []program into CPU memory at index 0
	// log.Println("Entered load().")
	cpu.Reset()
	cpu.Load(program, len(program))
	cpu.Preprocess(program, len(program))
	log.Println("Program loaded")
	dashboard.SetStatus("Program loaded.")
	dashboard.UpdateAll()
	go monitorCPUStatus() // Start CPU monitor in background
}

func run() {
	log.Println("Running loaded program, standby...")
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		log.Println("ERROR: No program loaded.")
		return
	}
	if !cpu.GetRunning() {
		go monitorCPUStatus()
		cpu.SetRunning(true)
	}
	go cpu.RunFromPC(len(program))
}

func step() {
	log.Println("Single-step in program.")
	result := cpu.VerifyProgramInMemory()
	if !result {
		dashboard.SetStatus("ERROR: No program loaded.")
		log.Println("ERROR: No program loaded.")
		return
	}
	cpu.SetRunning(true)
	if cpu.PC < len(program) {
		cpu.SetRunning(true)
		cpu.FetchInstruction(cpu.Memory)
		//log.Printf("PC = x%04x, SP = %d", cpu.PC, cpu.SP)
		dashboard.SetStatus(fmt.Sprintf("Step: PC = %d, SP = %d, S[0] = %d\n", cpu.PC, cpu.SP, cpu.Stack[0]))
		dashboard.UpdateAll()
		cpu.SetRunning(false)
		// log.Println("Sleep ", cpu.Clock, " seconds")
		// time.Sleep(time.Duration(cpu.Clock) * time.Second)
	} else {
		log.Println("End of memory reached, reset and load new program, or press halt to quit application.")
		log.Printf("End of memory, R0 = %d; PC = %d, SP = %d, S[0] = %d\n", cpu.Registers[0], cpu.PC, cpu.SP, cpu.Stack[0])
		dashboard.SetStatus(fmt.Sprintf("End of memory reached, reset and load new program, or press halt. PC = %d, SP = %d, S[0] = %d\n", cpu.PC, cpu.SP, cpu.Stack[0]))
		dashboard.UpdateAll()
		dashboard.UpdateAll()
		return
	}
}

func halt() {
	log.Println("Halt button pressed.")
	dashboard.SetStatus("Halt program.")
	dashboard.UpdateAll()
	cpu.SetRunning(false)
	os.Exit(0)
}

func reset() {
	log.Println("Reset button pressed.")
	cpu.Reset()
	dashboard.SetStatus("CPU and memory reset.")
	dashboard.UpdateAll()
}

func pause() {
	log.Println("Pause button pressed.")
	cpu.SetHalt(true)
	cpu.SetRunning(false)
	dashboard.SetStatus("CPU paused. Press Run or Step to continue current program.")
	dashboard.UpdateAll()
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
