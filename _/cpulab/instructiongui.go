package main

import (
	"fmt"
	"github.com/kierdavis/go/termdialog"
)

func addInstructionCallback(name string, arg interface{}) (close bool) {
	addInstructionDialog.SetValue("")

	instruction := &Instruction{
		Name:              name,
		MicroInstructions: make([]IMicroInstruction, 0),
	}

	addInstructionToGui(instruction)
	cpu.Instructions = append(cpu.Instructions, instruction)

	return true
}

func addInstructionToGui(instruction *Instruction) {
	instruction.editOption = instructionsDialog.AddOption(&termdialog.Option{"Edit " + instruction.Name, editInstructionOpenCallback, instruction})
	instruction.microInstructionsDialog = termdialog.NewSelectionDialog("Microinstructions", nil)
	instruction.microInstructionsDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addMicroInstructionDialog})
	instruction.microInstructionsDialog.AddOption(&termdialog.Option{"Close", nil, nil})
}

func updateEditInstructionDialogs() {
	editInstructionDialog.SetTitle("Edit Instruction " + currentlyEditingInstruction.Name)
	editInstructionDialog.GetOption(0).Text = fmt.Sprintf("Edit mnemonic (%s)", currentlyEditingInstruction.Name)
	editInstructionDialog.GetOption(1).Data = currentlyEditingInstruction.microInstructionsDialog

	editInstructionNameDialog.SetTitle("Edit Instruction " + currentlyEditingInstruction.Name)
	editInstructionNameDialog.SetValue(currentlyEditingInstruction.Name)
}

func editInstructionOpenCallback(option *termdialog.Option) (close bool) {
	currentlyEditingInstruction = option.Data.(*Instruction)
	updateEditInstructionDialogs()
	dialogStack.Open(editInstructionDialog)
	return false
}

func editInstructionNameCallback(name string, args interface{}) (close bool) {
	currentlyEditingInstruction.Name = name
	updateEditInstructionDialogs()
	currentlyEditingInstruction.editOption.Text = "Edit " + name
	return true
}

func deleteInstructionCallback(option *termdialog.Option) (close bool) {
	currentConfirmationCallback = deleteInstruction
	dialogStack.Open(confirmationDialog)
	return false
}

func deleteInstruction() {
	dialogStack.Close(editInstructionDialog)

	for i, instruction := range cpu.Instructions {
		if instruction == currentlyEditingInstruction {
			cpu.Instructions = append(cpu.Instructions[:i], cpu.Instructions[i+1:]...)
			instructionsDialog.RemoveOption(i + 2)
		}
	}
}
