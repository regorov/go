// Package lem1802_tb implements the LEM1802 display for the DCPU-16 using Termbox as a backend.
// 
// Device spec: http://0x10c.com/highnerd/rc_1/lem1802.txt
package lem1802_tb

import (
    "fmt"
    "github.com/kierdavis/go/dcpuem"
    "github.com/nsf/termbox-go"
    "sync"
    "time"
)

type LEM1802 struct {
    Mutex      sync.Mutex
    Em         *dcpuem.Emulator
    VideoMap   uint16
    ErrChan    chan error
    StopChan   chan bool
    BlinkCount uint8
}

func New() (d *LEM1802) {
    d = new(LEM1802)
    d.ErrChan = make(chan error, 16)
    d.StopChan = make(chan bool)
    return d
}

func (d *LEM1802) AssociateEmulator(em *dcpuem.Emulator) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    d.Em = em
}

func (d *LEM1802) ID() (id uint32) {
    return 0x7349f615
}

func (d *LEM1802) Version() (ver uint16) {
    return 0x1802
}

func (d *LEM1802) Manufacturer() (manu uint32) {
    return 0x1c6c8b36
}

func (d *LEM1802) Start() {
    go d.Run()

    // Also take the oppurtunity to set up the border:

    d.DrawBorder(7)
}

func (d *LEM1802) Stop() {
    d.StopChan <- true
}

func (d *LEM1802) GetError() (err error) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    if len(d.ErrChan) > 0 {
        return <-d.ErrChan
    }

    return nil
}

func (d *LEM1802) Run() {
    ticker := time.NewTicker(time.Second / 30.0)

    for {
        select {
        case <-d.StopChan:
            return

        case <-ticker.C:
            if d.VideoMap != 0 {
                d.Render()
            }

            d.BlinkCount++
        }
    }
}

func (d *LEM1802) Render() {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    addr := uint16(d.VideoMap)

    for y := 0; y < 12; y++ {
        for x := 0; x < 32; x++ {
            word := d.Em.MemoryLoad(addr)
            fg := (word >> 12) & 0xF
            bg := (word >> 8) & 0xF
            ch := rune(word & 0x7F)
            addr++

            if ch == 0 {
                ch = ' '
            }

            if (word&0x0080) != 0 && (d.BlinkCount&0x10) != 0 {
                fg, bg = bg, fg
            }

            termbox.SetCell(x+1, y+1, ch, Palette[fg], Palette[bg])
        }
    }

    termbox.Flush()
}

func (d *LEM1802) DrawBorder(color uint16) {
    fg := Palette[color]

    for i := 0; i < 12; i++ {
        termbox.SetCell(0, i+1, 0x2502, fg, 0)
        termbox.SetCell(33, i+1, 0x2502, fg, 0)
    }

    for i := 0; i < 32; i++ {
        termbox.SetCell(i+1, 0, 0x2500, fg, 0)
        termbox.SetCell(i+1, 13, 0x2500, fg, 0)
    }

    termbox.SetCell(0, 0, 0x250C, fg, 0)
    termbox.SetCell(33, 0, 0x2510, fg, 0)
    termbox.SetCell(0, 13, 0x2514, fg, 0)
    termbox.SetCell(33, 13, 0x2518, fg, 0)
}

func (d *LEM1802) Interrupt() (err error) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    a := d.Em.Regs[dcpuem.A]
    b := d.Em.Regs[dcpuem.B]

    switch a {
    case 0:
        d.VideoMap = b

    case 1:
        // Fonts not supported

    case 2:
        // Palette swapping not supported, although a rudimentary solution involving rounding the
        // RGB values to termbox.Color* constants could be implemented.

    case 3:
        d.DrawBorder(b & 0xF)

    case 4:
        // Fonts not supported.

    case 5:
        d.Em.MemoryStore(b, 0x0000)
        d.Em.MemoryStore(b+1, 0x0A00)
        d.Em.MemoryStore(b+2, 0x00A0)
        d.Em.MemoryStore(b+3, 0x05A0)
        d.Em.MemoryStore(b+4, 0x000A)
        d.Em.MemoryStore(b+5, 0x0A0A)
        d.Em.MemoryStore(b+6, 0x00AA)
        d.Em.MemoryStore(b+7, 0x0AAA)
        d.Em.MemoryStore(b+8, 0x0555)
        d.Em.MemoryStore(b+9, 0x0F55)
        d.Em.MemoryStore(b+10, 0x05F5)
        d.Em.MemoryStore(b+11, 0x0FF5)
        d.Em.MemoryStore(b+12, 0x055F)
        d.Em.MemoryStore(b+13, 0x0F5F)
        d.Em.MemoryStore(b+14, 0x05FF)
        d.Em.MemoryStore(b+15, 0x0FFF)

    default:
        return &dcpuem.Error{dcpuem.ErrInvalidHardwareCommand, fmt.Sprintf("Invalid command to LEM1802: 0x%04X", a), int(a)}
    }

    return nil
}

var Palette = [16]termbox.Attribute{
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
