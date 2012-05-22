package dcpuemlib

import (
    "errors"
)

// MiscOperand type constants.
const (
    MISC_SP   = iota // SP
    MISC_PC          // PC
    MISC_O           // O
    MISC_PUSH        // PUSH
    MISC_POP         // POP
)

// Opcode class constants.
const (
    OP_BASIC = iota // Basic opcodes
    OP_EXT          // Extended opcodes (o = 0)
)

// Maps register numbers to names.
var RegisterNames = []string{"A", "B", "C", "X", "Y", "Z", "I", "J"}

// Errors.
var (
    ErrInvalidOpcode = errors.New("Invalid opcode")
)
