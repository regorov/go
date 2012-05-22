package k750emlib

import (
    "fmt"
)

type Peripheral interface {
    ReadRegister(uint8) uint32
    WriteRegister(uint8, uint32)
    Interrupt(uint8)
    GetPendingInterruptsChannel() chan uint8
    Start()
    Stop()
}

type Emulator struct {
    PC          uint32
    SR          uint32
    SC          uint8
    Regs        [16]uint32
    Memory      []byte
    Peripherals []Peripheral
}

func NewEmulator() (em *Emulator) {
    em = new(Emulator)
    em.Memory = make([]byte, 1024)
    return em
}

func (em *Emulator) GrowMemory(newsize uint32) {
    m := make([]byte, newsize)
    copy(m, em.Memory)
    em.Memory = m
}

func (em *Emulator) MemoryLoad8(address uint32) (value uint8) {
    if address >= uint32(len(em.Memory)) {
        return 0
    }
    return em.Memory[address]
}

func (em *Emulator) MemoryLoad16(address uint32) (value uint16) {
    value = uint16(em.MemoryLoad8(address)) << 8
    value |= uint16(em.MemoryLoad8(address + 1))
    return value
}

func (em *Emulator) MemoryLoad32(address uint32) (value uint32) {
    value = uint32(em.MemoryLoad16(address)) << 16
    value |= uint32(em.MemoryLoad16(address + 2))
    return value
}

func (em *Emulator) MemoryStore8(address uint32, value uint8) {
    if address >= uint32(len(em.Memory)) {
        newsize := uint32(len(em.Memory)) + 1

        for address >= newsize {
            newsize *= 2
        }

        em.GrowMemory(newsize)
    }

    em.Memory[address] = value
}

func (em *Emulator) MemoryStore16(address uint32, value uint16) {
    em.MemoryStore8(address, uint8(value>>8))
    em.MemoryStore8(address+1, uint8(value))
}

func (em *Emulator) MemoryStore32(address uint32, value uint32) {
    em.MemoryStore16(address, uint16(value>>16))
    em.MemoryStore16(address+2, uint16(value))
}

func (em *Emulator) Fetch8() (value uint8) {
    value = em.MemoryLoad8(em.PC)
    em.PC += 1
    return value
}

func (em *Emulator) Fetch16() (value uint16) {
    value = em.MemoryLoad16(em.PC)
    em.PC += 2
    return value
}

func (em *Emulator) Fetch32() (value uint32) {
    value = em.MemoryLoad32(em.PC)
    em.PC += 4
    return value
}

func (em *Emulator) LoadOperand(key uint8) (operand Operand) {
    switch {
    case (key & 0x80) == 0x00:
        v := uint32(key & 0x7F)
        if v&0x40 != 0 {
            v |= 0xFFFFFF80
        }

        return &LiteralOperand{v}

    case (key & 0xF0) == 0x80:
        return &RegisterOperand{Register(key & 0x0F)}

    case (key & 0xF0) == 0x90:
        return &MemoryOperand{Mem8, Register(key & 0x0F), 0}

    case (key & 0xF0) == 0xA0:
        return &MemoryOperand{Mem16, Register(key & 0x0F), 0}

    case (key & 0xF0) == 0xB0:
        return &MemoryOperand{Mem32, Register(key & 0x0F), 0}

    case (key & 0xF0) == 0xC0:
        d := uint32(int16(em.Fetch16()))
        return &MemoryOperand{Mem8, Register(key & 0x0F), d}

    case (key & 0xF0) == 0xD0:
        d := uint32(int16(em.Fetch16()))
        return &MemoryOperand{Mem16, Register(key & 0x0F), d}

    case (key & 0xF0) == 0xE0:
        d := uint32(int16(em.Fetch16()))
        return &MemoryOperand{Mem32, Register(key & 0x0F), d}

    case (key & 0xFC) == 0xF0:
        bi := em.Fetch8()
        d := uint32(int16(em.Fetch16()))
        b := Register((bi >> 4) & 0x0F)
        i := Register(bi & 0xF)
        return &ArrayOperand{Mem8, b, i, Scale(key & 0x03), d}

    case (key & 0xFC) == 0xF4:
        bi := em.Fetch8()
        d := uint32(int16(em.Fetch16()))
        b := Register((bi >> 4) & 0x0F)
        i := Register(bi & 0xF)
        return &ArrayOperand{Mem16, b, i, Scale(key & 0x03), d}

    case (key & 0xFC) == 0xF8:
        bi := em.Fetch8()
        d := uint32(int16(em.Fetch16()))
        b := Register((bi >> 4) & 0x0F)
        i := Register(bi & 0xF)
        return &ArrayOperand{Mem32, b, i, Scale(key & 0x03), d}

    case key == 0xFC:
        return &PCOperand{}

    case key == 0xFD:
        return &SROperand{}

    case key == 0xFE:
        i := em.Fetch32()
        return &MemoryOperand{Mem32, NoRegister, i}

    case key == 0xFF:
        i := em.Fetch32()
        return &LiteralOperand{i}
    }

    return nil
}

func (em *Emulator) RunOne() (err error) {
    inst := em.Fetch8()

    switch inst {
    case 0x00:
        //return em.doNop()
    case 0x01:
        return em.doMov()
    case 0x02:
        return em.doNot()
    case 0x03:
        return em.doNeg()
    case 0x04:
        return em.doPush()
    case 0x05:
        return em.doPop()
    case 0x06:
        return em.doPusha()
    case 0x07:
        return em.doPopa()
    case 0x08:
        return em.doAb()
    case 0x09:
        return em.doOb()
    case 0x0A:
        return em.doXb()
    case 0x0B:
        return em.doPpn()
    case 0x0C:
        //return em.doInt()
    case 0x0D:
        //return em.doRih()
    case 0x0E:
        //return em.doJbc()
    case 0x0F:
        //return em.doJbs()
    case 0x10:
        //return em.doAdd()
    case 0x11:
        //return em.doSub()
    case 0x12:
        //return em.doAnd()
    case 0x13:
        //return em.doOr()
    case 0x14:
        //return em.doXor()
    case 0x15:
        //return em.doPpi()
    case 0x16:
        //return em.doPpr()
    case 0x17:
        //return em.doPpw()
    case 0x18:
        //return em.doJeq()
    case 0x19:
        //return em.doJne()
    case 0x1A:
        //return em.doJlt()
    case 0x1B:
        //return em.doJge()
    case 0x1C:
        //return em.doJac()
    case 0x1D:
        //return em.doJas()
    case 0x1E:
        //return em.doCall()
    case 0x1F:
        //return em.doReti()
    }

    switch {
    case inst&0xF0 == 0x20:
        //return em.doCb(inst & 0x0F)
    case inst&0xF0 == 0x30:
        //return em.doSb(inst & 0x0F)
    case inst&0xE0 == 0x40:
        //return em.doRll(inst & 0x1F)
    case inst&0xE0 == 0x60:
        //return em.doRlr(inst & 0x1F)
    case inst&0xE0 == 0x80:
        //return em.doShl(inst & 0x1F)
    case inst&0xE0 == 0xA0:
        //return em.doLshr(inst & 0x1F)
    case inst&0xE0 == 0xC0:
        //return em.doAshr(inst & 0x1F)
    case inst&0xFE == 0xE0:
        //return em.doLdb(inst & 0x01)
    case inst&0xFE == 0xE2:
        //return em.doStb(inst & 0x01)
    }

    return &Error{ErrInvalidOpcode, fmt.Sprintf("Invalid opcode 0x%02X", inst)}
}

func (em *Emulator) Push(v uint32) {
    em.Regs[SP] -= 4
    em.MemoryStore32(em.Regs[SP], v)
}

func (em *Emulator) Pop() (v uint32) {
    v = em.MemoryLoad32(em.Regs[SP])
    em.Regs[SP] += 4
    return v
}

func (em *Emulator) GetBit(bit uint8) (v bool) {
    if em.SR&(1<<bit) == 0 {
        return false
    }

    return true
}

func (em *Emulator) SetBit(bit uint8, v bool) {
    if v {
        em.SR |= (1 << bit)
    } else {
        em.SR &= ^(1 << bit)
    }
}

func (em *Emulator) loadOperands(operands ...*Operand) {
    keys := make([]uint8, len(operands))

    for i := 0; i < len(operands); i++ {
        keys[i] = em.Fetch8()
    }

    for i, key := range keys {
        *(operands[i]) = em.LoadOperand(key)
    }
}
