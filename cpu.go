package cpu

import (
	"strconv"
	"strings"
)

// The constants below are masks for corresponding machine instructions
const (
	MaskSet   = 0x00
	MaskAdd   = 0x20
	MaskSub   = 0x40
	MaskMul   = 0x60
	MaskPush  = 0x80
	MaskPop   = 0xa0
	MaskGoto  = 0xc0
	MaskLabel = 0xe0
)

// CPU is the central structure representing the processor with its resources
type CPU struct {
	Registers [17]int
	Labels    [16]int
	PC        int // Program counter
	Stack     [1000]int
	SP        int // Stack pointer
}

// FetchInstruction is a dispatcher function, which takes care of properly
// interpreting bytes as instructions and carrying those out
func (c *CPU) FetchInstruction(code []byte) {
	instruction := code[c.PC]
	c.PC++
	op := instruction & 0xe0
	switch op {
	case 0x00: // SET
		val := instruction & 0x1f
		c.Registers[0] = int(val)
	case 0x20: // ADD
		reg := (instruction&0x1e)>>1 + 1
		c.Registers[0] += c.Registers[reg]
	case 0x40: // SUB
		reg := (instruction&0x1e)>>1 + 1
		c.Registers[0] -= c.Registers[reg]
	case 0x60: // MUL
		reg := (instruction&0x1e)>>1 + 1
		c.Registers[0] *= c.Registers[reg]
	case 0x80: // PUSH
		opt := instruction & 0x01
		var reg byte
		if opt == 1 {
			reg = 0
		} else {
			reg = (instruction&0x1e)>>1 + 1
		}
		c.Stack[c.SP] = c.Registers[reg]
		c.SP++
	case 0xa0: // POP
		opt := instruction & 0x01
		var reg byte
		if opt == 1 {
			reg = 0
		} else {
			reg = (instruction&0x1e)>>1 + 1
		}
		c.SP--
		c.Registers[reg] = c.Stack[c.SP]
	case 0xc0: // GOTO
		opt := instruction & 0x01
		if opt == 1 { // R0 != 0
			if c.Registers[0] != 0 {
				c.PC = c.Labels[(instruction&0x1e)>>1]
			}
		} else { // R0 == 0
			if c.Registers[0] == 0 {
				c.PC = c.Labels[(instruction&0x1e)>>1]
			}
		}
	case 0xe0: // LABEL
		break
	}
}

// Preprocess takes care of parsing labels to allow forward references in the
// code
func (c *CPU) Preprocess(code []byte, codeLength int) {
	for i := 0; i < codeLength; i++ {
		if code[i]&0xe0 == 0xe0 {
			label := (code[i] & 0x1e) >> 1
			c.Labels[label] = i + 1
		}
	}
}

func (c *CPU) Reset() {
	c.PC = 0
	c.SP = 0
	for i := 0; i < 16; i++ {
		c.Labels[i] = 0
	}
	for i := 0; i < 17; i++ {
		c.Registers[i] = 0
	}
	for i := 0; i < 1000; i++ {
		c.Stack[i] = 0
	}
}

// Run resets the CPU, carries out a sequence of instruction and finally returns
// the contents of R0
func (c *CPU) Run(code []byte, codeLength int) int {
	c.Reset()
	c.Preprocess(code, codeLength)
	for c.PC < codeLength {
		c.FetchInstruction(code)
	}
	return c.Registers[0]
}

// Translate a symbolic instruction mnemonic into a byte
func asmToByte(s string) byte {
	parts := strings.Split(s, "_")
	var b byte
	switch parts[0] {
	case "set":
		i, _ := strconv.Atoi(parts[1])
		b = 0x00 + byte(i)
	case "add":
		i, _ := strconv.Atoi(parts[1])
		i = (i - 1) << 1
		b = 0x20 + byte(i)
	case "sub":
		i, _ := strconv.Atoi(parts[1])
		i = (i - 1) << 1
		b = 0x40 + byte(i)
	case "mul":
		i, _ := strconv.Atoi(parts[1])
		i = (i - 1) << 1
		b = 0x60 + byte(i)
	case "push":
		i, _ := strconv.Atoi(parts[1])
		if i == 0 {
			b = 0x80 + 0x01
		} else {
			i = (i - 1) << 1
			b = 0x80 + byte(i)
		}
	case "pop":
		i, _ := strconv.Atoi(parts[1])
		if i == 0 {
			b = 0xa0 + 0x01
		} else {
			i = (i - 1) << 1
			b = 0xa0 + byte(i)
		}
	case "goto":
		l, _ := strconv.Atoi(parts[1])
		o, _ := strconv.Atoi(parts[2])
		b = 0xc0 + byte((l<<1)+o)
	case "label":
		l, _ := strconv.Atoi(parts[1])
		l = l << 1
		b = 0xe0 + byte(l)
	}
	return b
}

// AsmCodeToBytes translates a full program from symbolic to machine form
func AsmCodeToBytes(code []string) []byte {
	bytes := make([]byte, len(code))
	for i, asm := range code {
		bytes[i] = asmToByte(asm)
	}

	return bytes
}
