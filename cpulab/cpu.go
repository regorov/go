package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/kierdavis/go/termdialog"
	"io"
)

type RamOperation uint8

const (
	RamRead RamOperation = iota
	RamWrite
)

type CPU struct {
	Registers     []*Register
	RegisterBanks []*RegisterBank
	Rams          []*Ram

	Instructions []*Instruction
}

type Register struct {
	Name       string
	Width      uint
	editOption *termdialog.Option `json:"-" xml:"-"`
}

type RegisterBank struct {
	Name       string
	Width      uint
	Depth      uint
	editOption *termdialog.Option `json:"-" xml:"-"`
}

type Ram struct {
	Name       string
	Width      uint
	Depth      uint
	editOption *termdialog.Option `json:"-" xml:"-"`
}

type Instruction struct {
	Name                    string
	MicroInstructions       []IMicroInstruction
	editOption              *termdialog.Option          `json:"-" xml:"-"`
	microInstructionsDialog *termdialog.SelectionDialog `json:"-" xml:"-"`
}

type IMicroInstruction interface {
	ToVerilog(io.Writer, string)
	String() string
}

// Implements IMicroInstruction
type MoveMicroInstruction struct {
	Destination IMicroDestination
	Source      IMicroSource
}

// Implements IMicroInstruction
type RamOperationMicroInstruction struct {
	Ram       *Ram
	Operation RamOperation
}

type IMicroDestination interface {
	ToVerilogExpr() string
	String() string
}

type IMicroSource interface {
	IMicroSimpleSource
}

type IMicroSimpleSource interface {
	ToVerilogExpr() string
	String() string
}

// Implements IMicroDestination, IMicroSource and IMicroSimpleSource
type MicroRegister struct {
	Register *Register
}

// Implements IMicroDestination and IMicroSource
type MicroRegisterBankIndexed struct {
	RegisterBank *RegisterBank
	Index        IMicroSimpleSource
}

// Implements IMicroDestination
type MicroRamAddr struct {
	Ram *Ram
}

// Implements IMicroDestination
type MicroRamInput struct {
	Ram *Ram
}

// Implements IMicroSource and IMicroSimpleSource
type MicroRamOutput struct {
	Ram *Ram
}

func (self *MoveMicroInstruction) ToVerilog(writer io.Writer, ind string) {
	fmt.Fprintf(writer, "%s%s <= %s;\n", ind, self.Destination.ToVerilogExpr(), self.Source.ToVerilogExpr())
}

func (self *MoveMicroInstruction) String() (str string) {
	return fmt.Sprintf("{%s <- %s}", self.Destination.String(), self.Source.String())
}

func (self *RamOperationMicroInstruction) ToVerilog(writer io.Writer, ind string) {
	if self.Operation == RamRead {
		fmt.Fprintf(writer, "%sram_read <= 1'd1;\n")
	} else if self.Operation == RamWrite {
		fmt.Fprintf(writer, "%sram_write <= 1'd1;\n")
	} else {
		showError(errors.New("Invalid value for Operation field"))
	}
}

func (self *RamOperationMicroInstruction) String() (str string) {
	if self.Operation == RamRead {
		str = fmt.Sprintf("{Read %s}", self.Ram.Name)
	} else if self.Operation == RamWrite {
		str = fmt.Sprintf("{Write %s}", self.Ram.Name)
	} else {
		showError(errors.New("Invalid value for Operation field"))
	}

	return str
}

func (self *MicroRegister) ToVerilogExpr() (str string) {
	return fmt.Sprintf("reg_%s", self.Register.Name)
}

func (self *MicroRegister) String() (str string) {
	return self.Register.Name
}

func (self *MicroRegisterBankIndexed) ToVerilogExpr() (str string) {
	return fmt.Sprintf("regbank_%s[%s]", self.RegisterBank.Name, self.Index.ToVerilogExpr())
}

func (self *MicroRegisterBankIndexed) String() (str string) {
	if self.Index == nil {
		str = fmt.Sprintf("%s[...]", self.RegisterBank.Name)
	} else {
		str = fmt.Sprintf("%s[%s]", self.RegisterBank.Name, self.Index.String())
	}
	return str
}

func (self *MicroRamAddr) ToVerilogExpr() (str string) {
	return fmt.Sprintf("ram_%s_addr", self.Ram.Name)
}

func (self *MicroRamAddr) String() (str string) {
	return fmt.Sprintf("%s.Address", self.Ram.Name)
}

func (self *MicroRamInput) ToVerilogExpr() (str string) {
	return fmt.Sprintf("ram_%s_data", self.Ram.Name)
}

func (self *MicroRamInput) String() (str string) {
	return fmt.Sprintf("%s.Input", self.Ram.Name)
}

func (self *MicroRamOutput) ToVerilogExpr() (str string) {
	return fmt.Sprintf("ram_%s_q", self.Ram.Name)
}

func (self *MicroRamOutput) String() (str string) {
	return fmt.Sprintf("%s.Output", self.Ram.Name)
}

/*
func (w *IMicroInstructionWrapper) UnmarshalJSON(data []byte) (err error) {
	var d map[string]json.RawMessage
	err = json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	var t InterfaceType
	err = json.Unmarshal(d["Type"], &t)
	if err != nil {
		return err
	}
	w.Type = t

	switch t {
	case MoveMicroInstructionType:
		var v *MoveMicroInstruction
		err = json.Unmarshal(d["V"], &v)
		w.V = v

	case RamOperationMicroInstructionType:
		var v *RamOperationMicroInstruction
		err = json.Unmarshal(d["V"], &v)
		w.V = v
	}

	return err
}
*/

func init() {
	// Implementers of IMicroInstruction
	gob.Register(&MoveMicroInstruction{})
	gob.Register(&RamOperationMicroInstruction{})

	// Implementers of IMicroDestination/IMicroSource
	gob.Register(&MicroRegister{})
	gob.Register(&MicroRegisterBankIndexed{})
	gob.Register(&MicroRamAddr{})
	gob.Register(&MicroRamInput{})
	gob.Register(&MicroRamOutput{})
}
