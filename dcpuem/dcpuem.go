// dcpuem.go - Constants, variables and the Emulator type.

package dcpuem

import (
    "fmt"
    "log"
    "time"
)

// Type *Emulator represents a DCPU-16 emulator.
type Emulator struct {
    // The A,B,C,X,Y,Z,I,J registers.
    Regs [8]uint16

    // The value of the program counter at the start of the instruction.
    LastPC uint16

    // The program counter.
    PC uint16

    // The stack pointer.
    SP uint16

    // The excess register.
    EX uint16

    // The interrupt address register.
    IA uint16

    // The main RAM.
    RAM []uint16

    // A FIFO queue of pending interrupts.
    Interrupts []uint16

    // List of connected hardware devices.
    Hardware []Device

    // Execution traces & debugging info will be logged to this logger.
    Logger *log.Logger

    // When set to false, the Run method will stop.
    Running bool

    // If set, the emulator will continue to skip instructions until a non-test instruction is
    // encountered.
    Skip bool

    // If set, interrupts will be queued instead of executed immediately. Internally, an interrupt
    // will be allowed to execute at the start of every instruction unless this is set to true.
    InterruptQueueing bool

    // The clock ticker. If this is not nil, execution will be limited by it.
    ClockTicker *time.Ticker
}

// Function NewEmulator creates, initialises and returns a new emulator.
func NewEmulator() (em *Emulator) {
    em = new(Emulator)
    em.Reset()

    em.Hardware = make([]Device, 0, 16)
    em.Logger = nil

    return em
}

// Function Reset resets the emulator's CPU state.
func (em *Emulator) Reset() {
    for i := 0; i < 8; i++ {
        em.Regs[i] = 0
    }

    em.LastPC = 0
    em.PC = 0
    em.SP = 0
    em.EX = 0
    em.IA = 0

    em.RAM = make([]uint16, 1024)

    em.Interrupts = make([]uint16, 0, 256)
    em.Running = false
    em.Skip = false
    em.InterruptQueueing = false
}

// Function Log logs information to the logger.
func (em *Emulator) Log(s string, args ...interface{}) {
    if em.Logger != nil {
        s = fmt.Sprintf(s, args...)
        em.Logger.Output(2, s)
    }
}

// Function LogInstruction logs an instruction execution to the logger.
func (em *Emulator) LogInstruction(s string, args ...interface{}) {
    if em.Logger != nil {
        s = fmt.Sprintf(s, args...)
        em.Logger.Output(2, fmt.Sprintf("[0x%04X] %s", em.LastPC, s))
    }
}

// Function DumpState dumps the state of the emulator to the logger.
func (em *Emulator) DumpState() {
    em.Log("A:  0x%04X   Y:  0x%04X", em.Regs[A], em.Regs[Y])
    em.Log("B:  0x%04X   Z:  0x%04X", em.Regs[B], em.Regs[Z])
    em.Log("C:  0x%04X   I:  0x%04X", em.Regs[C], em.Regs[I])
    em.Log("X:  0x%04X   J:  0x%04X", em.Regs[X], em.Regs[J])
    em.Log("PC: 0x%04X   SP: 0x%04X", em.PC, em.SP)
    em.Log("EX: 0x%04X   IA: 0x%04X", em.EX, em.IA)
}

// Type InstructionMode represents an instruction mode.
type InstructionMode uint8

// InstructionMode constants.
const (
    _ InstructionMode = iota
    Basic
    Special
)

// Type OperandMode represents an operand mode.
type OperandMode uint8

// OperandMode constants.
const (
    _ OperandMode = iota
    Literal
    Register
    Memory
    SP
    PC
    EX
)

// Type Operand represents an operand.
type Operand struct {
    Mode   OperandMode
    Info   uint16
    String string
}

var NilOperand = Operand{0, 0, ""}

type Error struct {
    Type    ErrorType
    Message string
    Arg     int
}

func (e *Error) Error() (str string) {
    return e.Message
}

type ErrorType uint16

const (
    _ ErrorType = iota
    ErrInvalidOpcode
    ErrInvalidOperand
    ErrCrashLoop
    ErrInvalidHardwareIndex
    ErrInvalidHardwareCommand
)

// Register names.
var RegisterNames = [8]string{"A", "B", "C", "X", "Y", "Z", "I", "J"}

// Register constants.
const (
    A = 0
    B = 1
    C = 2
    X = 3
    Y = 4
    Z = 5
    I = 6
    J = 7
)
