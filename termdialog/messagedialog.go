package termdialog

import (
    "github.com/nsf/termbox-go"
)

/*
  +-----------+
  |           |
  |  Title    |
  |           |
  |  Message  |
  |           |
  +-----------+
*/

type MessageDialog struct {
    BaseDialog
    message string
}

func NewMessageDialog(title string, message string) (dialog *MessageDialog) {
    dialog = &MessageDialog{
        message: message,
    }

    dialog.BaseDialog.title = title
    dialog.BaseDialog.metricsDirty = true
    dialog.BaseDialog.theme = DefaultTheme
    return dialog
}

func (dialog *MessageDialog) GetMessage() (message string) {
    return dialog.message
}

func (dialog *MessageDialog) SetMessage(message string) {
    dialog.message = message
    dialog.metricsDirty = true
}

func (dialog *MessageDialog) CalcMetrics() {
    windowWidth, windowHeight := termbox.Size()

    maxWidth := len(dialog.message)
    if len(dialog.BaseDialog.title) > maxWidth {
        maxWidth = len(dialog.BaseDialog.title)
    }

    dialog.width = 6 + maxWidth // 6 = "|  " + "  |"
    dialog.height = 7

    dialog.x = (windowWidth / 2) - (dialog.width / 2)
    dialog.y = (windowHeight / 2) - (dialog.height / 2)

    dialog.metricsDirty = false
}

func (dialog *MessageDialog) Open() {
    baseDialogOpen(dialog)

    DrawString(dialog.x+3, dialog.y+4, dialog.message, dialog.theme.InactiveItem.FG, dialog.theme.InactiveItem.BG)
}

func (dialog *MessageDialog) HandleEvent(event termbox.Event) (handled bool, close bool) {
    handled, close = baseDialogHandleEvent(dialog, event)
    if handled {
        return
    }

    switch event.Type {
    case termbox.EventKey:
        switch event.Key {
        case termbox.KeyEnter, termbox.KeySpace:
            return true, true
        }
    }
    return false, false
}
