# Introduction

This CPU simulator project is my way of learning the Go programming language. My interest in microprocessors led to my choice of a microprocessor simulator as the base of the project. I've found the best way to learn new programming languages is to read and work with other's code. I found a basic CPU simulator in Go by Wojciech S. Gac. The code was a single package and file with only seven instructions. Just reading Wojtek's code to understand the CPU and how the instructions worked gave me some solid understanding of Go syntax and semantics.

I decided I wanted to integrate a graphical user interface (GUI) dashboard with the simulator in order to watch the internal working of the CPU and experiment with the instructions. I selected the Fyne.io cross-platform library for its completeness and coverage of Linux, macOS, and Windows.

I added an extended instruction set as well as re-architected the memory and stack to make this more like an early microprocessor. The extended instruction set uses bit 4 as a flag to activate an extended instruction. Extended instructions may use additional bytes following the actual instruction byte to represent a memory address or value. I make use of goroutines and channels to enable live buttons to interact with the running CPU. Goroutines and channels help avoid blocking of the UI on long programs.

## Architecture description

The CPU has 17 registers, special register `R0`, which serves as an accumulator, among other things, and 16 general purpose registers, referred to as `R1`-`R16` in this document. Programs are sequences of 8-bit bytes, each byte encoding a single instruction, together with its arguments. It is assumed that after completion, the results of a program are stored in `R0`. This simple CPU uses a stack that starts at the highest even memory location available and counts downward as items are pushed onto the stack. When a PUSH is executed, the SP is first decrements by 2, and then the low byte of the target is pushed onto the stack. Next, the SP is decremented and the high byte is pushed. The SP is left pointing to the high byte of the last item pushed onto the stack. This results in a "Big Endian" storage scheme with the most significant byte at the lowest memory address.

The initial memory size is set upon initialization of the CPU structure. The PC is set to 0 and the SP is set to the highest available even address. All registers and memory are zeroed. Programs are loaded as sequence of bytes starting at address x0000. Programs should always terminate with a HALT instruction to avoid infinite loops and memory overrun errors.

## Instructions

In the description below the following symbolic bits are used:
`V` - value (bits encoding literal value)
`R` - register (registers are encoded according to the pattern: `0` - `R1`, `1` - `R2`, etc.)
`X` - ignored
`S` - switch bit (switches between two modes of an instruction)
`L` - label (used to define targets for GOTO jumps)

Note: The label `L` is not a normal machine code usage. In this context, it provides a simple means for using GOTO
instructions without having to reference a memory address. A proper assembler would translate a label in the assembly
code to a memory address to which a JUMP instruction would move the Program Counter (PC). The extended instruction set
to this simulator removes the need for the preprocessor and `LABEL` instruction. Instead, it provides for a CALL and
RETurn to implement subroutines using the stack to maintain state.

The table below describes the complete instruction set, together with bit patterns:

Instruction|Bit Pattern|Description
----------|----|-----
SET|000XVVVV|Set R0 to VVVV
ADD|001RRRRX|R0 += RRRR
SUB|010RRRRX|R0 -= RRRR
MUL|011RRRRX|R0 *= RRRR
PUSH|100RRRRS|if (S==1) {PUSH(R0)} else {PUSH(RRRR)}
POP|101RRRRS|if (S==1) {POP(R0)} else {POP(RRRR)}
GOTO|110LLLLS|if (S==1) {if (R0!=0) {GOTO LLLL}} else {IF (R0==0) {GOTO LLLL}}
LABEL|111LLLLX|Mark next instruction with label LLLL

### Extended Instruction Set
Instruction|Bit Pattern|Description
----------|----|-----
NOOP|00010000|PC++
HALT|00010001|CPU stops processing at current PC
STORE|00010010|R0 --> PC+1, PC+2 (big endian)
LOAD|00010011|R0 <-- PC+1, PC+2 (big endian)
SWAP|00010100|R0 <--> R1
CALL|00010101|SP-2, PC --> SP (big endian), PC <-- (PC+1,PC+2)
RET|00010110|PC <-- SP (big endian), SP+2
CMP|00010111|R0 compare R1, if equal, CMPFLAG true, else CMPFLAG false
XSET|00011000|R0 <-- Set R0 to value in next two bytes (big endian)

## Assembler

The code also contains a simple mnemonic translator (the function `AsmCodeToBytes`), which allows one to write programs using mnemonics instead of bare hex numbers. Below is a sample mnemonic definition of a program summing numbers from 1 to 10:

```go
sum1To10 := []string{
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
```
## GUI Dashboard

The basic CPU simulator devloped by Wojciech S. Gac ran only in a terminal. I selected this code because it was a good starting point around which I can learn Go and Fyne and build a functional GUI to run the simulator. My interests include learning Go and Fyne, but also in building useful CPU simulators to be used to learn how a CPU functions internally.

I want to be able to watch the changing registers and memory as the CPU executes each instruction. To do that, I need a comprehensive dashboard to control the CPU and display the results. While implementing the GUI, I added buttons to trigger the running of a program. There are buttons to step through a program as well as pausing and restarting it. I wanted to learn how to use Go structures and related methods to implement simple object-oriented features. The CPU becomes an OO structure with related methods. Likewise, the dashboard is another OO structure that supports methods for updating elements of the UI.

During the implementation of the UI, I realized that once a CPU entered the instruction execution loop, that function never returned until the program finished. So, even if I have a button to pause execution, the program doesn't see it until the Run() function returns. I learned about Go routines and channels to enable concurrent execution. This allows the Run() function to get a message from the dashboard that the pause button was pressed.

This re-architecting of the simulator now supports the introduction of more complex CPU designs. The present implentation has only eight instructions and does not have a NOP or HALT instruction. Running a program requires the simulator to know how long the program is in order to stop before running on indefinitely or generation a segmentation fault. Also, the simple CPU mixes bytes and ints, whereas a true CPU works on byte and words. A word is defined to be a specific number of bits, usually 16, 32, or 64. I changed the GUI display to show everything in hexadecimal rather than integers to better show what is happening inside the simulator. Also, I copy the program into memory space and run from there as a real CPU would do. The stack should also be placed in memory, and that will be a future update.

Here is a sample of what you will see on the dashboard
![Dashboard](./dashboard.png)

## Building
First be sure the latest version of golang is installed.
```
$ sudo rm -rf /usr/local/go && curl -sSL "https://go.dev/dl/go1.21.6.linux-arm64.tar.gz" | sudo tar -xz -C /usr/local
$ echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
$ source $HOME/.profile
$ go version
go version go1.21.6 linux/arm64
```
Clone the github repo for go-cpu-simulator
```
$ cd go-cpu-simulator
$ go mod tidy
$ go build simplesimulator
$ run ./simplesimulator
```

