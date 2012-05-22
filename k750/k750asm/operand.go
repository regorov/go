package main

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

type Operand interface {
    String() string
    Length() uint32
    ReduceLabel(map[string]uint32)
    SatisfiesType(OperandType) bool
    LiteralValue() uint32
    BitValue() uint8
    EncodeKey() byte
    EncodeExtra([]byte)
}

type LiteralOperand struct {
    Literal
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
    num Register
}

func (o *RegisterOperand) String() (str string) {
    return "%" + RegisterNames[o.num]
}

func (o *RegisterOperand) SatisfiesType(t OperandType) (result bool) {
    return t == DynamicType
}

func (o *RegisterOperand) Length() (length uint32) {
    return 0
}

func (o *RegisterOperand) ReduceLabel(labelMap map[string]uint32) {

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

type PCOperand struct {
}

func (o *PCOperand) String() (str string) {
    return "pc"
}

func (o *PCOperand) SatisfiesType(t OperandType) (result bool) {
    return t == DynamicType
}

func (o *PCOperand) Length() (length uint32) {
    return 0
}

func (o *PCOperand) ReduceLabel(labelMap map[string]uint32) {

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
