// specialhandlers.go - Handlers for special instructions.

package dcpuem

type specialHandler func(*Emulator, Operand) error

var specialHandlers = [32]specialHandler{
    nil,       // 00
    handleJSR, // 01
    nil,       // 02
    nil,       // 03
    nil,       // 04
    nil,       // 05
    nil,       // 06
    nil,       // 07
    nil,       // 08
    nil,       // 09
    nil,       // 0A
    nil,       // 0B
    nil,       // 0C
    nil,       // 0D
    nil,       // 0E
    nil,       // 0F
    handleHWN, // 10
    handleHWQ, // 11
    handleHWI, // 12
    nil,       // 13
    nil,       // 14
    nil,       // 15
    nil,       // 16
    nil,       // 17
    nil,       // 18
    nil,       // 19
    nil,       // 1A
    nil,       // 1B
    nil,       // 1C
    nil,       // 1D
    nil,       // 1E
    nil,       // 1F
}

func handleJSR(em *Emulator, a Operand) (err error) {
    em.Push(em.PC)
    em.PC = em.Load(a)
    em.LogInstruction("JSR %s ; destination = 0x%04X", a.String, em.PC)
    return nil
}

func handleHWN(em *Emulator, a Operand) (err error) {
    n := uint16(len(em.Hardware))
    err = em.Store(a, n)
    if err != nil {
        return err
    }

    em.LogInstruction("HWN %s ; returned %d", a.String, n)
    return nil
}

func handleHWQ(em *Emulator, a Operand) (err error) {
    n := em.Load(a)
    hw := em.Hardware[n]

    id := hw.ID()
    ver := hw.Version()
    manu := hw.Manufacturer()

    em.Regs[A] = uint16(id)
    em.Regs[B] = uint16(id >> 16)
    em.Regs[C] = ver
    em.Regs[X] = uint16(manu)
    em.Regs[Y] = uint16(manu >> 16)

    em.LogInstruction("HWQ %s ; Info about device %d: 0x%08X 0x%04X 0x%08X", a.String, n, id, ver, manu)
    return nil
}

func handleHWI(em *Emulator, a Operand) (err error) {
    n := em.Load(a)
    hw := em.Hardware[n]
    em.LogInstruction("HWI %s ; interrupting device %d", a.String, n)

    em.Log("Device interrupt handling beginning")
    err = hw.Interrupt(em)

    if err == nil {
        em.Log("Device interrupt handling completed successfully")
    } else {
        em.Log("Device interrupt handling completed with errors")
    }

    return err
}
