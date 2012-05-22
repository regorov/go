package k750emlib

func sizedLoad(em *Emulator, size MemorySize, addr uint32) (v uint32) {
    switch size {
    case Mem8:
        return uint32(em.MemoryLoad8(addr))

    case Mem16:
        return uint32(em.MemoryLoad16(addr))

    case Mem32:
        return em.MemoryLoad32(addr)
    }

    return 0
}

func sizedStore(em *Emulator, size MemorySize, addr uint32, v uint32) {
    switch size {
    case Mem8:
        em.MemoryStore8(addr, uint8(v))

    case Mem16:
        em.MemoryStore16(addr, uint16(v))

    case Mem32:
        em.MemoryStore32(addr, v)
    }
}

type Operand interface {
    Load(*Emulator) (uint32, error)
    Store(*Emulator, uint32) error
}

type LiteralOperand struct {
    Value uint32
}

func (o *LiteralOperand) Load(em *Emulator) (v uint32, err error) {
    return o.Value, nil
}

func (o *LiteralOperand) Store(em *Emulator, v uint32) (err error) {
    return &Error{ErrStoreToLiteral, "Attempt to store to a literal operand"}
}

type RegisterOperand struct {
    Reg Register
}

func (o *RegisterOperand) Load(em *Emulator) (v uint32, err error) {
    return em.Regs[o.Reg], nil
}

func (o *RegisterOperand) Store(em *Emulator, v uint32) (err error) {
    em.Regs[o.Reg] = v
    return nil
}

type MemoryOperand struct {
    Size    MemorySize
    Reg     Register
    Literal uint32
}

func (o *MemoryOperand) Addr(em *Emulator) (addr uint32) {
    addr = o.Literal
    if o.Reg != NoRegister {
        addr += em.Regs[o.Reg]
    }

    return addr
}

func (o *MemoryOperand) Load(em *Emulator) (v uint32, err error) {
    return sizedLoad(em, o.Size, o.Addr(em)), nil
}

func (o *MemoryOperand) Store(em *Emulator, v uint32) (err error) {
    sizedStore(em, o.Size, o.Addr(em), v)
    return nil
}

type ArrayOperand struct {
    Size  MemorySize
    Base  Register
    Index Register
    Scale Scale
    Disp  uint32
}

func (o *ArrayOperand) Addr(em *Emulator) (addr uint32) {
    base := em.Regs[o.Base]
    index := em.Regs[o.Index]

    return base + (index << (o.Scale + 1)) + o.Disp
}

func (o *ArrayOperand) Load(em *Emulator) (v uint32, err error) {
    return sizedLoad(em, o.Size, o.Addr(em)), nil
}

func (o *ArrayOperand) Store(em *Emulator, v uint32) (err error) {
    sizedStore(em, o.Size, o.Addr(em), v)
    return nil
}

type PCOperand struct {
}

func (o *PCOperand) Load(em *Emulator) (v uint32, err error) {
    return em.PC, nil
}

func (o *PCOperand) Store(em *Emulator, v uint32) (err error) {
    em.PC = v
    return nil
}

type SROperand struct {
}

func (o *SROperand) Load(em *Emulator) (v uint32, err error) {
    return em.SR, nil
}

func (o *SROperand) Store(em *Emulator, v uint32) (err error) {
    em.SR = v
    return nil
}
