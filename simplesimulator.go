package main

import (
	"fmt"
	"log"
	"os"

	"chrisriddick.net/cpusimple"
	"chrisriddick.net/dashboard"
	"fyne.io/fyne/v2"
)

const (
	memSize   = 128
	stackSize = 10
)

var (
	cpu     cpusimple.CPU
	program = []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1, 0x00, 0x05, 0x80, 0x01, 0x00, 0x05, 0x00, 0x05,
	}

/* 	program = []byte{
	0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
	0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
	0xc1, 0x80, 0xa1,
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

	cpu.Reset()
	cpu.Load(program, len(program))
	cpu.Preprocess(cpu.Memory[0:], len(program))

	var w fyne.Window = dashboard.New(&cpu, load, run, step, halt, reset)

	w.ShowAndRun()

}

func load() {
	// Loads code in []p into CPU memory at address a
	cpu.Load(program, len(program))
	dashboard.SetStatus("Loaded new program into memory")
	log.Println("Program loaded")
	dashboard.UpdateAll()
}

func run() {
	cpu.Run(program, len(program))
	log.Printf("Program finished. R0 = %d; PC = %d, SP = %d, S[0] = %d\n", cpu.Registers[0], cpu.PC, cpu.SP, cpu.Stack[0])
	log.Println("\nMemory:")
	for i := 0; i < len(cpu.Memory); i = i + 16 {
		fmt.Println(cpu.GetMemory(i))
	}
	log.Printf("\nStack:\n" + cpu.GetStack())
	dashboard.SetStatus(fmt.Sprintf("Program finished. R0 = %d; PC = %d, SP = %d, S[0] = %d\n", cpu.Registers[0], cpu.PC, cpu.SP, cpu.Stack[0]))
	dashboard.UpdateAll()
}

func step() {
	if cpu.PC < len(program)+1 {
		cpu.FetchInstruction(cpu.Memory)
		log.Printf("PC = %x, SP = %d", cpu.PC, cpu.SP)
		dashboard.SetStatus(fmt.Sprintf("Single step: PC = %d, SP = %d, S[0] = %d\n", cpu.PC, cpu.SP, cpu.Stack[0]))
		dashboard.UpdateAll()
	} else {
		log.Println("End of memory reached, reset and load new program, or press halt.")
		cpu.FetchInstruction(cpu.Memory)
		log.Printf("R0 = %d; PC = %d, SP = %d, S[0] = %d\n", cpu.Registers[0], cpu.PC, cpu.SP, cpu.Stack[0])
		dashboard.SetStatus(fmt.Sprintf("End of memory reached: PC = %d, SP = %d, S[0] = %d\n", cpu.PC, cpu.SP, cpu.Stack[0]))
		dashboard.SetStatus("End of memory reached, reset and load new program, or press halt.")
		dashboard.UpdateAll()
		return
	}
}

func halt() {
	log.Println("Halt button pressed")
	dashboard.SetStatus("Halt program.")
	dashboard.UpdateAll()
	os.Exit(0)
}

func reset() {
	log.Println("Reset button pressed")
	cpu.Reset()
	dashboard.SetStatus("CPU and memory reset")
	dashboard.UpdateAll()
}

func printAllMemory() {
	// Print contents of CPU Memory
	for i := 0; i < len(cpu.Memory); i = i + 16 {
		fmt.Println(cpu.GetMemory(i))
	}
}
