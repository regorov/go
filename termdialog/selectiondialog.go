package termdialog

import (
    "fmt"
    "github.com/nsf/termbox-go"
)

/*
  +---------+
  |         |
  |  Title  |
  |         |
  |  xxx    |
  |  yyy    |
  |  zzz    |
  |         |
  +---------+
*/

// Type Option represents an option in a selection dialog.
type Option struct {
    Text     string             // The text of the option.
    Callback func(*Option) bool // The callback.
    Data     interface{}        // Arbitary associated data that can be accessed by the callback.
}

// Type SelectionDialog represents a dialog with a number of selectable options.
type SelectionDialog struct {
    BaseDialog
    options       []*Option
    selectedIndex int
}

// Function NewSelectionDialog creates and returns a new selection dialog. The options argument can
// be nil, and an empty slice will be created (options can be added later with AddOption).
func NewSelectionDialog(title string, options []*Option) (dialog *SelectionDialog) {
    if options == nil {
        options = make([]*Option, 0)
    }

    dialog = &SelectionDialog{
        options:       options,
        selectedIndex: 0,
    }

    dialog.BaseDialog.title = title
    dialog.BaseDialog.metricsDirty = true
    dialog.BaseDialog.theme = DefaultTheme
    return dialog
}

// Function NOptions returns the number of options attached to this dialog.
func (dialog *SelectionDialog) NOptions() (num int) {
    return len(dialog.options)
}

func (dialog *SelectionDialog) GetOption(n int) (option *Option) {
    return dialog.options[n]
}

func (dialog *SelectionDialog) SetOption(n int, option *Option) {
    dialog.options[n] = option
    dialog.metricsDirty = true
}

func (dialog *SelectionDialog) AddOption(option *Option) (theSameOption *Option) {
    dialog.options = append(dialog.options, option)
    dialog.metricsDirty = true
    return option
}

func (dialog *SelectionDialog) RemoveOption(n int) {
    dialog.options = append(dialog.options[:n], dialog.options[n+1:]...)
    dialog.metricsDirty = true
}

func (dialog *SelectionDialog) ClearOptions() {
    dialog.options = make([]*Option, 0)
    dialog.metricsDirty = true
}

func (dialog *SelectionDialog) GetSelectedIndex() (selectedIndex int) {
    return dialog.selectedIndex
}

func (dialog *SelectionDialog) SetSelectedIndex(selectedIndex int) {
    dialog.selectedIndex = selectedIndex
}

func (dialog *SelectionDialog) CalcMetrics() {
    windowWidth, windowHeight := termbox.Size()

    maxWidth := 0
    for _, option := range dialog.options {
        if len(option.Text) > maxWidth {
            maxWidth = len(option.Text)
        }
    }

    maxWidth += 2 // Add the "* "
    if len(dialog.title) > maxWidth {
        maxWidth = len(dialog.title)
    }

    dialog.width = 6 + maxWidth             // 6 = "|  " + "  |"
    dialog.height = 6 + len(dialog.options) // 6 = Top border, Top padding, Title, Under-title padding, Bottom padding, Bottom border

    dialog.x = (windowWidth / 2) - (dialog.width / 2)
    dialog.y = (windowHeight / 2) - (dialog.height / 2)

    dialog.metricsDirty = false
}

func (dialog *SelectionDialog) Open() {
    baseDialogOpen(dialog)

    for i, option := range dialog.options {
        if i == dialog.selectedIndex {
            DrawString(dialog.x+3, dialog.y+4+i, fmt.Sprintf("* %s", option.Text), dialog.theme.ActiveItem.FG, dialog.theme.ActiveItem.BG)
        } else {
            DrawString(dialog.x+3, dialog.y+4+i, fmt.Sprintf("* %s", option.Text), dialog.theme.InactiveItem.FG, dialog.theme.InactiveItem.BG)
        }
    }
}

func (dialog *SelectionDialog) HandleEvent(event termbox.Event) (handled bool, close bool) {
    handled, close = baseDialogHandleEvent(dialog, event)
    if handled {
        return
    }

    maxOption := len(dialog.options) - 1

    switch event.Type {
    case termbox.EventKey:
        switch event.Key {
        case termbox.KeyArrowUp:
            if dialog.selectedIndex > 0 {
                dialog.selectedIndex--
                dialog.Open()
                termbox.Flush()
            }

            return true, false

        case termbox.KeyArrowDown:
            if dialog.selectedIndex < maxOption {
                dialog.selectedIndex++
                dialog.Open()
            }

            return true, false

        case termbox.KeyHome:
            dialog.selectedIndex = 0
            dialog.Open()

            return true, false

        case termbox.KeyEnd:
            dialog.selectedIndex = maxOption
            dialog.Open()

            return true, false

        case termbox.KeyEnter, termbox.KeySpace:
            option := dialog.options[dialog.selectedIndex]
            close = true
            if option.Callback != nil {
                close = option.Callback(option)
            }

            return true, close
        }
    }

    return false, false
}
