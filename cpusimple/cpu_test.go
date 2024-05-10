package cpusimple

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestSum1To10(t *testing.T) {
	// Below is a hexadecimal form of a summation program, summing numbers from
	// 1 to 10
	fmt.Println("TestSum1To10")
	code := []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1,
	}

	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	cpu := CPU{}
	cpu.CPUStatus = make(chan string)
	cpu.InitMemory(100)
	cpu.InitStack(100 - 1)
	cpu.SetClock(0)
	res := cpu.Run(code, uint16(len(code)))
	if res != 55 {
		t.Fatalf("Want: %d Got: %d", 55, res)
	}
}

func TestAlternativeSum1To10(t *testing.T) {
	// Modified summation program. Sums numbers in reverse order
	fmt.Println("TestAlternativeSum1To10")
	code := []byte{
		0x00, 0x81, 0xa0, // INIT R1
		0x0a, 0x81, 0xa2, // INIT R2
		0x01, 0x81, 0xa4, // INIT R3
		0xe0, // LABEL 0
		0x80, 0xa1, 0x22, 0x81, 0xa0, 0x82, 0xa1, 0x44, 0x81, 0xa2,
		0xc1, // GOTO 0
		0x80, 0xa1,
	}
	cpu := CPU{}
	cpu.InitMemory(100)
	cpu.InitStack(100 - 1)
	cpu.SetClock(0)
	res := cpu.Run(code, uint16(len(code)))
	if res != 55 {
		t.Fatalf("Want: %d Got: %d", 55, res)
	}
}

func TestMachineCodeGeneration(t *testing.T) {
	// Summation program in mnemonic form
	fmt.Println("TestMachineCodeGeneration")
	asmCode := []string{
		"set_0",
		"push_0",
		"pop_1",
		"set_10",
		"push_0",
		"pop_2",
		"set_1",
		"push_0",
		"pop_3",
		"label_0",
		"push_1",
		"pop_0",
		"add_2",
		"push_0",
		"pop_1",
		"push_2",
		"pop_0",
		"sub_3",
		"push_0",
		"pop_2",
		"goto_0_1",
		"push_1",
		"pop_0",
	}
	generatedCode := AsmCodeToBytes(asmCode)
	code := []byte{
		0x00, 0x81, 0xa0, // INIT R1
		0x0a, 0x81, 0xa2, // INIT R2
		0x01, 0x81, 0xa4, // INIT R3
		0xe0, // LABEL 0
		0x80, 0xa1, 0x22, 0x81, 0xa0, 0x82, 0xa1, 0x44, 0x81, 0xa2,
		0xc1, // GOTO 0
		0x80, 0xa1,
	}
	if !bytes.Equal(code, generatedCode) {
		t.Fatalf("Want: %#v Got %#v", code, generatedCode)
	}
	cpu := CPU{}
	cpu.InitMemory(100)
	cpu.InitStack(100 - 1)
	cpu.SetClock(0)
	res := cpu.Run(generatedCode, uint16(len(generatedCode)))
	if res != 55 {
		t.Fatalf("Want: %d Got: %d", 55, res)
	}
}

func TestSum1To100(t *testing.T) {
	fmt.Println("TestSum1To100")
	asmCode := []string{
		"set_0",
		"push_0",
		"pop_1",
		"set_10",
		"push_0",
		"pop_4",
		"mul_4",
		"push_0",
		"pop_2",
		"set_1",
		"push_0",
		"pop_3",
		"label_0",
		"push_1",
		"pop_0",
		"add_2",
		"push_0",
		"pop_1",
		"push_2",
		"pop_0",
		"sub_3",
		"push_0",
		"pop_2",
		"goto_0_1",
		"push_1",
		"pop_0",
	}
	generatedCode := AsmCodeToBytes(asmCode)
	cpu := CPU{}
	cpu.InitMemory(100)
	cpu.InitStack(100 - 1)
	cpu.SetClock(0)
	res := cpu.Run(generatedCode, uint16(len(generatedCode)))
	if res != 5050 {
		t.Fatalf("Want: %d Got: %d", 5050, res)
	}
}

func TestSequence1To15(t *testing.T) {
	fmt.Println("TestSequence1To15")
	asmCode := []string{
		"set_1",
		"push_0",
		"pop_1",
		"set_6",
		"push_0",
		"pop_5",
		"set_15",
		"push_0",
		"pop_2",
		"set_15",
		"mul_2",
		"add_5",
		"push_0",
		"pop_2",
		"set_1",
		"push_0",
		"pop_3",
		"set_0",
		"push_0",
		"pop_4",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_2",
		"mul_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_2",
		"mul_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_0",
		"sub_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
	}
	code := AsmCodeToBytes(asmCode)
	cpu := CPU{}
	cpu.InitMemory(100)
	cpu.InitStack(100 - 1)
	cpu.SetClock(0)
	res := cpu.Run(code, uint16(len(code)))
	if res != 75 {
		t.Fatalf("Want: %d Got: %d", 75, res)
	}
}

// TODO: Needs debugging and tuning
/*********
func TestSequence1To10000(t *testing.T) {
	fmt.Println("TestSequence1To10000")
	asmCode := []string{
		"set_1",
		"push_0",
		"pop_1",
		"set_6",
		"push_0",
		"pop_5",
		"set_30",
		"push_0",
		"pop_2",
		"set_22",
		"mul_2",
		"add_5",
		"push_0",
		"pop_2",
		"set_1",
		"push_0",
		"pop_3",
		"set_0",
		"push_0",
		"pop_4",
		"label_0",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_2",
		"mul_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_2",
		"mul_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_0",
		"sub_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_2",
		"pop_0",
		"sub_3",
		"goto_0_1",
		"push_0",
		"pop_2",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_2",
		"mul_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_4",
		"pop_0",
		"add_1",
		"push_0",
		"pop_4",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"push_1",
		"pop_0",
		"add_3",
		"push_0",
		"pop_1",
		"set_2",
		"mul_1",
		"add_4",
		"push_0",
		"pop_4",
		"push_4",
		"pop_0",
	}
	code := AsmCodeToBytes(asmCode)
	cpu := CPU{}
	cpu.InitMemory(100)
	cpu.InitStack(100-1)
	cpu.SetClock(0)
	res := cpu.Run(code, len(code))
	if res != 36678337 {
		t.Fatalf("Want: %d Got: %d", 36678337, res)
	}
}
*******/
