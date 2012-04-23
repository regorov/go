package main

import (
    "fmt"
    "github.com/kierdavis/go/termdialog"
    "strconv"
)

func addRamWidthCallback(name string, arg interface{}) (close bool) {
    addRamDialog.SetValue("")
    addRamWidthDialog.SetCallbackArg(name)
    dialogStack.Open(addRamWidthDialog)
    return true
}

func addRamDepthCallback(widthStr string, arg interface{}) (close bool) {
    name := arg.(string)
    width64, err := strconv.ParseUint(widthStr, 10, 0)
    if err != nil {
        panic(err)
    }

    width := uint(width64)
    addRamDepthDialog.SetCallbackArg(nameWidth{name, width})
    dialogStack.Open(addRamDepthDialog)
    return true
}

func addRamCallback(depthStr string, arg interface{}) (close bool) {
    nw := arg.(nameWidth)
    name := nw.name
    width := nw.width

    depth64, err := strconv.ParseUint(depthStr, 10, 0)
    if err != nil {
        panic(err)
    }

    depth := uint(depth64)

    ram := &Ram{
        Name:  name,
        Width: width,
        Depth: depth,
    }

    addRamToGui(ram)
    cpu.Rams = append(cpu.Rams, ram)

    return true
}

func addRamToGui(ram *Ram) {
    ram.EditOption = ramsDialog.AddOption(&termdialog.Option{"Edit " + ram.Name, editRamOpenCallback, ram})
}

func updateEditRamDialogs() {
    editRamDialog.SetTitle("Edit Register Bank " + currentlyEditingRam.Name)
    editRamDialog.GetOption(0).Text = fmt.Sprintf("Edit name (%s)", currentlyEditingRam.Name)
    editRamDialog.GetOption(1).Text = fmt.Sprintf("Edit width (%d)", currentlyEditingRam.Width)
    editRamDialog.GetOption(2).Text = fmt.Sprintf("Edit depth (%d)", currentlyEditingRam.Depth)

    editRamNameDialog.SetTitle("Edit Register Bank " + currentlyEditingRam.Name)
    editRamNameDialog.SetValue(currentlyEditingRam.Name)

    editRamWidthDialog.SetTitle("Edit Register Bank " + currentlyEditingRam.Name)
    editRamWidthDialog.SetValue(strconv.FormatUint(uint64(currentlyEditingRam.Width), 10))

    editRamDepthDialog.SetTitle("Edit Register Bank " + currentlyEditingRam.Name)
    editRamDepthDialog.SetValue(strconv.FormatUint(uint64(currentlyEditingRam.Depth), 10))
}

func editRamOpenCallback(option *termdialog.Option) (close bool) {
    currentlyEditingRam = option.Data.(*Ram)
    updateEditRamDialogs()
    dialogStack.Open(editRamDialog)
    return false
}

func editRamNameCallback(name string, arg interface{}) (close bool) {
    currentlyEditingRam.Name = name
    updateEditRamDialogs()
    currentlyEditingRam.EditOption.Text = "Edit " + name
    return true
}

func editRamWidthCallback(widthStr string, arg interface{}) (close bool) {
    width64, err := strconv.ParseUint(widthStr, 10, 0)
    if err != nil {
        panic(err)
    }

    currentlyEditingRam.Width = uint(width64)
    updateEditRamDialogs()
    return true
}

func editRamDepthCallback(depthStr string, arg interface{}) (close bool) {
    depth64, err := strconv.ParseUint(depthStr, 10, 0)
    if err != nil {
        panic(err)
    }

    currentlyEditingRam.Depth = uint(depth64)
    updateEditRamDialogs()
    return true
}

func deleteRamCallback(option *termdialog.Option) (close bool) {
    currentConfirmationCallback = deleteRam
    dialogStack.Open(confirmationDialog)
    return false
}

func deleteRam() {
    dialogStack.Close(editRamDialog)

    for i, ram := range cpu.Rams {
        if ram == currentlyEditingRam {
            cpu.Rams = append(cpu.Rams[:i], cpu.Rams[i+1:]...)
            ramsDialog.RemoveOption(i + 2)
        }
    }
}
