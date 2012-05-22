package k750emlib

type MemorySize uint8

const (
    Mem8  MemorySize = 8
    Mem16 MemorySize = 16
    Mem32 MemorySize = 32
)

type Scale uint8

const (
    Scale2  Scale = 0
    Scale4  Scale = 1
    Scale8  Scale = 2
    Scale16 Scale = 3
)

type Register uint8

const (
    V0 Register = iota
    V1
    V2
    V3
    V4
    V5
    V6
    V7
    A0
    A1
    A2
    A3
    Q0
    Q1
    SP
    AT

    NoRegister Register = 0xFF
)

type ErrNum uint8

const (
    ErrInvalidOpcode ErrNum = iota
    ErrStoreToLiteral
)

type Error struct {
    Num     ErrNum
    Message string
}

func (e *Error) Error() (msg string) {
    return e.Message
}

func xor(a bool, b bool) (x bool) {
    return (a || b) && !(a && b)
}
