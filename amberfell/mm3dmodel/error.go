package mm3dmodel

import (
    "fmt"
)

const (
    E_MAGIC_MISMATCH = iota
    E_UNIMPLEMENTED_VERSION
)

type Error struct {
    id int
    msg string
}

func NewError(id int, msg string) (e *Error) {
    return &Error{id: id, msg: msg}
}

func (e *Error) ID() (id int) {
    return e.id
}

func (e *Error) Error() (str string) {
    return fmt.Sprintf("[github.com/kierdavis/amberfell/mm3dmodel] %s", e.msg)
}
