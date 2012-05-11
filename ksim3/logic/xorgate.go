package logic

import (
    "github.com/kierdavis/go/ksim3"
)

type XORGate struct {
    Sim *ksim3.Simulator

    A ksim3.Node
    B ksim3.Node
    Q ksim3.Node

    N1 *NANDGate
    N2 *NANDGate
    N3 *NANDGate
    N4 *NANDGate
}

func NewXORGate(a, b, q ksim3.Node) (c *XORGate) {
    if a == nil {
        a = make(ksim3.Node)
    }

    if b == nil {
        b = make(ksim3.Node)
    }

    if q == nil {
        q = make(ksim3.Node)
    }

    x := make(ksim3.Node)
    y := make(ksim3.Node)
    z := make(ksim3.Node)

    return &XORGate{
        A:  a,
        B:  b,
        Q:  q,
        N1: NewNANDGate(a, b, x),
        N2: NewNANDGate(a, x, y),
        N3: NewNANDGate(b, x, z),
        N4: NewNANDGate(y, z, q),
    }
}

func (c *XORGate) Register(sim *ksim3.Simulator) (n int) {
    return 1 + c.N1.Register(sim) + c.N2.Register(sim) + c.N3.Register(sim) + c.N4.Register(sim)
}

func (c *XORGate) Run() {
    go c.N1.Run()
    go c.N2.Run()
    go c.N3.Run()
    go c.N4.Run()
}
