package lem1802_tb

import (
    "github.com/nsf/termbox-go"
)

var ColorMap = [16]termbox.Attribute{
    termbox.ColorBlack,
    termbox.ColorRed,
    termbox.ColorGreen,
    termbox.ColorYellow,
    termbox.ColorBlue,
    termbox.ColorMagenta,
    termbox.ColorCyan,
    termbox.ColorWhite,

    termbox.AttrBold | termbox.ColorBlack,
    termbox.AttrBold | termbox.ColorRed,
    termbox.AttrBold | termbox.ColorGreen,
    termbox.AttrBold | termbox.ColorYellow,
    termbox.AttrBold | termbox.ColorBlue,
    termbox.AttrBold | termbox.ColorMagenta,
    termbox.AttrBold | termbox.ColorCyan,
    termbox.AttrBold | termbox.ColorWhite,
}
