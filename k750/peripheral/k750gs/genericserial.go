package k750gs

import (
    "sync"
)

type GenericSerial struct {
    Port              chan byte
    Fifo              []byte
    Lock              sync.Mutex
    StopChan          chan bool
    InterruptNumber   uint8
    PendingInterrupts chan uint8
}

func NewGenericSerial(port chan byte) (pp *GenericSerial) {
    pp = new(GenericSerial)
    pp.Port = port
    pp.Fifo = make([]byte, 16)
    pp.StopChan = make(chan bool)
    pp.PendingInterrupts = make(chan uint8, 16)
    return pp
}

func (pp *GenericSerial) ReadRegister(reg uint8) (value uint32) {
    switch reg {
    case 0x00:
        value = 0xf04c5900

    case 0x10:
        pp.Lock.Lock()

        if len(pp.Fifo) > 0 {
            value = uint32(pp.Fifo[0])
            pp.Fifo = pp.Fifo[1:]
        }

        pp.Lock.Unlock()

    case 0x11:
        pp.Lock.Lock()
        value = uint32(len(pp.Fifo))
        pp.Lock.Unlock()

    case 0x12:
        value = 16

    case 0x20:
        pp.Lock.Lock()
        value = uint32(pp.InterruptNumber)
        pp.Lock.Unlock()

    case 0x30:
        value = 0
    }

    return value
}

func (pp *GenericSerial) WriteRegister(reg uint8, value uint32) {
    switch reg {
    case 0x10:
        pp.Port <- uint8(value)

    case 0x20:
        pp.Lock.Lock()
        pp.InterruptNumber = uint8(value)
        pp.Lock.Unlock()
    }
}

func (pp *GenericSerial) GetPendingInterruptsChannel() (ch chan uint8) {
    return pp.PendingInterrupts
}

func (pp *GenericSerial) RunService() {
    for {
        select {
        case b := <-pp.Port:
            pp.Lock.Lock()
            pp.Fifo = append(pp.Fifo, b)

            if pp.InterruptNumber != 0 {
                pp.PendingInterrupts <- pp.InterruptNumber
            }

            pp.Lock.Unlock()

        case <-pp.StopChan:
            pp.StopChan <- true
            return
        }
    }
}

func (pp *GenericSerial) Start() {
    go pp.RunService()
}

func (pp *GenericSerial) Stop() {
    pp.StopChan <- true
    <-pp.StopChan
}
