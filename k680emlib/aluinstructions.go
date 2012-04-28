package k680emlib

import ()

type aluInstructionHandler func(*Emulator, uint8, uint8, uint8)

var aluInstructions = map[uint8]aluInstructionHandler{
    0x0: handlePush,
    0x1: handlePop,
    0x2: handleAdd,
    0x3: handleSub,
    0x4: handleAnd,
    0x5: handleOr,
    0x6: handleXor,
    0x7: handleMul,
    0x8: handleMov,
    //0xA: handleNot,
    //0xB: handleNeg,
    //0xC: handleXeq,
    //0xD: handleXne,
    //0xE: handleXlt,
    //0xF: handleXge,
}

func handlePush(em *Emulator, a uint8, b uint8, d uint8) {
    em.Push(em.Regs[a])
    em.LogInstruction("push %s -- Value transferred was 0x%08X", RegisterNames[a], em.Regs[a])
}

func handlePop(em *Emulator, a uint8, b uint8, d uint8) {
    em.Regs[a] = em.Pop()
    em.LogInstruction("pop %s -- Value transferred was 0x%08X", RegisterNames[a], em.Regs[a])
}

func handleAdd(em *Emulator, a uint8, b uint8, d uint8) {
    x := em.Regs[a]
    y := em.Regs[b]
    z := x + y
    em.Regs[d] = z
    em.LogInstruction("add %s, %s, %s -- 0x%08X + 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], RegisterNames[b], x, y, z)
}

func handleSub(em *Emulator, a uint8, b uint8, d uint8) {
    x := em.Regs[a]
    y := em.Regs[b]
    z := x - y
    em.Regs[d] = z
    em.LogInstruction("sub %s, %s, %s -- 0x%08X - 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], RegisterNames[b], x, y, z)
}

func handleAnd(em *Emulator, a uint8, b uint8, d uint8) {
    x := em.Regs[a]
    y := em.Regs[b]
    z := x & y
    em.Regs[d] = z
    em.LogInstruction("and %s, %s, %s -- 0x%08X & 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], RegisterNames[b], x, y, z)
}

func handleOr(em *Emulator, a uint8, b uint8, d uint8) {
    x := em.Regs[a]
    y := em.Regs[b]
    z := x | y
    em.Regs[d] = z
    em.LogInstruction("or %s, %s, %s -- 0x%08X | 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], RegisterNames[b], x, y, z)
}

func handleXor(em *Emulator, a uint8, b uint8, d uint8) {
    x := em.Regs[a]
    y := em.Regs[b]
    z := x ^ y
    em.Regs[d] = z
    em.LogInstruction("xor %s, %s, %s -- 0x%08X ^ 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], RegisterNames[b], x, y, z)
}

func handleMul(em *Emulator, a uint8, b uint8, d uint8) {
    x := em.Regs[a] & 0x0000FFFF
    y := em.Regs[b] & 0x0000FFFF
    z := x * y
    em.Regs[d] = z
    em.LogInstruction("mul %s, %s, %s -- 0x%08X * 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], RegisterNames[b], x, y, z)
}

func handleMov(em *Emulator, a uint8, b uint8, d uint8) {
    em.Regs[d] = em.Regs[a]
    em.LogInstruction("mov %s, %s -- value transferred was 0x%08X", RegisterNames[d], RegisterNames[a], em.Regs[d])
}
