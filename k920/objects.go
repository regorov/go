package k920

import (
	"fmt"
	"strconv"
	"strings"
)

type Object interface {
	String() string
	Length() uint32
	Offset() uint32
	SetOffset(uint32)
}

type Instruction struct {
	Group    string
	Name     string
	Operands []Operand
	Offset_  uint32
}

func (inst *Instruction) String() (str string) {
	opstrs := make([]string, len(inst.Operands))

	for i, operand := range inst.Operands {
		opstrs[i] = operand.String()
	}

	return fmt.Sprintf("  %s.%s %s", inst.Group, inst.Name, strings.Join(opstrs, ", "))
}

func (inst *Instruction) Length() (length uint32) {
	return 4
}

func (inst *Instruction) Offset() (offset uint32) {
	return inst.Offset_
}

func (inst *Instruction) SetOffset(offset uint32) {
	inst.Offset_ = offset
}

type Label struct {
	Name    string
	Offset_ uint32
}

func (label *Label) String() (str string) {
	return fmt.Sprintf("%s:", label.Name)
}

func (label *Label) Length() (length uint32) {
	return 0
}

func (label *Label) Offset() (offset uint32) {
	return label.Offset_
}

func (label *Label) SetOffset(offset uint32) {
	label.Offset_ = offset
}

type Operand interface {
	String() string
}

type Register uint8

func (operand Register) String() (str string) {
	return "%" + RegisterNames[operand]
}

type BitReg uint8

func (operand BitReg) String() (str string) {
	return "%" + BitRegNames[operand]
}

type Integer int64

func (operand Integer) String() (str string) {
	return strconv.FormatInt(int64(operand), 10)
}

type LabelRef string

func (operand LabelRef) String() (str string) {
	return string(operand)
}
