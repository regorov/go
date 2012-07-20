package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/kierdavis/go/termgrid"
	"github.com/nsf/termbox-go"
	"math/big"
	"os"
	"strings"
	"time"
)

const Keys = "qweiop"

const PasswordLen = 30
const PaddingLen = 18

const KeyInterval = (time.Second * 1) / 2

func randInt(max int) (n int) {
	x, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}

	return int(x.Int64())
}

func randKeys() (ch chan byte) {
	ch = make(chan byte)

	go func() {
		a := randInt(len(Keys))
		ch <- Keys[a]

		for {
			b := randInt(len(Keys) - 1)

			if b >= a {
				b++
			}

			ch <- Keys[b]
			a = b
		}
	}()

	return ch
}

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer termbox.Close()

	grid := termgrid.NewGrid()
	numCols := len(Keys)
	colsStart := -(numCols / 2)
	colsEnd := numCols - (numCols / 2)
	numRows := 7
	colWidth := 13
	rowHeight := 7
	centreX := grid.Width() / 2
	centreY := grid.Height() / 2
	totalHeight := ((rowHeight - 1) * numRows) + 1
	startY := centreY - (totalHeight / 2)
	endY := centreY + (totalHeight - (totalHeight / 2))

	for i := colsStart; i < colsEnd; i++ {
		x := centreX + (i * (colWidth - 1))
		c := x + (colWidth / 2)

		for j := 1; j < numRows-1; j++ {
			grid.HozLine(x+1, c-1, startY+((rowHeight-1)*j))
			grid.HozLine(c+1, x+colWidth-2, startY+((rowHeight-1)*j))
		}

		//grid.VertLine(c, startY+(rowHeight/2), endY-(rowHeight/2))
	}

	grid.Draw(termgrid.Light, termbox.ColorWhite, termbox.ColorBlack)
	grid = termgrid.NewGrid()

	for i := colsStart; i < colsEnd; i++ {
		x := centreX + (i * (colWidth - 1))
		c := x + (colWidth / 2)

		grid.Box(x, startY, colWidth, totalHeight)
		grid.HozLine(x+1, c-1, endY-rowHeight)
		grid.HozLine(c+1, x+colWidth-2, endY-rowHeight)
	}

	//grid.HozLine(centreX-(3*(colWidth-1)), centreX+(3*(colWidth-1)), endY-rowHeight)
	grid.Draw(termgrid.Heavy, termbox.ColorWhite, termbox.ColorBlack)

	termbox.Flush()

	playChan := make(chan byte)

	go func() {
		keyChan := randKeys()
		pwd := ""

		for i := 0; i < PasswordLen; i++ {
			pwd += string(<-keyChan)
		}

		// Repeat entire process 7 times
		for i := 0; i < 7; i++ {
			// Repeat 5 times
			for j := 0; j < 5; j++ {
				// Play password 3 times
				for k := 0; k < 3; k++ {
					for l := 0; l < PasswordLen; l++ {
						playChan <- pwd[l]
						time.Sleep(KeyInterval)
					}
				}

				// Play a padding gap
				for k := 0; k < PaddingLen; k++ {
					playChan <- <-keyChan
					time.Sleep(KeyInterval)
				}
			}

			// ... then take a break
			time.Sleep(time.Second * 10)
		}

		close(playChan)
	}()

	eventChan := make(chan termbox.Event)

	go func() {
		for {
			eventChan <- termbox.PollEvent()
		}
	}()

	cols := make([][]byte, numCols)
	colLength := (endY - (rowHeight / 2)) - (startY + (rowHeight / 2))

	for i := 0; i < numCols; i++ {
		cols[i] = make([]byte, colLength)

		for j := 0; j < colLength; j++ {
			cols[i][j] = ' '
		}
	}

	// Ticks at the rate that keys move down the screen.
	// Using rowHeight Hz means that keys move 1 row per second.
	ticker := time.NewTicker(KeyInterval / time.Duration(rowHeight))

	inverted := false

	for {
		select {
		case event := <-eventChan:
			switch event.Type {
			case termbox.EventKey:
				i := strings.IndexRune(Keys, event.Ch)
				if i >= 0 {
					col := cols[i]
					part := col[:(rowHeight/2)+1]

					j := bytes.Index(part, []byte{byte(event.Ch)})
					if j < 0 {
						inverted = true
					} else {
						part[j] = ' '
					}

					break
				}

				switch event.Ch {
				case 0:
					switch event.Key {
					case termbox.KeyEsc:
						//return
					}
				}
			}

		case <-ticker.C:
			// Draw the columns

			for i, col := range cols {
				for j, cell := range col {
					j = colLength - 1 - j
					x := centreX + ((i + colsStart) * (colWidth - 1)) + (colWidth / 2)
					y := j + startY + (rowHeight / 2)
					termbox.SetCell(x, y, rune(cell), termbox.ColorWhite, termbox.ColorBlack)
				}
			}

			if inverted {
				buffer := termbox.CellBuffer()
				for i, cell := range buffer {
					cell.Fg |= termbox.AttrReverse
					buffer[i] = cell
				}
			}

			termbox.Flush()

			if inverted {
				buffer := termbox.CellBuffer()
				for i, cell := range buffer {
					cell.Fg &= ^termbox.AttrReverse
					buffer[i] = cell
				}

				inverted = false
			}

			// Shift the columns

			for i, col := range cols {
				if col[0] != ' ' {
					inverted = true
				}

				cols[i] = append(col[1:], ' ')
			}

		case key := <-playChan:
			for i, col := range cols {
				if key == Keys[i] {
					col[colLength-1] = key
					break
				}
			}
		}
	}
}
