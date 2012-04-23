package termdialog

import (
    "github.com/nsf/termbox-go"
    "strings"
)

/*
  +-------------------+
  |                   |
  |  Title            |
  |                   |
  |  Prompt ________  |
  |                   |
  +-------------------+
*/

type InputDialog struct {
    BaseDialog
    prompt     string
    valueWidth int
    value      string
    callback   func(string, interface{}) bool
    arg        interface{}
}

func NewInputDialog(title string, prompt string, valueWidth int, valueInit string, callback func(string, interface{}) bool, arg interface{}) (dialog *InputDialog) {
    dialog = &InputDialog{
        prompt:     prompt,
        valueWidth: valueWidth,
        value:      valueInit,
        callback:   callback,
        arg:        arg,
    }

    dialog.BaseDialog.title = title
    dialog.BaseDialog.metricsDirty = true
    dialog.BaseDialog.theme = DefaultTheme
    return dialog
}

func (dialog *InputDialog) GetPrompt() (prompt string) {
    return dialog.prompt
}

func (dialog *InputDialog) SetPrompt(prompt string) {
    dialog.prompt = prompt
    dialog.metricsDirty = true
}

func (dialog *InputDialog) GetValueWidth() (valueWidth int) {
    return dialog.valueWidth
}

func (dialog *InputDialog) SetValueWidth(valueWidth int) {
    dialog.valueWidth = valueWidth
    dialog.metricsDirty = true
}

func (dialog *InputDialog) GetValue() (value string) {
    return dialog.value
}

func (dialog *InputDialog) SetValue(value string) {
    dialog.value = value
}

func (dialog *InputDialog) GetCallback() (callback func(string, interface{}) bool) {
    return dialog.callback
}

func (dialog *InputDialog) SetCallback(callback func(string, interface{}) bool) {
    dialog.callback = callback
}

func (dialog *InputDialog) GetCallbackArg() (arg interface{}) {
    return dialog.arg
}

func (dialog *InputDialog) SetCallbackArg(arg interface{}) {
    dialog.arg = arg
}

func (dialog *InputDialog) CalcMetrics() {
    windowWidth, windowHeight := termbox.Size()

    maxWidth := len(dialog.prompt) + 1 + dialog.valueWidth
    if len(dialog.BaseDialog.title) > maxWidth {
        maxWidth = len(dialog.BaseDialog.title)
    }

    dialog.width = 6 + maxWidth // 6 = "|  " + "  |"
    dialog.height = 7

    dialog.x = (windowWidth / 2) - (dialog.width / 2)
    dialog.y = (windowHeight / 2) - (dialog.height / 2)

    dialog.metricsDirty = false
}

func (dialog *InputDialog) Open() {
    baseDialogOpen(dialog)

    DrawString(dialog.x+3, dialog.y+4, dialog.prompt, dialog.theme.InactiveItem.FG, dialog.theme.InactiveItem.BG)

    v := dialog.value + strings.Repeat("_", dialog.valueWidth-len(dialog.value))
    DrawString(dialog.x+4+len(dialog.prompt), dialog.y+4, v, dialog.theme.ActiveItem.FG, dialog.theme.ActiveItem.BG)
}

func (dialog *InputDialog) HandleEvent(event termbox.Event) (handled bool, close bool) {
    handled, close = baseDialogHandleEvent(dialog, event)
    if handled {
        return
    }

    switch event.Type {
    case termbox.EventKey:
        if event.Ch == 0 {
            switch event.Key {
            case termbox.KeyEnter:
                close = true
                if dialog.callback != nil {
                    close = dialog.callback(dialog.value, dialog.arg)
                }

                return true, close

            case termbox.KeyBackspace, termbox.KeyBackspace2:
                if len(dialog.value) > 0 {
                    dialog.value = dialog.value[:len(dialog.value)-1]
                    dialog.Open()
                }

                return true, false
            }

        } else {
            if len(dialog.value) < dialog.valueWidth {
                dialog.value += string(event.Ch)
                dialog.Open()
            }

            return true, false
        }
    }
    return false, false
}
