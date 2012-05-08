package k680emlib

import ()

type otherInstructionHandler func(*Emulator, uint8, uint16)

var otherInstructions = map[uint8]otherInstructionHandler{
    0x00: handleNop,
    0x01: handleLdi,
    0x02: handleJr,
    0x03: handleCr,
    //0x04: handleInt,
    //0x05: handleRih,
    0x06: handleLda,
    0x07: handleSta,
    0x08: handleRet,
    //0x09: handleReti,
    0x0A: handleLdl,
    0x0B: handleLdu,
    //0x0C: handlePr,
    //0x0D: handlePw,
    //0x0E: handlePc,
    //0x0F: handlePs,
    0x10: handleCps,
    0x11: handleJmcs,

    0x7F: handleHlt,
}

func handleNop(em *Emulator, a uint8, i uint16) {
    em.LogInstruction("nop")
}

func handleLdi(em *Emulator, a uint8, i uint16) {
    em.Regs[a] = parseSigned16(i)
    em.LogInstruction("ldi %s, %s0x%04X", RegisterNames[a], sign16(i), abs16(i))
}

func handleJr(em *Emulator, a uint8, i uint16) {
    em.PC = em.Regs[a]
    em.LogInstruction("jr %s -- 0x%08X", RegisterNames[a], em.PC)
}

func handleCr(em *Emulator, a uint8, i uint16) {
    em.Push(em.PC)
    em.PC = em.Regs[a]
    em.LogInstruction("cr %s -- 0x%08X", RegisterNames[a], em.PC)
}

func handleLda(em *Emulator, a uint8, i uint16) {
    scale := (i >> 10) & 0x03
    base := (i >> 5) & 0x1F
    index := i & 0x1F
    addr := em.Regs[base] + (em.Regs[index] << scale)

    switch scale {
    case 0:
        data := em.MemoryLoad(addr)
        em.Regs[a] = uint32(data)
        em.LogInstruction("lda %s, %s, 0, %s -- [0x%08X] = 0x%02X", RegisterNames[a], RegisterNames[base], RegisterNames[index], addr, data)

    case 1:
        data := em.MemoryLoadHalf(addr)
        em.Regs[a] = uint32(data)
        em.LogInstruction("lda %s, %s, 1, %s -- [0x%08X] = 0x%04X", RegisterNames[a], RegisterNames[base], RegisterNames[index], addr, data)

    case 2:
        data := em.MemoryLoadWord(addr)
        em.Regs[a] = data
        em.LogInstruction("lda %s, %s, 2, %s -- [0x%08X] = 0x%08X", RegisterNames[a], RegisterNames[base], RegisterNames[index], addr, data)

    case 3:
        data := em.MemoryLoadDouble(addr)
        em.Regs[a>>1] = uint32(data >> 32)
        em.Regs[(a>>1)+1] = uint32(data)
        em.LogInstruction("lda %s, %s, 3, %s -- [0x%08X] = 0x%016X", RegisterNames[a], RegisterNames[base], RegisterNames[index], addr, data)
    }
}

func handleSta(em *Emulator, a uint8, i uint16) {
    scale := (i >> 10) & 0x03
    base := (i >> 5) & 0x1F
    index := i & 0x1F
    addr := em.Regs[base] + (em.Regs[index] << scale)
    data := em.Regs[a]

    switch scale {
    case 0:
        em.MemoryStore(addr, uint8(data))
        em.LogInstruction("sta %s, 0, %s, %s -- [0x%08X] = 0x%02X", RegisterNames[base], RegisterNames[index], RegisterNames[a], addr, data)

    case 1:
        em.MemoryStoreHalf(addr, uint16(data))
        em.LogInstruction("sta %s, 1, %s, %s -- [0x%08X] = 0x%04X", RegisterNames[base], RegisterNames[index], RegisterNames[a], addr, data)

    case 2:
        em.MemoryStoreWord(addr, data)
        em.LogInstruction("sta %s, 2, %s, %s -- [0x%08X] = 0x%08X", RegisterNames[base], RegisterNames[index], RegisterNames[a], addr, data)

    case 3:
        data64 := uint64(em.Regs[a>>1]) << 32
        data64 |= uint64(em.Regs[(a>>1)+1])
        em.MemoryStoreDouble(addr, data64)
        em.LogInstruction("sta %s, 3, %s, %s -- [0x%08X] = 0x%016X", RegisterNames[base], RegisterNames[index], RegisterNames[a], addr, data)
    }
}

func handleRet(em *Emulator, a uint8, i uint16) {
    em.PC = em.Pop()
    em.LogInstruction("ret")
}

func handleLdl(em *Emulator, a uint8, i uint16) {
    em.Regs[a] = (em.Regs[a] & 0xFFFF0000) | uint32(i)
    em.LogInstruction("ldl %s, 0x%04X (now = 0x%08X)", RegisterNames[a], i, em.Regs[a])
}

func handleLdu(em *Emulator, a uint8, i uint16) {
    em.Regs[a] = (em.Regs[a] & 0x0000FFFF) | (uint32(i) << 16)
    em.LogInstruction("ldu %s, 0x%04X (now = 0x%08X)", RegisterNames[a], i, em.Regs[a])
}

func handleCps(em *Emulator, a uint8, i uint16) {
    s := (i >> 5) & 0x1F
    d := i & 0x1F
    sa := em.Regs[s]
    da := em.Regs[d]

    v := em.MemoryLoad(sa)
    em.MemoryStore(da, v)

    em.Regs[a]--
    em.Regs[s]++
    em.Regs[d]++
    em.LogInstruction("cps %s, %s, %s -- [0x%08X] -> [0x%08X], value was 0x%02X, count is now 0x%08X", RegisterNames[d], RegisterNames[s], RegisterNames[a], sa, da, v, em.Regs[a])
}

func handleJmcs(em *Emulator, a uint8, i uint16) {
    em.PC = em.Regs[a]
    em.Regs[4] = em.Regs[7]
    em.LogInstruction("jmcs %s -- PC = 0x%08X, CS/US = 0x%08X", RegisterNames[a], em.PC, em.Regs[4])
}

func handleHlt(em *Emulator, a uint8, i uint16) {
    em.Running = false
    em.LogInstruction("hlt")
}
