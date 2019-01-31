package cpu

import (
	"bytes"
	"testing"
)

func TestSum1To10(t *testing.T) {
	// Below is a hexadecimal form of a summation program, summing numbers from
	// 1 to 10
	code := []byte{
		0x00, 0x81, 0xa0, 0x0b, 0x81, 0xa2, 0x01, 0x81, 0xa4, 0xe0,
		0x80, 0xa1, 0x24, 0x81, 0xa0, 0x01, 0x24, 0x81, 0xa4, 0x42,
		0xc1, 0x80, 0xa1,
	}
	cpu := CPU{}
	res := cpu.Run(code, len(code))
	if res != 55 {
		t.Fatalf("Want: %d Got: %d", 55, res)
	}
}

func TestAlternativeSum1To10(t *testing.T) {
	// Modified summation program. Sums numbers in reverse order
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
	res := cpu.Run(code, len(code))
	if res != 55 {
		t.Fatalf("Want: %d Got: %d", 55, res)
	}
}

func TestMachineCodeGeneration(t *testing.T) {
	// Summation program in mnemonic form
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
	res := cpu.Run(generatedCode, len(generatedCode))
	if res != 55 {
		t.Fatalf("Want: %d Got: %d", 55, res)
	}
}

func TestSum1To100(t *testing.T) {
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
	res := cpu.Run(generatedCode, len(generatedCode))
	if res != 5050 {
		t.Fatalf("Want: %d Got: %d", 5050, res)
	}
}
