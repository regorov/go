package hydraulics

import (
    
)

// Type Pipe represents a pipe that can take any number of inputs and provides one output. It
// satisfies Component, Producer and Receiver
type Pipe struct {
    quantity int        // The amount of fluid currently in the pipe.
    capacity int        // The maxiumum capacity of the pipe.
    output Receiver     // The pipe's output component.
}

// Function NewPipe creates an returns a new Pipe with the specified capacity.
func NewPipe(capacity int) (c *Pipe) {
    c = new(Pipe)
    c.quantity = 0
    c.capacity = capacity
    c.output = nil
    return c
}

// Function GetOutput returns the component currently attached to this pipe's output.
func (c *Pipe) GetOutput() (output Receiver) {return c.output}

// Function SetOutput sets the component currently attached to this pipe's output.
func (c *Pipe) SetOutput(output Receiver) {c.output = output}

// Function GetQuantity returns the amount of fluid currently in the pipe.
func (c *Pipe) GetQuantity() (quantity int) {return c.quantity}

// Function SetQuantity sets the amount of fluid currently in the pipe.
func (c *Pipe) SetQuantity(quantity int) {c.quantity = imin(quantity, c.capacity)}

// Function AddQuantity adds the argument to the amountt of fluid currently in the pipe.
func (c *Pipe) AddQuantity(quantity int) {c.quantity = imin(c.quantity + quantity, c.capacity)}

// Function GetCapacity returns the maximum capacity of the pipe.
func (c *Pipe) GetCapacity() (capacity int) {return c.capacity}

// Function Flow runs one cycle of the pipe's simulation.
func (c *Pipe) Flow() {
    if c.output != nil {
        thisQuantity, outputQuantity := balanceYLimited(c.quantity, c.output.GetQuantity(),
            c.output.GetCapacity())
        
        c.quantity = thisQuantity
        c.output.SetQuantity(outputQuantity)
    }
}

