package k270emlib

import (
    "fmt"
)

var IOpcodes = [16]func(*Emulator, int){
    HandleJmp,
    HandleCall,
    HandleInt,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
    nil,
}

func HandleIOpcode(em *Emulator, a int, i int) {
    f := IOpcodes[a]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid I opcode 0x%X", a))
    } else {
        f(em, i)
    }
}

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

func HandleJmp(em *Emulator, i int) {
    jumprel(em, i, "jmp")
}

func HandleCall(em *Emulator, i int) {
    em.PushWord(em.pc)
    jumprel(em, i, "call")
}

func HandleInt(em *Emulator, i int) {
    em.SetInterruptsEnabled(false)
    em.PushWord(em.pc)
    em.pc = em.InterruptRegistryLoad(uint8(i))
}
