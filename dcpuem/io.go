// io.go - Interrupt services & hardware interface.

package dcpuem

// Function Interrupt pushes an interrupt request with the specified message onto the interrupt
// queue.
func (em *Emulator) Interrupt(message uint16) {
    em.Interrupts = append(em.Interrupts, message)
}

// Function ServiceInterrupt services at most one interrupt if the queue is not empty. Otherwise, it
// does nothing.
func (em *Emulator) ServiceInterrupt() {
    if len(em.Interrupts) > 0 && !em.InterruptQueueing {
        message := em.Interrupts[0]
        em.Interrupts = em.Interrupts[1:]

        if em.IA == 0 {
            em.Log("Interrupting with message 0x%04X - ignoring (IA is unset)", message)

        } else {
            em.Push(em.PC)
            em.Push(em.Regs[A])

            em.PC = em.IA
            em.Regs[A] = message

            em.InterruptQueueing = true

            em.Log("Interrupt with message 0x%04X - jumping to 0x%04X", message, em.IA)
        }
    }
}

// Function AttachDevice attaches the specified device to the CPU and returns its index.
func (em *Emulator) AttachDevice(device Device) (index int) {
    device.AssociateEmulator(em)
    em.Hardware = append(em.Hardware, device)
    return len(em.Hardware) - 1
}

// Function NumDevices returns the number of attached devices.
func (em *Emulator) NumDevices() (num int) {
    return len(em.Hardware)
}

// Function GetDevice returns the device with the specified index.
func (em *Emulator) GetDevice(index int) (device Device) {
    return em.Hardware[index]
}

// Function DetachDevice detaches the specified device from the CPU.
func (em *Emulator) DetachDevice(device Device) {
    device.AssociateEmulator(nil)

    for i, d := range em.Hardware {
        if d == device {
            em.Hardware = append(em.Hardware[:i], em.Hardware[i+1:]...)
            return
        }
    }
}

// Interface Device represents a hardware device.
type Device interface {
    // Function AssociateEmulator should store a reference to the given emulator object in the
    // device's content structure.
    AssociateEmulator(*Emulator)

    // Function ID should return the device's ID number, for hardware querying.
    ID() uint32

    // Function Version should return the device's version number, for hardware querying.
    Version() uint16

    // Function Manufacturer should return the device's manufacturer ID, for hardware querying.
    Manufacturer() uint32

    // Function Interrupt should handle an interrupt triggered by a HWI instruction.
    Interrupt() error

    // Function Start should start any background services (e.g. event loops) needed by the device,
    // ideally starting goroutines.
    Start()

    // Function Stop should stop all background services needed by the device, ideally stopping the
    // goroutines started by Start.
    Stop()
}
