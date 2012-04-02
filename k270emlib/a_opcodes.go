package k270emlib

import (
    "fmt"
)

var AOpcodes = [16]func(*Emulator, int){
    HandleNot,
    HandleNeg,
    HandlePush,
    HandlePop,
    HandleShl,
    HandleAshr,
    HandleLshr,
    HandleVOpcode,
    HandleShlc,
    HandleShrc,
    HandleJr,
    HandleCr,
    HandleLdsp,
    HandleStsp,
    nil,
    nil,
}

func HandleAOpcode(em *Emulator, a int, b int) {
    f := AOpcodes[b]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid A opcode 0x%X", b))
    } else {
        f(em, a)
    }
}

func HandleNot(em *Emulator, a int) {
    before := em.GetReg(a)
    after := ^before
    em.SetReg(a, after)
    em.LogInstruction("not %s -- ~0x%02X = 0x%02X", RegisterNames[a], before, after)
}

func HandleNeg(em *Emulator, a int) {
    before := em.GetReg(a)
    after := -before
    em.SetReg(a, after)
    em.LogInstruction("neg %s -- -0x%02X = 0x%02X", RegisterNames[a], before, after)
}

func HandlePush(em *Emulator, a int) {
    v := em.GetReg(a)
    em.Push(v)
    em.LogInstruction("push %s -- value transferred was 0x%02X", RegisterNames[a], v)
}

func HandlePop(em *Emulator, a int) {
    v := em.Pop()
    em.SetReg(a, v)
    em.LogInstruction("pop %s -- value transferred was 0x%02X", RegisterNames[a], v)
}

func HandleShl(em *Emulator, a int) {
    before := em.GetReg(a)
    after := before << 1
    em.SetReg(a, after)
    em.LogInstruction("shl %s -- 0x%02X << 1 = 0x%02X", RegisterNames[a], before, after)
}

func HandleAshr(em *Emulator, a int) {
    before := em.GetReg(a)
    after := (before >> 1) | (before & 0x80)
    em.SetReg(a, after)
    em.LogInstruction("ashr %s -- 0x%02X >> 1 = 0x%02X", RegisterNames[a], before, after)
}

func HandleLshr(em *Emulator, a int) {
    before := em.GetReg(a)
    after := before >> 1
    em.SetReg(a, after)
    em.LogInstruction("lshr %s -- 0x%02X >> 1 = 0x%02X", RegisterNames[a], before, after)
}

func HandleShlc(em *Emulator, a int) {
    before := em.GetReg(a)
    c := uint8(0)
    if em.GetCarry() {c = 1}
    
    after := (before << 1) | c
    em.SetReg(a, after)
    em.LogInstruction("shlc %s -- 0x%02X << 1 = 0x%02X", RegisterNames[a], before, after)
}

func HandleShrc(em *Emulator, a int) {
    before := em.GetReg(a)
    c := uint8(0)
    if em.GetCarry() {c = 1}
    
    after := (before >> 1) | (c << 7)
    em.SetReg(a, after)
    em.LogInstruction("shrc %s -- 0x%02X >> 1 = 0x%02X", RegisterNames[a], before, after)
}

func HandleJr(em *Emulator, a int) {
    em.pc = em.GetWordReg(a)
    em.LogInstruction("jr %s -- 0x%04X", WordRegisterNames[a], em.pc)
}

func HandleCr(em *Emulator, a int) {
    em.PushWord(em.pc)
    em.pc = em.GetWordReg(a)
    em.LogInstruction("cr %s -- 0x%04X", WordRegisterNames[a], em.pc)
}

func HandleLdsp(em *Emulator, a int) {
    em.SetWordReg(a, em.sp)
    em.LogInstruction("ldsp %s -- 0x%04X", WordRegisterNames[a], em.sp)
}

func HandleStsp(em *Emulator, a int) {
    em.sp = em.GetWordReg(a)
    em.LogInstruction("stsp %s -- 0x%04X", WordRegisterNames[a], em.sp)
}
