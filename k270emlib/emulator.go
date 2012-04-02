package k270emlib

import (
    "fmt"
    "io"
)

type Emulator struct {
    memory []uint8
    interruptRegistry []uint16
    ioports []uint8
    
    lastpc uint16
    pc uint16
    sp uint16
    regs [16]uint8
    c bool
    a bool
    i bool
    u bool
    sc uint8
    
    traceFile io.Writer
    running bool
}

func NewEmulator() (em *Emulator) {
    em = new(Emulator)
    em.traceFile = nil
    em.Reset()
    em.ResetMemory()
    
    return em
}

func (em *Emulator) SetTraceFile(traceFile io.Writer) {
    em.traceFile = traceFile
}

func (em *Emulator) Reset() {
    em.lastpc = 0
    em.pc = 0
    em.sp = 0
    em.c = false
    em.a = false
    em.i = false
    em.u = false
    em.sc = 0
    
    for i := 0; i < 16; i++ {
        em.regs[i] = 0
    }
}

func (em *Emulator) ResetMemory() {
    em.memory = make([]uint8, 1024)
    em.interruptRegistry = make([]uint16, 256)
    em.ioports = make([]uint8, 256)
}

func (em *Emulator) GrowMemory(newsize int) {
    if newsize == 0 {
        newsize = (cap(em.memory) + 1) * 2
    }
    
    m := make([]uint8, newsize)
    copy(m, em.memory)
    em.memory = m
}

func (em *Emulator) GetMemory() (memory []uint8) {
    return em.memory
}

func (em *Emulator) SetMemory(memory []uint8) {
    em.memory = memory
}

func (em *Emulator) MemoryLoad(address uint16) (value uint8) {
    if int(address) >= len(em.memory) {
        //panic("Memory address out of range")
        return 0
    }
    
    return em.memory[address]
}

func (em *Emulator) MemoryStore(address uint16, value uint8) {
    if int(address) >= len(em.memory) {
        newsize := cap(em.memory) + 1
        
        for int(address) >= newsize {
            newsize *= 2
        }
        
        em.GrowMemory(newsize)
    }
    
    em.memory[address] = value
}

func (em *Emulator) InterruptRegistryLoad(number uint8) (value uint16) {
    return em.interruptRegistry[number]
}

func (em *Emulator) InterruptRegistryStore(number uint8, value uint16) {
    em.interruptRegistry[number] = value
}

func (em *Emulator) LoadIOPort(number uint8) (value uint8) {
    return em.ioports[number]
}

func (em *Emulator) StoreIOPort(number uint8, value uint8) {
    em.ioports[number] = value
}

func (em *Emulator) GetPC() (value uint16) {
    return em.pc
}

func (em *Emulator) SetPC(value uint16) {
    em.pc = value
}

func (em *Emulator) GetReg(number int) (value uint8) {
    if number < 0 || number > 15 {
        panic("Register index must be between 0 and 15")
    }
    
    return em.regs[number]
}

func (em *Emulator) GetWordReg(number int) (value uint16) {
    if number < 0 || number > 15 {
        panic("Register index must be between 0 and 15")
    }
    
    if number & 1 == 1 {
        number--
    }
    
    return (uint16(em.regs[number]) << 8) | uint16(em.regs[number + 1])
}

func (em *Emulator) SetReg(number int, value uint8) {
    if number < 1 || number > 15 {
        panic("Register index must be between 1 and 15") // excluding zero reg (r0)
    }
    
    em.regs[number] = value
}

func (em *Emulator) SetWordReg(number int, value uint16) {
    if number < 2 || number > 15 {
        panic("Register index must be between 2 and 15") // excluding zero reg pair (r0:r1)
    }
    
    if number & 1 == 1 {
        number--
    }
    
    em.regs[number] = uint8(value >> 8)
    em.regs[number + 1] = uint8(value)
}

func (em *Emulator) GetCarry() (value bool) {
    return em.c
}

func (em *Emulator) SetCarry(value bool) {
    em.c = value
}

func (em *Emulator) GetAuthorised() (value bool) {
    return em.a
}

func (em *Emulator) SetAuthorised(value bool) {
    em.a = value
}

func (em *Emulator) GetInterruptsEnabled() (value bool) {
    return em.i
}

func (em *Emulator) SetInterruptsEnabled(value bool) {
    em.i = value
}

func (em *Emulator) GetUserMode() (value bool) {
    return em.u
}

func (em *Emulator) SetUserMode(value bool) {
    em.u = value
}

func (em *Emulator) Push(value uint8) {
    em.sp--
    em.MemoryStore(em.sp, value)
}

func (em *Emulator) Pop() (value uint8) {
    value = em.MemoryLoad(em.sp)
    em.sp++
    return value
}

func (em *Emulator) PushWord(value uint16) {
    em.Push(uint8(value))
    em.Push(uint8(value >> 8))
}

func (em *Emulator) PopWord() (value uint16) {
    high := uint16(em.Pop())
    low := uint16(em.Pop())
    return (high << 8) | low
}

func (em *Emulator) FetchWord() (word uint16) {
    high := em.MemoryLoad(em.pc)
    low := em.MemoryLoad(em.pc + 1)
    
    em.lastpc = em.pc
    em.pc += 2
    
    return (uint16(high) << 8) | uint16(low)
}

func (em *Emulator) RunOne() {
    word := em.FetchWord()
    
    o := int(word >> 12)
    a := int((word >> 8) & 0xF)
    i := int(word & 0xFF)
    
    HandleAIOpcode(em, o, a, i)
}

func (em *Emulator) Run() {
    em.running = true
    
    for em.running {
        em.RunOne()
    }
}

func (em *Emulator) LogInstruction(format string, args ...interface{}) {
    if em.traceFile != nil {
        format = fmt.Sprintf(format, args...)
        fmt.Fprintf(em.traceFile, "[0x%04X] %s\n", em.lastpc, format)
    }
}

func (em *Emulator) GetRunning() (running bool) {
    return em.running
}

func (em *Emulator) SetRunning(running bool) {
    em.running = running
}