package k270emlib

import (
    "fmt"
    "io"
)

type Character struct {
    Char uint8
    Attr uint8
}

type Emulator struct {
    memory []uint8
    videoMemory []Character
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
    timer uint32
    
    traceFile io.Writer
    running bool
    interruptQueue chan uint8
    portHandlers map[uint8][](func(*Emulator, uint8, uint8))
    pinHandlers map[uint][](func(*Emulator, uint, bool))
}

func NewEmulator() (em *Emulator) {
    em = new(Emulator)
    em.traceFile = nil
    em.running = false
    em.interruptQueue = make(chan uint8, 16)
    em.portHandlers = make(map[uint8][](func(*Emulator, uint8, uint8)))
    em.pinHandlers = make(map[uint][](func(*Emulator, uint, bool)))
    
    em.Reset()
    em.ResetMemory()
    
    em.RegisterPortHandler(P_TR, timerResetHandler)
    
    return em
}

func (em *Emulator) SetTraceFile(traceFile io.Writer) {
    em.traceFile = traceFile
}

func (em *Emulator) RegisterPortHandler(port uint8, handler func(*Emulator, uint8, uint8)) {
    handlers, ok := em.portHandlers[port]
    
    if !ok || handlers == nil {
        handlers = make([](func(*Emulator, uint8, uint8)), 0, 8)
    }
    
    handlers = append(handlers, handler)
    em.portHandlers[port] = handlers
}

func (em *Emulator) Reset() {
    em.lastpc = 0
    em.pc = 0
    em.sp = 0
    em.c = false
    em.a = false
    em.i = true
    em.u = false
    em.sc = 0
    em.timer = 0
    
    for i := 0; i < 16; i++ {
        em.regs[i] = 0
    }
}

func (em *Emulator) ResetMemory() {
    em.memory = make([]uint8, 1024)
    em.videoMemory = make([]Character, VMEM_SIZE)
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

func (em *Emulator) VideoMemoryLoad(address uint16) (value uint8) {
    if address & 0x8000 != 0 {
        address = address & 0x7fff
        
        if address < VMEM_SIZE {
            return em.videoMemory[address].Attr
        }
    
    } else {
        if address < VMEM_SIZE {
            return em.videoMemory[address].Char
        }
    }
    
    return 0
}

func (em *Emulator) VideoMemoryStore(address uint16, value uint8) {
    if address & 0x8000 != 0 {
        address = address & 0x7fff
        
        if address < VMEM_SIZE {
            em.videoMemory[address].Attr = value
        }
    
    } else {
        if address < VMEM_SIZE {
            em.videoMemory[address].Char = value
        }
    }
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
    em.triggerPortHandlers(number, value)
}

func (em *Emulator) GetPC() (value uint16) {
    return em.pc
}

func (em *Emulator) SetPC(value uint16) {
    em.pc = value
}

func (em *Emulator) GetReg(number int) (value uint8) {
    if number < 0 || number > 15 {
        panic(NewError(E_REG_INDEX_OUT_OF_RANGE, "Register index must be between 0 and 15"))
    }
    
    return em.regs[number]
}

func (em *Emulator) GetWordReg(number int) (value uint16) {
    if number < 0 || number > 15 {
        panic(NewError(E_REG_INDEX_OUT_OF_RANGE, "Register index must be between 0 and 15"))
    }
    
    if number & 1 == 1 {
        number--
    }
    
    return (uint16(em.regs[number]) << 8) | uint16(em.regs[number + 1])
}

func (em *Emulator) SetReg(number int, value uint8) {
    if number < 1 || number > 15 { // excluding zero reg (r0)
        panic(NewError(E_REG_INDEX_OUT_OF_RANGE, "Register index must be between 1 and 15"))
    }
    
    em.regs[number] = value
}

func (em *Emulator) SetWordReg(number int, value uint16) {
    if number < 2 || number > 15 { // excluding zero reg pair (r0:r1)
        panic(NewError(E_REG_INDEX_OUT_OF_RANGE, "Register index must be between 2 and 15"))
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
    
    if value {
        // Interrupts have been re-enabled, try and handle another
        
        select {
        case i := <-em.interruptQueue: // Interrupt waiting
            em.Interrupt(i)
        default: // No interrupt waiting
        }
    }
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
    em.timer += 4
    em.CheckTimer()
    
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

func (em *Emulator) Interrupt(i uint8) {
    fmt.Fprintf(em.traceFile, "Interrupt 0x%02X requested, ", i)
    
    if em.GetInterruptsEnabled() {
        addr := em.InterruptRegistryLoad(i)
        if addr != 0 {
            fmt.Fprintf(em.traceFile, "executing now (calling 0x%04X)\n", addr)
            em.SetInterruptsEnabled(false)
            em.PushWord(em.pc)
            em.pc = addr
        
        } else {
            fmt.Fprintf(em.traceFile, "discarding (interrupt not registered)\n")
        }
    
    } else {
        fmt.Fprintf(em.traceFile, "queueing\n")
        em.interruptQueue <- i
    }
}

func (em *Emulator) CheckTimer() {
    if em.timer % 256 == 0 {
        em.Interrupt(INT_T0)
    }
    
    if em.timer % 1024 == 0 {
        em.Interrupt(INT_T1)
    }
    
    if em.timer % 4096 == 0 {
        em.Interrupt(INT_T2)
    }
}

func (em *Emulator) GetDigitalOutput(pin uint) (value bool) {
    port := (pin >> 3) & 3
    bit := pin & 7
    
    dout := em.ioports[P_DOUT0 + port]
    dmode := em.ioports[P_DMODE0 + port]
    
    if ((dmode >> bit) & 1) == 1 {
        return ((dout >> bit) & 1) == 1
    
    } else {
        panic(NewError(E_INCORRECT_MODE, "Pin is not currently defined as an output"))
    }
    return false
}

func (em *Emulator) SetDigitalInput(pin uint, value bool, triggerHandlers bool) {
    port := (pin >> 3) & 3
    bit := pin & 7
    v := uint8(0)
    if value {v = 1}
    
    dmode := em.ioports[P_DMODE0 + port]
    
    if ((dmode >> bit) & 1) == 0 {
        din := em.ioports[P_DIN0 + port]
        din = (din & ^(1 << bit)) | (v << bit)
        em.ioports[P_DIN0 + port] = din
        
        if triggerHandlers {
            em.triggerPinHandlers(pin, value)
        }
    
    } else {
        panic(NewError(E_INCORRECT_MODE, "Pin is not currently defined as an input"))
    }
}

func (em *Emulator) triggerPortHandlers(port uint8, value uint8) {
    handlers, ok := em.portHandlers[port]
    
    if ok && handlers != nil {
        for _, handler := range handlers {
            handler(em, port, value)
        }
    }
}

func (em *Emulator) triggerPinHandlers(pin uint, value bool) {
    handlers, ok := em.pinHandlers[pin]
    
    if ok && handlers != nil {
        for _, handler := range handlers {
            handler(em, pin, value)
        }
    }
}

func timerResetHandler(em *Emulator, port uint8, value uint8) {
    em.timer = 0
}