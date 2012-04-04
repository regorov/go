package k270emlib

import (
    "fmt"
)

var AB2Opcodes = [16]func(*Emulator, int, int){
    HandleAOpcode,
    nil,
    HandleLdv,
    HandleStv,
    HandlePand,
    HandlePor,
    HandlePxor,
    HandlePclr,
    HandleLd,
    HandleLdInc,
    HandleLdDec,
    HandleIn,
    HandleSt,
    HandleStInc,
    HandleStDec,
    HandleOut,
}

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

func HandleLdv(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b)
    data := em.VideoMemoryLoad(addr)
    em.SetReg(a, data)
    em.LogInstruction("ldv %s, %s -- VMEM[0x%04X] = 0x%02X", RegisterNames[a], WordRegisterNames[b >> 1], addr, data)
}

func HandleStv(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b)
    em.VideoMemoryStore(addr, data)
    em.LogInstruction("stv %s, %s -- VMEM[0x%04X] = 0x%02X", WordRegisterNames[b >> 1], RegisterNames[a], addr, data)
}

func HandlePand(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    before := em.LoadIOPort(port)
    after := before & uint8(i)
    em.StoreIOPort(port, after)
    
    em.LogInstruction("pand %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X & 0x%02X = 0x%02X", RegisterNames[b], RegisterNames[a], port, before, before, i, after)
}

func HandlePor(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    before := em.LoadIOPort(port)
    after := before | uint8(i)
    em.StoreIOPort(port, after)
    
    em.LogInstruction("por %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X | 0x%02X = 0x%02X", RegisterNames[b], RegisterNames[a], port, before, before, i, after)
}

func HandlePxor(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    before := em.LoadIOPort(port)
    after := before ^ uint8(i)
    em.StoreIOPort(port, after)
    
    em.LogInstruction("pxor %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X ^ 0x%02X = 0x%02X", RegisterNames[b], RegisterNames[a], port, before, before, i, after)
}

func HandlePclr(em *Emulator, a int, b int) {
    i := em.GetReg(a)
    port := em.GetReg(b)
    before := em.LoadIOPort(port)
    after := before & (^uint8(i))
    em.StoreIOPort(port, after)
    
    em.LogInstruction("pclr %s, %s -- ports[0x%02X] = 0x%02X, 0x%02X & ~0x%02X = 0x%02X", RegisterNames[b], RegisterNames[a], port, before, before, i, after)
}

func HandleLd(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b)
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.LogInstruction("ld %s, %s -- [0x%04X] = 0x%02X", RegisterNames[a], WordRegisterNames[b >> 1], addr, data)
}

func HandleLdInc(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b)
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.SetWordReg(b, addr + 1)
    em.LogInstruction("ld %s, %s+ -- [0x%04X] = 0x%02X", RegisterNames[a], WordRegisterNames[b >> 1], addr, data)
}

func HandleLdDec(em *Emulator, a int, b int) {
    addr := em.GetWordReg(b) - 1
    data := em.MemoryLoad(addr)
    em.SetReg(a, data)
    em.SetWordReg(b, addr)
    em.LogInstruction("ld %s, -%s -- [0x%04X] = 0x%02X", RegisterNames[a], WordRegisterNames[b >> 1], addr, data)
}

func HandleIn(em *Emulator, a int, b int) {
    addr := em.GetReg(b)
    data := em.LoadIOPort(addr)
    em.SetReg(a, data)
    em.LogInstruction("in %s, %s -- ports[0x%02X] = 0x%02X", RegisterNames[a], RegisterNames[b], addr, data)
}

func HandleSt(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b)
    em.MemoryStore(addr, data)
    em.LogInstruction("st %s, %s -- [0x%04X] = 0x%02X", WordRegisterNames[b >> 1], RegisterNames[a], addr, data)
}

func HandleStInc(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b)
    em.MemoryStore(addr, data)
    em.SetWordReg(b, addr + 1)
    em.LogInstruction("st %s+, %s -- [0x%04X] = 0x%02X", WordRegisterNames[b >> 1], RegisterNames[a], addr, data)
}

func HandleStDec(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetWordReg(b) - 1
    em.MemoryStore(addr, data)
    em.SetWordReg(b, addr)
    em.LogInstruction("st -%s, %s -- [0x%04X] = 0x%02X", WordRegisterNames[b >> 1], RegisterNames[a], addr, data)
}

func HandleOut(em *Emulator, a int, b int) {
    data := em.GetReg(a)
    addr := em.GetReg(b)
    em.StoreIOPort(addr, data)
    em.LogInstruction("out %s, %s -- ports[0x%02X] = 0x%02X", RegisterNames[b >> 1], RegisterNames[a], addr, data)
}
