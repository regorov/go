package k680emlib

import ()

type aluImmInstructionHandler func(*Emulator, uint8, uint8, uint16)

var aluImmInstructions = map[uint8]aluImmInstructionHandler{
    0x0: handlePushi,
    0x1: handleSti,
    0x2: handleAddi,
    0x3: handleSubi,
    0x4: handleAndi,
    0x5: handleOri,
    0x6: handleXori,
    0x7: handleMuli,
    0x8: handleShl,
    0x9: handleShr,
    0xA: handleAshr,
}

func handlePushi(em *Emulator, a uint8, d uint8, i uint16) {
    em.Push(uint32(i))
    em.LogInstruction("pushi 0x%04X", i)
}

func handleSti(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a]
    data := parseSigned14(i)
    em.MemoryStoreWord(addr, data)
    em.LogInstruction("sti %s, %s0x%04X -- [0x%08X] = 0x%08X", RegisterNames[a], sign14(i), abs14(i), addr, data)
}

func handleAddi(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := x + uint32(i)
    em.Regs[d] = y
    em.LogInstruction("addi %s, %s, 0x%04X -- 0x%08X + 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleSubi(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := x - uint32(i)
    em.Regs[d] = y
    em.LogInstruction("subi %s, %s, 0x%04X -- 0x%08X - 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleAndi(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := x & uint32(i)
    em.Regs[d] = y
    em.LogInstruction("andi %s, %s, 0x%04X -- 0x%08X & 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleOri(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := x | uint32(i)
    em.Regs[d] = y
    em.LogInstruction("ori %s, %s, 0x%04X -- 0x%08X | 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleXori(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := x ^ uint32(i)
    em.Regs[d] = y
    em.LogInstruction("xori %s, %s, 0x%04X -- 0x%08X ^ 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleMuli(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a] & 0x0000FFFF
    y := x * uint32(i)
    em.Regs[d] = y
    em.LogInstruction("muli %s, %s, 0x%04X -- 0x%08X * 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleShl(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := x << uint32(i)
    em.Regs[d] = y
    em.LogInstruction("shl %s, %s, 0x%04X -- 0x%08X << 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleShr(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := x >> uint32(i)
    em.Regs[d] = y
    em.LogInstruction("shr %s, %s, 0x%04X -- 0x%08X >> 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}

func handleAshr(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]
    y := uint32(int32(x) >> uint32(i))
    em.Regs[d] = y
    em.LogInstruction("ashr %s, %s, 0x%04X -- 0x%08X >> 0x%08X = 0x%08X", RegisterNames[d], RegisterNames[a], i, x, i, y)
}
