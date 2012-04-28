package k680emlib

import (
    "fmt"
)

type InvalidOpcodeError struct {
    Word uint32
}

func (err *InvalidOpcodeError) Error() string {
    return fmt.Sprintf("Invalid opcode error: %032b", err.Word)
}

var RegisterNames = []string{
    "%z", "%sp", "%q0", "%q1", "%cs", "%ds", "%ss", "%us",
    "%a0", "%a1", "%a2", "%a3", "%k0", "%k1", "%k2", "%k3",
    "%v0", "%v1", "%v2", "%v3", "%v4", "%v5", "%v6", "%v7",
    "%t0", "%t1", "%t2", "%t3", "%t4", "%t5", "%t6", "%t7",
}

const (
    Z  = 0
    SP = 1
    Q0 = 2
    Q1 = 3
    CS = 4
    DS = 5
    SS = 6
    US = 7
    A0 = 8
    A1 = 9
    A2 = 10
    A3 = 11
    K0 = 12
    K1 = 13
    K2 = 14
    K3 = 15
    V0 = 16
    V1 = 17
    V2 = 18
    V3 = 19
    V4 = 20
    V5 = 21
    V6 = 22
    V7 = 23
    T0 = 24
    T1 = 25
    T2 = 26
    T3 = 27
    T4 = 28
    T5 = 29
    T6 = 30
    T7 = 31
)

func sign14(v uint16) (s string) {
    if v&0x2000 != 0 {
        return "-"
    }
    return ""
}

func signPlus14(v uint16) (s string) {
    if v&0x2000 != 0 {
        return "-"
    }
    return "+"
}

func abs14(v uint16) (x uint16) {
    if v&0x2000 != 0 {
        return 0x4000 - v
    }
    return v
}

func sign16(v uint16) (s string) {
    if v&0x8000 != 0 {
        return "-"
    }
    return ""
}

func signPlus16(v uint16) (s string) {
    if v&0x8000 != 0 {
        return "-"
    }
    return "+"
}

func abs16(v uint16) (x uint16) {
    if v&0x8000 != 0 {
        return -v
    }
    return v
}

func parseSigned14(x uint16) (y uint32) {
    y = uint32(x)
    if x&0x2000 != 0 {
        return y - 0x4000
    }
    return y
}

func parseSigned16(x uint16) (y uint32) {
    y = uint32(x)
    if x&0x8000 != 0 {
        return y - 0x10000
    }
    return y
}
