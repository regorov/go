// Package cellulose implements a generative musical sequencer based on Batuhan Bozkurt's Otomata:
// http://www.earslap.com/projectslab/otomata
package cellulose

import ()

const (
    NORTH = iota
    EAST
    SOUTH
    WEST
)

// 1st dimension is X coord
// 2nd dimension is Y coord
// 3rd dimension is index of entity within the cell
type Grid [][][]uint8

type NoteHandler func(int, int)

type Sequencer struct {
    Grid           Grid
    GridWidth      int
    GridHeight     int
    RowNoteHandler NoteHandler
    ColNoteHandler NoteHandler
}

func NewSequencer(gridWidth, gridHeight int, rowNoteHandler, colNoteHandler NoteHandler) (seq *Sequencer) {
    seq = new(Sequencer)
    seq.Grid = MakeGrid(gridWidth, gridHeight)
    seq.GridWidth = gridWidth
    seq.GridHeight = gridHeight
    seq.RowNoteHandler = rowNoteHandler
    seq.ColNoteHandler = colNoteHandler
    return seq
}

func (seq *Sequencer) Iterate() {
    gridWidth := seq.GridWidth
    gridHeight := seq.GridHeight
    grid := seq.Grid
    newGrid := MakeGrid(gridWidth, gridHeight)

    for x := 0; x < gridWidth; x++ {
        for y := 0; y < gridHeight; y++ {
            for _, cell := range grid[x][y] {
                switch cell {
                case NORTH:
                    if y == 0 {
                        InsertCell(newGrid, x, 1, SOUTH)
                        seq.ColNoteHandler(x, y)
                    } else {
                        InsertCell(newGrid, x, y-1, NORTH)
                    }

                case EAST:
                    if x == gridWidth-1 {
                        InsertCell(newGrid, gridWidth-2, y, WEST)
                        seq.RowNoteHandler(x, y)
                    } else {
                        InsertCell(newGrid, x+1, y, EAST)
                    }

                case SOUTH:
                    if y == gridHeight-1 {
                        InsertCell(newGrid, x, gridHeight-2, NORTH)
                        seq.ColNoteHandler(x, y)
                    } else {
                        InsertCell(newGrid, x, y+1, SOUTH)
                    }

                case WEST:
                    if x == 0 {
                        InsertCell(newGrid, 1, y, EAST)
                        seq.RowNoteHandler(x, y)
                    } else {
                        InsertCell(newGrid, x-1, y, WEST)
                    }
                }
            }
        }
    }

    for x := 0; x < gridWidth; x++ {
        for y := 0; y < gridHeight; y++ {
            s := newGrid[x][y]
            if len(s) > 1 {
                for i := 0; i < len(s)-1; i++ {
                    for j, cell := range s {
                        s[j] = (cell + 1) % 4
                    }
                }
            }
        }
    }

    seq.Grid = newGrid
}

func InsertCell(grid Grid, x int, y int, cell uint8) {
    s := grid[x][y]
    if s == nil {
        s = make([]uint8, 0, 4)
    }

    grid[x][y] = append(s, cell)
}

func MakeGrid(width int, height int) (grid Grid) {
    grid = make(Grid, width)
    for x := 0; x < width; x++ {
        grid[x] = make([][]uint8, height)
    }

    return grid
}
