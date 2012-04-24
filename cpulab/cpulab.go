package main

import (
	"encoding/gob"
	"fmt"
	"github.com/kierdavis/go/termdialog"
	"github.com/nsf/termbox-go"
	"os"
	"path/filepath"
)

var SavesDir = os.Getenv("HOME") + "/cpulab/saves/"

var cpu CPU

func refreshCPU() {
	cpu.Registers = make([]*Register, 0)
	cpu.RegisterBanks = make([]*RegisterBank, 0)
	cpu.Rams = make([]*Ram, 0)
	cpu.Instructions = make([]*Instruction, 0)
}

// UI structure

var dialogStack = termdialog.NewDialogStack()

var (
	registerEditOptions     []*termdialog.Option
	registerBankEditOptions []*termdialog.Option
	ramEditOptions          []*termdialog.Option
	instructionEditOptions  []*termdialog.Option
)

var (
	currentlyEditingRegister          *Register
	currentlyEditingRegisterBank      *RegisterBank
	currentlyEditingRam               *Ram
	currentlyEditingInstruction       *Instruction
	currentlyEditingMicroInstruction  IMicroInstruction
	currentMoveMicroInstructionSource IMicroSource
	currentConfirmationCallback       func()
)

var (
	mainMenuDialog                           *termdialog.SelectionDialog
	cpuDialog                                *termdialog.SelectionDialog
	structureDialog                          *termdialog.SelectionDialog
	registersDialog                          *termdialog.SelectionDialog
	editRegisterDialog                       *termdialog.SelectionDialog
	registerBanksDialog                      *termdialog.SelectionDialog
	editRegisterBankDialog                   *termdialog.SelectionDialog
	ramsDialog                               *termdialog.SelectionDialog
	editRamDialog                            *termdialog.SelectionDialog
	instructionsDialog                       *termdialog.SelectionDialog
	editInstructionDialog                    *termdialog.SelectionDialog
	addMicroInstructionDialog                *termdialog.SelectionDialog
	addMoveMicroInstructionDialog            *termdialog.SelectionDialog
	addMoveMicroInstructionDestinationDialog *termdialog.SelectionDialog
	addRamReadMicroInstructionDialog         *termdialog.SelectionDialog
	addRamWriteMicroInstructionDialog        *termdialog.SelectionDialog
	editMicroInstructionDialog               *termdialog.SelectionDialog
	confirmationDialog                       *termdialog.SelectionDialog

	addRegisterDialog           *termdialog.InputDialog
	addRegisterWidthDialog      *termdialog.InputDialog
	editRegisterNameDialog      *termdialog.InputDialog
	editRegisterWidthDialog     *termdialog.InputDialog
	addRegisterBankDialog       *termdialog.InputDialog
	addRegisterBankWidthDialog  *termdialog.InputDialog
	addRegisterBankDepthDialog  *termdialog.InputDialog
	editRegisterBankNameDialog  *termdialog.InputDialog
	editRegisterBankWidthDialog *termdialog.InputDialog
	editRegisterBankDepthDialog *termdialog.InputDialog
	addRamDialog                *termdialog.InputDialog
	addRamWidthDialog           *termdialog.InputDialog
	addRamDepthDialog           *termdialog.InputDialog
	editRamNameDialog           *termdialog.InputDialog
	editRamWidthDialog          *termdialog.InputDialog
	editRamDepthDialog          *termdialog.InputDialog
	addInstructionDialog        *termdialog.InputDialog
	editInstructionNameDialog   *termdialog.InputDialog

	// Future FileChooserDialogs
	openDialog *termdialog.InputDialog
	saveDialog *termdialog.InputDialog
)

type openDialogArg struct {
	dialog termdialog.Dialog
	arg    interface{}
}

type nameWidth struct {
	name  string
	width uint
}

func confirmCallback(option *termdialog.Option) (close bool) {
	confirmationDialog.SetSelectedIndex(0)
	currentConfirmationCallback()
	return true
}

func openDialogCallback(option *termdialog.Option) (close bool) {
	dialogStack.Open(option.Data.(termdialog.Dialog))
	return false
}

func openCallback(filename string, arg interface{}) (close bool) {
	filename = filepath.Join(SavesDir, filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	refreshCPU()

	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&cpu)
	if err != nil {
		panic(err)
	}

	refreshGUI()

	for _, register := range cpu.Registers {
		addRegisterToGui(register)
	}

	for _, registerBank := range cpu.RegisterBanks {
		addRegisterBankToGui(registerBank)
	}

	for _, ram := range cpu.Rams {
		addRamToGui(ram)
	}

	for _, instruction := range cpu.Instructions {
		addInstructionToGui(instruction)
	}

	return true
}

func saveCallback(filename string, arg interface{}) (close bool) {
	filename = filepath.Join(SavesDir, filename)
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, os.ModeDir|0755)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(cpu)
	if err != nil {
		panic(err)
	}

	return true
}

// If a add*ToGui function modifies a dialog's options, move its initialisation into here so they
// can be reset.

func refreshGUI() {
	registersDialog.ClearOptions()
	registerBanksDialog.ClearOptions()
	ramsDialog.ClearOptions()
	instructionsDialog.ClearOptions()
	addMoveMicroInstructionDialog.ClearOptions()
	addMoveMicroInstructionDestinationDialog.ClearOptions()

	registersDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addRegisterDialog})
	registersDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	registerBanksDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addRegisterBankDialog})
	registerBanksDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	ramsDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addRamDialog})
	ramsDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	instructionsDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addInstructionDialog})
	instructionsDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	addMoveMicroInstructionDialog.AddOption(&termdialog.Option{"Edit registers", openDialogCallback, registersDialog})
	addMoveMicroInstructionDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	addMoveMicroInstructionDestinationDialog.AddOption(&termdialog.Option{"Edit registers", openDialogCallback, registersDialog})
	addMoveMicroInstructionDestinationDialog.AddOption(&termdialog.Option{"Close", nil, nil})
}

func init() {
	mainMenuDialog = termdialog.NewSelectionDialog("Main Menu", nil)
	cpuDialog = termdialog.NewSelectionDialog("CPU", nil)
	structureDialog = termdialog.NewSelectionDialog("Structure", nil)
	registersDialog = termdialog.NewSelectionDialog("Registers", nil)
	editRegisterDialog = termdialog.NewSelectionDialog("Edit Register", nil)
	registerBanksDialog = termdialog.NewSelectionDialog("Register Banks", nil)
	editRegisterBankDialog = termdialog.NewSelectionDialog("Edit Register Bank", nil)
	ramsDialog = termdialog.NewSelectionDialog("RAMs", nil)
	editRamDialog = termdialog.NewSelectionDialog("Edit RAM", nil)
	instructionsDialog = termdialog.NewSelectionDialog("Instructions", nil)
	editInstructionDialog = termdialog.NewSelectionDialog("Edit Instruction", nil)
	addMicroInstructionDialog = termdialog.NewSelectionDialog("Add Microinstruction", nil)
	addMoveMicroInstructionDialog = termdialog.NewSelectionDialog("Add Register Transfer Microinstruction", nil)
	addMoveMicroInstructionDestinationDialog = termdialog.NewSelectionDialog("Add Register Transfer Microinstruction", nil)
	addRamReadMicroInstructionDialog = termdialog.NewSelectionDialog("Add RAM Read Microinstruction", nil)
	addRamWriteMicroInstructionDialog = termdialog.NewSelectionDialog("Add RAM Write Microinstruction", nil)
	editMicroInstructionDialog = termdialog.NewSelectionDialog("Edit Microinstruction", nil)
	confirmationDialog = termdialog.NewSelectionDialog("Are you sure?", nil)

	addRegisterDialog = termdialog.NewInputDialog("Add Register", "Name:", 12, "", addRegisterWidthCallback, nil)
	addRegisterWidthDialog = termdialog.NewInputDialog("Add Register", "Width:", 4, "", addRegisterCallback, nil)
	editRegisterNameDialog = termdialog.NewInputDialog("Edit Register", "Name:", 12, "", editRegisterNameCallback, nil)
	editRegisterWidthDialog = termdialog.NewInputDialog("Edit Register", "Width:", 4, "", editRegisterWidthCallback, nil)
	addRegisterBankDialog = termdialog.NewInputDialog("Add Register Bank", "Name:", 12, "", addRegisterBankWidthCallback, nil)
	addRegisterBankWidthDialog = termdialog.NewInputDialog("Add Register Bank", "Width:", 4, "", addRegisterBankDepthCallback, nil)
	addRegisterBankDepthDialog = termdialog.NewInputDialog("Add Register Bank", "Depth:", 8, "", addRegisterBankCallback, nil)
	editRegisterBankNameDialog = termdialog.NewInputDialog("Edit Register Bank", "Name:", 12, "", editRegisterBankNameCallback, nil)
	editRegisterBankWidthDialog = termdialog.NewInputDialog("Edit Register Bank", "Width:", 4, "", editRegisterBankWidthCallback, nil)
	editRegisterBankDepthDialog = termdialog.NewInputDialog("Edit Register Bank", "Depth:", 8, "", editRegisterBankDepthCallback, nil)
	addRamDialog = termdialog.NewInputDialog("Add RAM", "Name:", 12, "", addRamWidthCallback, nil)
	addRamWidthDialog = termdialog.NewInputDialog("Add RAM", "Width:", 4, "", addRamDepthCallback, nil)
	addRamDepthDialog = termdialog.NewInputDialog("Add RAM", "Depth:", 8, "", addRamCallback, nil)
	editRamNameDialog = termdialog.NewInputDialog("Edit RAM", "Name:", 12, "", editRamNameCallback, nil)
	editRamWidthDialog = termdialog.NewInputDialog("Edit RAM", "Width:", 4, "", editRamWidthCallback, nil)
	editRamDepthDialog = termdialog.NewInputDialog("Edit RAM", "Depth:", 8, "", editRamDepthCallback, nil)
	addInstructionDialog = termdialog.NewInputDialog("Add Instruction", "Mnemonic:", 12, "", addInstructionCallback, nil)
	editInstructionNameDialog = termdialog.NewInputDialog("Edit Instruction", "Mnemonic:", 12, "", editInstructionNameCallback, nil)

	openDialog = termdialog.NewInputDialog("Open from "+SavesDir, "Filename:", 20, "", openCallback, nil)
	saveDialog = termdialog.NewInputDialog("Save to "+SavesDir, "Filename:", 20, "", saveCallback, nil)

	mainMenuDialog.AddOption(&termdialog.Option{"Edit CPU", openDialogCallback, cpuDialog})
	mainMenuDialog.AddOption(&termdialog.Option{"Open", openDialogCallback, openDialog})
	mainMenuDialog.AddOption(&termdialog.Option{"Save", openDialogCallback, saveDialog})
	mainMenuDialog.AddOption(&termdialog.Option{"Exit", nil, nil})

	cpuDialog.AddOption(&termdialog.Option{"Structure", openDialogCallback, structureDialog})
	cpuDialog.AddOption(&termdialog.Option{"Instructions", openDialogCallback, instructionsDialog})
	cpuDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	structureDialog.AddOption(&termdialog.Option{"Registers", openDialogCallback, registersDialog})
	structureDialog.AddOption(&termdialog.Option{"Register banks", openDialogCallback, registerBanksDialog})
	structureDialog.AddOption(&termdialog.Option{"RAMs", openDialogCallback, ramsDialog})
	structureDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	editRegisterDialog.AddOption(&termdialog.Option{"Edit name", openDialogCallback, editRegisterNameDialog})
	editRegisterDialog.AddOption(&termdialog.Option{"Edit width", openDialogCallback, editRegisterWidthDialog})
	editRegisterDialog.AddOption(&termdialog.Option{"Delete", deleteRegisterCallback, nil})
	editRegisterDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	editRegisterBankDialog.AddOption(&termdialog.Option{"Edit name", openDialogCallback, editRegisterBankNameDialog})
	editRegisterBankDialog.AddOption(&termdialog.Option{"Edit width", openDialogCallback, editRegisterBankWidthDialog})
	editRegisterBankDialog.AddOption(&termdialog.Option{"Edit depth", openDialogCallback, editRegisterBankDepthDialog})
	editRegisterBankDialog.AddOption(&termdialog.Option{"Delete", deleteRegisterBankCallback, nil})
	editRegisterBankDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	editRamDialog.AddOption(&termdialog.Option{"Edit name", openDialogCallback, editRamNameDialog})
	editRamDialog.AddOption(&termdialog.Option{"Edit width", openDialogCallback, editRamWidthDialog})
	editRamDialog.AddOption(&termdialog.Option{"Edit depth", openDialogCallback, editRamDepthDialog})
	editRamDialog.AddOption(&termdialog.Option{"Delete", deleteRamCallback, nil})
	editRamDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	editInstructionDialog.AddOption(&termdialog.Option{"Edit mnemonic", openDialogCallback, editInstructionNameDialog})
	editInstructionDialog.AddOption(&termdialog.Option{"Edit microinstructions", openDialogCallback, nil})
	editInstructionDialog.AddOption(&termdialog.Option{"Delete", deleteInstructionCallback, nil})
	editInstructionDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	addMicroInstructionDialog.AddOption(&termdialog.Option{"Register transfer", addMoveMicroInstructionOpenCallback, nil})
	addMicroInstructionDialog.AddOption(&termdialog.Option{"RAM read", addRamReadMicroInstructionOpenCallback, nil})
	addMicroInstructionDialog.AddOption(&termdialog.Option{"RAM write", addRamWriteMicroInstructionOpenCallback, nil})
	addMicroInstructionDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	editMicroInstructionDialog.AddOption(&termdialog.Option{"Delete", deleteMicroInstructionCallback, nil})
	editMicroInstructionDialog.AddOption(&termdialog.Option{"Close", nil, nil})

	confirmationDialog.AddOption(&termdialog.Option{"No", nil, nil})
	confirmationDialog.AddOption(&termdialog.Option{"Yes", confirmCallback, nil})

	refreshGUI()
}

func die(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func main() {
	refreshCPU()

	err := termbox.Init()
	die(err)
	defer termbox.Close()

	dialogStack.Open(mainMenuDialog)

	dialogStack.Run()
}
