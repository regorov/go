package main

import (
    "fmt"
)

type Register uint8

const (
    V0 Register = 0
    V1 Register = 1
    V2 Register = 2
    V3 Register = 3
    V4 Register = 4
    V5 Register = 5
    V6 Register = 6
    V7 Register = 7
    A0 Register = 8
    A1 Register = 9
    A2 Register = 10
    A3 Register = 11
    Q0 Register = 12
    Q1 Register = 13
    SP Register = 14
    AT Register = 15

    NoRegister Register = 255
)

type MemSize uint8

const (
    Mem8  MemSize = 8
    Mem16 MemSize = 16
    Mem32 MemSize = 32
)

type OperandType uint8

const (
    DynamicType OperandType = iota
    LiteralType
    BitRegType
)

var RegisterNames = []string{
    "v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7",
    "a0", "a1", "a2", "a3", "q0", "q1", "sp", "at",
}

func (r Register) String() (str string) {
    if r == NoRegister {
        return "0"
    }

    return "%" + RegisterNames[r]
}

type Operand interface {
    String() string
    Length() uint32
    ReduceLabel(map[string]uint32)
    SetSize(MemSize)
    SatisfiesType(OperandType) bool
    LiteralValue() uint32
    BitValue() uint8
    EncodeKey() byte
    EncodeExtra([]byte)
}

type LiteralOperand struct {
    coord Coord
    Literal
}

func (o *LiteralOperand) SetSize(size MemSize) {

}

func (o *LiteralOperand) SatisfiesType(t OperandType) (result bool) {
    return t == DynamicType || t == LiteralType
}

func (o *LiteralOperand) LiteralValue() (v uint32) {
    return o.Literal.Value()
}

func (o *LiteralOperand) BitValue() (v uint8) {
    return 0
}

func (o *LiteralOperand) EncodeKey() (key byte) {
    if o.Length() == 4 {
        return 0xFF
    }

    return byte(o.Literal.Value() & 0x7F)
}

func (o *LiteralOperand) EncodeExtra(extra []byte) {
    if o.Length() == 4 {
        v := o.Literal.Value()
        extra[0] = byte(v >> 24)
        extra[1] = byte(v >> 16)
        extra[2] = byte(v >> 8)
        extra[3] = byte(v)
    }
}

type RegisterOperand struct {
    coord Coord
    num   Register
}

func (o *RegisterOperand) String() (str string) {
    return o.num.String()
}

func (o *RegisterOperand) Length() (length uint32) {
    return 0
}

func (o *RegisterOperand) ReduceLabel(labelMap map[string]uint32) {

}

func (o *RegisterOperand) SetSize(size MemSize) {

}

func (o *RegisterOperand) SatisfiesType(t OperandType) (result bool) {
    return t == DynamicType
}

func (o *RegisterOperand) LiteralValue() (v uint32) {
    return 0
}

func (o *RegisterOperand) BitValue() (v uint8) {
    return 0
}

func (o *RegisterOperand) EncodeKey() (key byte) {
    return 0x80 | byte(o.num)
}

func (o *RegisterOperand) EncodeExtra(extra []byte) {

}

type MemoryOperand struct {
    coord  Coord
    size   MemSize
    reg    Register
    disp   Literal
    length uint32
}

func (o *MemoryOperand) String() (str string) {
    return fmt.Sprintf("%d[%s + %s]", o.size, o.reg.String(), o.disp.String())
}

func (o *MemoryOperand) Length() (length uint32) {
    if o.length != 256 {
        return o.length
    }

    v := int32(o.disp.Value())

    if o.reg == NoRegister {
        if v < -0x8000 || v >= 0x8000 {
            errChan <- &AsmError{o.coord, fmt.Sprintf("Integer displacement out of range (-0x8000 to 0x7FFF): 0x%08X", v)}

        } else {
            length = 4
        }

    } else if !o.disp.Reduced() {
        length = 2

    } else if v == 0 {
        length = 0

    } else if v >= -0x8000 && v < 0x8000 {
        length = 2

    } else {
        errChan <- &AsmError{o.coord, fmt.Sprintf("Integer displacement out of range (-0x8000 to 0x7FFF): 0x%08X", v)}
    }

    o.length = length
    return length
}

func (o *MemoryOperand) ReduceLabel(labelMap map[string]uint32) {
    o.disp.ReduceLabel(labelMap)
}

func (o *MemoryOperand) SetSize(size MemSize) {
    o.size = size
}

func (o *MemoryOperand) SatisfiesType(t OperandType) (result bool) {
    return t == DynamicType
}

func (o *MemoryOperand) LiteralValue() (v uint32) {
    return 0
}

func (o *MemoryOperand) BitValue() (v uint8) {
    return 0
}

func (o *MemoryOperand) EncodeKey() (key byte) {
    if o.reg == NoRegister {
        return 0xFE
    }

    if o.Length() == 2 {
        switch o.size {
        case Mem8:
            return 0xC0 | byte(o.reg)
        case Mem16:
            return 0xD0 | byte(o.reg)
        case Mem32:
            return 0xE0 | byte(o.reg)
        }

    } else {
        switch o.size {
        case Mem8:
            return 0x90 | byte(o.reg)
        case Mem16:
            return 0xA0 | byte(o.reg)
        case Mem32:
            return 0xB0 | byte(o.reg)
        }
    }

    return 0
}

func (o *MemoryOperand) EncodeExtra(extra []byte) {
    if o.reg == NoRegister {
        v := o.disp.Value()
        extra[0] = byte(v >> 24)
        extra[1] = byte(v >> 16)
        extra[2] = byte(v >> 8)
        extra[3] = byte(v)

    } else if o.Length() == 2 {
        v := uint16(int32(o.disp.Value()))
        extra[0] = byte(v >> 8)
        extra[1] = byte(v)
    }
}

type PCOperand struct {
    coord Coord
}

func (o *PCOperand) String() (str string) {
    return "pc"
}

func (o *PCOperand) Length() (length uint32) {
    return 0
}

func (o *PCOperand) ReduceLabel(labelMap map[string]uint32) {

}

func (o *PCOperand) SetSize(size MemSize) {

}

func (o *PCOperand) SatisfiesType(t OperandType) (result bool) {
    return t == DynamicType
}

func (o *PCOperand) LiteralValue() (v uint32) {
    return 0
}

func (o *PCOperand) BitValue() (v uint8) {
    return 0
}

func (o *PCOperand) EncodeKey() (key byte) {
    return 0xFC
}

func (o *PCOperand) EncodeExtra(extra []byte) {

}
