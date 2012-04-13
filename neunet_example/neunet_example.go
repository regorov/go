// Command neunet_example exists to demonstrate the interface of the neunet package.
package main

import (
    "fmt"
    "github.com/kierdavis/go/neunet"
    "math/rand"
    "time"
)

type Case struct {
    Inputs []float64
    Outputs []float64
}

func ShowLayers(net *neunet.NeuralNet) {
    for i, layer := range net.Layers[1:] {
        if i == len(net.Layers) - 2 {
            fmt.Printf("O:\n")
        } else {
            fmt.Printf("H:\n")
        }
        
        for _, p_ := range layer {
            p := p_.(*neunet.ComputablePerceptron)
            fmt.Printf("  V:%v B:%v E:%v Ws:%v\n", p.Value, p.Bias, p.Error, p.Weights)
        }
    }
    
    fmt.Printf("\n\n")
}

func main() {
    // Seed the PRNG
    rand.Seed(time.Now().UnixNano())
    
    // We're modelling a shape recognition function.
    // The input is in the form of an 12x12 grid of pixels - each pixel is one of two states,
    // represented by 1.0 (on) or 0.0 (off).
    // We'll train it to recognise circles and squares.
    
    net := neunet.NewNeuralNet(neunet.NewParameters(10.0, 0.002))
    
    // We need 144 input perceptrons and 2 outputs, and lets add 100 hidden ones for extra processing:
    net.AddInputLayer(144)
    net.AddHiddenLayer(50)
    net.AddOutputLayer(2)
    
    for i := 0; i < 1000; i++ {
        for _, c := range TrainingData {
            error := net.Train(c.Inputs, c.Outputs)
            fmt.Printf("E: %.3f\n", error)
        }
    }
    
    for j, c := range TestData {
        fmt.Printf("Test data %d:\n", j)
        
        for i := 0; i < len(c.Inputs); i++ {
            net.SetInput(i, c.Inputs[i])
        }
        
        net.Propagate()
        
        for i := 0; i < len(c.Outputs); i++ {
            fmt.Printf("  Output %d: %.1f (err = %.1f)\n", i, net.GetOutput(i), (c.Outputs[i] - net.GetOutput(i)))
        }
        
        fmt.Printf("\n")
    }
}
