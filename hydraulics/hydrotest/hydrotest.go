// Command hydrotest demonstrates example usage of the hydraulics system.
package main

import (
    "fmt"
    "github.com/kierdavis/go/hydraulics"
    "time"
)

/*  General guidelines for capacities & thresholds:

 * Medium-power Sources should have a power of 100.
 * All components with a capacity should have it at 200, except for outputs and splitters which have 500.
 * All components with a threshold should have it at 100.
 */

const (
    MEDIUM_POWER    = 100
    CAPACITY        = 200
    OUTPUT_CAPACITY = 500
    THRESHOLD       = 100
)

func main() {
    source1 := hydraulics.NewSource(MEDIUM_POWER)
    source2 := hydraulics.NewSource(MEDIUM_POWER)
    pipe1 := hydraulics.NewPipe(CAPACITY)
    pipe2 := hydraulics.NewPipe(CAPACITY)
    gate := hydraulics.NewANDGate(CAPACITY, THRESHOLD)
    lightBulb := hydraulics.NewLightBulb(OUTPUT_CAPACITY, THRESHOLD)

    // Chain the components
    source1.SetOutput(gate.GetFeeder1())

    source2.SetOutput(pipe1)
    pipe1.SetOutput(pipe2)
    pipe2.SetOutput(gate.GetFeeder2())

    gate.SetOutput(lightBulb)

    for {
        source1.Flow()
        source2.Flow()
        pipe1.Flow()
        pipe2.Flow()
        gate.Flow()
        lightBulb.Flow()

        fmt.Printf("G.1: %d  P1: %d  P2: %d  G.2: %d  L: %d  G: %t  L: %t\n",
            gate.GetFeeder1().GetQuantity(),
            pipe1.GetQuantity(),
            pipe2.GetQuantity(),
            gate.GetFeeder2().GetQuantity(),
            lightBulb.GetQuantity(),
            gate.GetState(),
            lightBulb.GetState())

        time.Sleep(time.Second)
    }
}
