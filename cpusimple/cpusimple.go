package cpusimple

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
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
	Registers   [17]int
	Labels      [16]int
	PC          int // Program counter
	SP          int // Stack pointer
	Memory      []byte
	Stack       []int
	Clock       int64       // clock delay in seconds. If = 0, full speed
	HaltFlag    bool        // Halt flag used to stop CPU
	RunningFlag bool        // Indicates if CPU is executing a program
	CPUStatus   chan string // Channel for passing status to monitor goroutines
}

// FetchInstruction is a dispatcher function, which takes care of properly
// interpreting bytes as instructions and carrying those out
func (c *CPU) FetchInstruction(code []byte) {
	instruction := code[c.PC]
	c.PC++
	op := instruction & 0xe0
	switch op {
	case MaskSet: // SET
		val := instruction & 0x1f
		c.Registers[0] = int(val)
	case MaskAdd: // ADD
		reg := (instruction&0x1e)>>1 + 1
		c.Registers[0] += c.Registers[reg]
	case MaskSub: // SUB
		reg := (instruction&0x1e)>>1 + 1
		c.Registers[0] -= c.Registers[reg]
	case MaskMul: // MUL
		reg := (instruction&0x1e)>>1 + 1
		c.Registers[0] *= c.Registers[reg]
	case MaskPush: // PUSH
		opt := instruction & 0x01
		var reg byte
		if opt == 1 {
			reg = 0
		} else {
			reg = (instruction&0x1e)>>1 + 1
		}
		c.Stack[c.SP] = c.Registers[reg]
		c.SP++
	case MaskPop: // POP
		opt := instruction & 0x01
		var reg byte
		if opt == 1 {
			reg = 0
		} else {
			reg = (instruction&0x1e)>>1 + 1
		}
		c.SP--
		c.Registers[reg] = c.Stack[c.SP]
	case MaskGoto: // GOTO
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
	case MaskLabel: // LABEL
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
	c.RunningFlag = false
	c.HaltFlag = false

	for i := 0; i < len(c.Memory); i++ {
		c.Memory[i] = 0
	}
	for i := 0; i < 16; i++ {
		c.Labels[i] = 0
	}
	for i := 0; i < 17; i++ {
		c.Registers[i] = 0
	}
	for i := 0; i < len(c.Stack); i++ {
		c.Stack[i] = 0
	}
}

// Run resets the CPU, carries out a sequence of instruction and finally returns
// with a string indicating status
func (c *CPU) RunFromPC(codeLength int) {
	// Be sure that a program has been loaded
	for c.PC < codeLength {
		if !c.RunningFlag {
			c.CPUStatus <- "Program paused ..., press Step or Run to continue."
			log.Println("Program paused.")
			return
		}
		c.FetchInstruction(c.Memory[0:])
		// log.Println("RunFromPC: Sleep ", c.Clock, " seconds")
		time.Sleep(time.Duration(c.Clock) * time.Second)
		// log.Printf("RunFromPC: R0 = %d; PC = %d, SP = %d, S[0] = %d\n", c.Registers[0], c.PC, c.SP, c.Stack[0])
	}
	c.SetRunning(false)
	/* 	log.Printf("Program finished. R0 = %d; PC = %d, SP = %d, S[0] = %d\n", c.Registers[0], c.PC, c.SP, c.Stack[0])
	   	log.Println("\nMemory:\n" + c.GetAllMemory())
	   	log.Printf("\nStack:\n" + c.GetStack()) */
	c.CPUStatus <- "Program execution complete."
	log.Println("Program execution complete.")
}

// Be sure there is a program in memory
func (c *CPU) VerifyProgramInMemory() bool {
	// Be sure that a program has been loaded by testing  the first two bytes
	if (c.Memory[0] == 0) && (c.Memory[1] == 0) {
		return false
	}
	return true
}

// Load Memory with preprocessed program
func (c *CPU) Load(program []byte, programLength int) {
	for i := 0; i < programLength; i++ {
		c.Memory[i] = program[i]
	}
	// log.Println("Load: " + fmt.Sprintf("%02x ", c.Memory[99]))
}

// Translate a symbolic instruction mnemonic into a byte
func asmToByte(s string) byte {
	parts := strings.Split(s, "_")
	var b byte
	switch parts[0] {
	case "set":
		i, _ := strconv.Atoi(parts[1])
		b = MaskSet + byte(i)
	case "add":
		i, _ := strconv.Atoi(parts[1])
		i = (i - 1) << 1
		b = MaskAdd + byte(i)
	case "sub":
		i, _ := strconv.Atoi(parts[1])
		i = (i - 1) << 1
		b = MaskSub + byte(i)
	case "mul":
		i, _ := strconv.Atoi(parts[1])
		i = (i - 1) << 1
		b = MaskMul + byte(i)
	case "push":
		i, _ := strconv.Atoi(parts[1])
		if i == 0 {
			b = MaskPush + 0x01
		} else {
			i = (i - 1) << 1
			b = MaskPush + byte(i)
		}
	case "pop":
		i, _ := strconv.Atoi(parts[1])
		if i == 0 {
			b = MaskPop + 0x01
		} else {
			i = (i - 1) << 1
			b = MaskPop + byte(i)
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

// GetMemory returns a 16 byte formatted string starting at provided index
func (c *CPU) GetMemory(index int) string {
	var line string
	for i := index; i < index+16; i++ {
		line = line + fmt.Sprintf("%02x ", c.Memory[i])
	}
	return line
}

// GetAllMemory returns a 16 byte formatted string starting at provided index
func (c *CPU) GetAllMemory() string {
	var line string
	blocks := len(c.Memory) / 16
	remainder := len(c.Memory) % 16
	//log.Println("Blocks = ", blocks)
	//log.Println("Remainder = ", remainder)
	k := 0
	for j := 0; j < blocks; j++ {
		for i := k; i < k+16; i++ {
			line = line + fmt.Sprintf("%02x ", c.Memory[i])
		}
		line = line + "\n"
		k = k + 16
	}
	endBlock := blocks * 16
	for i := endBlock; i < endBlock+remainder; i++ {
		line = line + fmt.Sprintf("%02x ", c.Memory[i])
	}
	line = line + "\n"
	return line
}

// GetStack returns a formatted string of bytes beginning at SP(0)
func (c *CPU) GetStack() string {
	var s string
	//s = "\nStack\n"
	for i := 0; i < len(c.Stack); i++ {
		s = s + fmt.Sprintf("x%04x\n", c.Stack[i])
	}
	return s
}

// SetMemSize expands Memory slice to specified size and initializes to all zeros
func (c *CPU) SetMemSize(size int) {
	tempSlice := make([]byte, size)
	for i := 1; i < size; i++ {
		tempSlice[i] = 0
	}
	c.Memory = append(c.Memory, tempSlice...)
}

// SetStackSize expands Memory slice to specified size and initializes to all zeros
func (c *CPU) SetStackSize(size int) {
	tempStackSlice := make([]int, size)
	for i := 1; i < size; i++ {
		tempStackSlice[i] = 0
	}
	c.Stack = append(c.Stack, tempStackSlice...)
}

// Set CPU clock delay
func (c *CPU) SetClock(delay int64) {
	c.Clock = delay
}

// Set Halt Flag
func (c *CPU) SetHalt(flag bool) {
	c.HaltFlag = flag
}

// Get Halt Flag
func (c *CPU) GetHalt() bool {
	return c.HaltFlag
}

// Set Running Flag
func (c *CPU) SetRunning(flag bool) {
	c.RunningFlag = flag
}

// Get Running Flag
func (c *CPU) GetRunning() bool {
	return c.RunningFlag
}

// Initialize CPUStatus channel
func (c *CPU) InitChan() {
	c.CPUStatus = make(chan string) // Set up un-buffered channel for CPU status
}

func NewCPU() *CPU {
	return &CPU{}
}
