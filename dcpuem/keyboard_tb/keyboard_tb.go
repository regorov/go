// Package lem1802_tb implements the Generic Keyboard interface for the DCPU-16 using Termbox as a
// backend.
//
// Device spec: http://0x10c.com/highnerd/rc_1/keyboard.txt
package keyboard_tb

import (
    "fmt"
    "github.com/kierdavis/go/dcpuem"
    "github.com/nsf/termbox-go"
    "sync"
)

type Keyboard struct {
    Mutex            sync.Mutex
    Em               *dcpuem.Emulator
    Buffer           []uint16
    InterruptMessage uint16
}

func New() (d *Keyboard) {
    d = new(Keyboard)
    d.Buffer = make([]uint16, 16)
    return d
}

func (d *Keyboard) AssociateEmulator(em *dcpuem.Emulator) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    d.Em = em
}

func (d *Keyboard) ID() (id uint32) {
    return 0x30cf7406
}

func (d *Keyboard) Version() (ver uint16) {
    return 0
}

func (d *Keyboard) Manufacturer() (manu uint32) {
    return 0
}

func (d *Keyboard) Interrupt() (err error) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    a := d.Em.Regs[dcpuem.A]
    b := d.Em.Regs[dcpuem.B]

    switch a {
    case 0:
        d.Buffer = d.Buffer[:0]

    case 1:
        if len(d.Buffer) > 0 {
            d.Em.Regs[dcpuem.C] = d.Buffer[0]
            d.Buffer = d.Buffer[1:]
        } else {
            d.Em.Regs[dcpuem.C] = 0
        }

    case 2:
        // Key checking not implemented

    case 3:
        d.InterruptMessage = b

    default:
        return &dcpuem.Error{dcpuem.ErrInvalidHardwareCommand, fmt.Sprintf("Invalid command to Generic Keyboard: 0x%04X", a), int(a)}
    }

    return nil
}

func (d *Keyboard) Start() {

}

func (d *Keyboard) Stop() {

}

func (d *Keyboard) SendKey(key uint16) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    //termbox.SetCell(32, 12, rune(key), 0, 0)
    d.Buffer = append(d.Buffer, key)

    if d.InterruptMessage != 0 {
        d.Em.Interrupt(d.InterruptMessage)
    }
}

func (d *Keyboard) HandleEvent(ev termbox.Event) (handled bool) {
    switch ev.Type {
    case termbox.EventKey:
        if ev.Ch == 0 {
            switch ev.Key {
            case termbox.KeyBackspace, termbox.KeyBackspace2:
                d.SendKey(0x10)
                return true

            case termbox.KeyEnter:
                d.SendKey(0x11)
                return true

            case termbox.KeyInsert:
                d.SendKey(0x12)
                return true

            case termbox.KeyDelete:
                d.SendKey(0x13)
                return true

            case termbox.KeySpace:
                d.SendKey(' ')
                return true

            case termbox.KeyArrowUp:
                d.SendKey(0x80)
                return true

            case termbox.KeyArrowDown:
                d.SendKey(0x81)
                return true

            case termbox.KeyArrowLeft:
                d.SendKey(0x82)
                return true

            case termbox.KeyArrowRight:
                d.SendKey(0x83)
                return true
            }

        } else {
            d.SendKey(uint16(ev.Ch))
            return true
        }
    }

    return false
}
