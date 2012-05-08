package k680emlib

import ()

type jumpMemInstructionHandler func(*Emulator, uint8, uint8, uint16)

var jumpMemInstructions = map[uint8]jumpMemInstructionHandler{
    0x0: handleJmp,
    0x1: handleCall,
    0x2: handleJbc,
    0x3: handleJbs,
    0x4: handleJeq,
    0x5: handleJne,
    0x6: handleJlt,
    0x7: handleJge,
    0x8: handleLdb,
    0x9: handleStb,
    0xA: handleLdh,
    0xB: handleSth,
    0xC: handleLdw,
    0xD: handleStw,
    0xE: handleLds,
    0xF: handleSts,
}

func handleJmp(em *Emulator, a uint8, d uint8, i uint16) {
    em.PC += parseSigned14(i)
    em.LogInstruction("jmp .%s0x%04X -- to 0x%08X", signPlus14(i), abs14(i), em.PC)
}

func handleCall(em *Emulator, a uint8, d uint8, i uint16) {
    em.Push(em.PC)
    em.PC += parseSigned14(i)
    em.LogInstruction("call .%s0x%04X -- to 0x%08X", signPlus14(i), abs14(i), em.PC)
}

func handleJbc(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]

    if (x>>d)&1 == 0 {
        em.PC += parseSigned14(i)
        em.LogInstruction("jbc %s, %d, .%s0x%04X -- bitcleared(0x%08X, %d) is true, taking branch to 0x%08X", RegisterNames[a], d, signPlus14(i), abs14(i), x, d, em.PC)
    } else {
        em.LogInstruction("jbc %s, %d, .%s0x%04X -- bitcleared(0x%08X, %d) is false, not taking branch", RegisterNames[a], d, signPlus14(i), abs14(i), x, d)
    }
}

func handleJbs(em *Emulator, a uint8, d uint8, i uint16) {
    x := em.Regs[a]

    if (x>>d)&1 == 1 {
        em.PC += parseSigned14(i)
        em.LogInstruction("jbs %s, %d, .%s0x%04X -- bitset(0x%08X, %d) is true, taking branch to 0x%08X", RegisterNames[a], d, signPlus14(i), abs14(i), x, d, em.PC)
    } else {
        em.LogInstruction("jbs %s, %d, .%s0x%04X -- bitset(0x%08X, %d) is false, not taking branch", RegisterNames[a], d, signPlus14(i), abs14(i), x, d)
    }
}

func handleJeq(em *Emulator, a uint8, d uint8, i uint16) {
    if em.Regs[a] == em.Regs[d] {
        em.PC += parseSigned14(i)
        em.LogInstruction("jeq %s, %s, .%s0x%04X -- 0x%08X == 0x%08X is true, taking branch to 0x%08X", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d], em.PC)
    } else {
        em.LogInstruction("jeq %s, %s, .%s0x%04X -- 0x%08X == 0x%08X is false, not taking branch", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d])
    }
}

func handleJne(em *Emulator, a uint8, d uint8, i uint16) {
    if em.Regs[a] != em.Regs[d] {
        em.PC += parseSigned14(i)
        em.LogInstruction("jne %s, %s, .%s0x%04X -- 0x%08X != 0x%08X is true, taking branch to 0x%08X", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d], em.PC)
    } else {
        em.LogInstruction("jne %s, %s, .%s0x%04X -- 0x%08X != 0x%08X is false, not taking branch", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d])
    }
}

func handleJlt(em *Emulator, a uint8, d uint8, i uint16) {
    if em.Regs[a] < em.Regs[d] {
        em.PC += parseSigned14(i)
        em.LogInstruction("jlt %s, %s, .%s0x%04X -- 0x%08X < 0x%08X is true, taking branch to 0x%08X", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d], em.PC)
    } else {
        em.LogInstruction("jlt %s, %s, .%s0x%04X -- 0x%08X < 0x%08X is false, not taking branch", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d])
    }
}

func handleJge(em *Emulator, a uint8, d uint8, i uint16) {
    if em.Regs[a] >= em.Regs[d] {
        em.PC += parseSigned14(i)
        em.LogInstruction("jge %s, %s, .%s0x%04X -- 0x%08X >= 0x%08X is true, taking branch to 0x%08X", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d], em.PC)
    } else {
        em.LogInstruction("jge %s, %s, .%s0x%04X -- 0x%08X >= 0x%08X is false, not taking branch", RegisterNames[a], RegisterNames[d], signPlus14(i), abs14(i), em.Regs[a], em.Regs[d])
    }
}

func handleLdb(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a] + parseSigned14(i)
    data := em.MemoryLoad(addr)
    em.Regs[d] = uint32(data)
    em.LogInstruction("ldb %s, %s, %s0x%04X -- [0x%08X] = 0x%02X", RegisterNames[d], RegisterNames[a], sign14(i), abs14(i), addr, data)
}

func handleStb(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a] + parseSigned14(i)
    data := em.Regs[d]
    em.MemoryStore(addr, uint8(data))
    em.LogInstruction("stb %s, %s, %s0x%04X -- [0x%08X] = 0x%02X", RegisterNames[a], RegisterNames[d], sign14(i), abs14(i), addr, data)
}

func handleLdh(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a] + parseSigned14(i)
    data := em.MemoryLoadHalf(addr)
    em.Regs[d] = uint32(data)
    em.LogInstruction("ldh %s, %s, %s0x%04X -- [0x%08X] = 0x%04X", RegisterNames[d], RegisterNames[a], sign14(i), abs14(i), addr, data)
}

func handleSth(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a] + parseSigned14(i)
    data := em.Regs[d]
    em.MemoryStoreHalf(addr, uint16(data))
    em.LogInstruction("sth %s, %s, %s0x%04X -- [0x%08X] = 0x%04X", RegisterNames[a], RegisterNames[d], sign14(i), abs14(i), addr, data)
}

func handleLdw(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a] + parseSigned14(i)
    data := em.MemoryLoadWord(addr)
    em.Regs[d] = data
    em.LogInstruction("ldw %s, %s, %s0x%04X -- [0x%08X] = 0x%08X", RegisterNames[d], RegisterNames[a], sign14(i), abs14(i), addr, data)
}

func handleStw(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a] + parseSigned14(i)
    data := em.Regs[d]
    em.MemoryStoreWord(addr, data)
    em.LogInstruction("stw %s, %s, %s0x%04X -- [0x%08X] = 0x%08X", RegisterNames[a], RegisterNames[d], sign14(i), abs14(i), addr, data)
}

func handleLds(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a]
    data := em.MemoryLoad(addr)
    em.Regs[d] = uint32(data)
    em.Regs[a] = addr + 1
    em.Regs[i] = em.Regs[i] - 1

    em.LogInstruction("lds %s, %s, %s -- [0x%08X] = 0x%02X, count now = 0x%08X", RegisterNames[d], RegisterNames[a], RegisterNames[i], sign14(i), abs14(i), addr, data, em.Regs[i])
}

func handleSts(em *Emulator, a uint8, d uint8, i uint16) {
    addr := em.Regs[a] + parseSigned14(i)
    data := em.Regs[d]
    em.MemoryStore(addr, uint8(data))
    em.Regs[a] = addr + 1
    em.Regs[i] = em.Regs[i] - 1

    em.LogInstruction("sts %s, %s, %s -- [0x%08X] = 0x%02X, count now = 0x%08X", RegisterNames[a], RegisterNames[d], RegisterNames[i], sign14(i), abs14(i), addr, data, em.Regs[i])
}
