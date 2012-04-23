package main

import (
    "encoding/json"
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
    currentlyEditingRegister     *Register
    currentlyEditingRegisterBank *RegisterBank
    currentlyEditingRam          *Ram
    currentConfirmationCallback  func()
)

var (
    mainMenuDialog         *termdialog.SelectionDialog
    cpuDialog              *termdialog.SelectionDialog
    structureDialog        *termdialog.SelectionDialog
    registersDialog        *termdialog.SelectionDialog
    editRegisterDialog     *termdialog.SelectionDialog
    registerBanksDialog    *termdialog.SelectionDialog
    editRegisterBankDialog *termdialog.SelectionDialog
    ramsDialog             *termdialog.SelectionDialog
    editRamDialog          *termdialog.SelectionDialog
    confirmationDialog     *termdialog.SelectionDialog

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

    decoder := json.NewDecoder(f)
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

    encoder := json.NewEncoder(f)
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

    registersDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addRegisterDialog})
    registersDialog.AddOption(&termdialog.Option{"Close", nil, nil})

    registerBanksDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addRegisterBankDialog})
    registerBanksDialog.AddOption(&termdialog.Option{"Close", nil, nil})

    ramsDialog.AddOption(&termdialog.Option{"Add", openDialogCallback, addRamDialog})
    ramsDialog.AddOption(&termdialog.Option{"Close", nil, nil})
}

func init() {
    mainMenuDialog = termdialog.NewSelectionDialog("Main Menu", nil)
    cpuDialog = termdialog.NewSelectionDialog("CPU", nil)
    structureDialog = termdialog.NewSelectionDialog("Structure", nil)
    registersDialog = termdialog.NewSelectionDialog("Registers", nil)
    editRegisterDialog = termdialog.NewSelectionDialog("Edit Register", nil)
    registerBanksDialog = termdialog.NewSelectionDialog("Register banks", nil)
    editRegisterBankDialog = termdialog.NewSelectionDialog("Edit Register Bank", nil)
    ramsDialog = termdialog.NewSelectionDialog("RAMs", nil)
    editRamDialog = termdialog.NewSelectionDialog("Edit RAM", nil)
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

    openDialog = termdialog.NewInputDialog("Open from "+SavesDir, "Filename:", 20, "", openCallback, nil)
    saveDialog = termdialog.NewInputDialog("Save to "+SavesDir, "Filename:", 20, "", saveCallback, nil)

    mainMenuDialog.AddOption(&termdialog.Option{"Edit CPU", openDialogCallback, cpuDialog})
    mainMenuDialog.AddOption(&termdialog.Option{"Open", openDialogCallback, openDialog})
    mainMenuDialog.AddOption(&termdialog.Option{"Save", openDialogCallback, saveDialog})
    mainMenuDialog.AddOption(&termdialog.Option{"Exit", nil, nil})

    cpuDialog.AddOption(&termdialog.Option{"Structure", openDialogCallback, structureDialog})
    cpuDialog.AddOption(&termdialog.Option{"Close", nil, nil})

    structureDialog.AddOption(&termdialog.Option{"Registers", openDialogCallback, registersDialog})
    structureDialog.AddOption(&termdialog.Option{"Register Banks", openDialogCallback, registerBanksDialog})
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
