package main

import (
    "fmt"
    "time"
)

type Link chan int

func NewLink() (link Link) {
    return make(chan int, 3)
}

type Source struct {
    power int
    interval time.Duration
    output Link
}

func NewSource(power int, interval time.Duration) (source *Source) {
    source = new(Source)
    source.power = power
    source.interval = interval
    source.output = NewLink()
    return source
}

func (source *Source) GetOutput() (output Link) {return source.output}
func (source *Source) SetOutput(output Link) {source.output = output}

func (source *Source) Pump(power int) {
    source.output <- power
}

func (source *Source) Run() {
    fmt.Printf("Source started\n")
    for {
        //fmt.Printf("Pumping %d\n", source.power)
        source.output <- source.power
        time.Sleep(source.interval)
    }
}

type Pipe struct {
    input Link
    output Link
}

func NewPipe(input Link) (pipe *Pipe) {
    pipe = new(Pipe)
    pipe.input = input
    pipe.output = NewLink()
    return pipe
}

func (pipe *Pipe) Run() {
    for {
        pressure := <-pipe.input
        pipe.output <- pressure
    }
}

func (pipe *Pipe) GetOutput() (output Link) {return pipe.output}
func (pipe *Pipe) SetOutput(output Link) {pipe.output = output}

func (pipe *Pipe) GetInput() (input Link) {return pipe.input}
func (pipe *Pipe) SetInput(input Link) {pipe.input = input}

type Accumulator struct {
    input1 Link
    input2 Link
    output Link
}

func NewAccumulator(input1 Link, input2 Link) (accumulator *Accumulator) {
    accumulator = new(Accumulator)
    accumulator.input1 = input1
    accumulator.input2 = input2
    accumulator.output = NewLink()
    return accumulator
}

func (accumulator *Accumulator) Run() {
    for {
        totalPressure := 0
        
        select {
            case totalPressure += <-accumulator.input1 {
                
            }
        }
    }
}

type LightBulb struct {
    input Link
}

func NewLightBulb(input Link) (lightBulb *LightBulb) {
    lightBulb = new(LightBulb)
    lightBulb.input = input
    return lightBulb
}

func (lightBulb *LightBulb) Run() {
    for {
        fmt.Printf("Light bulb: %d\n", <-lightBulb.input)
    }
}

func (lightBulb *LightBulb) GetInput() (input Link) {return lightBulb.input}
func (lightBulb *LightBulb) SetInput(input Link) {lightBulb.input = input}

func main() {
    source := NewSource(100, time.Second)
    pipe1 := NewPipe(source.GetOutput())
    pipe2 := NewPipe(pipe1.GetOutput())
    lightBulb := NewLightBulb(pipe2.GetOutput())
    
    go source.Run()
    go pipe1.Run()
    go pipe2.Run()
    go lightBulb.Run()
    
    for {
        time.Sleep(time.Second)
    }
}
