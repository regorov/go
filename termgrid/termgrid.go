package termgrid

import (
	"github.com/nsf/termbox-go"
)

const (
	Top = 1 << iota
	Bottom
	Left
	Right
)

type BoxStyle struct {
	Hoz, Vert, Cross                       rune
	Top, Bottom, Left, Right               rune
	TLCorner, TRCorner, BLCorner, BRCorner rune
	TopTee, BottomTee, LeftTee, RightTee   rune
}

func (style *BoxStyle) Char(sides uint8) (ch rune) {
	switch sides {
	case Top:
		return style.Top
	case Bottom:
		return style.Bottom
	case Bottom | Top:
		return style.Vert
	case Left:
		return style.Left
	case Left | Top:
		return style.BRCorner
	case Left | Bottom:
		return style.TRCorner
	case Left | Bottom | Top:
		return style.RightTee
	case Right:
		return style.Right
	case Right | Top:
		return style.BLCorner
	case Right | Bottom:
		return style.TLCorner
	case Right | Bottom | Top:
		return style.LeftTee
	case Right | Left:
		return style.Hoz
	case Right | Left | Top:
		return style.BottomTee
	case Right | Left | Bottom:
		return style.TopTee
	case Right | Left | Bottom | Top:
		return style.Cross
	}

	return ' '
}

var Light = &BoxStyle{
	Hoz:       0x2500,
	Vert:      0x2502,
	Cross:     0x253C,
	Top:       0x2575,
	Bottom:    0x2577,
	Left:      0x2574,
	Right:     0x2576,
	TLCorner:  0x250C,
	TRCorner:  0x2510,
	BLCorner:  0x2514,
	BRCorner:  0x2518,
	TopTee:    0x252C,
	BottomTee: 0x2534,
	LeftTee:   0x251C,
	RightTee:  0x2524,
}

var Heavy = &BoxStyle{
	Hoz:       0x2501,
	Vert:      0x2503,
	Cross:     0x254B,
	Top:       0x2579,
	Bottom:    0x257B,
	Left:      0x2578,
	Right:     0x257A,
	TLCorner:  0x250F,
	TRCorner:  0x2513,
	BLCorner:  0x2517,
	BRCorner:  0x251B,
	TopTee:    0x2533,
	BottomTee: 0x253B,
	LeftTee:   0x2523,
	RightTee:  0x252B,
}

type Grid [][]uint8

func NewGrid() (grid Grid) {
	width, height := termbox.Size()

	grid = make(Grid, width)

	for i := 0; i < width; i++ {
		grid[i] = make([]uint8, height)
	}

	return grid
}

func (grid Grid) Width() (width int) {
	return len(grid)
}

func (grid Grid) Height() (height int) {
	return len(grid[0])
}

func (grid Grid) HozLine(startX, endX, y int) {
	for x := startX; x < endX; x++ {
		grid[x][y] |= Right
	}

	for x := endX; x > startX; x-- {
		grid[x][y] |= Left
	}
}

func (grid Grid) VertLine(x, startY, endY int) {
	for y := startY; y < endY; y++ {
		grid[x][y] |= Bottom
	}

	for y := endY; y > startY; y-- {
		grid[x][y] |= Top
	}
}

func (grid Grid) Box(startX, startY, width, height int) {
	endX := startX + width - 1
	endY := startY + height - 1

	grid.HozLine(startX, endX, startY)
	grid.HozLine(startX, endX, endY)
	grid.VertLine(startX, startY, endY)
	grid.VertLine(endX, startY, endY)
}

func (grid Grid) Draw(style *BoxStyle, fg, bg termbox.Attribute) {
	for x, col := range grid {
		for y, cell := range col {
			if cell != 0 {
				termbox.SetCell(x, y, style.Char(cell), fg, bg)
			}
		}
	}
}
