package main

import (
    "fmt"
    "github.com/kierdavis/go/termdialog"
    "strconv"
)

func addRegisterWidthCallback(name string, arg interface{}) (close bool) {
    addRegisterDialog.SetValue("")
    addRegisterWidthDialog.SetCallbackArg(name)
    dialogStack.Open(addRegisterWidthDialog)
    return true
}

func addRegisterCallback(widthStr string, arg interface{}) (close bool) {
    name := arg.(string)
    width64, err := strconv.ParseUint(widthStr, 10, 0)
    if err != nil {
        panic(err)
    }

    width := uint(width64)

    //termdialog.Debug(0, "%s (%d)", name, width)
    register := &Register{
        Name:  name,
        Width: width,
    }

    addRegisterToGui(register)
    cpu.Registers = append(cpu.Registers, register)

    return true
}

func addRegisterToGui(register *Register) {
    register.EditOption = registersDialog.AddOption(&termdialog.Option{"Edit " + register.Name, editRegisterOpenCallback, register})
}

func updateEditRegisterDialogs() {
    editRegisterDialog.SetTitle("Edit Register " + currentlyEditingRegister.Name)
    editRegisterDialog.GetOption(0).Text = fmt.Sprintf("Edit name (%s)", currentlyEditingRegister.Name)
    editRegisterDialog.GetOption(1).Text = fmt.Sprintf("Edit width (%d)", currentlyEditingRegister.Width)

    editRegisterNameDialog.SetTitle("Edit Register " + currentlyEditingRegister.Name)
    editRegisterNameDialog.SetValue(currentlyEditingRegister.Name)

    editRegisterWidthDialog.SetTitle("Edit Register " + currentlyEditingRegister.Name)
    editRegisterWidthDialog.SetValue(strconv.FormatUint(uint64(currentlyEditingRegister.Width), 10))
}

func editRegisterOpenCallback(option *termdialog.Option) (close bool) {
    currentlyEditingRegister = option.Data.(*Register)
    updateEditRegisterDialogs()
    dialogStack.Open(editRegisterDialog)
    return false
}

func editRegisterNameCallback(name string, arg interface{}) (close bool) {
    currentlyEditingRegister.Name = name
    updateEditRegisterDialogs()
    currentlyEditingRegister.EditOption.Text = "Edit " + name
    return true
}

func editRegisterWidthCallback(widthStr string, arg interface{}) (close bool) {
    width64, err := strconv.ParseUint(widthStr, 10, 0)
    if err != nil {
        panic(err)
    }

    currentlyEditingRegister.Width = uint(width64)
    updateEditRegisterDialogs()
    return true
}

func deleteRegisterCallback(option *termdialog.Option) (close bool) {
    currentConfirmationCallback = deleteRegister
    dialogStack.Open(confirmationDialog)
    return false
}

func deleteRegister() {
    dialogStack.Close(editRegisterDialog)

    for i, register := range cpu.Registers {
        if register == currentlyEditingRegister {
            cpu.Registers = append(cpu.Registers[:i], cpu.Registers[i+1:]...)
            registersDialog.RemoveOption(i + 2)
        }
    }
}
