package hydraulics

import (
    
)

// Type ANDGate represents a two-input AND gate (proper name coming soon). It implements Component and
// Producer, and has two built-in input feeders.
type ANDGate struct {
    feeder1 *Feeder
    feeder2 *Feeder
    threshold int
    output Receiver
    state bool
}

// Function NewANDGate creates and returns a new ANDGate.
func NewANDGate(feederCapacity int, threshold int) (c *ANDGate) {
    c = new(ANDGate)
    c.feeder1 = NewFeeder(feederCapacity)
    c.feeder2 = NewFeeder(feederCapacity)
    c.threshold = threshold
    c.output = nil
    return c
}

// Function GetFeeder1 returns this component's feeder 1.
func (c *ANDGate) GetFeeder1() (feeder *Feeder) {return c.feeder1}

// Function GetFeeder2 returns this component's feeder 2.
func (c *ANDGate) GetFeeder2() (feeder *Feeder) {return c.feeder2}

// Function GetThreshold returns this component's threshold level.
func (c *ANDGate) GetThreshold() (threshold int) {return c.threshold}

// Function SetThreshold sets this component's threshold level.
func (c *ANDGate) SetThreshold(threshold int) {c.threshold = threshold}

// Function GetOutput returns the component currently attached to this component's output.
func (c *ANDGate) GetOutput() (output Receiver) {return c.output}

// Function SetOutput sets the component currently attached to this component's output.
func (c *ANDGate) SetOutput(output Receiver) {c.output = output}

// Function GetState returns the most recently calculated state of the gate (true = on, false=off).
func (c *ANDGate) GetState() (state bool) {return c.state}

// Function Flow runs one cycle of the component's simulation.
func (c *ANDGate) Flow() {
    active1 := c.feeder1.GetQuantity() >= c.threshold
    active2 := c.feeder2.GetQuantity() >= c.threshold
    
    if active1 && active2 && c.output != nil { // Open the gate
        feeder1Quantity, feeder2Quantity, outputQuantity := balance3ZLimited(c.feeder1.GetQuantity(), c.feeder2.GetQuantity(), c.output.GetQuantity(), c.output.GetCapacity())
        c.feeder1.SetQuantity(feeder1Quantity)
        c.feeder2.SetQuantity(feeder2Quantity)
        c.output.SetQuantity(outputQuantity)
        c.state = true
    } else {
        c.state = false
    }
}
