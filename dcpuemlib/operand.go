package dcpuemlib

import (
    "fmt"
)

// Interface Operand defines methods that the different types of operand should implement.
type Operand interface {
    // Function Load should load the value from the operand and return it.
    Load(em *Emulator) (value uint16)
    
    // Function Store should store the value `value` to the operand.
    Store(em *Emulator, value uint16)
    
    // Function String should return a string representation of the operand, preferably in the
    //  format accepted by the assembler.
    String() string
}

// =================== LiteralOperand =====================

// Type LiteralOperand represents a literal value.
type LiteralOperand struct {
    // The value that this LiteralOperand represents.
    Value uint16
}

// Function NewLiteralOperand creates and returns a new LiteralOperand with the value `value`.
func NewLiteralOperand(value uint16) (operand *LiteralOperand) {
    operand = new(LiteralOperand)
    operand.Value = value
    return operand
}

// Function LiteralOperand.Load returns the operand's value.
func (operand *LiteralOperand) Load(em *Emulator) (value uint16) {
    return operand.Value
}

// Function LiteralOperand.Store will panic, as you cannot store a value to a literal.
func (operand *LiteralOperand) Store(em *Emulator, value uint16) {
    panic("Assigning to a literal operand!")
}

// Function LiteralOperand.String returns a hexidecimal representation of the value.
func (operand *LiteralOperand) String() (str string) {
    return fmt.Sprintf("0x%04X", operand.Value)
}

// ================== RegisterOperand =====================

// Type RegisterOperand represents a value stored in a general-purpose register.
type RegisterOperand struct {
    // The number of the register.
    Number uint8
}

// Function NewRegisterOperand creates and returns a new RegisterOperand with the number `number`.
//  It will panic if the number is out of range (the range is 0 to 7, inclusive).
func NewRegisterOperand(number uint8) (operand *RegisterOperand) {
    if number >= 8 {
        panic("Register index out of range")
    }
    
    operand = new(RegisterOperand)
    operand.Number = number
    return operand
}

// Function RegisterOperand.Load returns the value in the Emulator's register numbered Number. 
func (operand *RegisterOperand) Load(em *Emulator) (value uint16) {
    return em.Registers[operand.Number]
}

// Function RegisterOperand.Store stores `value` into the Emulator's register numbered Number.
func (operand *RegisterOperand) Store(em *Emulator, value uint16) {
    em.Registers[operand.Number] = value
}

// Function RegisterOperand.String returns the name of the register represented by the operand.
func (operand *RegisterOperand) String() (str string) {
    return RegisterNames[operand.Number]
}

// =================== MemoryOperand ======================

// Type MemoryOperand represents a value stored in the Emulator's RAM.
type MemoryOperand struct {
    // The address in memory.
    Address uint16
}

// Function NewMemoryOperand creates and returns a new MemoryOperand with the address `address`.
func NewMemoryOperand(address uint16) (operand *MemoryOperand) {
    operand = new(MemoryOperand)
    operand.Address = address
    return operand
}

// Function MemoryOperand.Load returns the value in the Emulator's RAM at address `Address`.
func (operand *MemoryOperand) Load(em *Emulator) (value uint16) {
    return em.MemoryLoad(operand.Address)
}

// Function MemoryOperand.Store stores the value `value` into the Emulator's RAM at address
//  `Address`.
func (operand *MemoryOperand) Store(em *Emulator, value uint16) {
    em.MemoryStore(operand.Address, value )
}

// Function MemoryOperand.String returns a string representation of the operand.
func (operand *MemoryOperand) String() (str string) {
    return fmt.Sprintf("[0x%04X]", operand.Address)
}

// =================== MiscOperand ========================

// Type MiscOperand represents one of:
//   * the stack pointer
//   * the program counter
//   * the overflow register
//   * a value pushed onto the stack
//   * a value popped off the stack
type MiscOperand struct {
    Type uint8
}

// Function NewSPOperand creates and returns a new MiscOperand referring to the stack pointer.
func NewSPOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_SP
    return operand
}

// Function NewPCOperand creates and returns a new MiscOperand referring to the program counter.
func NewPCOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_PC
    return operand
}

// Function NewOOperand creates and returns a new MiscOperand referring to the overflow register.
func NewOOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_O
    return operand
}

// Function NewPushOperand creates and returns a new MiscOperand referring to a value pushed onto
//  the stack.
func NewPushOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_PUSH
    return operand
}

// Function NewPopOperand creates and returns a new MiscOperand referring to a value popped off the
//  stack.
func NewPopOperand() (operand *MiscOperand) {
    operand = new(MiscOperand)
    operand.Type = MISC_POP
    return operand
}

// Function MiscOperand.Load loads the value from the operand.
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

// Function MiscOperand.Store stores the value `value` to the operand.
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

// Function MiscOperand.String returns a string representation of the operand.
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
