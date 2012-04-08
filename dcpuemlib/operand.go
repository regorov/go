package dcpuemlib

import (
    "fmt"
)

type Operand interface {
    Load(em *Emulator) (value uint16)
    Store(em *Emulator, value uint16)
    String() string
}

// =================== LiteralOperand =====================

type LiteralOperand struct {
    Value uint16
}

func NewLiteralOperand(value uint16) (operand *LiteralOperand) {
    operand = new(LiteralOperand)
    operand.Value = value
    return operand
}

func (operand *LiteralOperand) Load(em *Emulator) (value uint16) {
    return operand.Value
}

func (operand *LiteralOperand) Store(em *Emulator, value uint16) {
    panic("Assigning to a literal operand!")
}

func (operand *LiteralOperand) String() (str string) {
    return fmt.Sprintf("0x%04X", operand.Value)
}

// ================== RegisterOperand =====================

type RegisterOperand struct {
    Number uint8
}

func NewRegisterOperand(number uint8) (operand *RegisterOperand) {
    if number >= 8 {
        panic("Register index out of range")
    }
    
    operand = new(RegisterOperand)
    operand.Number = number
    return operand
}

func (operand *RegisterOperand) Load(em *Emulator) (value uint16) {
    return em.Registers[operand.Number]
}

func (operand *RegisterOperand) Store(em *Emulator, value uint16) {
    em.Registers[operand.Number] = value
}

func (operand *RegisterOperand) String() (str string) {
    return RegisterNames[operand.Number]
}

// =================== MemoryOperand ======================

type MemoryOperand struct {
    Address uint16
}

func NewMemoryOperand(address uint16) (operand *MemoryOperand) {
    operand = new(MemoryOperand)
    operand.Address = address
    return operand
}

func (operand *MemoryOperand) Load(em *Emulator) (value uint16) {
    return em.MemoryLoad(operand.Address)
}

func (operand *MemoryOperand) Store(em *Emulator, value uint16) {
    em.MemoryStore(operand.Address, value )
}

func (operand *MemoryOperand) String() (str string) {
    return fmt.Sprintf("[0x%04X]", operand.Address)
}

// =================== MiscOperand ========================

type MiscOperand struct {
    Type uint8
}

func NewSPOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_SP
    return operand
}

func NewPCOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_PC
    return operand
}

func NewOOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_O
    return operand
}

func NewPushOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_PUSH
    return operand
}

func NewPopOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_POP
    return operand
}

func (operand *MiscOperand) Load(em *Emulator) (value uint16) {
    switch operand.Type {
    case MISC_SP:
        return em.SP
    
    case MISC_PC:
        return em.PC
    
    case MISC_O:
        return em.O
    
    case MISC_PUSH:
        panic("Can't load from PUSH")
    
    case MISC_POP:
        return em.Pop()
    }
    
    panic("Invalid MISC type")
}

func (operand *MiscOperand) Store(em *Emulator, value uint16) {
    switch operand.Type {
    case MISC_SP:
        em.SP = value
    
    case MISC_PC:
        if value == em.LastPC {
            em.LogInstruction("'Crash loop' detected, halting execution")
            em.Running = false
        
        } else {
            em.PC = value
        }
    
    case MISC_O:
        em.O = value
    
    case MISC_PUSH:
        em.Push(value)
    
    case MISC_POP:
        panic("Can't store to POP")
    
    default:
        panic("Invalid MISC type")
    }
}

func (operand *MiscOperand) String() (str string) {
    switch operand.Type {
    case MISC_SP:
        return "SP"
    case MISC_PC:
        return "PC"
    case MISC_O:
        return "O"
    case MISC_PUSH:
        return "PUSH"
    case MISC_POP:
        return "POP"
    }
    
    panic("Invalid MISC type")
}
