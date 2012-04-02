package k270emlib

import (
    "fmt"
)

var VOpcodes = [16]func(*Emulator){
    HandleRet,
    HandleReti,
    HandlePusha,
    HandlePopa,
    HandleTgc,
    HandleTgi,
    HandleSwu,
    HandleHlt,
    HandleIfc,
    HandleIfa,
    HandleIfi,
    HandleIfu,
    HandleIfnc,
    HandleIfna,
    HandleIfni,
    HandleIfnu,
}

func HandleVOpcode(em *Emulator, a int) {
    f := VOpcodes[a]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid V opcode 0x%X", a))
    } else {
        f(em)
    }
}

func HandleRet(em *Emulator) {
    em.pc = em.PopWord()
    em.LogInstruction("ret")
}

func HandleReti(em *Emulator) {
    em.pc = em.PopWord()
    em.SetInterruptsEnabled(true)
    em.LogInstruction("reti")
}

func HandlePusha(em *Emulator) {
    em.sc = (em.sc - 1) & 0xF
    em.Push(em.GetReg(int(em.sc)))
    em.SetCarry(em.sc == 0)
    em.LogInstruction("pusha -- %s", RegisterNames[em.sc])
}

func HandlePopa(em *Emulator) {
    em.SetReg(int(em.sc), em.Pop())
    em.sc = (em.sc + 1) & 0xF
    em.SetCarry(em.sc == 0)
    em.LogInstruction("popa -- %s", RegisterNames[em.sc])
}

func HandleTgc(em *Emulator) {
    em.SetCarry(!em.GetCarry())
    em.LogInstruction("tgc -- C now = %t", em.GetCarry())
}

func HandleTgi(em *Emulator) {
    em.SetInterruptsEnabled(!em.GetInterruptsEnabled())
    em.LogInstruction("tgi -- I now = %t", em.GetInterruptsEnabled())
}

func HandleSwu(em *Emulator) {
    em.SetUserMode(true)
    em.LogInstruction("swu")
}

func HandleHlt(em *Emulator) {
    em.SetRunning(false)
    em.LogInstruction("hlt")
}

func HandleIfc(em *Emulator) {
    if em.GetCarry() {
        em.LogInstruction("ifc -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifc -- skipping next")
    }
}

func HandleIfa(em *Emulator) {
    if em.GetAuthorised() {
        em.LogInstruction("ifa -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifa -- skipping next")
    }
}

func HandleIfi(em *Emulator) {
    if em.GetInterruptsEnabled() {
        em.LogInstruction("ifi -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifi -- skipping next")
    }
}

func HandleIfu(em *Emulator) {
    if em.GetUserMode() {
        em.LogInstruction("ifu -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifu -- skipping next")
    }
}

func HandleIfnc(em *Emulator) {
    if !em.GetCarry() {
        em.LogInstruction("ifnc -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifnc -- skipping next")
    }
}

func HandleIfna(em *Emulator) {
    if !em.GetAuthorised() {
        em.LogInstruction("ifna -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifna -- skipping next")
    }
}

func HandleIfni(em *Emulator) {
    if !em.GetInterruptsEnabled() {
        em.LogInstruction("ifni -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifni -- skipping next")
    }
}

func HandleIfnu(em *Emulator) {
    if !em.GetUserMode() {
        em.LogInstruction("ifnu -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifnu -- skipping next")
    }
}