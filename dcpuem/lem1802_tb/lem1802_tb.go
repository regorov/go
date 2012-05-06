// Package lem1802_tb implements the LEM1802 display for the DCPU-16 using Termbox as a backend.
// 
// Example usage:
// 
//      var display *lem1802_tb.LEM1802
//      var em *dcpuem.Emulator
//      var err error
//      var ev termbox.Event
//      
//      // Create an emulator.
//      em = dcpuem.NewEmulator()
//      
//      // Create a display & attach it to the emulator.
//      display = lem1802_tb.New()
//      em.AttachDevice(display)
//      
//      // Initialise Termbox.
//      err = termbox.Init()
//      if err != nil {panic(err)}
//      defer termbox.Close()
//      
//      // Start the display's service.
//      display.Start()
//      defer display.Stop()
//      
//      // Start the emulator.
//      go em.Run()
//      
//      // Termbox event loop.
//      mainloop:
//      for {
//          ev = termbox.PollEvent()
//          
//          switch ev.Type {
//          case termbox.EventKey:
//              switch ev.Key {
//              case termbox.KeyEsc:
//                  break mainloop
//              }
//          }
//      }
package lem1802_tb

import (
    "github.com/kierdavis/go/dcpuem"
    "github.com/nsf/termbox-go"
    "sync"
    "time"
)

type LEM1802 struct {
    Mutex      sync.Mutex
    Em         *dcpuem.Emulator
    VideoMap   uint16
    FontMap    uint16
    PaletteMap uint16
    ErrChan    chan error
    StopChan   chan bool
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

            termbox.SetCell(x, y, ch, ColorMap[fg], ColorMap[bg])
        }
    }

    termbox.Flush()
}

func (d *LEM1802) Interrupt() (err error) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    switch d.Em.Regs[dcpuem.A] {
    case 0: // MEM_MAP_SCREEN
        d.VideoMap = d.Em.Regs[dcpuem.B]
    }

    return nil
}
