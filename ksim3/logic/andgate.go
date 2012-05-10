package logic

import (
    "github.com/kierdavis/go/ksim3"
)

type ANDGate struct {
    Sim *ksim3.Simulator
    A   ksim3.Node
    B   ksim3.Node
    Q   ksim3.Node
}

func NewANDGate(a, b, q ksim3.Node) (c *ANDGate) {
    if a == nil {
        a = make(ksim3.Node)
    }

    if b == nil {
        b = make(ksim3.Node)
    }

    if q == nil {
        q = make(ksim3.Node)
    }

    c = &ANDGate{
        A: a,
        B: b,
        Q: q,
    }

    return c
}

func (c *ANDGate) Register(sim *ksim3.Simulator) (n int) {
    c.Sim = sim
    return 1
}

func (c *ANDGate) Run() {
    var a bool
    var b bool

    for {
        select {
        case a = <-c.A:
            c.Q <- a && b

        case b = <-c.B:
            c.Q <- a && b

        case <-c.Sim.StopChan:
            c.Sim.ComponentStopped()
            return
        }
    }
}
