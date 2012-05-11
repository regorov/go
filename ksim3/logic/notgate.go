package logic

import (
    "github.com/kierdavis/go/ksim3"
)

type NOTGate struct {
    Sim *ksim3.Simulator

    A ksim3.Node
    Q ksim3.Node
}

func NewNOTGate(a, q ksim3.Node) (c *NOTGate) {
    if a == nil {
        a = make(ksim3.Node)
    }

    if q == nil {
        q = make(ksim3.Node)
    }

    c = &NOTGate{
        A: a,
        Q: q,
    }

    return c
}

func (c *NOTGate) Register(sim *ksim3.Simulator) (n int) {
    c.Sim = sim
    return 1
}

func (c *NOTGate) Run() {
    var a bool

    for {
        select {
        case a = <-c.A:
            c.Q <- !a

        case <-c.Sim.StopChan:
            c.Sim.ComponentStopped()
            return
        }
    }
}
