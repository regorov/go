package main

import (
	//    "fmt"
	"github.com/kierdavis/go/termdialog"
)

func addMoveMicroInstructionOpenCallback(option *termdialog.Option) (close bool) {
	dialogStack.Open(addMoveMicroInstructionDialog)
	return true
}

func addRamReadMicroInstructionOpenCallback(option *termdialog.Option) (close bool) {
	dialogStack.Open(addRamReadMicroInstructionDialog)
	return true
}

func addRamWriteMicroInstructionOpenCallback(option *termdialog.Option) (close bool) {
	dialogStack.Open(addRamWriteMicroInstructionDialog)
	return true
}

func addMoveMicroInstructionDestinationCallback(option *termdialog.Option) (close bool) {
	currentMoveMicroInstructionSource = option.Data.(IMicroSource)
	dialogStack.Open(addMoveMicroInstructionDestinationDialog)
	return true
}

func addMoveMicroInstructionCallback(option *termdialog.Option) (close bool) {
	mmi := &MoveMicroInstruction{
		Destination: option.Data.(IMicroDestination),
		Source:      currentMoveMicroInstructionSource,
	}

	addMicroInstructionToGui(mmi)
	currentlyEditingInstruction.MicroInstructions = append(currentlyEditingInstruction.MicroInstructions, IMicroInstruction(mmi))
	return true
}

func addMicroInstructionToGui(mi IMicroInstruction) {
	currentlyEditingInstruction.microInstructionsDialog.AddOption(&termdialog.Option{"Edit " + mi.String(), editMicroInstructionOpenCallback, mi})
}

func updateEditMicroInstructionDialogs() {
	editMicroInstructionDialog.SetTitle("Edit " + currentlyEditingMicroInstruction.String())
}

func editMicroInstructionOpenCallback(option *termdialog.Option) (close bool) {
	currentlyEditingMicroInstruction = option.Data.(IMicroInstruction)
	updateEditMicroInstructionDialogs()
	dialogStack.Open(editMicroInstructionDialog)
	return false
}

func deleteMicroInstructionCallback(option *termdialog.Option) (close bool) {
	currentConfirmationCallback = deleteMicroInstruction
	dialogStack.Open(confirmationDialog)
	return false
}

func deleteMicroInstruction() {
	dialogStack.Close(editMicroInstructionDialog)

	for i, mi := range currentlyEditingInstruction.MicroInstructions {
		if mi == currentlyEditingMicroInstruction {
			currentlyEditingInstruction.MicroInstructions = append(currentlyEditingInstruction.MicroInstructions[:i], currentlyEditingInstruction.MicroInstructions[i+1:]...)
			currentlyEditingInstruction.microInstructionsDialog.RemoveOption(i + 2)
		}
	}
}
