package termdialog

import (
    "github.com/nsf/termbox-go"
)

type DialogStack struct {
    dialogs []Dialog
}

func NewDialogStack() (dialogStack *DialogStack) {
    return &DialogStack{
        dialogs: make([]Dialog, 0),
    }
}

func (dialogStack *DialogStack) Open(dialog Dialog) {
    dialog.Open()
    dialogStack.dialogs = append(dialogStack.dialogs, dialog)
    //return dialog
}

func (dialogStack *DialogStack) Close(dialog Dialog) {
    dialog.Close()

    for i, d := range dialogStack.dialogs {
        if d == dialog {
            dialogStack.dialogs = append(dialogStack.dialogs[:i], dialogStack.dialogs[i+1:]...)
        }
    }
}

func (dialogStack *DialogStack) CloseTop() {
    dialog := dialogStack.dialogs[len(dialogStack.dialogs)-1]
    dialog.Close()
    dialogStack.dialogs = dialogStack.dialogs[:len(dialogStack.dialogs)-1]
    //return dialog
}

func (dialogStack *DialogStack) Run() {
    windowWidth, windowHeight := termbox.Size()
    Fill(0, 0, windowWidth, windowHeight, ' ', DefaultTheme.Screen.FG, DefaultTheme.Screen.BG)

    for len(dialogStack.dialogs) > 0 {
        for _, dialog := range dialogStack.dialogs {
            dialog.Open()
        }
        termbox.Flush()

        activeDialog := dialogStack.dialogs[len(dialogStack.dialogs)-1]
        event := termbox.PollEvent()

        _, close := activeDialog.HandleEvent(event)
        if close {
            dialogStack.Close(activeDialog)
        }
    }
}
