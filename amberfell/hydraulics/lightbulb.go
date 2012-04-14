package hydraulics

import (
    
)

// Type LightBulb represents a light bulb. It implements Component and Receiver.
type LightBulb struct {
    quantity int    // The amount of fluid currently in the pipe.
    capacity int    // The maxiumum capacity of the pipe.
    threshold int   // The amount of fluid required for the light bulb to trigger.
}

// Function NewLightBulb creates and returns a new LightBulb.
func NewLightBulb(capacity int, threshold int) (c *LightBulb) {
    c = new(LightBulb)
    c.quantity = 0
    c.capacity = capacity
    c.threshold = threshold
    return c
}

// Function GetQuantity returns the amount of fluid currently in the light bulb.
func (c *LightBulb) GetQuantity() (quantity int) {return c.quantity}

// Function SetQuantity sets the amount of fluid currently in the light bulb.
func (c *LightBulb) SetQuantity(quantity int) {c.quantity = imin(quantity, c.capacity)}

// Function AddQuantity adds the argument to the amount of fluid currently in the light bulb.
func (c *LightBulb) AddQuantity(quantity int) {c.quantity = imin(c.quantity + quantity, c.capacity)}

// Function GetCapacity returns the maximum capacity of the light bulb.
func (c *LightBulb) GetCapacity() (capacity int) {return c.capacity}

// Function Flow runs one cucle of the light bulb's simulation
func (c *LightBulb) Flow() {
    if c.quantity >= c.threshold {
        c.quantity -= c.threshold
    }
}

