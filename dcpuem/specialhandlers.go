// specialhandlers.go - Handlers for special instructions.

package dcpuem

import (
    "fmt"
)

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
    handleINT, // 08
    handleIAG, // 09
    handleIAS, // 0A
    handleRFI, // 0B
    handleIAQ, // 0C
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

func handleINT(em *Emulator, a Operand) (err error) {
    em.Interrupt(em.Load(a))
    em.LogInstruction("INT %s", a.String)
    return nil
}

func handleIAG(em *Emulator, a Operand) (err error) {
    err = em.Store(a, em.IA)
    if err != nil {
        return err
    }

    em.LogInstruction("IAG %s ; value transferred was 0x%04X", a.String, em.IA)
    return nil
}

func handleIAS(em *Emulator, a Operand) (err error) {
    em.IA = em.Load(a)
    em.LogInstruction("IAS %s ; value transferred was 0x%04X", a.String, em.IA)
    return nil
}

func handleRFI(em *Emulator, a Operand) (err error) {
    em.InterruptQueueing = false
    em.Regs[A] = em.Pop()
    em.PC = em.Pop()
    em.LogInstruction("RFI")
    return nil
}

func handleIAQ(em *Emulator, a Operand) (err error) {
    em.InterruptQueueing = em.Load(a) != 0
    em.LogInstruction("IAQ %s ; interrupts now queued: %t", a.String, em.InterruptQueueing)
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

    if uint16(len(em.Hardware)) <= n {
        return &Error{ErrInvalidHardwareIndex, fmt.Sprintf("Hardware index out of range: 0x%04X", n), int(n)}
    }

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

    if uint16(len(em.Hardware)) <= n {
        return &Error{ErrInvalidHardwareIndex, fmt.Sprintf("Hardware index out of range: 0x%04X", n), int(n)}
    }

    hw := em.Hardware[n]
    em.LogInstruction("HWI %s ; interrupting device %d", a.String, n)

    em.Log("Device interrupt handling beginning")
    err = hw.Interrupt()

    if err == nil {
        em.Log("Device interrupt handling completed successfully")
    } else {
        em.Log("Device interrupt handling completed with errors")
    }

    return err
}
