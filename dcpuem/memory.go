// memory.go - The dynamically-expanding memory implementation.

package dcpuem

// Function GrowMemory requests that the underlying size of the RAM be increased.
func (em *Emulator) GrowMemory(newsize int) {
    if newsize == 0 {
        newsize = (cap(em.RAM) + 1) * 2
    }

    if newsize > 0x10000 {
        newsize = 0x10000
    }

    m := make([]uint16, newsize)
    copy(m, em.RAM)
    em.RAM = m
}

// Function LoadProgram loads a slice of words into the emulator's RAM starting at the specified
// address.
func (em *Emulator) LoadProgram(program []uint16, address uint16) {
    end := int(address) + len(program)

    if len(em.RAM) < end {
        em.GrowMemory(end)
    }

    copy(em.RAM[address:], program)
}

// Function LoadProgramBytesBE loads a slice of bytes into the emulator's RAM, interpreting each
// pair of bytes as a big-endian word.
func (em *Emulator) LoadProgramBytesBE(program []byte, address uint16) {
    end := int(address) + (len(program) * 2)

    if len(em.RAM) < end {
        em.GrowMemory(end)
    }

    for i := 0; i < len(program)/2; i++ {
        high := uint16(program[i*2])
        low := uint16(program[(i*2)+1])
        em.RAM[address+uint16(i)] = (high << 8) | low
    }
}

// Function LoadProgramBytesBE loads a slice of bytes into the emulator's RAM, interpreting each
// pair of bytes as a little-endian word.
func (em *Emulator) LoadProgramBytesLE(program []byte, address uint16) {
    end := int(address) + (len(program) * 2)

    if len(em.RAM) < end {
        em.GrowMemory(end)
    }

    for i := 0; i < len(program)/2; i++ {
        low := uint16(program[i*2])
        high := uint16(program[(i*2)+1])
        em.RAM[address+uint16(i)] = (high << 8) | low
    }
}

// Function MemoryLoad returns the value in the emulator's RAM at the specified address, or 0 if it
// is greater that the size of the underlying storage.
func (em *Emulator) MemoryLoad(address uint16) (value uint16) {
    if int(address) >= len(em.RAM) {
        return 0
    } else {
        return em.RAM[address]
    }

    return 0
}

// Function MemoryStore stores a value into the emulator's RAM at the specified address, calling
// GrowMemory if needed.
func (em *Emulator) MemoryStore(address uint16, value uint16) {
    if int(address) >= len(em.RAM) {
        newsize := cap(em.RAM) + 1

        for int(address) >= newsize {
            newsize *= 2
        }

        em.GrowMemory(newsize)
    }

    em.RAM[address] = value
}

// Function Push stores the value into the emulator's RAM at the address specified by the stack
// pointer, then increments the stack pointer.
func (em *Emulator) Push(value uint16) {
    em.SP--
    em.MemoryStore(em.SP, value)
}

// Function Pop decrements the stack pointer, the loads and returns the value in the emulator's RAM
// at the address specified by the stack pointer.
func (em *Emulator) Pop() (value uint16) {
    value = em.MemoryLoad(em.SP - 1)
    em.SP++
    return value
}

// Function FetchWord fetches the next program word from the emulator's RAM.
func (em *Emulator) FetchWord() (word uint16) {
    word = em.MemoryLoad(em.PC)
    em.PC++
    return word
}
