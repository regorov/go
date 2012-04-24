package main

import (
	"fmt"
	"github.com/kierdavis/go/termdialog"
	"strconv"
)

func addRegisterBankWidthCallback(name string, arg interface{}) (close bool) {
	addRegisterBankDialog.SetValue("")
	addRegisterBankWidthDialog.SetCallbackArg(name)
	dialogStack.Open(addRegisterBankWidthDialog)
	return true
}

func addRegisterBankDepthCallback(widthStr string, arg interface{}) (close bool) {
	name := arg.(string)
	width64, err := strconv.ParseUint(widthStr, 10, 0)
	if err != nil {
		panic(err)
	}

	width := uint(width64)
	addRegisterBankDepthDialog.SetCallbackArg(nameWidth{name, width})
	dialogStack.Open(addRegisterBankDepthDialog)
	return true
}

func addRegisterBankCallback(depthStr string, arg interface{}) (close bool) {
	nw := arg.(nameWidth)
	name := nw.name
	width := nw.width

	depth64, err := strconv.ParseUint(depthStr, 10, 0)
	if err != nil {
		panic(err)
	}

	depth := uint(depth64)

	registerBank := &RegisterBank{
		Name:  name,
		Width: width,
		Depth: depth,
	}

	addRegisterBankToGui(registerBank)
	cpu.RegisterBanks = append(cpu.RegisterBanks, registerBank)

	return true
}

func addRegisterBankToGui(registerBank *RegisterBank) {
	registerBank.editOption = registerBanksDialog.AddOption(&termdialog.Option{"Edit " + registerBank.Name, editRegisterBankOpenCallback, registerBank})

	asMicroRegisterBankIndexed := &MicroRegisterBankIndexed{
		Register: register,
		Index:    nil,
	}

	asMicroSource := IMicroSource(asMicroRegisterBankIndexed)
	addMoveMicroInstructionDialog.AddOption(&termdialog.Option{"Use " + asMicroRegisterBankIndexed.String() + " as source operand", addMoveMicroInstructionDestinationCallback, asMicroSource})

	asMicroDestination := IMicroDestination(asMicroRegisterBankIndexed)
	addMoveMicroInstructionDestinationDialog.AddOption(&termdialog.Option{"Use " + asMicroRegisterBankIndexed.String() + " as destination operand", addMoveMicroInstructionCallback, asMicroDestination})
}

func updateEditRegisterBankDialogs() {
	editRegisterBankDialog.SetTitle("Edit Register Bank " + currentlyEditingRegisterBank.Name)
	editRegisterBankDialog.GetOption(0).Text = fmt.Sprintf("Edit name (%s)", currentlyEditingRegisterBank.Name)
	editRegisterBankDialog.GetOption(1).Text = fmt.Sprintf("Edit width (%d)", currentlyEditingRegisterBank.Width)
	editRegisterBankDialog.GetOption(2).Text = fmt.Sprintf("Edit depth (%d)", currentlyEditingRegisterBank.Depth)

	editRegisterBankNameDialog.SetTitle("Edit Register Bank " + currentlyEditingRegisterBank.Name)
	editRegisterBankNameDialog.SetValue(currentlyEditingRegisterBank.Name)

	editRegisterBankWidthDialog.SetTitle("Edit Register Bank " + currentlyEditingRegisterBank.Name)
	editRegisterBankWidthDialog.SetValue(strconv.FormatUint(uint64(currentlyEditingRegisterBank.Width), 10))

	editRegisterBankDepthDialog.SetTitle("Edit Register Bank " + currentlyEditingRegisterBank.Name)
	editRegisterBankDepthDialog.SetValue(strconv.FormatUint(uint64(currentlyEditingRegisterBank.Depth), 10))
}

func editRegisterBankOpenCallback(option *termdialog.Option) (close bool) {
	currentlyEditingRegisterBank = option.Data.(*RegisterBank)
	updateEditRegisterBankDialogs()
	dialogStack.Open(editRegisterBankDialog)
	return false
}

func editRegisterBankNameCallback(name string, arg interface{}) (close bool) {
	currentlyEditingRegisterBank.Name = name
	updateEditRegisterBankDialogs()
	currentlyEditingRegisterBank.editOption.Text = "Edit " + name
	return true
}

func editRegisterBankWidthCallback(widthStr string, arg interface{}) (close bool) {
	width64, err := strconv.ParseUint(widthStr, 10, 0)
	if err != nil {
		panic(err)
	}

	currentlyEditingRegisterBank.Width = uint(width64)
	updateEditRegisterBankDialogs()
	return true
}

func editRegisterBankDepthCallback(depthStr string, arg interface{}) (close bool) {
	depth64, err := strconv.ParseUint(depthStr, 10, 0)
	if err != nil {
		panic(err)
	}

	currentlyEditingRegisterBank.Depth = uint(depth64)
	updateEditRegisterBankDialogs()
	return true
}

func deleteRegisterBankCallback(option *termdialog.Option) (close bool) {
	currentConfirmationCallback = deleteRegisterBank
	dialogStack.Open(confirmationDialog)
	return false
}

func deleteRegisterBank() {
	dialogStack.Close(editRegisterBankDialog)

	for i, registerBank := range cpu.RegisterBanks {
		if registerBank == currentlyEditingRegisterBank {
			cpu.RegisterBanks = append(cpu.RegisterBanks[:i], cpu.RegisterBanks[i+1:]...)
			registerBanksDialog.RemoveOption(i + 2)
		}
	}
}
