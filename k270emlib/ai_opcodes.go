package k270emlib

import (
    "fmt"
)

var AIOpcodes = [16]func(*Emulator, int, int){
    HandleNop,
    HandleIOpcode,
    HandleRih,
    nil,
    HandleAB1Opcode,
    HandleAB2Opcode,
    HandleAdci,
    HandleSbci,
    HandleAddi,
    HandleSubi,
    HandleAndi,
    HandleOri,
    HandleXori,
    HandleLdi,
    HandleLdd,
    HandleStd,
}

func HandleAIOpcode(em *Emulator, o int, a int, i int) {
    f := AIOpcodes[o]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid AI opcode 0x%X", o))
    } else {
        f(em, a, i)
    }
}

func HandleNop(em *Emulator, a int, i int) {
    em.LogInstruction("nop")
}

func HandleRih(em *Emulator, a int, i int) {
    em.InterruptRegistryStore(uint8(i), em.GetWordReg(a))
    em.LogInstruction("rih 0x%02X, %s", i, RegisterNames[a])
}

func HandleAdci(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value + i + c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("adci %s, 0x%02X -- 0x%02X + 0x%02X + %d = 0x%02X, carry = %t", RegisterNames[a], i, a_value, i, c, r, em.GetCarry())
}

func HandleSbci(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value - i - c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("sbci %s, 0x%02X -- 0x%02X - 0x%02X - %d = 0x%02X, carry = %t", RegisterNames[a], i, a_value, i, c, r, em.GetCarry())
}

func HandleAddi(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value + i
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("addi %s, 0x%02X -- 0x%02X + 0x%02X = 0x%02X, carry = %t", RegisterNames[a], i, a_value, i, r, em.GetCarry())
}

func HandleSubi(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value - i
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("subi %s, 0x%02X -- 0x%02X - 0x%02X = 0x%02X, carry = %t", RegisterNames[a], i, a_value, i, r, em.GetCarry())
}

func HandleAndi(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value & i
    em.SetReg(a, uint8(r))
    em.LogInstruction("andi %s, 0x%02X -- 0x%02X & 0x%02X = 0x%02X", RegisterNames[a], i, a_value, i, r)
}

func HandleOri(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value | i
    em.SetReg(a, uint8(r))
    em.LogInstruction("ori %s, 0x%02X -- 0x%02X | 0x%02X = 0x%02X", RegisterNames[a], i, a_value, i, r)
}

func HandleXori(em *Emulator, a int, i int) {
    a_value := int(em.GetReg(a))
    r := a_value ^ i
    em.SetReg(a, uint8(r))
    em.LogInstruction("xori %s, 0x%02X -- 0x%02X ^ 0x%02X = 0x%02X", RegisterNames[a], i, a_value, i, r)
}

func HandleLdi(em *Emulator, a int, i int) {
    em.SetReg(a, uint8(i))
    em.LogInstruction("ldi %s, 0x%02X", RegisterNames[a], i)
}

func HandleLdd(em *Emulator, a int, i int) {
    v := em.MemoryLoad(uint16(i))
    em.SetReg(a, v)
    em.LogInstruction("ldd %s, 0x%02X -- [0x%04X] = 0x%02X", RegisterNames[a], i, i, v)
}

func HandleStd(em *Emulator, a int, i int) {
    v := em.GetReg(a)
    em.MemoryStore(uint16(i), v)
    em.LogInstruction("std 0x%02X, %s -- [0x%04X] = 0x%02X", i, RegisterNames[a], i, v)
}
