// Package clock implements the Generic Clock interface for the DCPU-16.
// 
// Device spec: http://0x10c.com/highnerd/rc_1/clock.txt
package clock

import (
    "fmt"
    "github.com/kierdavis/go/dcpuem"
    "sync"
    "time"
)

// Type Clock represents a Generic Clock device.
type Clock struct {
    // The governing mutex.
    Mutex sync.Mutex

    // The associated emulator.
    Em *dcpuem.Emulator

    // The internal ticker.
    Ticker *time.Ticker

    // Whether ticking has been enabled or not.
    TickingEnabled bool

    // The number of ticks since the last time the clock was configured.
    Ticks uint16

    // Send true to stop the ticker watcher.
    StopChan chan bool

    // The assigned message used for interrupts.
    InterruptMessage uint16
}

// Function New creates and returns a new clock.
func New() (d *Clock) {
    d = new(Clock)
    d.StopChan = make(chan bool)
    return d
}

// Function AssociateEmulator associates the specified emulator with the device.
func (d *Clock) AssociateEmulator(em *dcpuem.Emulator) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    d.Em = em
}

// Function ID returns the ID number of the device.
func (d *Clock) ID() (id uint32) {
    return 0x12d0b402
}

// Function Version returns the version number of the device. Since this is a generic device, this
// function simply returns 0.
func (d *Clock) Version() (ver uint16) {
    return 0
}

// Function Manufacturer returns the manufacturer ID of the device. Since this is a generic device,
// this function simple returns 0.
func (d *Clock) Manufacturer() (manu uint32) {
    return 0
}

// Function Interrupt handles a hardware interrupt.
func (d *Clock) Interrupt() (err error) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    a := d.Em.Regs[dcpuem.A]
    b := d.Em.Regs[dcpuem.B]

    switch a {
    case 0:
        if b == 0 {
            d.TickingEnabled = false
            d.Ticker = nil

        } else {
            if d.Ticker != nil {
                d.Ticker.Stop()
            }

            d.TickingEnabled = true
            d.Ticker = time.NewTicker(time.Second / time.Duration(b))
            d.Ticks = 0
        }

    case 1:
        d.Em.Regs[dcpuem.C] = d.Ticks

    case 2:
        d.InterruptMessage = b

    default:
        return &dcpuem.Error{dcpuem.ErrInvalidHardwareCommand,
            fmt.Sprintf("Invalid command to Generic Clock: 0x%04X", a), int(a)}
    }

    return nil
}

// Function Start starts the ticker watcher.
func (d *Clock) Start() {
    go d.Run()
}

// Function Stop stops the ticker watcher (by sending true along StopChan).
func (d *Clock) Stop() {
    d.StopChan <- true
}

// Function Run is the main body of the ticker watcher, which counts the ticks emitted by the
// associated ticker and triggers interrupts when necessary.
func (d *Clock) Run() {
    for {
        d.Mutex.Lock()

        select {
        case <-d.Ticker.C:
            d.Tick()

        case <-d.StopChan:
            return
        }
    }
}

// Function Tick records a tick by incrementing Ticks and sending an interrupt to the emulator
// (if enabled).
func (d *Clock) Tick() {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    if d.TickingEnabled {
        d.Ticks++

        if d.InterruptMessage != 0 {
            d.Em.Interrupt(d.InterruptMessage)
        }
    }
}
