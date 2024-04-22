package main

import (
	"fmt"
	"log"

	"chrisriddick.net/cpusimple"
	"chrisriddick.net/dashboard"

	"fyne.io/fyne/v2"
)

const (
	memSize   = 128
	stackSize = 1000
)

func main() {

	program := []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1,
	}

	cpu := cpusimple.CPU{}
	cpu.SetMemSize(memSize)
	cpu.SetStackSize(stackSize)

	/*
	*
	* Main processing loop, exit by pressing Halt
	*
	 */

	var w fyne.Window = dashboard.New(&cpu)

	w.ShowAndRun()

	for {
		cpu.Reset()
		cpu.Load(program, len(program))
		cpu.Preprocess(cpu.Memory[0:], len(program))
		res := cpu.RunProgram(len(program))
		// Print contents of CPU Memory
		for i := 0; i < len(cpu.Memory); i = i + 16 {
			fmt.Println(cpu.GetMemory(i))
		}
		log.Println("R0 = ", res)
	}
}
