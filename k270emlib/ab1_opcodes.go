package k270emlib

import (
    "fmt"
)

var AB1Opcodes = [16]func(*Emulator, int, int){
    HandleIfbc,
    HandleIfbs,
    HandleIfbcp,
    HandleIfbsp,
    HandleIfeq,
    HandleIfne,
    HandleIflt,
    HandleIfge,
    HandleAdd,
    HandleSub,
    HandleAnd,
    HandleOr,
    HandleXor,
    HandleMov,
    HandleAdc,
    HandleSbc,
}

func HandleAB1Opcode(em *Emulator, a int, i int) {
    q := i >> 4
    b := i & 0xF
    
    f := AB1Opcodes[q]
    
    if f == nil {
        panic(fmt.Sprintf("Invalid AB1 opcode 0x%X", q))
    } else {
        f(em, a, b)
    }
}

func HandleIfbc(em *Emulator, a int, b int) {
    v := em.GetReg(a)
    bit := (v >> uint(b)) & 1
    
    if bit == 0 {
        em.LogInstruction("ifbc %s, %d -- bit(0x%02X, %d) = 0, executing next", RegisterNames[a], b, v, b)
    } else {
        em.pc += 2
        em.LogInstruction("ifbc %s, %d -- bit(0x%02X, %d) = 1, skipping next", RegisterNames[a], b, v, b)
    }
}

func HandleIfbs(em *Emulator, a int, b int) {
    v := em.GetReg(a)
    bit := (v >> uint(b)) & 1
    
    if bit == 1 {
        em.LogInstruction("ifbs %s, %d -- bit(0x%02X, %d) = 1, executing next", RegisterNames[a], b, v, b)
    } else {
        em.pc += 2
        em.LogInstruction("ifbs %s, %d -- bit(0x%02X, %d) = 0, skipping next", RegisterNames[a], b, v, b)
    }
}

func HandleIfbcp(em *Emulator, a int, b int) {
    p := em.GetReg(a)
    
    if em.getPortAccess(p) {
        v := em.LoadIOPort(p)
        bit := (v >> uint(b)) & 1
        
        if bit == 0 {
            em.LogInstruction("ifbc %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 0, executing next", RegisterNames[a], b, p, v, v, b)
        } else {
            em.pc += 2
            em.LogInstruction("ifbc %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 1, skipping next", RegisterNames[a], b, p, v, v, b)
        }
    
    } else {
        em.LogInstruction("ifbc %s, %d -- not authorised", RegisterNames[a], b)
    }
}

func HandleIfbsp(em *Emulator, a int, b int) {
    p := em.GetReg(a)
    
    if em.getPortAccess(p) {
        v := em.LoadIOPort(p)
        bit := (v >> uint(b)) & 1
        
        if bit == 1 {
            em.LogInstruction("ifbs %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 1, executing next", RegisterNames[a], b, p, v, v, b)
        } else {
            em.pc += 2
            em.LogInstruction("ifbs %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 0, skipping next", RegisterNames[a], b, p, v, v, b)
        }
    
    } else {
        em.LogInstruction("ifbs %s, %d -- not authorised", RegisterNames[a], b)
    }
}

func HandleIfeq(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value == b_value {
        em.LogInstruction("ifeq %s, %s -- 0x%02X == 0x%02X, executing next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    } else {
        em.pc += 2
        em.LogInstruction("ifeq %s, %s -- 0x%02X != 0x%02X, skipping next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    }
}

func HandleIfne(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value != b_value {
        em.LogInstruction("ifne %s, %s -- 0x%02X != 0x%02X, executing next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    } else {
        em.pc += 2
        em.LogInstruction("ifne %s, %s -- 0x%02X == 0x%02X, skipping next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    }
}

func HandleIflt(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value < b_value {
        em.LogInstruction("iflt %s, %s -- 0x%02X < 0x%02X, executing next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    } else {
        em.pc += 2
        em.LogInstruction("iflt %s, %s -- 0x%02X >= 0x%02X, skipping next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    }
}

func HandleIfge(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value >= b_value {
        em.LogInstruction("ifge %s, %s -- 0x%02X >= 0x%02X, executing next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    } else {
        em.pc += 2
        em.LogInstruction("ifge %s, %s -- 0x%02X < 0x%02X, skipping next", RegisterNames[a], RegisterNames[b], a_value, b_value)
    }
}

func HandleAdd(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    r := a_value + b_value
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("add %s, %s -- 0x%02X + 0x%02X = 0x%02X, carry = %t", RegisterNames[a], RegisterNames[b], a_value, b_value, r, em.GetCarry())
}

func HandleSub(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    r := a_value - b_value
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("sub %s, %s -- 0x%02X - 0x%02X = 0x%02X, carry = %t", RegisterNames[a], RegisterNames[b], a_value, b_value, r, em.GetCarry())
}

func HandleAnd(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    r := a_value & b_value
    em.SetReg(a, r)
    em.LogInstruction("and %s, %s -- 0x%02X & 0x%02X = 0x%02X", RegisterNames[a], RegisterNames[b], a_value, b_value, r)
}

func HandleOr(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    r := a_value | b_value
    em.SetReg(a, r)
    em.LogInstruction("or %s, %s -- 0x%02X | 0x%02X = 0x%02X", RegisterNames[a], RegisterNames[b], a_value, b_value, r)
}

func HandleXor(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    r := a_value ^ b_value
    em.SetReg(a, r)
    em.LogInstruction("xor %s, %s -- 0x%02X ^ 0x%02X = 0x%02X", RegisterNames[a], RegisterNames[b], a_value, b_value, r)
}

func HandleMov(em *Emulator, a int, b int) {
    v := em.GetReg(b)
    em.SetReg(a, v)
    em.LogInstruction("mov %s, %s -- value transferred was 0x%02X", RegisterNames[a], RegisterNames[b], v)
}

func HandleAdc(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value + b_value + c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("adc %s, %s -- 0x%02X + 0x%02X + %d = 0x%02X, carry = %t", RegisterNames[a], RegisterNames[b], a_value, b_value, c, r, em.GetCarry())
}

func HandleSbc(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value - b_value - c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("sbc %s, %s -- 0x%02X - 0x%02X - %d = 0x%02X, carry = %t", RegisterNames[a], RegisterNames[b], a_value, b_value, c, r, em.GetCarry())
}
