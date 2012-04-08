package main

import (
    "fmt"
    "github.com/kierdavis/go/ihex"
    "log"
    "strings"
)

func calcOffsets(asm Node) (labels map[string]int) {
    offset := 0
    labels = make(map[string]int)
    
    for _, node := asm.GetChildNodes() {
        inst, ok := node.(*Instruction)
        
        if ok {
            inst.Name = strings.ToLower(inst.Name)
            length := 0
            
            if inst.Name == "jsr" {
                if len(inst.Operands) < 1 {
                    error(fmt.Sprintf("Not enough operands for %s (expected 1, got %d)", inst.Name, len(inst.Operands)))
                } else if len(inst.Operands) > 1 {
                    error(fmt.Sprintf("Too many operands for %s (expected 1, got %d)", inst.Name, len(inst.Operands)))
                }
                
                length = 1 + inst.Operands[0].(*Operand).GetSize()
            
            } else if inst.Name == "set"
                 || inst.Name == "add"
                 || inst.Name == "sub"
                 || inst.Name == "mul"
                 || inst.Name == "div"
                 || inst.Name == "mod"
                 || inst.Name == "shl"
                 || inst.Name == "shr"
                 || inst.Name == "and"
                 || inst.Name == "bor"
                 || inst.Name == "xor"
                 || inst.Name == "ife"
                 || inst.Name == "ifn"
                 || inst.Name == "ifg"
                 || inst.Name == "ifb" {
                
                if len(inst.Operands) < 2 {
                    error(fmt.Sprintf("Not enough operands for %s (expected 2, got %d)", inst.Name, len(inst.Operands)))
                } else if len(inst.Operands) > 2 {
                    error(fmt.Sprintf("Too many operands for %s (expected 2, got %d)", inst.Name, len(inst.Operands)))
                }
                
                length = 1 + inst.Operands[0].(*Operand).GetSize() + inst.Operands[1].(*Operand).GetSize()
            
            } else if inst.Name == "word" {
                length = len(inst.Operands)
            
            } else if inst.Name == "bss" {
                if len(inst.Operands) < 1 {
                    error(fmt.Sprintf("Not enough operands for %s (expected 1, got %d)", inst.Name, len(inst.Operands)))
                } else if len(inst.Operands) > 1 {
                    error(fmt.Sprintf("Too many operands for %s (expected 1, got %d)", inst.Name, len(inst.Operands)))
                }
                
                length = evaluateOperandPart(inst.Operands[0].X)
            
            } else {
                error(fmt.Sprintf("Invalid instruction: %s", inst.Name))
            }
            
            inst.Length = length
            inst.Offset = offset
            offset += length
        }
        
        label, ok := node.(*Label)
        
        if ok {
            labels[inst.Name] = offset
        }
    }
    
    return labels
}

func evaluateOperandPart(part Node, labels map[string]int) (v int) {
    if part == nil {return 0}
    
    constant, ok := part.(*Constant)
    if ok {return constant.Value}
    
    labelref, ok := part.(*LabelRef)
    if ok {return labels[labelref.Name]}
    
    panic("Misplaced node in operand part field: " + part.GetName())
}

func evaluateOperand(operand *Operand, labels map[string]int) (f int, x int, y int) {
    return operand.Format, evaluateOperandPart(operand.X), evaluateOperandPart(operand.Y)
}

func encodeOperand(operand *Operand, labels map[string]int) (v int, hasExtraWord bool, extraWord int) {
    switch operand.Format {
    case O_REG:
        return evaluateOperandPart(operand.X), false, 0
    case O_MEM:
        return 0x08 | evaluateOperandPart(operand.X), false, 0
    case O_MEMDISP:
        return 0x10 | evaluateOperandPart(operand.X), true, evaluateOperandPart(operand.Y)
    case O_POP:
        return 0x18, false, 0
    case O_PEEK:
        return 0x19, false, 0
    case O_PUSH:
        return 0x1A, false, 0
    case O_SP:
        return 0x1B, false, 0
    case O_PC:
        return 0x1C, false, 0
    case O_O:
        return 0x1D, false, 0
    case O_MEMIMM:
        return 0x1E, true, evaluateOperandPart(operand.X)
    case O_IMM:
        v := evaluateOperandPart(operand.Y)
        if v < 0x20 {
            return 0x20 | v
        } else {
            return 0x1F, true, v
        }
    }
}

func encode(asm Node, labels map[string]int, ix *ihex.IHex) {
    for _, node := asm.GetChildNodes() {
        inst, ok := node.(*Instruction)
        
        if ok {
            for _, opnode := inst.Operands {
                operand, ok := node.(*Operand)
                if !ok {panic("Non-Operand in Operands field")}
                
                v, hasExtraWord, extraWord := encodeOperand(operand, labels)
            }
        }
    }
}