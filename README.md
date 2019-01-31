# Introduction

This CPU simulator is a re-implementation in Go of my solution to a programming contest held at my workplace recently. It's also the result of my growing interest in trying to understand low-level details of how processors carry out computation. The architecture of this particular processor is very simplistic, allowing easy understanding of its internals, at the cost of making actual programming cumbersome and tedious.

## Architecture description

The CPU has 17 registers, special register `R0`, which serves as an accumulator, among other things, and 16 general purpose registers, referred to as `R1`-`R16` in this document. Programs are sequences of 8-bit bytes, each byte encoding a single instruction, together with its arguments. The CPU also has a stack of size 1000. It is assumed that after completion the results of a program are stored in `R0`.

## Instructions

In the description below the following symbolic bits are used:
`V` - value (bits encoding literal value)
`R` - register (registers are encoded according to the pattern: `0` - `R1`, `1` - `R2`, etc.)
`X` - ignored
`S` - switch bit (switches between two modes of an instruction)
`L` - label (used to define targets for GOTO jumps)

The table below describes the complete instruction set, together with bit patterns:

Instruction|Bit Pattern|Description
----------|----|-----
SET|000VVVVV|Set R0 to VVVVV
ADD|001RRRRX|R0 += RRRR
SUB|010RRRRX|R0 -= RRRR
MUL|011RRRRX|R0 *= RRRR
PUSH|100RRRRS|if (S==1) {PUSH(R0)} else {PUSH(RRRR)}
POP|101RRRRS|if (S==1) {POP(R0)} else {POP(RRRR)}
GOTO|110LLLLS|if (S==1) {if (R0!=0) {GOTO LLLL}} else {IF (R0==0) {GOTO LLLL}}
LABEL|111LLLLX|Mark next instruction with label LLLL

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
