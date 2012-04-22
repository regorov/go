package termdialog

import (
    "github.com/nsf/termbox-go"
)

type Style struct {
    FG termbox.Attribute
    BG termbox.Attribute
}

type Theme struct {
    Screen       Style
    Border       Style
    Dialog       Style
    Title        Style
    InactiveItem Style
    ActiveItem   Style
}

var DefaultTheme = &Theme{
    Screen:       Style{termbox.ColorBlack, termbox.ColorBlack},
    Border:       Style{termbox.ColorWhite, termbox.ColorBlack},
    Dialog:       Style{termbox.ColorWhite, termbox.ColorWhite},
    Title:        Style{termbox.ColorBlack | termbox.AttrUnderline, termbox.ColorWhite},
    InactiveItem: Style{termbox.ColorBlack, termbox.ColorWhite},
    ActiveItem:   Style{termbox.ColorWhite, termbox.ColorBlack},
}
