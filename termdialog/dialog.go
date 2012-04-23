package termdialog

import (
    "github.com/nsf/termbox-go"
)

type Dialog interface {
    GetTitle() string
    SetTitle(string)
    GetMetricsDirty() bool
    GetWidth() int
    GetHeight() int
    GetX() int
    GetY() int
    GetTheme() *Theme
    SetTheme(*Theme)

    CalcMetrics()
    Open()
    Close()
    HandleEvent(termbox.Event) (bool, bool)
}

type BaseDialog struct {
    title        string
    metricsDirty bool
    width        int
    height       int
    x            int
    y            int
    theme        *Theme
}

func (dialog *BaseDialog) GetTitle() (title string) {
    return dialog.title
}

func (dialog *BaseDialog) SetTitle(title string) {
    dialog.title = title
    dialog.metricsDirty = true
}

func (dialog *BaseDialog) GetMetricsDirty() (dirty bool) {
    return dialog.metricsDirty
}

func (dialog *BaseDialog) GetWidth() (width int) {
    return dialog.width
}

func (dialog *BaseDialog) GetHeight() (height int) {
    return dialog.height
}

func (dialog *BaseDialog) GetX() (x int) {
    return dialog.x
}

func (dialog *BaseDialog) GetY() (y int) {
    return dialog.y
}

func (dialog *BaseDialog) GetTheme() (theme *Theme) {
    return dialog.theme
}

func (dialog *BaseDialog) SetTheme(theme *Theme) {
    dialog.theme = theme
}

func (dialog *BaseDialog) Close() {
    if dialog.metricsDirty {
        dialog.CalcMetrics()
    }

    Fill(dialog.x, dialog.y, dialog.width, dialog.height, ' ', dialog.theme.Screen.FG, dialog.theme.Screen.BG)
}

func (dialog *BaseDialog) CalcMetrics() {

}

func baseDialogOpen(dialog Dialog) {
    if dialog.GetMetricsDirty() {
        dialog.CalcMetrics()
    }

    title := dialog.GetTitle()
    x := dialog.GetX()
    y := dialog.GetY()
    width := dialog.GetWidth()
    height := dialog.GetHeight()
    theme := dialog.GetTheme()

    if theme.HasShadow {
        Fill(x+theme.ShadowOffsetX, y+theme.ShadowOffsetY, width, height, ' ', theme.Shadow.FG, theme.Shadow.BG)
    }

    DrawBox(x, y, width, height, theme.Border.FG, theme.Border.BG)
    Fill(x+1, y+1, width-2, height-2, ' ', theme.Dialog.FG, theme.Dialog.BG)

    DrawString(x+3, y+2, title, theme.Title.FG, theme.Title.BG)
}

func baseDialogHandleEvent(dialog Dialog, event termbox.Event) (handled bool, close bool) {
    switch event.Type {
    case termbox.EventKey:
        switch event.Key {
        case termbox.KeyEsc:
            return true, true
        }
    }

    return false, false
}
