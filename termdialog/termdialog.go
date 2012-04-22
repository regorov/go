package termdialog

import (
    "github.com/nsf/termbox-go"
)

const (
    BOX_HOZ       rune = 0x2500
    BOX_VERT      rune = 0x2502
    BOX_CORNER_TL rune = 0x250C
    BOX_CORNER_TR rune = 0x2510
    BOX_CORNER_BL rune = 0x2514
    BOX_CORNER_BR rune = 0x2518
    BOX_TEE_L     rune = 0x251C
    BOX_TEE_R     rune = 0x2524
    BOX_TEE_T     rune = 0x252C
    BOX_TEE_B     rune = 0x2534
    BOX_CROSS     rune = 0x253C
)

func drawBox(x int, y int, width int, height int, fg termbox.Attribute, bg termbox.Attribute) {
    xmax := x + width - 1
    ymax := y + height - 1

    termbox.SetCell(x, y, BOX_CORNER_TL, fg, bg)
    termbox.SetCell(xmax, y, BOX_CORNER_TR, fg, bg)
    termbox.SetCell(x, ymax, BOX_CORNER_BL, fg, bg)
    termbox.SetCell(xmax, ymax, BOX_CORNER_BR, fg, bg)

    for i := x + 1; i <= xmax-1; i++ {
        termbox.SetCell(i, y, BOX_HOZ, fg, bg)
        termbox.SetCell(i, ymax, BOX_HOZ, fg, bg)
    }

    for i := y + 1; i <= ymax-1; i++ {
        termbox.SetCell(x, i, BOX_VERT, fg, bg)
        termbox.SetCell(xmax, i, BOX_VERT, fg, bg)
    }
}

func drawString(x int, y int, str string, fg termbox.Attribute, bg termbox.Attribute) {
    startX := x

    for _, c := range str {
        if c == '\r' {
            x = startX
        } else if c == '\n' {
            y++
        } else {
            termbox.SetCell(x, y, c, fg, bg)
            x++
        }
    }
}

func fill(x int, y int, width int, height int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
    for i := 0; i < width; i++ {
        for j := 0; j < height; j++ {
            termbox.SetCell(x+i, y+j, ch, fg, bg)
        }
    }
}

type Dialog interface {
    GetTitle() string
    SetTitle(string)
    CalcMetrics()
    Open()
    Close()
    HandleEvent() (bool, bool)
}
