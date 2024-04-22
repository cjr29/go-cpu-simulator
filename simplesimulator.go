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
	stackSize = 1000
)

var (
	cpu     cpusimple.CPU
	program = []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1,
	}
)

func main() {

	// cpu := cpusimple.CPU{}
	cpu.SetMemSize(memSize)
	cpu.SetStackSize(stackSize)

	var w fyne.Window = dashboard.New(&cpu, load, run, step, halt, reset)

	cpu.Reset()
	cpu.Load(program, len(program))
	cpu.Preprocess(cpu.Memory[0:], len(program))

	// Print contents of CPU Memory
	for i := 0; i < len(cpu.Memory); i = i + 16 {
		fmt.Println(cpu.GetMemory(i))
	}

	w.ShowAndRun()
}

func load() {
	// Loads code in []p into CPU memory at address a
	cpu.Load(program, len(program))
}

func run() {
	log.Println("Enter run() callback function.")
	res := cpu.Run(program, len(program))
	log.Printf("R0 = %d; PC = %d\n", res, cpu.PC)
	for i := 0; i < len(cpu.Memory); i = i + 16 {
		fmt.Println(cpu.GetMemory(i))
	}
}

func step() {
	if cpu.PC < len(program) {
		log.Println("Enter step() callback.")
		cpu.FetchInstruction(cpu.Memory)
		log.Println("PC = ", cpu.PC)
	} else {
		log.Println("End of memory reached, reset and load new program, or press halt.")
		return
	}
}

func halt() {
	log.Println("Halt button pressed")
	os.Exit(0)
}

func reset() {
	log.Println("Reset button pressed")
	cpu.Reset()
}
