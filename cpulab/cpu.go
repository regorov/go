package main

import (
    "github.com/kierdavis/go/termdialog"
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
    EditOption *termdialog.Option `json:"-"`
}

type RegisterBank struct {
    Name       string
    Width      uint
    Depth      uint
    EditOption *termdialog.Option `json:"-"`
}

type Ram struct {
    Name       string
    Width      uint
    Depth      uint
    EditOption *termdialog.Option `json:"-"`
}

type Instruction struct {
    //Microinstructions []*Microinstruction
}
