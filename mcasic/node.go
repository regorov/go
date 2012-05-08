package main

import (
    "math/rand"
)

type Input struct {
    Node   *Node
    Weight float64
}

type Node struct {
    Inputs []*Input
    Bias   float64
    Value  float64
}

func NewNode(nodes []*Node) (node *Node) {
    if nodes == nil {
        inputs := make([]*Input, 0)

    } else {
        inputs := make([]*Input, len(nodes))
        for i, n := range nodes {
            inputs[i] = &Input{n, 0.0}
        }
    }

    node = &Node{
        Inputs: inputs,
        Bias:   0.0,
        Value:  0.0,
    }

    return node
}

func (node *Node) Fill(genome chan float64) {
    node.Bias = <-genome

    for _, input := range node.Inputs {
        input.Weight = <-genome
    }
}

func (node *Node) Calculate(c chan bool) {
    value = node.Bias

    for _, input := range node.Inputs {
        value += input.Node.Value * input.Weight
    }

    input.Value = value
    ch <- true
}
