package hydraulics

import (
    
)

// Type Source represents an (as far as this package understands) infinite source of fluid. It
// satisfies Component and Producer.
type Source struct {
    power int           // The amount of fluid produced every cycle.
    output Receiver     // The pipe's output component.
}

// Function NewSource creates and returns a new Source.
func NewSource(power int) (c *Source) {
    c = new(Source)
    c.power = power
    c.output = nil
    return c
}

// Function GetPower returns the source's power.
func (c *Source) GetPower() (power int) {return c.power}

// Function SetPower sets the source's power.
func (c *Source) SetPower(power int) {c.power = power}

// Function GetOutput returns the component currently attached to this source's output.
func (c *Source) GetOutput() (output Receiver) {return c.output}

// Function SetOutput sets the component currently attached to this source's output.
func (c *Source) SetOutput(output Receiver) {c.output = output}

// Function Flow runs one cycle of the source's simulation.
func (c *Source) Flow() {
    if c.output != nil {
        c.output.AddQuantity(c.power)
    }
}
