package k920asmlib

import (
	. "github.com/kierdavis/go/k920"
	"log"
	"reflect"
)

type Pipe chan Object

func verifyInstruction(inst *Instruction, types string) (ok bool) {
	if len(inst.Operands) < len(types) {
		log.Printf("Not enough operands for %s.%s: expected %d, got %d", inst.Group, inst.Name, len(types), len(inst.Operands))
		return false
	}

	if len(inst.Operands) > len(types) {
		log.Printf("Too many operands for %s.%s: expected %d, got %d", inst.Group, inst.Name, len(types), len(inst.Operands))
		return false
	}

	for i, t := range types {
		operand := inst.Operands[i]
		operandType := reflect.TypeOf(operand)

		switch t {
		case 'r':
			_, ok := operand.(Register)
			if !ok {
				log.Printf("Invalid type for operand %d to %s.%s: expected Register, got %s", operandType.Name())
				return false
			}

		case 'b':
			_, ok := operand.(BitReg)
			if !ok {
				log.Printf("Invalid type for operand %d to %s.%s: expected BitReg, got %s", operandType.Name())
				return false
			}

		case 'i':
			_, ok := operand.(Integer)
			if !ok {
				log.Printf("Invalid type for operand %d to %s.%s: expected Integer, got %s", operandType.Name())
				return false
			}

		case 'l':
			_, ok1 := operand.(Integer)
			_, ok2 := operand.(LabelRef)

			if !(ok1 || ok2) {
				log.Printf("Invalid type for operand %d to %s.%s: expected Integer or LabelRef, got %s", operandType.Name())
				return false
			}
		}
	}

	return true
}

func verifyStdInstruction(inst *Instruction, out Pipe) {
	switch inst.Name {
	case "nop", "ret", "reti", "pusha", "popa", "swu", "hlt":
		if verifyInstruction(inst, "") {
			out <- inst
		}

	case "jr", "cr", "push", "pop":
		if verifyInstruction(inst, "r") {
			out <- inst
		}

	case "cb", "sb":
		if verifyInstruction(inst, "b") {
			out <- inst
		}

	case "ab", "ob", "xb":
		if verifyInstruction(inst, "bbb") {
			out <- inst
		}

	case "mov", "not", "neg":
		if verifyInstruction(inst, "rr") {
			out <- inst
		}

	case "bld":
		if verifyInstruction(inst, "bri") {
			out <- inst
		}

	case "bst":
		if verifyInstruction(inst, "rib") {
			out <- inst
		}

	case "beq", "bne", "blt", "bge", "bbc", "bbs":
		if verifyInstruction(inst, "brr") {
			out <- inst
		}

	case "jmp", "call", "rjmp", "rcall":
		if verifyInstruction(inst, "l") {
			out <- inst
		}

	case "ldl", "ldu":
		if verifyInstruction(inst, "rl") {
			out <- inst
		}

	case "ldi":
		if verifyInstruction(inst, "rl") {
			i, ok := inst.Operands[1].(Integer)
			if ok && (i < -0x8000 || i >= 0x8000) {
				out <- &Instruction{Group: "std", Name: "ldu", Operands: []Operand{(i >> 16) & 0xFFFF}}
				out <- &Instruction{Group: "std", Name: "ldl", Operands: []Operand{i & 0xFFFF}}

			} else {
				out <- inst
			}
		}

	case "jb", "jnb":
		if verifyInstruction(inst, "bl") {
			out <- inst
		}

	case "ldb", "ldb+", "ldb-", "ldh", "ldh+", "ldh-", "ldw", "ldw+", "ldw-":
		if len(inst.Operands) < 3 {
			inst.Operands = append(inst.Operands, Integer(0))
		}

		if verifyInstruction(inst, "rrl") {
			out <- inst
		}

	case "addi", "adci", "andi", "ori", "xori", "jac", "jas", "jeq", "jne", "jlt", "jge":
		if verifyInstruction(inst, "rrl") {
			out <- inst
		}

	case "jgt", "jle":
		if verifyInstruction(inst, "rrl") {
			inst.Operands[0], inst.Operands[1] = inst.Operands[1], inst.Operands[0]
			out <- inst
		}

	case "jeqi", "jnei", "jlti", "jgti", "jlei", "jgei":
		if verifyInstruction(inst, "ril") {
			out <- inst
		}

	case "stb", "stb+", "stb-", "sth", "sth+", "sth-", "stw", "stw+", "stw-":
		if len(inst.Operands) < 3 {
			x := inst.Operands[1]
			inst.Operands[1] = Integer(0)
			inst.Operands = append(inst.Operands, x)
		}

		if verifyInstruction(inst, "rlr") {
			out <- inst
		}

	case "add", "sub", "adc", "sbc", "and", "or", "xor":
		if verifyInstruction(inst, "rrr") {
			out <- inst
		}

	case "asl", "lsl", "lsr", "rol", "ror":
		if verifyInstruction(inst, "rri") {
			out <- inst
		}

	default:
		log.Printf("Invalid instruction std.%s\n", inst.Name)
	}
}

func verifyMmuInstruction(inst *Instruction, out Pipe) {
	switch inst.Name {
	default:
		log.Printf("Invalid instruction mmu.%s\n", inst.Name)
	}
}

func verifySerialInstruction(inst *Instruction, out Pipe) {
	switch inst.Name {
	case "sel", "send", "int", "sq", "rq", "recv":
		if verifyInstruction(inst, "r") {
			out <- inst
		}

	case "seli", "sendi", "inti":
		if verifyInstruction(inst, "l") {
			out <- inst
		}

	default:
		log.Printf("Invalid instruction serial.%s\n", inst.Name)
	}
}

func Verify(in Pipe) (out Pipe) {
	out = make(Pipe)

	go func() {
		defer close(out)

		for iobj := range in {
			switch obj := iobj.(type) {
			case *Instruction:
				switch obj.Group {
				case "std":
					verifyStdInstruction(obj, out)
				case "mmu":
					verifyMmuInstruction(obj, out)
				case "serial":
					verifySerialInstruction(obj, out)
				default:
					log.Printf("Invalid instruction group %s\n", obj.Group)
				}

			default:
				out <- iobj
			}
		}
	}()

	return out
}

func CalcLabels(in Pipe) (out Pipe, labels map[string]uint32) {
	out = make(Pipe)
	labels = make(map[string]uint32)

	go func() {
		defer close(out)

		var offset uint32

		for iobj := range in {
			iobj.SetOffset(offset)

			switch obj := iobj.(type) {
			case *Label:
				labels[obj.Name] = offset
			}

			offset += iobj.Length()
			out <- iobj
		}
	}()

	return out, labels
}

func EvalLabels(in Pipe, labels map[string]uint32) (out Pipe, labels_ map[string]uint32) {
	out = make(Pipe)

	go func() {
		defer close(out)

		buffer := make([]Object, 0)

		for {
			obj := <-in
			inst, ok := obj.(*Instruction)
			if ok {
				for _, operand := range inst.Operands {
					label, ok := operand.(LabelRef)
					if ok {
						for {
							_, ok := labels[string(label)]
							if ok {
								break
							}

							o := <-in
							if o == nil {
								log.Fatalf("Label not defined: %s", label)
							}

							buffer = append(buffer, o)
						}
					}
				}
			}

			out <- obj

			if len(buffer) > 0 {
				for _, o := range buffer {
					out <- o
				}

				buffer = buffer[:0]
			}
		}
	}()

	return out, labels
}
