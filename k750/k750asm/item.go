package main

import (
    "fmt"
    "strings"
)

type OperandMode uint8

const (
    OperandModeNone OperandMode = iota
    OperandModeD
    OperandModeDD
    OperandModeDDD
    OperandModeDDL
    OperandModeB
    OperandModeBD
    OperandModeBB
    OperandModeBDL
    OperandModeDLB
)

var OperandTypeLookup = [][]OperandType{
    {},
    {DynamicType},
    {DynamicType, DynamicType},
    {DynamicType, DynamicType, DynamicType},
    {DynamicType, DynamicType, LiteralType},
    {BitRegType},
    {BitRegType, DynamicType},
    {BitRegType, BitRegType},
    {BitRegType, DynamicType, LiteralType},
    {DynamicType, LiteralType, BitRegType},
}

var LengthLookup = []uint32{
    1,
    2,
    3,
    4,
    3,
    1,
    3,
    2,
    3,
    3,
}

type Item interface {
    String() string
    VerifyAndReduce()
    Encode(map[string]uint32, []byte)
    Label() (string, bool)
    Length() uint32
    Offset() uint32
    SetOffset(uint32)
    Encoded() []byte
}

type Instruction struct {
    coord       Coord
    name        string
    operands    []Operand
    operandMode OperandMode
    length      uint32
    offset      uint32
    encoded     []byte
}

func (item *Instruction) String() (str string) {
    if len(item.operands) == 0 {
        return item.name
    }

    operandStrings := make([]string, len(item.operands))

    for i, operand := range item.operands {
        operandStrings[i] = operand.String()
    }

    return fmt.Sprintf("%s %s", item.name, strings.Join(operandStrings, ", "))
}

func (item *Instruction) VerifyAndReduce() {
    defer waitGroup.Done()

    var operandMode OperandMode

    switch item.name {
    case "nop", "popa", "pusha", "ret", "reti":
        operandMode = OperandModeNone

    case "call", "jmp", "int", "pop", "ppn", "push":
        operandMode = OperandModeD

    case "mov", "neg", "not", "ppi", "rih":
        operandMode = OperandModeDD

    case "add", "and", "jac", "jas", "jeq", "jge", "jgt", "jle", "jlt", "jne", "or", "ppr", "ppw", "sub", "xor":
        operandMode = OperandModeDDD

    case "rll", "rlr", "shl", "lshr", "ashr":
        operandMode = OperandModeDDL

    case "cb", "sb":
        operandMode = OperandModeB

    case "jbc", "jbs":
        operandMode = OperandModeBD

    case "ab", "ob", "xb":
        operandMode = OperandModeBB

    case "ldb":
        operandMode = OperandModeBDL

    case "stb":
        operandMode = OperandModeDLB

    default:
        errChan <- &AsmError{item.coord, fmt.Sprintf("Invalid instruction name: %s", item.name)}
        return
    }

    operandTypes := OperandTypeLookup[operandMode]
    length := LengthLookup[operandMode]

    if len(item.operands) != len(operandTypes) {
        errChan <- &AsmError{item.coord, fmt.Sprintf("Invalid number of operands (expected %d, got %d)", len(operandTypes), len(item.operands))}
        return
    }

    for i, o := range item.operands {
        t := operandTypes[i]

        if !o.SatisfiesType(t) {
            errChan <- &AsmError{item.coord, fmt.Sprintf("Invalid type for operand %d (0-indexed) to %s", i, item.name)}
            return
        }

        if t == DynamicType {
            length += o.Length()
        }
    }

    switch item.name {
    case "jmp":
        item.name = "mov"
        item.operands = []Operand{&PCOperand{coord: item.coord}, item.operands[0]}
        operandMode = OperandModeDD
        length = LengthLookup[operandMode] + item.operands[1].Length()

    case "ret":
        item.name = "pop"
        item.operands = []Operand{&PCOperand{coord: item.coord}}
        operandMode = OperandModeD
        length = LengthLookup[operandMode]

    case "jgt":
        item.name = "jlt"
        item.operands = []Operand{item.operands[1], item.operands[0], item.operands[2]}
        // Same mode & length

    case "jle":
        item.name = "jge"
        item.operands = []Operand{item.operands[1], item.operands[0], item.operands[2]}
        // Same mode & length
    }

    item.operandMode = operandMode
    item.length = length
}

func (item *Instruction) Encode(labelMap map[string]uint32, buffer []byte) {
    defer waitGroup.Done()

    operands := item.operands

    for _, operand := range operands {
        operand.ReduceLabel(labelMap)
    }

    var opcode byte

    switch item.name {
    case "nop":
        opcode = 0x00
    case "mov":
        opcode = 0x01
    }

    //buffer := make([]byte, item.length)
    buffer[0] = opcode

    switch item.operandMode {
    case OperandModeNone:

    case OperandModeD:
        buffer[1] = operands[0].EncodeKey()

        operands[0].EncodeExtra(buffer[2:])

    case OperandModeDD:
        buffer[1] = operands[0].EncodeKey()
        buffer[2] = operands[1].EncodeKey()

        operands[0].EncodeExtra(buffer[3:])
        l := operands[0].Length()
        operands[1].EncodeExtra(buffer[3+l:])

    case OperandModeDDD:
        buffer[1] = operands[0].EncodeKey()
        buffer[2] = operands[1].EncodeKey()
        buffer[3] = operands[2].EncodeKey()

        operands[0].EncodeExtra(buffer[4:])
        l := operands[0].Length()
        operands[1].EncodeExtra(buffer[4+l:])
        l += operands[1].Length()
        operands[2].EncodeExtra(buffer[4+l:])

    case OperandModeDDL:
        buffer[0] |= byte(operands[2].LiteralValue()) & 0x1F
        buffer[1] = operands[0].EncodeKey()
        buffer[2] = operands[1].EncodeKey()

        operands[0].EncodeExtra(buffer[3:])
        l := operands[0].Length()
        operands[1].EncodeExtra(buffer[3+l:])

    case OperandModeB:
        buffer[0] |= byte(operands[0].BitValue()) & 0x0F

    case OperandModeBD:
        buffer[1] = operands[1].EncodeKey()
        buffer[2] = byte(operands[0].BitValue()) & 0x0F

        operands[1].EncodeExtra(buffer[3:])

    case OperandModeBB:
        x := byte(operands[0].BitValue()) & 0x0F
        y := byte(operands[1].BitValue()) & 0x0F

        buffer[1] = (x << 4) | y

    case OperandModeBDL:
        j := byte(operands[2].LiteralValue()) & 0x1F
        x := byte(operands[0].BitValue()) & 0x0F

        buffer[0] |= j >> 4
        buffer[1] = operands[1].EncodeKey()
        buffer[2] = ((j & 0x0F) << 4) | x

        operands[1].EncodeExtra(buffer[3:])

    case OperandModeDLB:
        j := byte(operands[1].LiteralValue()) & 0x1F
        x := byte(operands[2].BitValue()) & 0x0F

        buffer[0] |= j >> 4
        buffer[1] = operands[0].EncodeKey()
        buffer[2] = ((j & 0x0F) << 4) | x

        operands[0].EncodeExtra(buffer[3:])
    }

    item.encoded = buffer
}

func (item *Instruction) Label() (label string, ok bool) {
    return "", false
}

func (item *Instruction) Length() (length uint32) {
    return item.length
}

func (item *Instruction) Offset() (offset uint32) {
    return item.offset
}

func (item *Instruction) SetOffset(offset uint32) {
    item.offset = offset
}

func (item *Instruction) Encoded() (encoded []byte) {
    return item.encoded
}

type Label struct {
    coord  Coord
    name   string
    offset uint32
}

func (item *Label) String() (str string) {
    return item.name + ":"
}

func (item *Label) VerifyAndReduce() {
    defer waitGroup.Done()
}

func (item *Label) Encode(labelMap map[string]uint32, buffer []byte) {
    defer waitGroup.Done()
}

func (item *Label) Label() (label string, ok bool) {
    return item.name, true
}

func (item *Label) Length() (length uint32) {
    return 0
}

func (item *Label) Offset() (offset uint32) {
    return item.offset
}

func (item *Label) SetOffset(offset uint32) {
    item.offset = offset
}

func (item *Label) Encoded() (encoded []byte) {
    return []byte{}
}
