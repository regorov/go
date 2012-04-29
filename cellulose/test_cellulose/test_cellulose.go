package main

import (
    "github.com/banthar/Go-SDL/mixer"
    "github.com/banthar/Go-SDL/sdl"
    "github.com/kierdavis/go/amberfell/cellulose"
    "github.com/kierdavis/go/resourcemanager"
    "github.com/nsf/termbox-go"
    "math/rand"
    "path/filepath"
    "time"
)

var ResourceManager *resourcemanager.ResourceManager
var RowSounds []*mixer.Chunk
var ColSounds []*mixer.Chunk

func playRow(x int, y int) {
    RowSounds[y%len(RowSounds)].PlayChannel(-1, 0)
}

func playCol(x int, y int) {
    ColSounds[x%len(ColSounds)].PlayChannel(-1, 0)
}

func runSeq() {
    seq := cellulose.NewSequencer(9, 9, playRow, playCol)

    rand.Seed(time.Now().UnixNano())

    n := rand.Intn(seq.GridWidth) + seq.GridWidth
    for i := 0; i < n; i++ {
        x := rand.Intn(5) + 2
        y := rand.Intn(5) + 2
        dir := uint8(rand.Intn(4))
        cellulose.InsertCell(seq.Grid, x, y, dir)
    }

    for i := 0; i < seq.GridWidth; i++ {
        termbox.SetCell(i, seq.GridHeight, '-', 0, 0)
    }

    for i := 0; i < seq.GridHeight; i++ {
        termbox.SetCell(seq.GridWidth, i, '|', 0, 0)
    }

    for {
        seq.Iterate()

        for x := 0; x < seq.GridWidth; x++ {
            for y := 0; y < seq.GridHeight; y++ {
                var ch rune

                s := seq.Grid[x][y]

                switch len(s) {
                case 0:
                    ch = ' '
                case 1:
                    switch s[0] {
                    case cellulose.NORTH:
                        ch = '^'
                    case cellulose.EAST:
                        ch = '>'
                    case cellulose.SOUTH:
                        ch = 'v'
                    case cellulose.WEST:
                        ch = '<'
                    }
                default:
                    ch = 'o'
                }

                termbox.SetCell(x, y, ch, 0, 0)
            }
        }

        termbox.Flush()
        time.Sleep(time.Second / 6)
    }
}

func main() {
    ResourceManager = resourcemanager.NewResourceManager("github.com/kierdavis/go/amberfell/cellulose/test_cellulose")

    if mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT, mixer.DEFAULT_CHANNELS, 4096) != 0 {
        panic(sdl.GetError())
    }
    defer mixer.CloseAudio()

    soundFiles, err := filepath.Glob(ResourceManager.GetFilename("sound/out/*.wav"))
    if err != nil {
        panic(err)
    }

    RowSounds = make([]*mixer.Chunk, 0)
    ColSounds = make([]*mixer.Chunk, 0)

    l2 := len(soundFiles) / 2

    for _, filename := range soundFiles[:l2] {
        RowSounds = append(RowSounds, mixer.LoadWAV(filename))
    }

    for _, filename := range soundFiles[l2:] {
        ColSounds = append(ColSounds, mixer.LoadWAV(filename))
    }

    err = termbox.Init()
    if err != nil {
        panic(err)
    }
    defer termbox.Close()

    go runSeq()

loop:
    for {
        switch ev := termbox.PollEvent(); ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            case termbox.KeyEsc:
                break loop
            }
        }
    }
}
