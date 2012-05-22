package main

import (
    "fmt"
)

type Literal interface {
    String() string
    Length() uint32
    ReduceLabel(map[string]uint32)
    Value() uint32
}

type ConstantLiteral struct {
    value uint32
}

func (l *ConstantLiteral) String() (str string) {
    return fmt.Sprintf("0x%x", l.value)
}

func (l *ConstantLiteral) Length() (length uint32) {
    // Can either be packed into a sign-extended 7-bit inline (0) or a 32-bit extra (4)

    sv := int32(l.value)
    if sv >= -0x80 && sv < 0x80 {
        return 0
    }

    return 4
}

func (l *ConstantLiteral) ReduceLabel(labelMap map[string]uint32) {

}

func (l *ConstantLiteral) Value() (value uint32) {
    return l.value
}

type LabelLiteral struct {
    coord   Coord
    name    string
    value   uint32
    reduced bool
}

func (l *LabelLiteral) String() (str string) {
    if l.reduced {
        return fmt.Sprintf("%s(0x%08X)", l.name, l.value)
    }
    return l.name
}

func (l *LabelLiteral) Length() (length uint32) {
    // We don't know how big the actual value is but it's nearly always bigger than 128. Assume 32-bit extra.
    return 4
}

func (l *LabelLiteral) ReduceLabel(labelMap map[string]uint32) {
    value, ok := labelMap[l.name]

    if !ok {
        errChan <- &AsmError{l.coord, fmt.Sprintf("Label '%s' not defined", l.name)}
    }

    l.value = value
    l.reduced = true
}

func (l *LabelLiteral) Value() (value uint32) {
    if l.reduced {
        return l.value
    }

    return 0
}
