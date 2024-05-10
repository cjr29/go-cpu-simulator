package cpusimple

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
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

	// Extended instrcution set
	MaskExtended = 0xf0
	NOOP         = 0x00 // No operation, move to next PC
	HALT         = 0x01 // Stop CPU execution at current PC
	STORE        = 0x02 // Store contents of R0 at memory location addressed by next two bytes (big endian)
	LOAD         = 0x03 // Load R0 with two bytes from memory location addressed by next two bytes
	SWAP         = 0x04 // Exchange contents of R0 with R1
	SUB          = 0x05 // Jump to subroutine at PC,PC++, pushing PC+2 onto stack
	RET          = 0x06 // Return from subroutine, popping PC from stack
	CMP          = 0x07 // Compare contents of R0 with contents of R1 and set CPU Flag to true if matched or false if not
)

var logger *log.Logger

// CPU is the central structure representing the processor with its resources
type CPU struct {
	Registers   [17]uint16
	Labels      [16]uint16
	PC          uint16 // Program counter
	SP          uint16 // Stack pointer
	Flag        bool   // Processor flag
	Memory      []byte
	StackHead   uint16 // Starting index of stack in Memory array
	StackSize   uint16
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
	opt := instruction & MaskExtended
	if opt == 0x10 {
		// Handle extended instruction
		logger.Println("Extended instruction set op code. Handle it.")
		c.ProcessExtendedOpCode(instruction)
		return
	}
	// Not using Extended Instruction Set, Proceed with original instruction set
	op := instruction & 0xe0
	switch op {
	case MaskSet: // SET
		val := instruction & 0x1f
		c.Registers[0] = uint16(val)
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
		//		c.Stack[c.SP] = c.Registers[reg]
		c.pushRegOnStack(reg)
		//		c.SP++
	case MaskPop: // POP
		opt := instruction & 0x01
		var reg byte
		if opt == 1 {
			reg = 0
		} else {
			reg = (instruction&0x1e)>>1 + 1
		}
		c.popRegFromStack(reg)
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

func (c *CPU) ProcessExtendedOpCode(instruction byte) {
	logger.Println("ProcessExtendedOpCode")
	op := instruction & 0x0f
	switch op {
	case HALT:
		logger.Println("HALT instruction")
		c.RunningFlag = false
		c.HaltFlag = true
	case NOOP:
		logger.Println("NOOP instruction")
	case STORE:
		// Stores the two bytes of R0 at the location specified by the next two bytes from the PC
		// Use Big-Endian for storing.
		c.Memory[c.PC] = byte(c.Registers[0] >> 8)
		c.Memory[c.PC+1] = byte(c.Registers[0] & 0x00ff)
		c.PC = c.PC + 2
		logger.Println("STORE instruction")
	case LOAD:
		logger.Println("LOAD instruction")
		// Loads the two bytes starting at PC into R0
		hibyte := uint16(c.Memory[c.PC])
		lobyte := uint16(c.Memory[c.PC+1])
		c.Registers[0] = (hibyte << 8) + lobyte
		c.PC = c.PC + 2
	case SWAP:
		logger.Println("SWAP instruction")
		// Exchange the contents of R0 with R1 and vice versa
		temp := c.Registers[0]
		c.Registers[0] = c.Registers[1]
		c.Registers[1] = temp
	case SUB:
		logger.Println("SUB instruction")
		// Jump to subroutine at address pointed to by PC,PC++, pushing PC+2 onto stack

	case RET:
		logger.Println("RET instruction")
		// Return from subroutine by popping the return address off the stack and setting PC to that value
	case CMP:
		logger.Println("CMP instruction")
		// Compare contents of R0 with contents of R1 and set CPU Flag to true if matched or false if not
		if c.Registers[0] == c.Registers[1] {
			c.Flag = true
		} else {
			c.Flag = false
		}
	default:
		logger.Println("Undefined extended instruction")
	}
}

// Preprocess takes care of parsing labels to allow forward references in the
// code
func (c *CPU) Preprocess(code []byte, codeLength uint16) {
	var i uint16
	for i = 0; i < codeLength; i++ {
		if code[i]&0xe0 == 0xe0 {
			label := (code[i] & 0x1e) >> 1
			c.Labels[label] = i + 1
		}
	}
}

func (c *CPU) Reset() {
	c.PC = 0
	c.SP = c.StackHead + 2
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
}

// Run resets the CPU, carries out a sequence of instruction and finally returns
// with a string indicating status
func (c *CPU) RunFromPC(codeLength int) {
	for c.PC < uint16(len(c.Memory)) {
		if !c.RunningFlag {
			logger.Println("CPU exiting run loop.")
			c.CPUStatus <- "CPU exiting run loop."
			return
		}
		c.FetchInstruction(c.Memory[0:])
		// logger.Println("RunFromPC: Sleep ", c.Clock, " seconds")
		time.Sleep(time.Duration(c.Clock) * time.Second)
		// log.Printf("RunFromPC: R0 = %d; PC = %d, SP = %d, S[0] = %d\n", c.Registers[0], c.PC, c.SP, c.Stack[0])
	}
	c.SetRunning(false)
	c.SetHalt(true)
	// log.Printf("Program finished. R0 = %d; PC = %d, SP = %d, S[0] = %d\n", c.Registers[0], c.PC, c.SP, c.Stack[0])
	logger.Println("End of memory fault. Execution halted.")
	c.CPUStatus <- "End of memory fault. Execution halted."
}

// DO NOT DELETE! Used by go tests
// Run resets the CPU, carries out a sequence of instruction and finally returns
// the contents of R0
func (c *CPU) Run(code []byte, codeLength uint16) uint16 {
	c.Reset()
	c.Preprocess(code, codeLength)
	for c.PC < codeLength {
		c.FetchInstruction(code)
	}
	return c.Registers[0]
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
}

// Translate a symbolic instruction mnemonic into a byte
func asmToByte(s string) byte {
	//logger.Println("Asm: " + s)
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
func (c *CPU) GetMemory(index uint16) string {
	var line string
	var i uint16
	for i = index; i < index+16; i++ {
		line = line + fmt.Sprintf("%02x ", c.Memory[i])
	}
	return line
}

// GetAllMemory returns a 16 byte formatted string starting at 0000
func (c *CPU) GetAllMemory() string {
	var line string
	blocks := len(c.Memory) / 16
	remainder := len(c.Memory) % 16
	// Send header line with memory locations
	line = "       00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f\n"
	//rowHead := 0
	k := 0
	for j := 0; j < blocks; j++ {
		//rowHead = k & 0xff00
		line = line + fmt.Sprintf("%04x:  ", k)
		//		line = line + fmt.Sprintf("%04d:  ", rowHead)
		for i := k; i < k+16; i++ {
			line = line + fmt.Sprintf("%02x ", c.Memory[i])
		}
		line = line + "\n"
		k = k + 16
	}
	if k >= len(c.Memory) {
		return line
	}
	endBlock := blocks * 16
	line = line + fmt.Sprintf("%04x:  ", k)
	for i := endBlock; i < endBlock+remainder; i++ {
		line = line + fmt.Sprintf("%02x ", c.Memory[i])
	}
	line = line + "\n"
	return line
}

// GetStack returns a formatted string of 16 words (big-endian) beginning at SP down to Head
func (c *CPU) GetStack() string {
	var s string
	for i := c.SP; i <= c.StackHead; i = i + 2 {
		s = s + fmt.Sprintf("%04x: x%04x\n", i, binary.BigEndian.Uint16(c.Memory[i:]))
	}
	return s
}

// GetRegisters returns a formatted string of register values
func (c *CPU) GetRegisters() string {
	var s string
	for i := 0; i < len(c.Registers); i++ {
		s = s + fmt.Sprintf("R%02d: x%04x\n", i, c.Registers[i])
	}
	return s
}

// InitMemory expands Memory slice to specified size and initializes to all zeros
func (c *CPU) InitMemory(size uint16) {
	tempSlice := make([]byte, size)
	var i uint16
	for i = 1; i < size; i++ {
		tempSlice[i] = 0
	}
	c.Memory = append(c.Memory, tempSlice...)
}

// InitStack initializes the SP to the specified address
func (c *CPU) InitStack(loc uint16) {
	head := loc
	if head%2 != 0 {
		head = head - 1
	}
	c.SP = head
	c.StackHead = head
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

func NewCPU() *CPU {
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &CPU{}
}

/**********************
*
* Internal functions to implement instruction set
*
***********************/

// Pushes the two bytes from specified register onto stack in Big Endian format
func (c *CPU) pushRegOnStack(reg byte) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b[0:], c.Registers[reg])
	c.SP-- // Move SP to first available position
	// Push Hi byte
	c.Memory[c.SP] = b[1] // Hi byte
	c.SP--
	c.Memory[c.SP] = b[0] // Lo byte
	// SP now points to MSB of value
}

// Pops the two bytes from the stack into the specified register using the Big Endian format
func (c *CPU) popRegFromStack(reg byte) uint16 {
	// SP currently points to last value at top of stack
	rval := binary.BigEndian.Uint16(c.Memory[c.SP:])
	c.Registers[reg] = rval
	c.SP = c.SP + 2
	return rval
}
