package termdialog

import (
    "github.com/nsf/termbox-go"
)

// Type Style represents a pair of foreground and background attributes.
type Style struct {
    FG termbox.Attribute
    BG termbox.Attribute
}

// Type Theme represents a GUI theme.
type Theme struct {
    Screen       Style // The style for the empty background region.
    Shadow       Style // The style for the shadow of dialogs (if enabled).
    Border       Style // The style for the border of dialogs.
    Dialog       Style // The style for the empty background of dialogs.
    Title        Style // The style for the title text of dialogs.
    InactiveItem Style // The style for inactive items and static text on dialogs.
    ActiveItem   Style // The style for active items and widgets that can be interacted with.

    HasShadow     bool // Whether to display a shadow behind dialogs.
    ShadowOffsetX int  // The X offset of the shadow, relative to the dialog's coordinates.
    ShadowOffsetY int  // The Y offset of the shadow, relative to the dialog's coordinates.
}

// Variable DefaultTheme is the default theme.
var DefaultTheme = &Theme{
    Screen:       Style{termbox.ColorBlack, termbox.ColorBlack},
    Shadow:       Style{termbox.ColorBlack, termbox.ColorBlack},
    Border:       Style{termbox.ColorWhite, termbox.ColorBlack},
    Dialog:       Style{termbox.ColorWhite, termbox.ColorWhite},
    Title:        Style{termbox.ColorBlack | termbox.AttrUnderline, termbox.ColorWhite},
    InactiveItem: Style{termbox.ColorBlack, termbox.ColorWhite},
    ActiveItem:   Style{termbox.ColorWhite, termbox.ColorRed},

    HasShadow:     false,
    ShadowOffsetX: 2,
    ShadowOffsetY: 1,
}
