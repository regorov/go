package termdialog

import (
    "fmt"
    "github.com/nsf/termbox-go"
)

/*
  +------------+
  |            |
  |  Title     |
  |            |
  |  00 - xxx  |
  |  01 - yyy  |
  |  02 - zzz  |
  |            |
  +------------+
*/

type Option struct {
    Text     string
    Callback func(Option)
}

type SelectionDialog struct {
    title         string
    options       []Option
    metricsDirty  bool
    width         int
    height        int
    x             int
    y             int
    selectedIndex int
    theme         *Theme
}

func NewSelectionDialog(title string, options []Option) (dialog *SelectionDialog) {
    if options == nil {
        options = make([]Option, 0)
    }

    return &SelectionDialog{
        title:         title,
        options:       options,
        metricsDirty:  true,
        width:         0,
        height:        0,
        x:             0,
        y:             0,
        selectedIndex: 0,
        theme:         DefaultTheme,
    }
}

func (dialog *SelectionDialog) GetTitle() (title string) {
    return dialog.title
}

func (dialog *SelectionDialog) SetTitle(title string) {
    dialog.title = title
    dialog.metricsDirty = true
}

func (dialog *SelectionDialog) NOptions() (num int) {
    return len(dialog.options)
}

func (dialog *SelectionDialog) GetOption(n int) (option Option) {
    return dialog.options[n]
}

func (dialog *SelectionDialog) SetOption(n int, option Option) {
    dialog.options[n] = option
    dialog.metricsDirty = true
}

func (dialog *SelectionDialog) AddOption(option Option) {
    dialog.options = append(dialog.options, option)
    dialog.metricsDirty = true
}

func (dialog *SelectionDialog) RemoveOption(n int) {
    dialog.options = append(dialog.options[:n], dialog.options[n+1:]...)
    dialog.metricsDirty = true
}

func (dialog *SelectionDialog) GetMetricsDirty() (dirty bool) {
    return dialog.metricsDirty
}

func (dialog *SelectionDialog) GetWidth() (width int) {
    return dialog.width
}

func (dialog *SelectionDialog) GetHeight() (height int) {
    return dialog.height
}

func (dialog *SelectionDialog) GetX() (x int) {
    return dialog.x
}

func (dialog *SelectionDialog) GetY() (y int) {
    return dialog.y
}

func (dialog *SelectionDialog) GetSelectedIndex() (selectedIndex int) {
    return dialog.selectedIndex
}

func (dialog *SelectionDialog) SetSelectedIndex(selectedIndex int) {
    dialog.selectedIndex = selectedIndex
}

func (dialog *SelectionDialog) GetTheme() (theme *Theme) {
    return dialog.theme
}

func (dialog *SelectionDialog) SetTheme(theme *Theme) {
    dialog.theme = theme
}

func (dialog *SelectionDialog) CalcMetrics() {
    windowWidth, windowHeight := termbox.Size()

    maxWidth := 0
    for _, option := range dialog.options {
        if len(option.Text) > maxWidth {
            maxWidth = len(option.Text)
        }
    }

    maxWidth += 5 // Add the "00 - "
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
    if dialog.metricsDirty {
        dialog.CalcMetrics()
    }

    drawBox(dialog.x, dialog.y, dialog.width, dialog.height, dialog.theme.Border.FG, dialog.theme.Border.BG)
    fill(dialog.x+1, dialog.y+1, dialog.width-2, dialog.height-2, ' ', dialog.theme.Dialog.FG, dialog.theme.Dialog.BG)

    drawString(dialog.x+3, dialog.y+2, dialog.title, dialog.theme.Title.FG, dialog.theme.Title.BG)

    for i, option := range dialog.options {
        if i == dialog.selectedIndex {
            drawString(dialog.x+3, dialog.y+4+i, fmt.Sprintf("%2d - %s", i+1, option.Text), dialog.theme.ActiveItem.FG, dialog.theme.ActiveItem.BG)
        } else {
            drawString(dialog.x+3, dialog.y+4+i, fmt.Sprintf("%2d - %s", i+1, option.Text), dialog.theme.InactiveItem.FG, dialog.theme.InactiveItem.BG)
        }
    }
}

func (dialog *SelectionDialog) Close() {
    if dialog.metricsDirty {
        dialog.CalcMetrics()
    }

    fill(dialog.x, dialog.y, dialog.width, dialog.height, ' ', dialog.theme.Screen.FG, dialog.theme.Screen.BG)
}

func (dialog *SelectionDialog) HandleEvent(event termbox.Event) (handled bool, close bool) {
    switch event.Type {
    case termbox.EventKey:
        switch event.Key {
        case termbox.KeyEsc:
            return true, true

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
                termbox.Flush()
            }

            return true, false

        case termbox.KeyHome:
            dialog.selectedIndex = 0
            dialog.Open()
            termbox.Flush()

            return true, false

        case termbox.KeyEnd:
            dialog.selectedIndex = maxOption
            dialog.Open()
            termbox.Flush()

            return true, false

        case termbox.KeyEnter, termbox.KeySpace:
            option := dialog.options[dialog.selectedIndex]
            if option.Callback != nil {
                option.Callback(option)
            }

            return true, true
        }
    }

    return false, false
}
