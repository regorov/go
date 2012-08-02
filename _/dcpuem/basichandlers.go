// basichandlers.go - Handlers for basic instructions.

package dcpuem

type basicHandler func(*Emulator, Operand, Operand) error

var basicHandlers = [32]basicHandler{
    nil,       // 00
    handleSET, // 01
    handleADD, // 02
    handleSUB, // 03
    handleMUL, // 04
    handleMLI, // 05
    handleDIV, // 06
    handleDVI, // 07
    handleMOD, // 08
    handleMDI, // 09
    handleAND, // 0A
    handleBOR, // 0B
    handleXOR, // 0C
    handleSHR, // 0D
    handleASR, // 0E
    handleSHL, // 0F
    handleIFB, // 10
    handleIFC, // 11
    handleIFE, // 12
    handleIFN, // 13
    handleIFG, // 14
    handleIFA, // 15
    handleIFL, // 16
    handleIFU, // 17
    nil,       // 18
    nil,       // 19
    handleADX, // 1A
    handleSBX, // 1B
    nil,       // 1C
    nil,       // 1D
    handleSTI, // 1E
    handleSTD, // 1F
}

func handleSET(em *Emulator, a Operand, b Operand) (err error) {
    v := em.Load(a)
    err = em.Store(b, v)
    if err != nil {
        return err
    }

    em.LogInstruction("SET %s, %s ; value transferred was 0x%04X", b.String, a.String, v)
    return nil
}

func handleADD(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))
    result := x + y

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16(result >> 16)
    em.LogInstruction("ADD %s, %s ; 0x%04X + 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleSUB(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))
    result := x - y

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16(result >> 16)
    em.LogInstruction("SUB %s, %s ; 0x%04X - 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleMUL(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))
    result := x * y

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16(result >> 16)
    em.LogInstruction("MUL %s, %s ; 0x%04X * 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleMLI(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))
    result := uint32(int32(x) * int32(y))

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16(result >> 16)
    em.LogInstruction("MLI %s, %s ; 0x%04X * 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleDIV(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))

    var result uint32

    if y == 0 {
        result = 0
    } else {
        result = x / y
    }

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16((x << 16) / y)
    em.LogInstruction("DIV %s, %s ; 0x%04X / 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleDVI(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))

    var result uint32

    if y == 0 {
        result = 0
    } else {
        result = uint32(int32(x) / int32(y))
    }

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16(uint32(int32(x<<16) / int32(y)))
    em.LogInstruction("DVI %s, %s ; 0x%04X / 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleMOD(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))

    var result uint32

    if y == 0 {
        result = 0
    } else {
        result = x % y
    }

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.LogInstruction("MOD %s, %s ; 0x%04X %% 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleMDI(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))

    var result uint32

    if y == 0 {
        result = 0
    } else {
        result = uint32(int32(x) % int32(y))
    }

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.LogInstruction("MDI %s, %s ; 0x%04X %% 0x%04X -> 0x%08X", b.String, a.String, x, y, result)
    return nil
}

func handleAND(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)
    result := x & y

    err = em.Store(b, result)
    if err != nil {
        return err
    }

    em.LogInstruction("AND %s, %s ; 0x%04X & 0x%04X -> 0x%04X", b.String, a.String, x, y, result)
    return nil
}

func handleBOR(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)
    result := x | y

    err = em.Store(b, result)
    if err != nil {
        return err
    }

    em.LogInstruction("BOR %s, %s ; 0x%04X | 0x%04X -> 0x%04X", b.String, a.String, x, y, result)
    return nil
}

func handleXOR(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)
    result := x ^ y

    err = em.Store(b, result)
    if err != nil {
        return err
    }

    em.LogInstruction("XOR %s, %s ; 0x%04X ^ 0x%04X -> 0x%04X", b.String, a.String, x, y, result)
    return nil
}

func handleSHR(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)
    result := x >> y

    err = em.Store(b, result)
    if err != nil {
        return err
    }

    em.EX = uint16((uint32(x) << 16) >> y)
    em.LogInstruction("SHR %s, %s ; 0x%04X >> 0x%04X -> 0x%04X, EX = 0x%04X", b.String, a.String, x, y, result, em.EX)
    return nil
}

func handleASR(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)
    result := uint16(int16(x) >> y)
    // If the operand is signed, an arith shift is used instead of logical.

    err = em.Store(b, result)
    if err != nil {
        return err
    }

    em.EX = uint16(int32(uint32(x)<<16) >> y)
    em.LogInstruction("ASR %s, %s ; 0x%04X >>> 0x%04X -> 0x%04X, EX = 0x%04X", b.String, a.String, x, y, result, em.EX)
    return nil
}

func handleSHL(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)
    result := x << y

    err = em.Store(b, result)
    if err != nil {
        return err
    }

    em.EX = uint16(uint32(x<<y) >> 16)
    em.LogInstruction("SHL %s, %s ; 0x%04X << 0x%04X -> 0x%04X, EX = 0x%04X", b.String, a.String, x, y, result, em.EX)
    return nil
}

func handleIFB(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if (x & y) != 0 {
        em.LogInstruction("IFB %s, %s ; test (0x%04X & 0x%04X) != 0 passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFB %s, %s ; test (0x%04X & 0x%04X) != 0 failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleIFC(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if (x & y) == 0 {
        em.LogInstruction("IFC %s, %s ; test (0x%04X & 0x%04X) == 0 passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFC %s, %s ; test (0x%04X & 0x%04X) == 0 failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleIFE(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if x == y {
        em.LogInstruction("IFE %s, %s ; test 0x%04X == 0x%04X passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFE %s, %s ; test 0x%04X == 0x%04X failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleIFN(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if x != y {
        em.LogInstruction("IFN %s, %s ; test 0x%04X != 0x%04X passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFN %s, %s ; test 0x%04X != 0x%04X failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleIFG(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if x > y {
        em.LogInstruction("IFG %s, %s ; test 0x%04X > 0x%04X passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFG %s, %s ; test 0x%04X > 0x%04X failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleIFA(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if int16(x) > int16(y) {
        em.LogInstruction("IFA %s, %s ; test 0x%04X > 0x%04X (signed) passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFA %s, %s ; test 0x%04X > 0x%04X (signed) failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleIFL(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if x < y {
        em.LogInstruction("IFL %s, %s ; test 0x%04X < 0x%04X passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFL %s, %s ; test 0x%04X < 0x%04X failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleIFU(em *Emulator, a Operand, b Operand) (err error) {
    x := em.Load(b)
    y := em.Load(a)

    if int16(x) < int16(y) {
        em.LogInstruction("IFU %s, %s ; test 0x%04X < 0x%04X (signed) != 0 passed, continuing execution", b.String, a.String, x, y)
    } else {
        em.Skip = true
        em.LogInstruction("IFU %s, %s ; test 0x%04X < 0x%04X (signed) != 0 failed, skipping", b.String, a.String, x, y)
    }

    return nil
}

func handleADX(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))
    ex := uint32(em.EX)
    result := x + y + ex

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16(result >> 16)
    em.LogInstruction("ADX %s, %s ; 0x%04X + 0x%04X + 0x%04X -> 0x%08X", b.String, a.String, x, y, ex, result)
    return nil
}

func handleSBX(em *Emulator, a Operand, b Operand) (err error) {
    x := uint32(em.Load(b))
    y := uint32(em.Load(a))
    ex := uint32(em.EX)
    result := x - y + ex

    err = em.Store(b, uint16(result))
    if err != nil {
        return err
    }

    em.EX = uint16(result >> 16)
    em.LogInstruction("SBX %s, %s ; 0x%04X - 0x%04X + 0x%04X -> 0x%08X", b.String, a.String, x, y, ex, result)
    return nil
}

func handleSTI(em *Emulator, a Operand, b Operand) (err error) {
    v := em.Load(a)
    err = em.Store(b, v)
    if err != nil {
        return err
    }

    em.Regs[I]++
    em.Regs[J]++

    em.LogInstruction("STI %s, %s ; value transferred was 0x%04X, I = 0x%04X, J = 0x%04X", b.String, a.String, v, em.Regs[I], em.Regs[J])
    return nil
}

func handleSTD(em *Emulator, a Operand, b Operand) (err error) {
    v := em.Load(a)
    err = em.Store(b, v)
    if err != nil {
        return err
    }

    em.Regs[I]--
    em.Regs[J]--

    em.LogInstruction("STD %s, %s ; value transferred was 0x%04X, I = 0x%04X, J = 0x%04X", b.String, a.String, v, em.Regs[I], em.Regs[J])
    return nil
}
