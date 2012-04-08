package k270emlib

import (
    
)

// Standard interrupt numbers.
const (
    INT_T0 = 0x10   // T0 - Timer 0 fired
    INT_T1 = 0x11   // T1 - Timer 1 fired
    INT_T2 = 0x12   // T2 - Timer 2 fired
)

// Standard port numbers
const (
    P_TR = 0x10     // TR - Reset timer
    
    P_DIN0 = 0x20   // DIN0 - Digital inputs
    P_DIN1 = 0x21   // DIN1 - Digital inputs
    P_DIN2 = 0x22   // DIN2 - Digital inputs
    P_DIN3 = 0x23   // DIN3 - Digital inputs
    P_DOUT0 = 0x24  // DOUT0 - Digital outputs
    P_DOUT1 = 0x25  // DOUT1 - Digital outputs
    P_DOUT2 = 0x26  // DOUT2 - Digital outputs
    P_DOUT3 = 0x27  // DOUT3 - Digital outputs
    P_DMODE0 = 0x28 // DMODE0 - Digital pin mode
    P_DMODE1 = 0x29 // DMODE1 - Digital pin mode
    P_DMODE2 = 0x2A // DMODE2 - Digital pin mode
    P_DMODE3 = 0x2B // DMODE3 - Digital pin mode
)

// Error ID constants
const (
    E_REG_INDEX_OUT_OF_RANGE = iota     // Invalid register index
    E_INCORRECT_MODE                    // Attempted to call GetDigitalOutput on an input pin, or
                                        // SetDigitalInput on an output pin.
)

// Height of the character display, in characters
const VMEM_HEIGHT = 48

// Width of the character display, in characters
const VMEM_WIDTH = 128

// Total size of the character display, in characters
const VMEM_SIZE = VMEM_HEIGHT * VMEM_WIDTH

// Maps register numbers to register names
var RegisterNames = []string{
    "z",  "q",  "k0", "k1", "a0", "a1", "a2", "a3",
    "v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7",
}

// Maps register number (divided by 2) to word register names
var WordRegisterNames = []string{
    "z:q",   "k0:k1", "a0:a1", "a2:a3",
    "v0:v1", "v2:v3", "v4:v5", "v6:v7",
}

// Type Error represents an emulation error. It implements the error interface.
type Error struct {
    ID int              // The error ID; one of the E_* constants
    Message string      // An associated textual message
}

// Function NewError creates and returns a new Error.
func NewError(id int, message string) (err *Error) {
    err = new(Error)
    err.ID = id
    err.Message = message
    return err
}

// Function Error.Error returns the error's message.
func (err *Error) Error() (str string) {
    return err.Message
}
