package k270emlib

import (
    "fmt"
)

// Variable IOpcodes is an array of functions that handle opcodes in the I class. Indexing into this
// array with a 4-bit number returns a function that will handle that opcode.
var IOpcodes = [16]func(*Emulator, int){
    HandleJmp,  // 0000 0
    HandleCall, // 0001 1
    HandleInt,  // 0010 2
    nil,        // 0011 3
    nil,        // 0100 4
    nil,        // 0101 5
    nil,        // 0110 6
    nil,        // 0111 7
    nil,        // 1000 8
    nil,        // 1001 9
    nil,        // 1010 10
    nil,        // 1011 11
    nil,        // 1100 12
    nil,        // 1101 13
    nil,        // 1110 14
    nil,        // 1111 15
}

// Function HandleIOpcode distributes the handling of an I opcode to the appropriate opcode handler.
func HandleIOpcode(em *Emulator, a int, i int) {
    f := IOpcodes[a]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid I opcode 0x%X", a))
    } else {
        f(em, i)
    }
}

// Function jumprel implements relative jumping and includes a call to Emulator.LogInstruction. `i`
// is the distance to jump (in 8-bit two's complement form), `opname` is the name of the
// instruction.
func jumprel(em *Emulator, i int, opname string) {
    if i & 0x80 != 0 {
        offset := (i - 0x100) * 2
        em.pc += uint16(offset)
        em.LogInstruction("%s .-0x%04X", opname, -offset)
    
    } else {
        offset := i * 2
        em.pc += uint16(offset)
        em.LogInstruction("%s .+0x%04X", opname, offset)
    }
}

// Function HandleJmp handles a JMP instruction.
func HandleJmp(em *Emulator, i int) {
    jumprel(em, i, "jmp")
}

// Function HandleCall handles a CALL instruction.
func HandleCall(em *Emulator, i int) {
    em.PushWord(em.pc)
    jumprel(em, i, "call")
}

// Function HandleInt handles an INT instruction.
func HandleInt(em *Emulator, i int) {
    em.Interrupt(uint8(i))
}
