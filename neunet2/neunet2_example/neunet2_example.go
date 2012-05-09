package main

import (
    "github.com/kierdavis/go/neunet2"
)

func main() {
    m := neunet2.NewMatrix(3, 3)

    for i := 0; i < 9; i++ {
        m.Data[i/3][i%3] = float32(i)
    }

    m.Show()

    m.SMul(2.0).Show()

    m.HMul(m).Show()

    m.KMul(m).Show()
}
