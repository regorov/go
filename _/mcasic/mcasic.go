package main

import (
    "fmt"
)

func test(net *Net, a byte, b byte) {
    inputs := []byte{a, b}
    outputs := []byte{0}
    net.Propagate(inputs, outputs)

    fmt.Printf("%d ^ %d -> %d\n", a, b, outputs[0])
}

func main() {
    net := NewNet()

    net.AddLayer(2)
    net.AddLayer(5)
    net.AddLayer(1)

    test(net, 0, 0)
    test(net, 0, 1)
    test(net, 1, 0)
    test(net, 1, 1)
}
