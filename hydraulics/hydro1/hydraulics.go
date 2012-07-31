// Command hydrotest demonstrates example usage of the hydraulics system.
package main

import (
    "fmt"
    "time"
)

func main() {
    source := NewSource(100)
    splitter := NewSplitter(500)
    pipe1 := NewPipe(200)
    pipe2 := NewPipe(200)
    pipe3 := NewPipe(200)
    lightBulb1 := NewLightBulb(500, 60)
    lightBulb2 := NewLightBulb(500, 60)
    
    source.SetOutput(splitter)
    splitter.SetOutput1(lightBulb1)
    splitter.SetOutput2(pipe1)
    pipe1.SetOutput(pipe2)
    pipe2.SetOutput(pipe3)
    pipe3.SetOutput(lightBulb2)
    
    for {
        source.Flow()
        splitter.Flow()
        lightBulb1.Flow("lefty")
        pipe1.Flow()
        pipe2.Flow()
        pipe3.Flow()
        lightBulb2.Flow("righty")
        
        fmt.Printf("S: %d  L1: %d  P1: %d  P2: %d  P3: %d  L2: %d\n", splitter.GetQuantity(), lightBulb1.GetQuantity(), pipe1.GetQuantity(), pipe2.GetQuantity(), pipe3.GetQuantity(), lightBulb2.GetQuantity())
        time.Sleep(time.Second)
    }
}