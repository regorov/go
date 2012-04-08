package k270emlib

import (
    "fmt"
)

// Variable AOpcodes is an array of functions that handle opcodes in the A class. Indexing into this
// array with a 4-bit number returns a function that will handle that opcode.
var AOpcodes = [16]func(*Emulator, int){
    HandleNot,     // 0000 0
    HandleNeg,     // 0001 1
    HandlePush,    // 0010 2
    HandlePop,     // 0011 3
    HandleShl,     // 0100 4
    HandleAshr,    // 0101 5
    HandleLshr,    // 0110 6
    HandleVOpcode, // 0111 7
    HandleShlc,    // 1000 8
    HandleShrc,    // 1001 9
    HandleJr,      // 1010 10
    HandleCr,      // 1011 11
    HandleLdsp,    // 1100 12
    HandleStsp,    // 1101 13
    HandleRtl,     // 1110 14
    HandleRtr,     // 1111 15
}

// Function HandleAOpcode distributes the handling of an A opcode to the appropriate opcode handler.
func HandleAOpcode(em *Emulator, a int, b int) {
    f := AOpcodes[b]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid A opcode 0x%X", b))
    } else {
        f(em, a)
    }
}

// Function HandleNot handles a NOT instruction.
func HandleNot(em *Emulator, a int) {
    before := em.GetReg(a)
    after := ^before
    em.SetReg(a, after)
    em.LogInstruction("not %s -- ~0x%02X = 0x%02X", RegisterNames[a], before, after)
}

// Function HandleNeg handles a NEG instruction.
func HandleNeg(em *Emulator, a int) {
    before := em.GetReg(a)
    after := -before
    em.SetReg(a, after)
    em.LogInstruction("neg %s -- -0x%02X = 0x%02X", RegisterNames[a], before, after)
}

// Function HandlePush handles a PUSH instruction.
func HandlePush(em *Emulator, a int) {
    v := em.GetReg(a)
    em.Push(v)
    em.LogInstruction("push %s -- value transferred was 0x%02X", RegisterNames[a], v)
}

// Function HandlePop handles a POP instruction.
func HandlePop(em *Emulator, a int) {
    v := em.Pop()
    em.SetReg(a, v)
    em.LogInstruction("pop %s -- value transferred was 0x%02X", RegisterNames[a], v)
}

// Function HandleShl handles a SHL instruction.
func HandleShl(em *Emulator, a int) {
    before := em.GetReg(a)
    after := before << 1
    em.SetCarry(before & 0x80 != 0)
    em.SetReg(a, after)
    em.LogInstruction("shl %s -- 0x%02X << 1 = 0x%02X", RegisterNames[a], before, after)
}

// Function HandleAshr handles an ASHR instruction.
func HandleAshr(em *Emulator, a int) {
    before := em.GetReg(a)
    after := (before >> 1) | (before & 0x80)
    em.SetCarry(before & 1 != 0)
    em.SetReg(a, after)
    em.LogInstruction("ashr %s -- 0x%02X >> 1 = 0x%02X", RegisterNames[a], before, after)
}

// Function HandleLshr handles a LSHR instruction.
func HandleLshr(em *Emulator, a int) {
    before := em.GetReg(a)
    after := before >> 1
    em.SetCarry(before & 1 != 0)
    em.SetReg(a, after)
    em.LogInstruction("lshr %s -- 0x%02X >> 1 = 0x%02X", RegisterNames[a], before, after)
}

// Function HandleShlc handles a SHLC instruction.
func HandleShlc(em *Emulator, a int) {
    before := em.GetReg(a)
    c := uint8(0)
    if em.GetCarry() {c = 1}
    
    after := (before << 1) | c
    em.SetCarry(before & 0x80 != 0)
    em.SetReg(a, after)
    em.LogInstruction("shlc %s -- 0x%02X << 1 = 0x%02X", RegisterNames[a], before, after)
}

// Function HandleShrc handles a SHRC instruction.
func HandleShrc(em *Emulator, a int) {
    before := em.GetReg(a)
    c := uint8(0)
    if em.GetCarry() {c = 1}
    
    after := (before >> 1) | (c << 7)
    em.SetCarry(before & 1 != 0)
    em.SetReg(a, after)
    em.LogInstruction("shrc %s -- 0x%02X >> 1 = 0x%02X", RegisterNames[a], before, after)
}

// Function HandleJr handles a JR instruction.
func HandleJr(em *Emulator, a int) {
    em.pc = em.GetWordReg(a)
    em.LogInstruction("jr %s -- 0x%04X", WordRegisterNames[a >> 1], em.pc)
}

// Function HandleCr handles a CR instruction.
func HandleCr(em *Emulator, a int) {
    em.PushWord(em.pc)
    em.pc = em.GetWordReg(a)
    em.LogInstruction("cr %s -- 0x%04X", WordRegisterNames[a >> 1], em.pc)
}

// Function HandleLdsp handles a LDSP instruction.
func HandleLdsp(em *Emulator, a int) {
    em.SetWordReg(a, em.sp)
    em.LogInstruction("ldsp %s -- 0x%04X", WordRegisterNames[a >> 1], em.sp)
}

// Function HandleStsp handles a STSP instruction.
func HandleStsp(em *Emulator, a int) {
    em.sp = em.GetWordReg(a)
    em.LogInstruction("stsp %s -- 0x%04X", WordRegisterNames[a >> 1], em.sp)
}

// Function HandleRtl handles a RTL instruction.
func HandleRtl(em *Emulator, a int) {
    before := em.GetReg(a)
    after := ((before << 1) & 0xFE) | (before >> 7)
    em.SetReg(a, after)
    em.LogInstruction("rtl %s -- 0x%02X <<< 1 = 0x%02X", RegisterNames[a], before, after)
}

// Function HandleRtr handles a RTR instruction.
func HandleRtr(em *Emulator, a int) {
    before := em.GetReg(a)
    after := (before >> 1) | ((before << 7) & 0x01)
    em.SetReg(a, after)
    em.LogInstruction("rtr %s -- 0x%02X >>> 1 = 0x%02X", RegisterNames[a], before, after)
}
