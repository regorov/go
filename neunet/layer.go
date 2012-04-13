package neunet

import (
    
)

// Type Layer represents a layer of perceptrons.
type Layer []Perceptron

// Function NewInputLayer creates and returns a new layer of input perceptrons.
func NewInputLayer(size int) (layer Layer) {
    layer = make(Layer, size)
    
    for i := 0; i < size; i++ {
        layer[i] = NewInputPerceptron()
    }
    
    return layer
}

// Function NewComputableLayer creates and returns a new layer of computable perceptrons. See
// the documentation for NewNeuralNet for information on the `h` and `gain` parameters.
func NewComputableLayer(parentLayer Layer, size int, params *Parameters) (layer Layer) {
    layer = make(Layer, size)
    
    for i := 0; i < size; i++ {
        layer[i] = NewComputablePerceptron(parentLayer, params)
    }
    
    return layer
}
