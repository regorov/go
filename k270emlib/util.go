package k270emlib

import (
    
)

const (
    INT_T0 = 0x10
    INT_T1 = 0x11
    INT_T2 = 0x12
)

const (
    P_TR = 0x10
    
    P_DIN0 = 0x20
    P_DIN1 = 0x21
    P_DIN2 = 0x22
    P_DIN3 = 0x23
    P_DOUT0 = 0x24
    P_DOUT1 = 0x25
    P_DOUT2 = 0x26
    P_DOUT3 = 0x27
    P_DMODE0 = 0x28
    P_DMODE1 = 0x29
    P_DMODE2 = 0x2A
    P_DMODE3 = 0x2B
)

const (
    E_REG_INDEX_OUT_OF_RANGE = iota
    E_INCORRECT_MODE
)

const VMEM_HEIGHT = 48
const VMEM_WIDTH = 128
const VMEM_SIZE = VMEM_HEIGHT * VMEM_WIDTH

var RegisterNames = []string{
    "z",  "q",  "k0", "k1", "a0", "a1", "a2", "a3",
    "v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7",
}

var WordRegisterNames = []string{
    "z:q",   "k0:k1", "a0:a1", "a2:a3",
    "v0:v1", "v2:v3", "v4:v5", "v6:v7",
}

type Error struct {
    ID int
    Message string
}

func NewError(id int, message string) (err *Error) {
    err = new(Error)
    err.ID = id
    err.Message = message
    return err
}

func (err *Error) Error() (str string) {
    return err.Message
}
