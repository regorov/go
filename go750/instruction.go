package main

import (
	"fmt"
	"strings"
)

type Instruction struct {
	Name     string
	Operands []Operand
}

func NewInstruction(name string, operands ...Operand) (inst Instruction) {
	return Instruction{name, operands}
}

func (inst Instruction) String() (str string) {
	opStrs := make([]string, len(inst.Operands))

	for i, operand := range inst.Operands {
		opStrs[i] = operand.String()
	}

	return inst.Name + " " + strings.Join(opStrs, ", ")
}

type Operand interface {
	String() string
}

type Integer interface {
	String() string
}

type Register uint8

func (operand Register) String() (str string) {
	return "%" + RegisterNames[operand]
}

type Literal int64

func (operand Literal) String() (str string) {
	if operand < 0 {
		return fmt.Sprintf("-0x%X", -operand)
	}

	return fmt.Sprintf("0x%X", operand)
}

type Label string

func (operand Label) String() (str string) {
	return string(operand)
}

type MemRef struct {
	Size    Size
	Base    Register
	Index   Register
	Scale   Scale
	Disp    Integer
	IsArray bool
}

func (operand MemRef) String() (str string) {
	if operand.IsArray {
		return fmt.Sprintf("%s[%s, %s, %s, %s]", operand.Size, operand.Base, operand.Index, operand.Scale, operand.Disp)
	}

	if operand.Disp != nil {
		return fmt.Sprintf("%s[%s + %s]", operand.Size, operand.Base, operand.Disp)
	}

	return fmt.Sprintf("%s[%s]", operand.Size, operand.Base)
}

type LitMemRef struct {
	Addr Integer
}

func (operand LitMemRef) String() (str string) {
	return fmt.Sprintf("32[%s]", operand.Addr)
}

type BitReg uint8

func (operand BitReg) String() (str string) {
	return fmt.Sprintf("%%b%d", operand)
}

type otherOperand string

func (operand otherOperand) String() (str string) {
	return string(operand)
}

var PC = otherOperand("PC")
var SR = otherOperand("SR")

type Size uint8

func (size Size) String() (str string) {
	switch size {
	case Byte:
		return "8"
	case Half:
		return "16"
	case Word:
		return "32"
	}

	return ""
}

type Scale uint8

func (scale Scale) String() (str string) {
	switch scale {
	case Scale2:
		return "2"
	case Scale4:
		return "4"
	case Scale8:
		return "8"
	case Scale16:
		return "16"
	}

	return ""
}

const (
	V0 Register = iota
	V1
	V2
	V3
	V4
	V5
	V6
	V7
	A0
	A1
	A2
	A3
	Q0
	Q1
	SP
	KT
)

var RegisterNames = []string{
	"v0",
	"v1",
	"v2",
	"v3",
	"v4",
	"v5",
	"v6",
	"v7",
	"a0",
	"a1",
	"a2",
	"a3",
	"q0",
	"q1",
	"sp",
	"kt",
}

const (
	Byte Size = iota
	Half
	Word
)

const (
	Scale2 Scale = iota
	Scale4
	Scale8
	Scale16
)
