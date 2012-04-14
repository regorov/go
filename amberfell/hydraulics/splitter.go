package hydraulics

import (
    
)

// Type Splitter represents a splitter that can take any number of inputs and provides two outputs.
// Half of the input is directed to each output. It satisfies Component, Producer? and Receiver.
type Splitter struct {
    quantity int        // The amount of fluid currently in the splitter.
    capacity int        // The maxiumum capacity of the splitter.
    output1 Receiver    // One of the splitter's output components.
    output2 Receiver    // The other one of the splitter's output components.
}

// Function NewSplitter creates an returns a new Splitter with the specified capacity.
func NewSplitter(capacity int) (c *Splitter) {
    c = new(Splitter)
    c.quantity = 0
    c.capacity = capacity
    c.output1 = nil
    c.output2 = nil
    return c
}

// Function GetOutput1 returns the component currently attached to this splitter's output 1.
func (c *Splitter) GetOutput1() (output Receiver) {return c.output1}

// Function GetOutput2 returns the component currently attached to this splitter's output 2.
func (c *Splitter) GetOutput2() (output Receiver) {return c.output2}

// Function SetOutput1 sets the component currently attached to this splitter's output 1.
func (c *Splitter) SetOutput1(output Receiver) {c.output1 = output}

// Function SetOutput2 sets the component currently attached to this splitter's output 2.
func (c *Splitter) SetOutput2(output Receiver) {c.output2 = output}

// Function GetQuantity returns the amount of fluid currently in the splitter.
func (c *Splitter) GetQuantity() (quantity int) {return c.quantity}

// Function SetQuantity sets the amount of fluid currently in the splitter.
func (c *Splitter) SetQuantity(quantity int) {c.quantity = imin(quantity, c.capacity)}

// Function AddQuantity adds the argument to the amountt of fluid currently in the splitter.
func (c *Splitter) AddQuantity(quantity int) {c.quantity = imin(c.quantity + quantity, c.capacity)}

// Function GetCapacity returns the maximum capacity of the splitter.
func (c *Splitter) GetCapacity() (capacity int) {return c.capacity}

// Function Flow runs one cycle of the splitter's simulation.
func (c *Splitter) Flow() {
    if c.output1 != nil {
        if c.output2 != nil {
            thisQuantity, output1Quantity, output2Quantity := balance3YZLimited(c.quantity, c.output1.GetQuantity(), c.output2.GetQuantity(), c.output1.GetCapacity(), c.output2.GetCapacity())
            
            c.quantity = thisQuantity
            c.output1.SetQuantity(output1Quantity)
            c.output2.SetQuantity(output2Quantity)
        
        }
    }
}
