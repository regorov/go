package k750emlib

func (em *Emulator) doMov() (err error) {
    var a, b Operand
    em.loadOperands(&a, &b)

    x, err := b.Load(em)
    if err != nil {
        return err
    }

    err = a.Store(em, x)
    if err != nil {
        return err
    }

    return nil
}

func (em *Emulator) doNot() (err error) {
    var a, b Operand
    em.loadOperands(&a, &b)

    x, err := b.Load(em)
    if err != nil {
        return err
    }

    err = a.Store(em, ^x)
    if err != nil {
        return err
    }

    return nil
}

func (em *Emulator) doNeg() (err error) {
    var a, b Operand
    em.loadOperands(&a, &b)

    x, err := b.Load(em)
    if err != nil {
        return err
    }

    err = a.Store(em, (^x)+1)
    if err != nil {
        return err
    }

    return nil
}

func (em *Emulator) doPush() (err error) {
    var a Operand
    em.loadOperands(&a)

    x, err := a.Load(em)
    if err != nil {
        return err
    }

    em.Push(x)
    return nil
}

func (em *Emulator) doPop() (err error) {
    var a Operand
    em.loadOperands(&a)

    x := em.Pop()

    err = a.Store(em, x)
    if err != nil {
        return err
    }

    return nil
}

func (em *Emulator) doPusha() (err error) {
    em.SC = (em.SC - 1) & 0x0F

    if em.SC == 14 { // SP
        em.Push(0)
    } else {
        em.Push(em.Regs[em.SC])
    }

    em.SetBit(15, em.SC == 0)
    return nil
}

func (em *Emulator) doPopa() (err error) {
    if em.SC == 14 { // SP
        em.Pop()
    } else {
        em.Regs[em.SC] = em.Pop()
    }

    em.SC = (em.SC + 1) & 0x0F
    em.SetBit(15, em.SC == 0)
    return nil
}

func (em *Emulator) doAb() (err error) {
    xy := em.Fetch8()
    y := xy & 0xF
    x := (xy >> 4) & 0xF

    em.SetBit(x, em.GetBit(x) && em.GetBit(y))
    return nil
}

func (em *Emulator) doOb() (err error) {
    xy := em.Fetch8()
    y := xy & 0xF
    x := (xy >> 4) & 0xF

    em.SetBit(x, em.GetBit(x) || em.GetBit(y))
    return nil
}

func (em *Emulator) doXb() (err error) {
    xy := em.Fetch8()
    y := xy & 0xF
    x := (xy >> 4) & 0xF

    em.SetBit(x, xor(em.GetBit(x), em.GetBit(y)))

    return nil
}

func (em *Emulator) doPpn() (err error) {
    var a Operand
    em.loadOperands(&a)

    err = a.Store(em, uint32(len(em.Peripherals)))
    if err != nil {
        return err
    }

    return nil
}
