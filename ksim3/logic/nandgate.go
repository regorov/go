package logic

import (
    "github.com/kierdavis/go/ksim3"
)

type NANDGate struct {
    Sim *ksim3.Simulator

    A ksim3.Node
    B ksim3.Node
    Q ksim3.Node

    And *ANDGate
    Not *NOTGate
}

func NewNANDGate(a, b, q ksim3.Node) (c *NANDGate) {
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

    return &NANDGate{
        A:   a,
        B:   b,
        Q:   q,
        And: NewANDGate(a, b, x),
        Not: NewNOTGate(x, q),
    }
}

func (c *NANDGate) Register(sim *ksim3.Simulator) (n int) {
    return c.And.Register(sim) + c.Not.Register(sim)
}

func (c *NANDGate) Run() {
    go c.And.Run()
    go c.Not.Run()
}
