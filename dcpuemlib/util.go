package dcpuemlib

import (
    
)

// Misc operand constants
const (
    MISC_SP = iota  // SP
    MISC_PC         // PC
    MISC_O          // O
    MISC_PUSH       // PUSH
    MISC_POP        // POP
)

// Opcode type constants
const (
    OP_BASIC = iota     // Basic opcodes
    OP_EXT              // Extended opcodes (o = 0)
)

var RegisterNames = []string{"A", "B", "C", "X", "Y", "Z", "I", "J"}
