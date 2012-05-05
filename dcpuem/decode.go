// decode.go - Instruction decoding.

package dcpuem

import (
    "fmt"
)

// Function DecodeWord splits an instruction word into a, b, instruction mode and opcode components.
func (em *Emulator) DecodeWord(word uint16) (a uint16, b uint16, opcode uint16) {
    a = (word >> 10) & 0x3F
    b = (word >> 5) & 0x1F
    opcode = word & 0x1F

    return a, b, opcode
}

// Function DecodeOperand takes an operand specifier and returns the operand.
func (em *Emulator) DecodeOperand(x uint16, isA bool) (operand Operand, err error) {
    var mode OperandMode
    var info uint16
    var str string

    if (x & 0x38) == 0x00 {
        mode = Register
        info = x & 0x07
        str = RegisterNames[info]

    } else if (x & 0x38) == 0x08 {
        num := x & 0x07

        mode = Memory
        info = em.Regs[num]
        str = fmt.Sprintf("[%s]", RegisterNames[num])

    } else if (x & 0x38) == 0x10 {
        num := x & 0x07
        offset := em.FetchWord()

        mode = Memory
        info = em.Regs[num] + offset
        str = fmt.Sprintf("[%s+0x%04X]", RegisterNames[num], offset)

    } else if x == 0x18 {
        if isA {
            mode = Memory
            info = em.SP
            str = "[SP++]"
            em.SP++

        } else {
            em.SP--
            mode = Memory
            info = em.SP
            str = "[--SP]"
        }

    } else if x == 0x19 {
        mode = Memory
        info = em.SP
        str = "[SP]"

    } else if x == 0x1A {
        offset := em.FetchWord()

        mode = Memory
        info = em.SP + offset
        str = fmt.Sprintf("[SP+0x%04X]", offset)

    } else if x == 0x1B {
        mode = SP
        str = "SP"

    } else if x == 0x1C {
        mode = PC
        str = "PC"

    } else if x == 0x1D {
        mode = EX
        str = "EX"

    } else if x == 0x1E {
        mode = Memory
        info = em.FetchWord()
        str = fmt.Sprintf("[0x%04X]", info)

    } else if x == 0x1F {
        mode = Literal
        info = em.FetchWord()
        str = fmt.Sprintf("0x%04X", info)

    } else if (x&0x20) == 0x20 && isA {
        mode = Literal
        info = (x & 0x1F) - 1 // - 1 because literals are -1..30, not 0..31
        str = fmt.Sprintf("0x%04X", info)

    } else {
        return NilOperand, ErrInvalidOperand
    }

    return Operand{mode, info, str}, nil
}

// Function RunOne runs one instruction.
func (em *Emulator) RunOne() (err error) {
    em.ServiceInterrupt()

    em.LastPC = em.PC
    word := em.FetchWord()
    aSpec, bSpec, opcode := em.DecodeWord(word)

    a, err := em.DecodeOperand(aSpec, true)
    if err != nil {
        return err
    }

    var b Operand

    if opcode != 0 {
        b, err = em.DecodeOperand(bSpec, false)
        if err != nil {
            return err
        }
    }

    if em.Skip {
        if opcode < 0x10 || opcode > 0x17 { // I.e. not a IF* instruction.
            em.Skip = false
        }

        return nil
    }

    if opcode == 0 {
        handler := specialHandlers[bSpec]
        if handler == nil {
            em.Log("Invalid special opcode: 0x%02X", bSpec)
            return ErrInvalidOpcode
        }

        return handler(em, a)

    } else {
        handler := basicHandlers[opcode]
        if handler == nil {
            em.Log("Invalid basic opcode: 0x%02X", bSpec)
            return ErrInvalidOpcode
        }

        return handler(em, a, b)
    }

    return nil
}

// Function Run sets the running flag to true, then runs until it is false.
func (em *Emulator) Run() (err error) {
    em.Running = true

    for em.Running {
        err = em.RunOne()

        if err == ErrCrashLoop {
            em.Running = false
        } else if err != nil {
            return err
        }
    }

    return nil
}

// Function Load loads the value of an operand.
func (em *Emulator) Load(operand Operand) (value uint16) {
    switch operand.Mode {
    case Literal:
        return operand.Info

    case Register:
        return em.Regs[operand.Info]

    case Memory:
        return em.MemoryLoad(operand.Info)

    case SP:
        return em.SP

    case PC:
        return em.PC

    case EX:
        return em.EX
    }

    return 0
}

// Function Store stores the value of an operand.
func (em *Emulator) Store(operand Operand, value uint16) (err error) {
    switch operand.Mode {
    case Literal:
        //return ErrStoringToLiteral
        // Assignments to literals are supposed to fail silently
        em.Logger.Printf("Notice: ignoring assignment to literal operand.")

    case Register:
        em.Regs[operand.Info] = value

    case Memory:
        em.MemoryStore(operand.Info, value)

    case SP:
        em.SP = value

    case PC:
        if value == em.LastPC {
            em.Logger.Printf("Crash loop detected - halting")
            return ErrCrashLoop
        }

        em.PC = value

    case EX:
        em.EX = value
    }

    return nil
}
