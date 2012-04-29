// Package cellulose implements a generative musical sequencer based on Batuhan Bozkurt's Otomata:
// http://www.earslap.com/projectslab/otomata
package cellulose

const (
    NORTH = iota // The value of a cell that is heading north.
    EAST         // The value of a cell that is heading east.
    SOUTH        // The value of a cell that is heading south.
    WEST         // The value of a cell that is heading west.
)

// Type Grid represents a grid of cells. The 1st dimension is the X coordinate, the 2nd dimension is
// the Y coordinate and the 3rd dimension is the index of the cell within this grid location (each
// location can contain multiple cells).
type Grid [][][]uint8

// Function NewGrid creates and returns a new grid.
func NewGrid(width int, height int) (grid Grid) {
    grid = make(Grid, width)
    for x := 0; x < width; x++ {
        grid[x] = make([][]uint8, height)
    }

    return grid
}

// Type NoteHandler represents a function that is called when a cell bounces off an edge. The first
// parameter is the X coordinate and the second is the Y coordinate.
type NoteHandler func(int, int)

// Type Sequencer represents a sequencer.
type Sequencer struct {
    Grid           Grid
    GridWidth      int
    GridHeight     int
    RowNoteHandler NoteHandler
    ColNoteHandler NoteHandler
}

// Function NewSequencer creates and returns a new sequencer, given the specified grid size and
// row/column bounce handlers.
func NewSequencer(gridWidth int, gridHeight int, rowNoteHandler NoteHandler, colNoteHandler NoteHandler) (seq *Sequencer) {
    seq = new(Sequencer)
    seq.Grid = NewGrid(gridWidth, gridHeight)
    seq.GridWidth = gridWidth
    seq.GridHeight = gridHeight
    seq.RowNoteHandler = rowNoteHandler
    seq.ColNoteHandler = colNoteHandler
    return seq
}

// Function Iterate runs an iteration of the simulation, moving all cells as dictated by their
// direction.
func (seq *Sequencer) Iterate() {
    gridWidth := seq.GridWidth
    gridHeight := seq.GridHeight
    grid := seq.Grid
    newGrid := NewGrid(gridWidth, gridHeight)

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

// Function InsertCell inserts a cell into a grid at the specified coordinates.
func InsertCell(grid Grid, x int, y int, cell uint8) {
    s := grid[x][y]
    if s == nil {
        s = make([]uint8, 0, 4)
    }

    grid[x][y] = append(s, cell)
}
