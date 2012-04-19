package k270emlib

import (
    "fmt"
)

// Variable AJOpcodes is an array of functions that handle opcodes in the AJ class. Indexing into
// this array with a 2-bit number returns a function that will handle that opcode.
var AJOpcodes = [4]func(*Emulator, int, int){
    HandleNop,       // 00 0
    nil,             // 01 1
    HandleLds,       // 10 2
    HandleSts,       // 11 3
}

// Function HandleAJOpcode distributes the handling of an AJ opcode to the appropriate opcode
// handler.
func HandleAJOpcode(em *Emulator, a int, i int) {
    q := (i >> 6) & 0x3
    j := i & 0x3F
    
    f := AJOpcodes[q]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid AJ opcode 0x%X", q))
    } else {
        f(em, a, j)
    }
}

// Function HandleNop handles a NOP instruction.
func HandleNop(em *Emulator, a int, j int) {
    em.LogInstruction("nop")
    
    em.timer += 4;
}

// Function HandleLds handles a LDS instruction.
func HandleLds(em *Emulator, a int, j int) {
    addr := em.sp + uint16(j)
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.LogInstruction("lds %s, 0x%02X -- [0x%04X] = 0x%02X", RegisterNames[a], j, addr, data)
    
    em.timer += 6;
}

// Function HandleSts handles a STS instruction.
func HandleSts(em *Emulator, a int, j int) {
    addr := em.sp + uint16(j)
    data := em.GetReg(a)
    em.MemoryStore(addr, data)
    em.LogInstruction("sts 0x%02X, %s -- [0x%04X] = 0x%02X", j, RegisterNames[a], addr, data)
    
    em.timer += 5;
}
