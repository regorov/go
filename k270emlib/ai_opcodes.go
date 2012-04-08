package k270emlib

import (
    "fmt"
)

// Variable AIOpcodes is an array of functions that handle opcodes in the AI class. Indexing into
// this array with a 4-bit number returns a function that will handle that opcode.
var AIOpcodes = [16]func(*Emulator, int, int){
    HandleNop,       // 0000 0
    HandleIOpcode,   // 0001 1
    HandleRih,       // 0010 2
    nil,             // 0011 3
    HandleAB1Opcode, // 0100 4
    HandleAB2Opcode, // 0101 5
    HandleAdci,      // 0110 6
    HandleSbci,      // 0111 7
    HandleAddi,      // 1000 8
    HandleSubi,      // 1001 9
    HandleAndi,      // 1010 10
    HandleOri,       // 1011 11
    HandleXori,      // 1100 12
    HandleLdi,       // 1101 13
    HandleLdd,       // 1110 14
    HandleStd,       // 1111 15
}

// Function HandleAIOpcode distributes the handling of an AI opcode to the appropriate opcode
// handler.
func HandleAIOpcode(em *Emulator, o int, a int, i int) {
    f := AIOpcodes[o]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid AI opcode 0x%X", o))
    } else {
        f(em, a, i)
    }
}

// Function HandleNop handles a NOP instruction.
func HandleNop(em *Emulator, a int, i int) {
    em.LogInstruction("nop")
}

// Function HandleRih handles a RIH instruction.
func HandleRih(em *Emulator, a int, i int) {
    if i < 0x80 && em.GetUserMode() {
        em.SetAuthorised(false)
    } else {
        em.SetAuthorised(true)
        em.InterruptRegistryStore(uint8(i), em.GetWordReg(a))
    }
    
    em.LogInstruction("rih 0x%02X, %s -- A = %t", i, RegisterNames[a], em.GetAuthorised())
}

// Function HandleAdci handles an ADCI instruction.
func HandleAdci(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value + i + c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("adci %s, 0x%02X -- 0x%02X + 0x%02X + %d = 0x%02X, carry = %t",
        RegisterNames[a], i, a_value, i, c, r, em.GetCarry())
}

// Function HandleSbci handles an SBCI instruction.
func HandleSbci(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value - i - c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("sbci %s, 0x%02X -- 0x%02X - 0x%02X - %d = 0x%02X, carry = %t",
        RegisterNames[a], i, a_value, i, c, r, em.GetCarry())
}

// Function HandleAddi handles an ADDI instruction.
func HandleAddi(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value + i
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("addi %s, 0x%02X -- 0x%02X + 0x%02X = 0x%02X, carry = %t", RegisterNames[a],
        i, a_value, i, r, em.GetCarry())
}

// Function HandleSubi handles a SUBI instruction.
func HandleSubi(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value - i
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("subi %s, 0x%02X -- 0x%02X - 0x%02X = 0x%02X, carry = %t", RegisterNames[a],
        i, a_value, i, r, em.GetCarry())
}

// Function HandleAndi handles an ANDI instruction.
func HandleAndi(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value & i
    em.SetReg(a, uint8(r))
    em.LogInstruction("andi %s, 0x%02X -- 0x%02X & 0x%02X = 0x%02X", RegisterNames[a], i, a_value,
        i, r)
}

// Function HandleOri handles an ORI instruction.
func HandleOri(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value | i
    em.SetReg(a, uint8(r))
    em.LogInstruction("ori %s, 0x%02X -- 0x%02X | 0x%02X = 0x%02X", RegisterNames[a], i, a_value, i,
        r)
}

// Function HandleXori handles an XORI instruction.
func HandleXori(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value ^ i
    em.SetReg(a, uint8(r))
    em.LogInstruction("xori %s, 0x%02X -- 0x%02X ^ 0x%02X = 0x%02X", RegisterNames[a], i, a_value,
        i, r)
}

// Function HandleLdi handles a LDI instruction.
func HandleLdi(em *Emulator, a int, i int) {
    em.SetReg(a, uint8(i))
    em.LogInstruction("ldi %s, 0x%02X", RegisterNames[a], i)
}

// Function HandleLdd handles a LDD instruction.
func HandleLdd(em *Emulator, a int, i int) {
    v := em.MemoryLoad(uint16(i))
    em.SetReg(a, v)
    em.LogInstruction("ldd %s, 0x%02X -- [0x%04X] = 0x%02X", RegisterNames[a], i, i, v)
}

// Function HandleStd handles a STD instruction.
func HandleStd(em *Emulator, a int, i int) {
    v := em.GetReg(a)
    em.MemoryStore(uint16(i), v)
    em.LogInstruction("std 0x%02X, %s -- [0x%04X] = 0x%02X", i, RegisterNames[a], i, v)
}
