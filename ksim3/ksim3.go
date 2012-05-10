package ksim3

import (
    "sync"
)

type Node chan bool

type Component interface {
    Register(*Simulator) int
    Run()
}

type Simulator struct {
    ComponentCount int
    RootComponents []Component
    WaitGroup      sync.WaitGroup
    StopChan       chan bool
}

func NewSimulator() (sim *Simulator) {
    sim = new(Simulator)
    sim.RootComponents = make([]Component, 0)
    sim.StopChan = make(chan bool)
    return sim
}

func (sim *Simulator) Add(c Component) {
    sim.ComponentCount += c.Register(sim)
    sim.RootComponents = append(sim.RootComponents, c)
}

func (sim *Simulator) Start() {
    sim.WaitGroup.Add(sim.ComponentCount)

    for _, c := range sim.RootComponents {
        go c.Run()
    }
}

func (sim *Simulator) Stop() {
    for i := 0; i < sim.ComponentCount; i++ {
        sim.StopChan <- true
    }

    sim.WaitGroup.Wait()
}

func (sim *Simulator) ComponentStopped() {
    sim.WaitGroup.Done()
}
