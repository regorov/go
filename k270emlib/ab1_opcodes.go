package k270emlib

import (
    "fmt"
)

// Variable AB1Opcodes is an array of functions that handle opcodes in the AB1 class. Indexing into
// this array with a 4-bit number returns a function that will handle that opcode.
var AB1Opcodes = [16]func(*Emulator, int, int){
    HandleIfbc,  // 0000 0
    HandleIfbs,  // 0001 1
    HandleIfbcp, // 0010 2
    HandleIfbsp, // 0011 3
    HandleIfeq,  // 0100 4
    HandleIfne,  // 0101 5
    HandleIflt,  // 0110 6
    HandleIfge,  // 0111 7
    HandleAdd,   // 1000 8
    HandleSub,   // 1001 9
    HandleAnd,   // 1010 10
    HandleOr,    // 1011 11
    HandleXor,   // 1100 12
    HandleMov,   // 1101 13
    HandleAdc,   // 1110 14
    HandleSbc,   // 1111 15
}

// Function HandleAB1Opcode distributes the handling of an AB1 opcode to the appropriate opcode
// handler.
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

// Function HandleIfbc handles an IFBC instruction.
func HandleIfbc(em *Emulator, a int, b int) {
    v := em.GetReg(a)
    bit := (v >> uint(b)) & 1
    
    if bit == 0 {
        em.LogInstruction("ifbc %s, %d -- bit(0x%02X, %d) = 0, executing next", RegisterNames[a], b,
            v, b)
    } else {
        em.pc += 2
        em.LogInstruction("ifbc %s, %d -- bit(0x%02X, %d) = 1, skipping next", RegisterNames[a], b,
            v, b)
    }
}

// Function HandleIfbs handles an IFBS instruction.
func HandleIfbs(em *Emulator, a int, b int) {
    v := em.GetReg(a)
    bit := (v >> uint(b)) & 1
    
    if bit == 1 {
        em.LogInstruction("ifbs %s, %d -- bit(0x%02X, %d) = 1, executing next", RegisterNames[a], b,
            v, b)
    } else {
        em.pc += 2
        em.LogInstruction("ifbs %s, %d -- bit(0x%02X, %d) = 0, skipping next", RegisterNames[a], b,
            v, b)
    }
}

// Function HandleIfbcp handles an IFBCP instruction.
func HandleIfbcp(em *Emulator, a int, b int) {
    p := em.GetReg(a)
    
    if em.getPortAccess(p) {
        v := em.LoadIOPort(p)
        bit := (v >> uint(b)) & 1
        
        if bit == 0 {
            em.LogInstruction(
                "ifbc %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 0, executing next",
                RegisterNames[a], b, p, v, v, b)
        
        } else {
            em.pc += 2
            em.LogInstruction(
                "ifbc %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 1, skipping next",
                RegisterNames[a], b, p, v, v, b)
        }
    
    } else {
        em.LogInstruction("ifbc %s, %d -- not authorised", RegisterNames[a], b)
    }
}

// Function HandleIfbsp handles an IFBSP instruction.
func HandleIfbsp(em *Emulator, a int, b int) {
    p := em.GetReg(a)
    
    if em.getPortAccess(p) {
        v := em.LoadIOPort(p)
        bit := (v >> uint(b)) & 1
        
        if bit == 1 {
            em.LogInstruction(
                "ifbs %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 1, executing next",
                RegisterNames[a], b, p, v, v, b)
        
        } else {
            em.pc += 2
            em.LogInstruction(
                "ifbs %s, %d -- ports[0x%02X] = 0x%02X, bit(0x%02X, %d) = 0, skipping next",
                RegisterNames[a], b, p, v, v, b)
        }
    
    } else {
        em.LogInstruction("ifbs %s, %d -- not authorised", RegisterNames[a], b)
    }
}

// Function HandleIfeq handles an IFEQ instruction.
func HandleIfeq(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value == b_value {
        em.LogInstruction("ifeq %s, %s -- 0x%02X == 0x%02X, executing next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    
    } else {
        em.pc += 2
        em.LogInstruction("ifeq %s, %s -- 0x%02X != 0x%02X, skipping next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    }
}

// Function HandleIfne handles an IFNE instruction.
func HandleIfne(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value != b_value {
        em.LogInstruction("ifne %s, %s -- 0x%02X != 0x%02X, executing next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    
    } else {
        em.pc += 2
        em.LogInstruction("ifne %s, %s -- 0x%02X == 0x%02X, skipping next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    }
}

// Function HandleIflt handles an IFLT instruction.
func HandleIflt(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value < b_value {
        em.LogInstruction("iflt %s, %s -- 0x%02X < 0x%02X, executing next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    
    } else {
        em.pc += 2
        em.LogInstruction("iflt %s, %s -- 0x%02X >= 0x%02X, skipping next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    }
}

// Function HandleIfge handles an IFGE instruction.
func HandleIfge(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    
    if a_value >= b_value {
        em.LogInstruction("ifge %s, %s -- 0x%02X >= 0x%02X, executing next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    
    } else {
        em.pc += 2
        em.LogInstruction("ifge %s, %s -- 0x%02X < 0x%02X, skipping next", RegisterNames[a],
            RegisterNames[b], a_value, b_value)
    }
}

// Function HandleAdd handles an ADD instruction.
func HandleAdd(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    r := a_value + b_value
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("add %s, %s -- 0x%02X + 0x%02X = 0x%02X, carry = %t", RegisterNames[a],
        RegisterNames[b], a_value, b_value, r, em.GetCarry())
}

// Function HandleSub handles a SUB instruction.
func HandleSub(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    r := a_value - b_value
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("sub %s, %s -- 0x%02X - 0x%02X = 0x%02X, carry = %t", RegisterNames[a],
        RegisterNames[b], a_value, b_value, r, em.GetCarry())
}

// Function HandleAnd handles an AND instruction.
func HandleAnd(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    r := a_value & b_value
    em.SetReg(a, r)
    em.LogInstruction("and %s, %s -- 0x%02X & 0x%02X = 0x%02X", RegisterNames[a], RegisterNames[b],
        a_value, b_value, r)
}

// Function HandleOr handles an OR instruction.
func HandleOr(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    r := a_value | b_value
    em.SetReg(a, r)
    em.LogInstruction("or %s, %s -- 0x%02X | 0x%02X = 0x%02X", RegisterNames[a], RegisterNames[b],
        a_value, b_value, r)
}

// Function HandleXor handles an XOR instruction.
func HandleXor(em *Emulator, a int, b int) {
    a_value := em.GetReg(a)
    b_value := em.GetReg(b)
    r := a_value ^ b_value
    em.SetReg(a, r)
    em.LogInstruction("xor %s, %s -- 0x%02X ^ 0x%02X = 0x%02X", RegisterNames[a], RegisterNames[b],
        a_value, b_value, r)
}

// Function HandleMov handles a MOV instruction.
func HandleMov(em *Emulator, a int, b int) {
    v := em.GetReg(b)
    em.SetReg(a, v)
    em.LogInstruction("mov %s, %s -- value transferred was 0x%02X", RegisterNames[a],
        RegisterNames[b], v)
}

// Function HandleAdc handles an ADC instruction.
func HandleAdc(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value + b_value + c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("adc %s, %s -- 0x%02X + 0x%02X + %d = 0x%02X, carry = %t", RegisterNames[a],
        RegisterNames[b], a_value, b_value, c, r, em.GetCarry())
}

// Function HandleSbc handles an SBC instruction.
func HandleSbc(em *Emulator, a int, b int) {
    a_value := int(em.GetReg(a))
    b_value := int(em.GetReg(b))
    c := 0
    if em.GetCarry() {c = 1}
    
    r := a_value - b_value - c
    em.SetCarry(r & 0x100 != 0)
    r &= 0xFF
    em.SetReg(a, uint8(r))
    em.LogInstruction("sbc %s, %s -- 0x%02X - 0x%02X - %d = 0x%02X, carry = %t", RegisterNames[a],
        RegisterNames[b], a_value, b_value, c, r, em.GetCarry())
}
