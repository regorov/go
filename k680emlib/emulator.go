// Package k680emlib provides a backend for including a K680 emulator in your project.
package k680emlib

import (
    "fmt"
    "io"
)

// Type Emulator represents a K680 emulator.
type Emulator struct {
    Regs      [32]uint32
    LastPC    uint32
    PC        uint32
    Memory    []byte
    TraceFile io.Writer
    Running   bool
}

// Function NewEmulator creates and returns a new emulator.
func NewEmulator() (em *Emulator) {
    em = new(Emulator)

    em.TraceFile = nil
    em.Reset()

    return em
}

// Function Reset resets the state of the emulator's registers and memory.
func (em *Emulator) Reset() {
    for i := 0; i < 32; i++ {
        em.Regs[i] = 0
    }

    em.LastPC = 0
    em.PC = 0
    em.Memory = make([]byte, 1024)
}

// Function GrowMemory expands the size of the main RAM to be at least the size specified.
func (em *Emulator) GrowMemory(newsize int) {
    if newsize == 0 {
        newsize = (cap(em.Memory) + 1) * 2
    }

    m := make([]byte, newsize)
    copy(m, em.Memory)
    em.Memory = m
}

// Function MemoryLoad returns the byte at the specified location in the emulator's RAM.
func (em *Emulator) MemoryLoad(address uint32) (value byte) {
    if address >= uint32(len(em.Memory)) {
        return 0
    }

    return em.Memory[address]
}

// Function MemoryLoadHalf returns the halfword at the specified location in the emulator's RAM.
func (em *Emulator) MemoryLoadHalf(address uint32) (value uint16) {
    a := uint16(em.MemoryLoad(address))
    b := uint16(em.MemoryLoad(address + 1))
    return (a << 8) | b
}

// Function MemoryLoadWord returns the word at the specified location in the emulator's RAM.
func (em *Emulator) MemoryLoadWord(address uint32) (value uint32) {
    a := uint32(em.MemoryLoadHalf(address))
    b := uint32(em.MemoryLoadHalf(address + 2))
    return (a << 16) | b
}

// Function MemoryLoadDouble returns the doubleword at the specified location in the emulator's RAM.
func (em *Emulator) MemoryLoadDouble(address uint32) (value uint64) {
    a := uint64(em.MemoryLoadWord(address))
    b := uint64(em.MemoryLoadWord(address + 4))
    return (a << 32) | b
}

// Function MemoryStore stores the value to the specified location in the emulator's RAM.
func (em *Emulator) MemoryStore(address uint32, value byte) {
    if address >= uint32(len(em.Memory)) {
        newsize := cap(em.Memory) + 1

        for int(address) >= newsize {
            newsize *= 2
        }

        em.GrowMemory(newsize)
    }

    em.Memory[address] = value
}

// Function MemoryStoreHalf stores the halfword to the specified location in the emulator's RAM.
func (em *Emulator) MemoryStoreHalf(address uint32, value uint16) {
    em.MemoryStore(address, uint8(value>>8))
    em.MemoryStore(address+1, uint8(value))
}

// Function MemoryStoreWord stores the word to the specified location in the emulator's RAM.
func (em *Emulator) MemoryStoreWord(address uint32, value uint32) {
    em.MemoryStoreHalf(address, uint16(value>>16))
    em.MemoryStoreHalf(address+2, uint16(value))
}

// Function MemoryStoreDouble stores the doubleword to the specified location in the emulator's RAM.
func (em *Emulator) MemoryStoreDouble(address uint32, value uint64) {
    em.MemoryStoreWord(address, uint32(value>>32))
    em.MemoryStoreWord(address+4, uint32(value))
}

// Function Push pushes the value onto the stack.
func (em *Emulator) Push(value uint32) {
    em.Regs[SP] -= 4
    em.MemoryStoreWord(em.Regs[SP], value)
}

// Function Pop pops a value off the stack.
func (em *Emulator) Pop() (value uint32) {
    value = em.MemoryLoadWord(em.Regs[SP])
    em.Regs[SP] += 4
    return value
}

// Function LoadProgram loads the specified program into RAM.
func (em *Emulator) LoadProgram(program []byte, offset uint32) {
    end := len(program) + int(offset)
    if len(em.Memory) < end {
        em.GrowMemory(end)
    }

    copy(em.Memory[offset:], program)
}

// Function FetchWord fetches a program word and increments the program counter.
func (em *Emulator) FetchWord() (word uint32) {
    em.LastPC = em.PC
    word = em.MemoryLoadWord(em.PC)
    em.PC += 4
    return word
}

// Function DecodeInstruction extracts the mode, exection condition, opcode and a-index from an
// instruction word.
func (em *Emulator) DecodeInstruction(word uint32) (mode uint8, xc uint8, opcode uint8, a uint8) {
    mode = uint8(word>>30) & 0x03
    xc = uint8(word>>28) & 0x03
    opcode = uint8(word>>24) & 0x0F
    a = uint8(word>>19) & 0x1F
    return mode, xc, opcode, a
}

// Function DecodeOther extracts the extended opcode and immediate from an other-type instruction
// word.
func (em *Emulator) DecodeOther(word uint32) (opext uint8, i uint16) {
    opext = uint8(word>>16) & 0x07
    i = uint16(word & 0xFFFF)
    return opext, i
}

// Function DecodeALU extracts the b-index and d-index from an ALU-type instruction word.
func (em *Emulator) DecodeALU(word uint32) (b uint8, d uint8) {
    b = uint8(word>>14) & 0x1F
    d = uint8(word>>9) & 0x1F
    return b, d
}

// Function DecodeJMI extracts the d-index and immediate from a jump, memory or ALU/immediate
// instruction word.
func (em *Emulator) DecodeJMI(word uint32) (d uint8, i uint16) {
    d = uint8(word>>14) & 0x1F
    i = uint16(word & 0x3FFF)
    return d, i
}

// Function LogInstruction formats a message and logs it to the trace file.
func (em *Emulator) LogInstruction(format string, args ...interface{}) {
    if em.TraceFile != nil {
        format = fmt.Sprintf(format, args...)
        fmt.Fprintf(em.TraceFile, "[0x%08X] %s\n", em.LastPC, format)
    }
}

// Function RunOne runs one instruction.
func (em *Emulator) RunOne() (err error) {
    word := em.FetchWord()
    mode, _, opcode, a := em.DecodeInstruction(word)

    switch mode {
    case 0:
        opext, i := em.DecodeOther(word)
        handler, ok := otherInstructions[(opcode<<3)|opext]
        if ok {
            handler(em, a, i)
        } else {
            return &InvalidOpcodeError{word}
        }

    case 1:
        b, d := em.DecodeALU(word)
        handler, ok := aluInstructions[opcode]
        if ok {
            handler(em, a, b, d)
        } else {
            return &InvalidOpcodeError{word}
        }

    case 2:
        d, i := em.DecodeJMI(word)
        handler, ok := jumpMemInstructions[opcode]
        if ok {
            handler(em, a, d, i)
        } else {
            return &InvalidOpcodeError{word}
        }

    case 3:
        d, i := em.DecodeJMI(word)
        handler, ok := aluImmInstructions[opcode]
        if ok {
            handler(em, a, d, i)
        } else {
            return &InvalidOpcodeError{word}
        }

    default:
        return &InvalidOpcodeError{word}
    }

    return nil
}

// Function Run runs until a halting condition is encountered.
func (em *Emulator) Run() (err error) {
    em.Running = true

    for em.Running {
        err = em.RunOne()
        if err != nil {
            return err
        }
    }

    return nil
}

// Function DumpState dumps the state of the processor to stdout.
func (em *Emulator) DumpState() {
    fmt.Printf("PC: 0x%08X/%d\n", em.PC, em.PC)
    fmt.Printf("\n")

    for i := 0; i < 32; i++ {
        v := em.Regs[i]
        fmt.Printf("%3s: 0x%08X/%d\n", RegisterNames[i], v, v)
    }

    fmt.Printf("\n")
}
