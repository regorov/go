// Package k270emlib emulates a K270 processor. Typical usage:
// 
//     em := k270emlib.NewEmulator()    // Create an emulator
//     em.LoadProgram(myprogram)        // Load a program
//     em.SetTraceFile(os.Stdout)       // (optional) log instructions to stdout
//     em.Run()                         // Go!
//
package k270emlib

// BUG(kierdavis): Pin handlers are not fully implemented yet.

import (
    "fmt"
    "io"
    "sync"
)

// Type Character represents a character in the video memory.
type Character struct {
    Char uint8  // The character code
    Attr uint8  // The attribute byte
}

// Type Emulator represents a K270 processor.
type Emulator struct {
    memory []uint8              // The main memory.
    videoMemory []Character     // The video memory.
    interruptRegistry []uint16  // The interrupt registry.
    ioports []uint8             // The I/O ports.
    
    lastpc uint16   // The address of the last instruction executed.
    pc uint16       // The program counter.
    sp uint16       // The stack pointer.
    regs [16]uint8  // The 16 GP registers.
    c bool          // The C (carry) flag.
    a bool          // The A (authorised) flag.
    i bool          // The I (interrupts enabled) flag
    u bool          // The U (user mode) flag.
    sc uint8        // The internal stack counter, used for PUSHA and POPA instructions.
    timer uint32    // The system timer.
    
    traceFile io.Writer                                         // A file that a debug trace will be
                                                                // written to.
    running bool                                                // When this is set to false, Run()
                                                                // stops.
    interruptQueue chan uint8                                   // The interrupt queue.
    portHandlers map[uint8][](func(*Emulator, uint8, uint8))    // A map of port numbers to handler
                                                                // functions.
    pinHandlers map[uint][](func(*Emulator, uint, bool))        // A map of pin numbers to handler
                                                                // functions.
    getKey func() byte                                          // Returns a character from the
                                                                // keyboard if input is requested
                                                                // (from reading the KBDK I/O port)
    
    Mutex sync.Mutex    // The global mutex. This is locked during RunOne(). Run() temporarily
                        // unlocks it after every instruction, to allow other goroutines to modify
                        // the emulator's properties.
}

// Function NewEmulator creates, initialises and returns a new Emulator.
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

// Function Emulator.SetTraceFile sets the trace file to `traceFile`. Executed instructions and
// other information is logged to this writer. If `tracefile` is nil, logging is disabled (the
// default).
func (em *Emulator) SetTraceFile(traceFile io.Writer) {
    em.traceFile = traceFile
}

// Function Emulator.RegisterPortHandler sets up `handler` to be run whenever the I/O port `port` is
// modified.
func (em *Emulator) RegisterPortHandler(port uint8, handler func(*Emulator, uint8, uint8)) {
    handlers, ok := em.portHandlers[port]
    
    if !ok || handlers == nil {
        handlers = make([](func(*Emulator, uint8, uint8)), 0, 8)
    }
    
    handlers = append(handlers, handler)
    em.portHandlers[port] = handlers
}

// Function Emulator.Reset clears the values of all of the emulator's registers and flags.
func (em *Emulator) Reset() {
    em.lastpc = 0
    em.pc = 0
    em.sp = 0
    em.c = false
    em.a = false
    em.i = false
    em.u = false
    em.sc = 0
    em.timer = 0
    
    for i := 0; i < 16; i++ {
        em.regs[i] = 0
    }
}

// Function Emulator.ResetMemory resets all the memory's of the emulator, including the main RAM,
// the video RAM, the interrupt registry and the I/O ports.
func (em *Emulator) ResetMemory() {
    em.memory = make([]uint8, 1024)
    em.videoMemory = make([]Character, VMEM_SIZE)
    em.interruptRegistry = make([]uint16, 256)
    em.ioports = make([]uint8, 256)
}

// Function Emulator.GrowMemory expands the size of the main RAM to be at least the size specified
// by `newsize`.
func (em *Emulator) GrowMemory(newsize int) {
    if newsize == 0 {
        newsize = (cap(em.memory) + 1) * 2
    }
    
    m := make([]uint8, newsize)
    copy(m, em.memory)
    em.memory = m
}

// Function Emulator.GetMemory returns the emulator's RAM. This is the same object that is attached
// to the emulator, not a copy.
func (em *Emulator) GetMemory() (memory []uint8) {
    return em.memory
}

// Function Emulator.SetMemory sets the emulator's RAM to `memory`. The new value completely
// overrides the old value; it is not copied.
func (em *Emulator) SetMemory(memory []uint8) {
    em.memory = memory
}

// Function Emulator.MemoryLoad loads and returns the value at address `address` from the
// emulator's RAM.
func (em *Emulator) MemoryLoad(address uint16) (value uint8) {
    if int(address) >= len(em.memory) {
        return 0
    }
    
    return em.memory[address]
}

// Function Emulator.MemoryStore stores `value` to address `address` in the emulator's RAM, calling
// GrowMemory if needed.
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

// Function Emulator.VideoMemoryLoad loads and returns the value at address `address` in the
// emulator's video RAM.
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

// Function Emulator.VideoMemoryStore stores the value `value` to the emulator's video RAM at
// address `address`.
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

// Function Emulator.GetVideoMemory returns the emulator's video RAM (not a copy!).
func (em *Emulator) GetVideoMemory() (vmem []Character) {
    return em.videoMemory
}

// Function Emulator.InterruptRegistryLoad returns the address of the handler for the interrupt
// numbered `number`.
func (em *Emulator) InterruptRegistryLoad(number uint8) (value uint16) {
    return em.interruptRegistry[number]
}

// Function Emulator.InterruptRegistryStore sets the address of the handler for the interrupt
// numbered `number` to `value`.
func (em *Emulator) InterruptRegistryStore(number uint8, value uint16) {
    em.interruptRegistry[number] = value
}

// Function Emulator.LoadIOPort returns the value at the I/O port numbered `number`.
func (em *Emulator) LoadIOPort(number uint8) (value uint8) {
    // The keyboard is a special case
    if number == P_KBDK {
        return em.getKey()
    }
    
    return em.ioports[number]
}

// Function Emulator.StoreIOPort stores the value `value` to the I/O port numbered `number`.
func (em *Emulator) StoreIOPort(number uint8, value uint8) {
    em.ioports[number] = value
    em.triggerPortHandlers(number, value)
}

// Function Emulator.GetPC returns the emulator's program counter.
func (em *Emulator) GetPC() (value uint16) {
    return em.pc
}

// Function Emulator.SetPC sets the emulator's program counter to `value`.
func (em *Emulator) SetPC(value uint16) {
    em.pc = value
}

// Function Emulator.GetSP returns the emulator's stack pointer.
func (em *Emulator) GetSP() (value uint16) {
    return em.sp
}

// Function Emulator.SetSP sets the emulator's stack pointer to `value`.
func (em *Emulator) SetSP(value uint16) {
    em.sp = value
}

// Function Emulator.GetReg returns the value of the register numbered `number`. It raises a
// runtime panic if `number` is out of range. The panic value is an instance of `Error`, and the ID
// is E_REG_INDEX_OUT_OF_RANGE.
func (em *Emulator) GetReg(number int) (value uint8) {
    if number < 0 || number > 15 {
        panic(NewError(E_REG_INDEX_OUT_OF_RANGE, "Register index must be between 0 and 15"))
    }
    
    return em.regs[number]
}

// Function Emulator.GetWordReg returns the value of the 16-bit register pair numbered `number`. It
// raises a runtime panic if `number` is out of range. The panic value is an instance of `Error`,
// and the ID is E_REG_INDEX_OUT_OF_RANGE.
func (em *Emulator) GetWordReg(number int) (value uint16) {
    if number < 0 || number > 15 {
        panic(NewError(E_REG_INDEX_OUT_OF_RANGE, "Register index must be between 0 and 15"))
    }
    
    if number & 1 == 1 {
        number--
    }
    
    return (uint16(em.regs[number]) << 8) | uint16(em.regs[number + 1])
}

// Function Emulator.SetReg sets the value of the register numbered `number` to `value`. It raises a
// runtime panic if `number` is out of range. The panic value is an instance of `Error`, and the ID
// is E_REG_INDEX_OUT_OF_RANGE.
func (em *Emulator) SetReg(number int, value uint8) {
    if number < 1 || number > 15 { // excluding zero reg (r0)
        panic(NewError(E_REG_INDEX_OUT_OF_RANGE, "Register index must be between 1 and 15"))
    }
    
    em.regs[number] = value
}

// Function Emulator.SetWordReg sets value of the 16-bit register pair numbered `number` to `value`.
// It raises a runtime panic if `number` is out of range. The panic value is an instance of `Error`,
// and the ID is E_REG_INDEX_OUT_OF_RANGE.
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

// Function Emulator.GetCarry returns the C (carry) flag.
func (em *Emulator) GetCarry() (value bool) {
    return em.c
}

// Function Emulator.SetCarry sets the C (carry) flag to `value`.
func (em *Emulator) SetCarry(value bool) {
    em.c = value
}

// Function Emulator.GetAuthorised returns the A (authorised) flag.
func (em *Emulator) GetAuthorised() (value bool) {
    return em.a
}

// Function Emulator.SetAuthorised sets the A (authorised) flag to `value`.
func (em *Emulator) SetAuthorised(value bool) {
    em.a = value
}

// Function Emulator.GetInterruptsEnabled returns the I (interrupts enabled) flag.
func (em *Emulator) GetInterruptsEnabled() (value bool) {
    return em.i
}

// Function Emulator.SetInterruptsEnabled sets the I (interrupts enabled) flag to `value`. Also, if
// there are interrupts waiting in the queue it will pop one off and get ready to execute it.
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

// Function Emulator.GetUserMode returns the U (user mode) flag.
func (em *Emulator) GetUserMode() (value bool) {
    return em.u
}

// Function Emulator.SetUserMode sets the U (user mode) value to `value`.
func (em *Emulator) SetUserMode(value bool) {
    em.u = value
}

// Function Emulator.Push decrements the stack pointer, then stores `value` into RAM at the address
// specified by the stack pointer.
func (em *Emulator) Push(value uint8) {
    em.sp--
    em.MemoryStore(em.sp, value)
}

// Function Emulator.Pop loads and returns the value in RAM at the address specified by the stack
// pointer before incrementing the stack pointer.
func (em *Emulator) Pop() (value uint8) {
    value = em.MemoryLoad(em.sp)
    em.sp++
    return value
}

// Function Emulator.PushWord pushes the 16-bit word `value` onto the stack.
func (em *Emulator) PushWord(value uint16) {
    em.Push(uint8(value))
    em.Push(uint8(value >> 8))
}

// Function Emulator.PopWord pops a 16-bit word off the stack and returns it.
func (em *Emulator) PopWord() (value uint16) {
    high := uint16(em.Pop())
    low := uint16(em.Pop())
    return (high << 8) | low
}

// Function Emulator.FetchWord loads and returns the next program word.
func (em *Emulator) FetchWord() (word uint16) {
    high := em.MemoryLoad(em.pc)
    low := em.MemoryLoad(em.pc + 1)
    
    em.lastpc = em.pc
    em.pc += 2
    
    return (uint16(high) << 8) | uint16(low)
}

// Function Emulator.RunOne loads the next program word and executes it (by distributing it to
// HandleAIOpcode). It locks the mutex during the execution of this function.
func (em *Emulator) RunOne() {
    em.Mutex.Lock()
    em.timer += 4
    em.CheckTimer()
    
    word := em.FetchWord()
    
    o := int(word >> 12)
    a := int((word >> 8) & 0xF)
    i := int(word & 0xFF)
    
    HandleAIOpcode(em, o, a, i)
    em.Mutex.Unlock()
}

// Function Emulator.Run sets the running flag to true, then runs instructions until it is set to
// false. The repeated unlocking and locking of the mutex is to allow external goroutines to access
// the emulator too.
func (em *Emulator) Run() {
    em.running = true
    
    for em.running {
        em.RunOne()
    }
}

// Function Emulator.LogInstruction works like a fmt.Fprintf to the traceFile attribute, except that
// it only executes if the traceFile attribute is not nil.
func (em *Emulator) LogInstruction(format string, args ...interface{}) {
    if em.traceFile != nil {
        format = fmt.Sprintf(format, args...)
        fmt.Fprintf(em.traceFile, "[0x%04X] %s\n", em.lastpc, format)
    }
}

// Function Emulator.GetRunning returns the state of the running flag.
func (em *Emulator) GetRunning() (running bool) {
    return em.running
}

// Function Emulator.SetRunning sets the running flag to `running`.
func (em *Emulator) SetRunning(running bool) {
    em.running = running
}

// Function Emulator.Interrupt triggers the interrupt numbered `i`. One of three things can occur:
// 
// * Interrupts are disabled, in which case the event is pushed onto the interrupt queue.
// 
// * The interrupt is registered in the registry, in which case interrupts are disabled, PC is
// pushed onto the stack and the interrupt's address is jumped. Additionally, if the interrupt
// number is a system interrupt (i.e. 0x00 - 0x7F), the Emulator switches to SYS (system) mode.
// 
// * The interrupt is not registered, in which case nothing happens.
func (em *Emulator) Interrupt(i uint8) {
    if em.traceFile != nil {fmt.Fprintf(em.traceFile, "Interrupt 0x%02X requested, ", i)}
    
    addr := em.InterruptRegistryLoad(i)
    if addr != 0 {
        if em.GetInterruptsEnabled() {
            if em.traceFile != nil {fmt.Fprintf(em.traceFile, "executing now (calling 0x%04X)\n",
                addr)}
            
            em.SetInterruptsEnabled(false)
            em.PushWord(em.pc)
            em.pc = addr
            
            if i < 0x80 {
                em.SetUserMode(false)
            }
        
        } else {
            if em.traceFile != nil {fmt.Fprintf(em.traceFile, "queueing\n")}
            em.interruptQueue <- i
        }
    
    } else {
        if em.traceFile != nil {fmt.Fprintf(em.traceFile,
            "discarding (interrupt not registered)\n")}
        
        em.SetInterruptsEnabled(true) // This triggers the emulator to execute another interrupt
                                      // off the queue.
    }
}

// Function Emulator.CheckTimer will check the system timer and trigger the T0, T1 and/or T2
// interrupts if appropriate.
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

// Function Emulator.GetDigitalOutput will return the state of the digital output pin numbered
// `pin`.
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

// Function Emulator.SetDigitalInput will set the state of the digital input pin numbered `pin` to
// `value`. Additionally, if `triggerHandlers` is true, it will trigger any pin handlers
// associated with the pin.
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

// Function Emulator.SetGetKey sets up the emulator so that `getKey` is called if input is
// requested (i.e. port KBDK is read). `getKey` should return a single byte, which will be used as
// the read character.
func (em *Emulator) SetGetKey(getKey func() byte) {
    em.getKey = getKey
}

// Function Emulator.triggerPortHandlers triggers any handlers associated with I/O port `port`. The
// new value of the port, `value`, is passed to the handlers.
func (em *Emulator) triggerPortHandlers(port uint8, value uint8) {
    handlers, ok := em.portHandlers[port]
    
    if ok && handlers != nil {
        for _, handler := range handlers {
            handler(em, port, value)
        }
    }
}

// Function Emulator.triggerPinHandlers triggers any handlers associated with digital pin `pin`. The
// new value of the pin, `value`, is passed to the handlers.
func (em *Emulator) triggerPinHandlers(pin uint, value bool) {
    handlers, ok := em.pinHandlers[pin]
    
    if ok && handlers != nil {
        for _, handler := range handlers {
            handler(em, pin, value)
        }
    }
}

// Function Emulator.getPortAccess returns whether the port numbered `port` is accessible in the
// current user mode, setting the A (authorised) flag to this value.
func (em *Emulator) getPortAccess(port uint8) (authorised bool) {
    if port < 0x80 && em.GetUserMode() {
        em.SetAuthorised(false)
        return false
    
    } else {
        em.SetAuthorised(true)
        return true
    }
    
    return false
}

// Function timerResetHandler is an internal port handler that is attached to the TR port. It
// handles resetting the system timer.
func timerResetHandler(em *Emulator, port uint8, value uint8) {
    em.timer = 0
}
