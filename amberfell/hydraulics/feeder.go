package hydraulics

import (
    
)

// Type Feeder is an internal component that is used to implement components with more that one
// distinct input. It implements Component and Reciever.
type Feeder struct {
    quantity int
    capacity int
}

// Function NewFeeder creates and returns a new feeder object.
func NewFeeder(capacity int) (c *Feeder) {
    c = new(Feeder)
    c.quantity = 0
    c.capacity = capacity
    return c
}

// Function GetQuantity returns the amount of fluid currently in the feeder.
func (c *Feeder) GetQuantity() (quantity int) {return c.quantity}

// Function SetQuantity sets the amount of fluid currently in the feeder.
func (c *Feeder) SetQuantity(quantity int) {c.quantity = imin(quantity, c.capacity)}

// Function AddQuantity adds the argument to the amount of fluid currently in the feeder.
func (c *Feeder) AddQuantity(quantity int) {c.quantity = imin(c.quantity + quantity, c.capacity)}

// Function GetCapacity returns the maximum capacity of the feeder.
func (c *Feeder) GetCapacity() (capacity int) {return c.capacity}

// Function SetCapacity sets the maximum capacity of the feeder.
func (c *Feeder) SetCapacity(capacity int) {c.capacity = capacity}

// Function Flow runs one cycle of the feeder's simulation.
func (c *Feeder) Flow() {
    // Does nothing.
}
