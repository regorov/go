// Package neunet provides a basic implementation of an artificial neural net.
package neunet

import (
    
)

// Type Parameters holds the global H and gain parameters. See the documentation for NewParameters
// for information on the parameters.
type Parameters struct {
    H float64           // The coefficient for the sigmoid function (passed to
                        // NewComputablePerceptron)
    Gain float64        // The learning gain (passed to NewComutablePerceptron)
}

// Function NewParameters creates and returns a new Parameters object. `h` is the coefficient for
// the sigmoid function and `gain` is the learning gain. If you are unsure what to put for these
// values, 1.0 and 5.0 tend to be a good starting point.
func NewParameters(h float64, gain float64) (params *Parameters) {
    params = new(Parameters)
    params.H = h
    params.Gain = gain
    return params
}

// Type NeuralNet represents a neural network.
type NeuralNet struct {
    Layers []Layer      // The layers in the network.
    InputLayer Layer    // The input layer (must be in Layers as well).
    OutputLayer Layer   // The output layer (must be in Layers as well).
    Params *Parameters  // The H and gain parameters.
    prevLayer Layer     // The last layer added.
}

// Function NewNeuralNet creates and returns a new NeuralNet.
func NewNeuralNet(params *Parameters) (net *NeuralNet) {
    net = new(NeuralNet)
    net.Layers = make([]Layer, 0, 3)
    net.InputLayer = nil
    net.OutputLayer = nil
    net.Params = params
    net.prevLayer = nil
    
    return net
}

// Function NeuralNet.AddInputLayer creates and adds a new input layer to the net.
func (net *NeuralNet) AddInputLayer(size int) (layer Layer) {
    layer = NewInputLayer(size)
    
    net.Layers = append(net.Layers, layer)
    net.InputLayer = layer
    net.prevLayer = layer
    
    return layer
}

// Function NeuralNet.AddHiddenLayer creates and adds a new hidden layer to the net.
func (net *NeuralNet) AddHiddenLayer(size int) (layer Layer) {
    layer = NewComputableLayer(net.prevLayer, size, net.Params)
    
    net.Layers = append(net.Layers, layer)
    net.prevLayer = layer
    
    return layer
}

// Function NeuralNet.AddOutputLayer creates and adds a new output layer to the net.
func (net *NeuralNet) AddOutputLayer(size int) (layer Layer) {
    layer = NewComputableLayer(net.prevLayer, size, net.Params)
    
    net.Layers = append(net.Layers, layer)
    net.OutputLayer = layer
    net.prevLayer = layer
    
    return layer
}

// Function NeuralNet.SetInput sets the value of the input perceptron numbered `num` to `value`.
func (net *NeuralNet) SetInput(num int, value float64) {
    p := net.InputLayer[num]
    p.(*InputPerceptron).SetValue(value)
}

// Function NeuralNet.GetOutput returns the value of the output perceptron numbered `num`.
func (net *NeuralNet) GetOutput(num int) (value float64) {
    p := net.OutputLayer[num]
    return p.GetValue()
}

// Function NeuralNet.Propagate calls Compute on all perceptrons.
func (net *NeuralNet) Propagate() {
    for _, layer := range net.Layers {
        for _, p := range layer {
            p.Compute()
        }
    }
}

// Function NeuralNet.CalculateOutputError calculates the error between the `expectedOutputs`
// parameter and the output perceptrons.
func (net *NeuralNet) CalculateOutputError(expectedOutputs []float64) (error float64) {
    error = 0.0
    
    for i, p := range net.OutputLayer {
        thisError := expectedOutputs[i] - p.GetValue()
        error += 0.5 * thisError * thisError
    }
    
    return error
}

// Function NeuralNet.BackPropagate updates the weights based on the error
func (net *NeuralNet) BackPropagate(expectedOutputs []float64) {
    for i, p_ := range net.OutputLayer {
        p := p_.(*ComputablePerceptron)
        v := p.GetValue()
        error := v * (1.0 - v) * (expectedOutputs[i] - v)
        p.AdjustWeights(error)
    }
    
    for i := len(net.Layers) - 2; i > 0; i-- {
        layer := net.Layers[i]
        
        for _, p_ := range layer {
            p := p_.(*ComputablePerceptron)
            v := p.GetValue()
            
            var sum float64 = 0.0
            for _, output_ := range net.Layers[i + 1] {
                output := output_.(*ComputablePerceptron)
                sum += output.GetIncomingWeight(p_) * output.Error
            }
            
            error := v * (1.0 - v) * sum
            p.AdjustWeights(error)
        }
    }
}

// Function NeuralNet.Train performs a single training cycle on the training set `tset`, returning
// the error distance.
func (net *NeuralNet) Train(inputs []float64, expectedOutputs []float64) (error float64) {
    for i, p := range net.InputLayer {
        p.(*InputPerceptron).SetValue(inputs[i])
    }
    
    net.Propagate()
    error = net.CalculateOutputError(expectedOutputs)
    net.BackPropagate(expectedOutputs)
    return error
}
