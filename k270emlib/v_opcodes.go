package k270emlib

import (
    "fmt"
)

// Variable VOpcodes is an array of functions that handle opcodes in the V class. Indexing into this
// array with a 4-bit number returns a function that will handle that opcode.
var VOpcodes = [16]func(*Emulator){
    HandleRet,   // 0000 0
    HandleReti,  // 0001 1
    HandlePusha, // 0010 2
    HandlePopa,  // 0011 3
    HandleTgc,   // 0100 4
    HandleTgi,   // 0101 5
    HandleSwu,   // 0110 6
    HandleHlt,   // 0111 7
    HandleIfc,   // 1000 8
    HandleIfa,   // 1001 9
    HandleIfi,   // 1010 10
    HandleIfu,   // 1011 11
    HandleIfnc,  // 1100 12
    HandleIfna,  // 1101 13
    HandleIfni,  // 1110 14
    HandleIfnu,  // 1111 15
}

// Function HandleVOpcode distributes the handling of an V opcode to the appropriate opcode handler.
func HandleVOpcode(em *Emulator, a int) {
    f := VOpcodes[a]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid V opcode 0x%X", a))
    } else {
        f(em)
    }
}

// Function HandleRet handles a RET instruction.
func HandleRet(em *Emulator) {
    em.pc = em.PopWord()
    em.LogInstruction("ret")
}

// Function HandleReti handles a RETI instruction.
func HandleReti(em *Emulator) {
    em.pc = em.PopWord()
    em.SetInterruptsEnabled(true)
    em.LogInstruction("reti")
}

// Function HandlePusha handles a PUSHA instruction.
func HandlePusha(em *Emulator) {
    em.sc = (em.sc - 1) & 0xF
    em.Push(em.GetReg(int(em.sc)))
    em.SetCarry(em.sc == 0)
    em.LogInstruction("pusha -- %s", RegisterNames[em.sc])
}

// Function HandlePopa handles a POPA instruction.
func HandlePopa(em *Emulator) {
    em.SetReg(int(em.sc), em.Pop())
    em.sc = (em.sc + 1) & 0xF
    em.SetCarry(em.sc == 0)
    em.LogInstruction("popa -- %s", RegisterNames[em.sc])
}

// Function HandleTgc handles a TGC instruction.
func HandleTgc(em *Emulator) {
    em.SetCarry(!em.GetCarry())
    em.LogInstruction("tgc -- C now = %t", em.GetCarry())
}

// Function HandleTgi handles a TGI instruction.
func HandleTgi(em *Emulator) {
    em.SetInterruptsEnabled(!em.GetInterruptsEnabled())
    em.LogInstruction("tgi -- I now = %t", em.GetInterruptsEnabled())
}

// Function HandleSwu handles a SWU instruction.
func HandleSwu(em *Emulator) {
    em.SetUserMode(true)
    em.LogInstruction("swu")
}

// Function HandleHlt handles a HLT instruction.
func HandleHlt(em *Emulator) {
    em.SetRunning(false)
    em.LogInstruction("hlt")
}

// Function HandleIfc handles an IFC instruction.
func HandleIfc(em *Emulator) {
    if em.GetCarry() {
        em.LogInstruction("ifc -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifc -- skipping next")
    }
}

// Function HandleIfa handles an IFA instruction.
func HandleIfa(em *Emulator) {
    if em.GetAuthorised() {
        em.LogInstruction("ifa -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifa -- skipping next")
    }
}

// Function HandleIfi handles an IFI instruction.
func HandleIfi(em *Emulator) {
    if em.GetInterruptsEnabled() {
        em.LogInstruction("ifi -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifi -- skipping next")
    }
}

// Function HandleIfu handles an IFU instruction.
func HandleIfu(em *Emulator) {
    if em.GetUserMode() {
        em.LogInstruction("ifu -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifu -- skipping next")
    }
}

// Function HandleIfnc handles an IFNC instruction.
func HandleIfnc(em *Emulator) {
    if !em.GetCarry() {
        em.LogInstruction("ifnc -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifnc -- skipping next")
    }
}

// Function HandleIfna handles an IFNA instruction.
func HandleIfna(em *Emulator) {
    if !em.GetAuthorised() {
        em.LogInstruction("ifna -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifna -- skipping next")
    }
}

// Function HandleIfni handles an IFNI instruction.
func HandleIfni(em *Emulator) {
    if !em.GetInterruptsEnabled() {
        em.LogInstruction("ifni -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifni -- skipping next")
    }
}

// Function HandleIfnu handles an IFNU instruction.
func HandleIfnu(em *Emulator) {
    if !em.GetUserMode() {
        em.LogInstruction("ifnu -- executing next")
    } else {
        em.pc += 2
        em.LogInstruction("ifnu -- skipping next")
    }
}
