package hmd2043

type HMU1440 struct {
    sectors     map[uint16][]uint16
    writeLocked bool
}

func NewHMU1440() (disk *HMU1440) {
    disk = new(HMU1440)
    disk.sectors = make(map[uint16][]uint16)
    return disk
}

func (disk *HMU1440) SectorSize() (sectorSize uint16) {
    return 512
}

func (disk *HMU1440) NumSectors() (numSectors uint16) {
    return 1440
}

func (disk *HMU1440) WriteLocked() (writeLocked bool) {
    return disk.writeLocked
}

func (disk *HMU1440) SetWriteLocked(writeLocked bool) {
    disk.writeLocked = writeLocked
}

func (disk *HMU1440) ReadSector(number uint16, buffer []uint16) {
    sector, ok := disk.sectors[number]

    if ok {
        copy(buffer, sector)

    } else {
        for i := 0; i < 512; i++ {
            buffer[i] = 0
        }
    }
}

func (disk *HMU1440) WriteSector(number uint16, buffer []uint16) {
    sector, ok := disk.sectors[number]

    if !ok {
        sector = make([]uint16, 512)
        disk.sectors[number] = sector
    }

    copy(sector, buffer)
}
