package k270emlib

import (
    "fmt"
)

// Variable AB2Opcodes is an array of functions that handle opcodes in the AB2 class. Indexing into
// this array with a 4-bit number returns a function that will handle that opcode.
var AB2Opcodes = [16]func(*Emulator, int, int){
    HandleAOpcode, // 0000 0
    nil,           // 0001 1
    HandleLdv,     // 0010 2
    HandleStv,     // 0011 3
    HandleIn,      // 0100 4
    HandleOut,     // 0101 5
    nil,           // 0110 6
    nil,           // 0111 7
    HandleLd,      // 1000 8
    HandleLdInc,   // 1001 9
    HandleLdDec,   // 1010 10
    HandleLdOne,   // 1011 11
    HandleSt,      // 1100 12
    HandleStInc,   // 1101 13
    HandleStDec,   // 1110 14
    HandleStOne,   // 1111 15
}

// Function HandleAB2Opcode distributes the handling of an AB2 opcode to the appropriate opcode
// handler.
func HandleAB2Opcode(em *Emulator, a int, i int) {
    q := i >> 4
    b := i & 0xF
    
    f := AB2Opcodes[q]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid AB2 opcode 0x%X", q))
    } else {
        f(em, a, b)
    }
}

// Function HandleLdv handles a LDV instruction.
func HandleLdv(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b)
    data := em.VideoMemoryLoad(addr)
    em.SetReg(a, data)
    em.LogInstruction("ldv %s, %s -- VMEM[0x%04X] = 0x%02X", RegisterNames[a],
        WordRegisterNames[b >> 1], addr, data)
}

// Function HandleStv handles a STV instruction.
func HandleStv(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b)
    em.VideoMemoryStore(addr, data)
    em.LogInstruction("stv %s, %s -- VMEM[0x%04X] = 0x%02X", WordRegisterNames[b >> 1],
        RegisterNames[a], addr, data)
}

// Function HandlePand handles a PAND instruction (obsolete).
func HandlePand(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    
    if em.getPortAccess(port) {
        before := em.LoadIOPort(port)
        after := before & uint8(i)
        em.StoreIOPort(port, after)
        
        em.LogInstruction("pand %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X & 0x%02X = 0x%02X",
            RegisterNames[b], RegisterNames[a], port, before, before, i, after)
    
    } else {
        em.LogInstruction("pand %s, %s -- not authorised", RegisterNames[b], RegisterNames[a])
    }
}

// Function HandlePor handles a POR instruction (obsolete).
func HandlePor(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    
    if em.getPortAccess(port) {
        before := em.LoadIOPort(port)
        after := before | uint8(i)
        em.StoreIOPort(port, after)
        
        em.LogInstruction("por %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X | 0x%02X = 0x%02X",
            RegisterNames[b], RegisterNames[a], port, before, before, i, after)
    
    } else {
        em.LogInstruction("pand %s, %s -- not authorised", RegisterNames[b], RegisterNames[a])
    }
}

// Function HandlePxor handles a PXOR instruction (obsolete).
func HandlePxor(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    
    if em.getPortAccess(port) {
        before := em.LoadIOPort(port)
        after := before ^ uint8(i)
        em.StoreIOPort(port, after)
        
        em.LogInstruction("pxor %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X ^ 0x%02X = 0x%02X",
            RegisterNames[b], RegisterNames[a], port, before, before, i, after)
    
    } else {
        em.LogInstruction("pand %s, %s -- not authorised", RegisterNames[b], RegisterNames[a])
    }
}

// Function HandlePclr handles a PCLR instruction (obsolete).
func HandlePclr(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    
    if em.getPortAccess(port) {
        before := em.LoadIOPort(port)
        after := before & (^uint8(i))
        em.StoreIOPort(port, after)
        
        em.LogInstruction("pclr %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X & ~0x%02X = 0x%02X",
            RegisterNames[b], RegisterNames[a], port, before, before, i, after)
    
    } else {
        em.LogInstruction("pand %s, %s -- not authorised", RegisterNames[b], RegisterNames[a])
    }
}

// Function HandleIn handles an IN instruction.
func HandleIn(em *Emulator, a int, b int) {
    addr := em.GetReg(b)
    
    if em.getPortAccess(addr) {
        data := em.LoadIOPort(addr)
        em.SetReg(a, data)
        em.LogInstruction("in %s, %s -- ports[0x%02X] = 0x%02X", RegisterNames[a],
            RegisterNames[b], addr, data)
    
    } else {
        em.LogInstruction("in %s, %s -- not authorised", RegisterNames[a], RegisterNames[b])
    }
}

// Function HandleOut handles an OUT instruction.
func HandleOut(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetReg(b)
    
    if em.getPortAccess(addr) {
        em.StoreIOPort(addr, data)
        em.LogInstruction("out %s, %s -- ports[0x%02X] = 0x%02X", RegisterNames[b],
            RegisterNames[a], addr, data)
    
    } else {
        em.LogInstruction("out %s, %s -- not authorised", RegisterNames[b], RegisterNames[a])
    }
}

// Function HandleLd handles a LD instruction.
func HandleLd(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b)
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.LogInstruction("ld %s, %s -- [0x%04X] = 0x%02X", RegisterNames[a], WordRegisterNames[b >> 1],
        addr, data)
}

// Function HandleLdInc handles a LD+ instruction.
func HandleLdInc(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b)
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.SetWordReg(b, addr + 1)
    em.LogInstruction("ld %s, %s+ -- [0x%04X] = 0x%02X", RegisterNames[a],
        WordRegisterNames[b >> 1], addr, data)
}

// Function HandleLdDec handles a -LD instruction.
func HandleLdDec(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b) - 1
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.SetWordReg(b, addr)
    em.LogInstruction("ld %s, -%s -- [0x%04X] = 0x%02X", RegisterNames[a],
        WordRegisterNames[b >> 1], addr, data)
}

// Function HandleLdOne handles a LD+1 instruction.
func HandleLdOne(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b) + 1
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.LogInstruction("ld %s, %s+1 -- [0x%04X] = 0x%02X", RegisterNames[a],
        WordRegisterNames[b >> 1], addr, data)
}

// Function HandleSt handles a ST instruction.
func HandleSt(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b)
    em.MemoryStore(addr, data)
    em.LogInstruction("st %s, %s -- [0x%04X] = 0x%02X", WordRegisterNames[b >> 1], RegisterNames[a],
        addr, data)
}

// Function HandleStInc handles a ST+ instruction.
func HandleStInc(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b)
    em.MemoryStore(addr, data)
    em.SetWordReg(b, addr + 1)
    em.LogInstruction("st %s+, %s -- [0x%04X] = 0x%02X", WordRegisterNames[b >> 1],
        RegisterNames[a], addr, data)
}

// Function HandleStDec handles a -ST instruction.
func HandleStDec(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b) - 1
    em.MemoryStore(addr, data)
    em.SetWordReg(b, addr)
    em.LogInstruction("st -%s, %s -- [0x%04X] = 0x%02X", WordRegisterNames[b >> 1],
        RegisterNames[a], addr, data)
}

// Function HandleStOne handles a ST+1 instruction.
func HandleStOne(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b) + 1
    em.MemoryStore(addr, data)
    em.LogInstruction("st %s+1, %s -- [0x%04X] = 0x%02X", WordRegisterNames[b >> 1],
        RegisterNames[a], addr, data)
}
