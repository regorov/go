// Package dcpuemlib emulates a DCPU-16 processor.
package dcpuemlib

import (
    "fmt"
    "io"
)

type Emulator struct {
    LastPC uint16
    Registers [8]uint16
    SP uint16
    PC uint16
    O uint16
    
    RAM []uint16
    
    SkipNext bool
    Running bool
    TraceFile io.Writer
}

func NewEmulator() (em *Emulator) {
    em = new(Emulator)
    em.Reset()
    em.ResetMemory()
    em.SkipNext = false
    em.Running = true
    em.TraceFile = nil
    return em
}

func (em *Emulator) Reset() {
    em.SP = 0
    em.PC = 0
    em.O = 0
    
    for i := 0; i < 8; i++ {
        em.Registers[i] = 0
    }
}

func (em *Emulator) ResetMemory() {
    em.RAM = make([]uint16, 1024)
}

func (em *Emulator) GrowMemory(newsize int) {
    if newsize == 0 {
        newsize = (cap(em.RAM) + 1) * 2
    }
    
    if newsize > 0x10000 {
        newsize = 0x10000
    }
    
    m := make([]uint16, newsize)
    copy(m, em.RAM)
    em.RAM = m
}

func (em *Emulator) LoadProgram(program []uint16) {
    if len(em.RAM) < len(program) {
        em.GrowMemory(len(program))
    }
    
    copy(em.RAM, program)
}

func (em *Emulator) LoadProgramBytesBE(program []byte) {
    if len(em.RAM) < (len(program) * 2) {
        em.GrowMemory(len(program) * 2)
    }
    
    for i := 0; i < len(program) / 2; i++ {
        high := uint16(program[i * 2])
        low := uint16(program[(i * 2) + 1])
        em.RAM[i] = (high << 8) | low
    }
}

func (em *Emulator) LoadProgramBytesLE(program []byte) {
    if len(em.RAM) < (len(program) * 2) {
        em.GrowMemory(len(program) * 2)
    }
    
    for i := 0; i < len(program) / 2; i++ {
        low := uint16(program[i * 2])
        high := uint16(program[(i * 2) + 1])
        em.RAM[i] = (high << 8) | low
    }
}

func (em *Emulator) MemoryLoad(address uint16) (value uint16) {
    if address >= uint16(len(em.RAM)) {
        return 0
    } else {
        return em.RAM[address]
    }
    
    return 0
}

func (em *Emulator) MemoryStore(address uint16, value uint16) {
    if int(address) >= len(em.RAM) {
        newsize := cap(em.RAM) + 1
        
        for int(address) >= newsize {
            newsize *= 2
        }
        
        em.GrowMemory(newsize)
    }
    
    em.RAM[address] = value
}

func (em *Emulator) Push(value uint16) {
    em.MemoryStore(em.SP, value)
    em.SP++
}

func (em *Emulator) Pop() (value uint16) {
    em.SP--
    return em.MemoryLoad(em.SP)
}

func (em *Emulator) FetchWord() (word uint16) {
    word = em.MemoryLoad(em.PC)
    em.PC++
    return word
}

func (em *Emulator) DecodeOperand(n uint8) (operand Operand) {
    if n < 0x08 { // register
        return NewRegisterOperand(n)
    
    } else if n < 0x10 { // [register]
        return NewMemoryOperand(em.Registers[n & 0x7])
    
    } else if n < 0x18 { // [next word + register]
        addr := em.FetchWord() + em.Registers[n & 0x7]
        return NewMemoryOperand(addr)
    
    } else if n == 0x18 { // POP
        return NewPopOperand()
    
    } else if n == 0x19 { // PEEK
        return NewMemoryOperand(em.SP)
    
    } else if n == 0x1A { // PUSH
        return NewPushOperand()
    
    } else if n == 0x1B { // SP
        return NewSPOperand()
    
    } else if n == 0x1C { // PC
        return NewPCOperand()
    
    } else if n == 0x1D { // O
        return NewOOperand()
    
    } else if n == 0x1E { // [next word]
        return NewMemoryOperand(em.FetchWord())
    
    } else if n == 0x1F { // next word
        return NewLiteralOperand(em.FetchWord())
    
    } else {
        return NewLiteralOperand(uint16(n & 0x1F))
    }
    
    return nil
}

func (em *Emulator) DecodeInstruction(word uint16) (fmt uint8, opcode uint8, dest Operand, src Operand) {
    o := word & 0x000F
    a := (word & 0x03F0) >> 4
    b := (word & 0xFC00) >> 10
    
    if o == 0 { // Opcode is in the a field
        dest = em.DecodeOperand(uint8(b))
        return OP_EXT, uint8(a), dest, nil
    
    } else {
        dest = em.DecodeOperand(uint8(a)) // Important that we handle A first
        src = em.DecodeOperand(uint8(b))
        return OP_BASIC, uint8(o), dest, src
    }
    
    return 0, 0, nil, nil
}

func (em *Emulator) LogInstruction(format string, args ...interface{}) {
    if em.TraceFile != nil {
        format = fmt.Sprintf(format, args...)
        fmt.Fprintf(em.TraceFile, "[0x%04X] %s\n", em.LastPC, format)
    }
}

func (em *Emulator) RunOne() {
    em.LastPC = em.PC
    word := em.FetchWord()
    fmt, opcode, dest, src := em.DecodeInstruction(word)
    
    // Make sure operands have fetched extra program words before we skip
    if em.SkipNext {
        em.SkipNext = false
        return
    }
    
    switch fmt {
    case OP_BASIC:
        switch opcode {
        case 0x1: // SET
            v := src.Load(em)
            dest.Store(em, v)
            em.LogInstruction("SET %s, %s -- value transferred was 0x%04X", dest.String(), src.String(), v)
        
        case 0x2: // ADD
            d := dest.Load(em)
            s := src.Load(em)
            v := uint32(d) + uint32(s)
            em.O = uint16(v >> 16)
            dest.Store(em, uint16(v))
            em.LogInstruction("ADD %s, %s -- 0x%04X + 0x%04X = 0x%08X", dest.String(), src.String(), d, s, v)
        
        case 0x3: // SUB
            d := dest.Load(em)
            s := src.Load(em)
            v := uint32(d) - uint32(s)
            em.O = uint16(v >> 16)
            dest.Store(em, uint16(v))
            em.LogInstruction("SUB %s, %s -- 0x%04X - 0x%04X = 0x%08X", dest.String(), src.String(), d, s, v)
        
        case 0x4: // MUL
            d := dest.Load(em)
            s := src.Load(em)
            v := uint32(d) * uint32(s)
            em.O = uint16(v >> 16)
            dest.Store(em, uint16(v))
            em.LogInstruction("MUL %s, %s -- 0x%04X * 0x%04X = 0x%08X", dest.String(), src.String(), d, s, v)
        
        case 0x5: // DIV
            d := dest.Load(em)
            s := src.Load(em)
            if s == 0 {
                em.O = 0
                dest.Store(em, 0)
                em.LogInstruction("DIV %s, %s -- 0x%04X / 0x%04X = DIV/0!", dest.String(), src.String(), d, s)
            
            } else {
                v := uint32(d) / uint32(s)
                em.O = uint16(v >> 16)
                dest.Store(em, uint16(v))
                em.LogInstruction("DIV %s, %s -- 0x%04X / 0x%04X = 0x%04X", dest.String(), src.String(), d, s, v)
            }
        
        case 0x6: // MOD
            d := dest.Load(em)
            s := src.Load(em)
            
            if s == 0 {
                dest.Store(em, 0)
                em.LogInstruction("MOD %s, %s -- 0x%04X %% 0x%04X = DIV/0", dest.String(), src.String(), d, s)
            
            } else {
                v := d % s
                dest.Store(em, v)
                em.LogInstruction("MOD %s, %s -- 0x%04X %% 0x%04X = 0x%04X", dest.String(), src.String(), d, s, v)
            }
        
        case 0x7: // SHL
            d := dest.Load(em)
            s := src.Load(em)
            v := uint32(d) << uint32(s)
            em.O = uint16(v >> 16)
            dest.Store(em, uint16(v))
            em.LogInstruction("SHL %s, %s -- 0x%04X << 0x%04X = 0x%08X", dest.String(), src.String(), d, s, v)
        
        case 0x8: // SHR
            d := dest.Load(em)
            s := src.Load(em)
            v := uint32(d) >> uint32(s)
            em.O = uint16(v >> 16)
            dest.Store(em, uint16(v))
            em.LogInstruction("SHR %s, %s -- 0x%04X >> 0x%04X = 0x%04X", dest.String(), src.String(), d, s, v)
        
        case 0x9: // AND
            d := dest.Load(em)
            s := src.Load(em)
            v := d & s
            dest.Store(em, v)
            em.LogInstruction("AND %s, %s -- 0x%04X & 0x%04X = 0x%04X", dest.String(), src.String(), d, s, v)
        
        case 0xA: // BOR
            d := dest.Load(em)
            s := src.Load(em)
            v := d | s
            dest.Store(em, v)
            em.LogInstruction("BOR %s, %s -- 0x%04X | 0x%04X = 0x%04X", dest.String(), src.String(), d, s, v)
        
        case 0xB: // XOR
            d := dest.Load(em)
            s := src.Load(em)
            v := d ^ s
            dest.Store(em, v)
            em.LogInstruction("XOR %s, %s -- 0x%04X ^ 0x%04X = 0x%04X", dest.String(), src.String(), d, s, v)
        
        case 0xC: // IFE
            d := dest.Load(em)
            s := src.Load(em)
            if d == s {
                em.LogInstruction("IFE %s, %s -- 0x%04X == 0x%04X, executing next", dest.String(), src.String(), d, s)
            } else {
                em.SkipNext = true
                em.LogInstruction("IFE %s, %s -- 0x%04X != 0x%04X, skipping next", dest.String(), src.String(), d, s)
            }
        
        case 0xD: // IFN
            d := dest.Load(em)
            s := src.Load(em)
            if d != s {
                em.LogInstruction("IFN %s, %s -- 0x%04X != 0x%04X, executing next", dest.String(), src.String(), d, s)
            } else {
                em.SkipNext = true
                em.LogInstruction("IFN %s, %s -- 0x%04X == 0x%04X, skipping next", dest.String(), src.String(), d, s)
            }
        
        case 0xE: // IFG
            d := dest.Load(em)
            s := src.Load(em)
            if d > s {
                em.LogInstruction("IFG %s, %s -- 0x%04X > 0x%04X, executing next", dest.String(), src.String(), d, s)
            } else {
                em.SkipNext = true
                em.LogInstruction("IFG %s, %s -- 0x%04X <= 0x%04X, skipping next", dest.String(), src.String(), d, s)
            }
        
        case 0xF: // IFB
            d := dest.Load(em)
            s := src.Load(em)
            if d & s != 0 {
                em.LogInstruction("IFB %s, %s -- 0x%04X & 0x%04X != 0, executing next", dest.String(), src.String(), d, s)
            } else {
                em.SkipNext = true
                em.LogInstruction("IFB %s, %s -- 0x%04X & 0x%04X == 0, skipping next", dest.String(), src.String(), d, s)
            }
        }
    
    case OP_EXT:
        switch opcode {
        case 0x01: // JSR
            em.Push(em.PC)
            em.PC = dest.Load(em)
            em.LogInstruction("JSR %s -- dest = %04X", dest.String(), em.PC)
        }
    }
}

func (em *Emulator) Run() {
    em.Running = true
    
    for em.Running {
        em.RunOne()
    }
}

func (em *Emulator) DumpState() {
    fmt.Printf("A: 0x%04X   Y: 0x%04X\n", em.Registers[0], em.Registers[4])
    fmt.Printf("B: 0x%04X   Z: 0x%04X\n", em.Registers[1], em.Registers[5])
    fmt.Printf("C: 0x%04X   I: 0x%04X\n", em.Registers[2], em.Registers[6])
    fmt.Printf("X: 0x%04X   J: 0x%04X\n", em.Registers[3], em.Registers[7])
    fmt.Printf("SP: 0x%04X\nPC: 0x%04X\nO: 0x%04X\n", em.SP, em.PC, em.O)
}
