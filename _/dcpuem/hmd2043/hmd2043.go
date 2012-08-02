// Package hmd2043 implements the unofficial Harold Media Drive, conforming to version 1.1 of the
// specification.
// 
// Device spec: https://gist.github.com/2495578
package hmd2043

import (
    "github.com/kierdavis/go/dcpuem"
    "sync"
    "time"
)

// Device flags.
const (
    NonBlocking          uint16 = 1 << 0
    MediaStatusInterrupt uint16 = 1 << 1
)

// Interrupt submessage IDs.
const (
    MediaStatus   uint16 = 0x0001
    ReadComplete  uint16 = 0x0002
    WriteComplete uint16 = 0x0003
)

// Error IDs.
const (
    ErrorNone          uint16 = 0x0000
    ErrorNoMedia       uint16 = 0x0001
    ErrorInvalidSector uint16 = 0x0002
    ErrorPending       uint16 = 0x0003
    ErrorWriteLocked   uint16 = 0x0004
)

type Disk interface {
    SectorSize() uint16
    NumSectors() uint16
    WriteLocked() bool
    ReadSector(uint16, []uint16)
    WriteSector(uint16, []uint16)
}

type HMD2043 struct {
    Mutex            sync.Mutex
    Em               *dcpuem.Emulator
    Disk             Disk
    Flags            uint16
    LastInterrupt    uint16
    InterruptMessage uint16
}

func New() (d *HMD2043) {
    d = new(HMD2043)
    d.InterruptMessage = 0xFFFF
    return d
}

func (d *HMD2043) Insert(disk Disk) {
    d.Disk = disk
}

func (d *HMD2043) Eject() (disk Disk) {
    disk = d.Disk
    d.Disk = nil
    return disk
}

func (d *HMD2043) AssociateEmulator(em *dcpuem.Emulator) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    d.Em = em
}

func (d *HMD2043) ID() (id uint32) {
    return 0x74fa4cae
}

func (d *HMD2043) Version() (ver uint16) {
    return 0x07c2
}

func (d *HMD2043) Manufacturer() (manu uint32) {
    return 0x21544948
}

func (d *HMD2043) Start() {

}

func (d *HMD2043) Stop() {

}

func (d *HMD2043) Interrupt() (err error) {
    d.Mutex.Lock()
    defer d.Mutex.Unlock()

    a := d.Em.Regs[dcpuem.A]
    b := d.Em.Regs[dcpuem.B]
    c := d.Em.Regs[dcpuem.C]
    x := d.Em.Regs[dcpuem.X]

    errcode := ErrorNone

    switch a {
    case 0x0000: // QUERY_MEDIA_PRESENT
        if d.Disk == nil {
            d.Em.Regs[dcpuem.B] = 0
        } else {
            d.Em.Regs[dcpuem.B] = 1
        }

    case 0x0001: // QUERY_MEDIA_PARAMETERS
        if d.Disk == nil {
            errcode = ErrorNoMedia

        } else {
            d.Em.Regs[dcpuem.B] = d.Disk.SectorSize()
            d.Em.Regs[dcpuem.C] = d.Disk.NumSectors()

            if d.Disk.WriteLocked() {
                d.Em.Regs[dcpuem.X] = 1
            } else {
                d.Em.Regs[dcpuem.X] = 0
            }
        }

    case 0x0002: // QUERY_DEVICE_FLAGS
        d.Em.Regs[dcpuem.B] = d.Flags

    case 0x0003: // UPDATE_DEVICE_FLAGS
        d.Flags = b

    case 0x0004: // QUERY_INTERRUPT_TYPE
        d.Em.Regs[dcpuem.B] = d.LastInterrupt

    case 0x0005: // SET_INTERRUPT_MESSAGE
        d.InterruptMessage = b

    case 0x0010: // READ_SECTORS
        if d.Disk == nil {
            errcode = ErrorNoMedia

        } else if b+c > d.Disk.NumSectors() {
            errcode = ErrorInvalidSector

        } else {
            if d.Flags&NonBlocking != 0 {
                d.ReadSectorsAsync(b, c, x)
            } else {
                d.ReadSectors(b, c, x)
            }
        }

    case 0x0011: // WRITE_SECTORS
        if d.Disk == nil {
            errcode = ErrorNoMedia

        } else if d.Disk.WriteLocked() {
            errcode = ErrorWriteLocked

        } else if b+c > d.Disk.NumSectors() {
            errcode = ErrorInvalidSector

        } else {
            if d.Flags&NonBlocking != 0 {
                d.WriteSectorsAsync(b, c, x)
            } else {
                d.WriteSectors(b, c, x)
            }
        }

    case 0xFFFF: // QUERY_MEDIA_QUALITY
        _, ok := d.Disk.(*HMU1440)

        if ok {
            d.Em.Regs[dcpuem.B] = 0x7FFF
        } else {
            d.Em.Regs[dcpuem.B] = 0xFFFF
        }
    }

    d.Em.Regs[dcpuem.A] = errcode

    return nil
}

func (d *HMD2043) ReadSectors(startSector uint16, numSectors uint16, addr uint16) {
    // 48 kW per second
    // 512 bytes per sector (generally)
    // = 93.75 sectors per second
    // Therefore we must copy sectors at 93.75 Hz

    ticker := time.NewTicker(time.Second / 94)
    sectorSize := d.Disk.SectorSize()

    for i := uint16(0); i < numSectors; i++ {
        d.Disk.ReadSector(startSector+i, d.Em.RAM[addr:addr+sectorSize])
        addr += sectorSize

        <-ticker.C
    }
}

func (d *HMD2043) ReadSectorsAsync(startSector uint16, numSectors uint16, startAddr uint16) {
    go func() {
        d.ReadSectors(startSector, numSectors, startAddr)

        d.LastInterrupt = ReadComplete
        d.Em.Interrupt(d.InterruptMessage)
    }()
}

func (d *HMD2043) WriteSectors(startSector uint16, numSectors uint16, addr uint16) {
    // 48 kW per second
    // 512 bytes per sector (generally)
    // = 93.75 sectors per second
    // Therefore we must copy sectors at 93.75 Hz

    ticker := time.NewTicker(time.Second / 93)
    sectorSize := d.Disk.SectorSize()

    for i := uint16(0); i < numSectors; i++ {
        d.Disk.WriteSector(startSector+i, d.Em.RAM[addr:addr+sectorSize])
        addr += sectorSize

        <-ticker.C
    }
}

func (d *HMD2043) WriteSectorsAsync(startSector uint16, numSectors uint16, startAddr uint16) {
    go func() {
        d.WriteSectors(startSector, numSectors, startAddr)

        d.LastInterrupt = WriteComplete
        d.Em.Interrupt(d.InterruptMessage)
    }()
}
