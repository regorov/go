package neunet

import (
    //"fmt"

    "math"
    "math/rand"
)

// Interface Perceptron represents a perceptron.
type Perceptron interface {
    Compute() float64
    GetValue() float64
}

// Type InputPerceptron represents an input perceptron.
type InputPerceptron struct {
    Value float64   // The perceptron's value
}

// Function NewInputPerceptron creates and returns a new InputPerceptron.
func NewInputPerceptron() (p *InputPerceptron) {
    p = new(InputPerceptron)
    p.Value = 0.0
    return p
}

// Function InputPerceptron.Compute does nothing, but exists to satisfy the Perceptron interface.
func (p *InputPerceptron) Compute() (value float64) {
    return p.Value
}

// Function InputPerceptron.SetValue sets the value of an input perceptron to `value`.
func (p *InputPerceptron) SetValue(value float64) {
    p.Value = value
}

// Function InputPerceptron.GetValue returns the perceptron's value.
func (p *InputPerceptron) GetValue() (value float64) {
    return p.Value
}

// Type ComputablePerceptron represents a perceptron who's output is computed from the outputs of
// the previous layer of perceptrons.
type ComputablePerceptron struct {
    ParentLayer Layer   // The previous layer.
    Weights []float64   // The perceptron-specific weights associated with the previous layer's
                        // perceptrons.
    Bias float64        // The bias weight.
    Value float64       // The perceptron's most recently calculated output value.
    Error float64       // The current error value (for training).
    Params *Parameters  // The parameters for computation.
}

// Function NewComputablePerceptron creates and returns a new ComputablePerceptron.
func NewComputablePerceptron(parentLayer Layer, params *Parameters) (p *ComputablePerceptron) {
    p = new(ComputablePerceptron)
    p.ParentLayer = parentLayer
    p.Weights = make([]float64, len(parentLayer))
    p.Bias = (rand.Float64() * 2.0) - 1.0
    p.Value = 0.0
    p.Error = 0.0
    p.Params = params
    
    for i := 0; i < len(p.Weights); i++ {
        p.Weights[i] = (rand.Float64() * 2.0) - 1.0
    }
    
    return p
}

// Function ComputablePerceptron.Compute uses the outputs of the parent layer to compute the new
// state of the perceptron. It also returns this value.
func (p *ComputablePerceptron) Compute() (value float64) {
    var total float64 = p.Bias
    
    for i := 0; i < len(p.ParentLayer); i++ {
        total += p.ParentLayer[i].GetValue() * p.Weights[i]
    }
    
    p.Value = 1.0 / (1.0 + math.Exp(-total * p.Params.H))
    //fmt.Printf("Sigmoid(%.1f) = %.1f\n", total, p.Value)
    return p.Value
}

// Function ComputablePerceptron.GetValue returns the current state of the perceptron.
func (p *ComputablePerceptron) GetValue() (value float64) {
    return p.Value
}

// Function ComputablePerceptron.AdjustWeights adjusts the input weights based on the error
// distance `error`.
func (p *ComputablePerceptron) AdjustWeights(error float64) {
    for i, input := range p.ParentLayer {
        p.Weights[i] += p.Params.Gain * error * input.GetValue()
    }
    
    p.Error = error
}

// Function ComputablePerceptron.GetIncomingWeight finds the weight of the input that has arrived
// from the given perceptron `other`.
func (p *ComputablePerceptron) GetIncomingWeight(other Perceptron) (weight float64) {
    for i, input := range p.ParentLayer {
        if input == other {
            return p.Weights[i]
        }
    }
    
    return 0.0
}
