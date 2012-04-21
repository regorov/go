package main

import (
    "flag"
    "fmt"
    "github.com/kierdavis/go/amberfell/mm3dmodel"
    "os"
)

func die(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
        os.Exit(1)
    }
}

func printMetadata(model *mm3dmodel.Model) {
    fmt.Printf("Metadata:\n")
    for key, value := range model.Metadata() {
        fmt.Printf("  %s: %s\n", key, value)
    }
    
    if len(model.Metadata()) == 0 {
        fmt.Printf("  None\n\n")
    } else {
        fmt.Printf("\n")
    }
}

func printVertices(model *mm3dmodel.Model) {
    fmt.Printf("Vertices (%d):\n" % model.NVertices())
    for i := 0; i < model.NVertices(); i++ {
        vertex := model.Vertex(i)
        fmt.Printf("  %4d: Flags: 0x%04X  X: %.3f  Y: %.3f  Z: %.3f\n", i, vertex.Flags(), vertex.X(), vertex.Y(), vertex.Z())
    }
    
    if model.NVertices() == 0 {
        fmt.Printf("  None\n\n")
    } else {
        fmt.Printf("\n")
    }
}

func main() {
    flag.Parse()
    
    if flag.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "Not enough arguments.\n\nUsage: %s <file.mm3d>\n", os.Args[0])
        os.Exit(2)
    }
    
    fname := flag.Arg(0)
    f, err := os.Open(fname); die(err)
    defer f.Close()
    
    model, err := mm3dmodel.Read(f); die(err)
    
    fmt.Printf("Version: 0x%02X:0x%02X\n", model.MajorVersion(), model.MinorVersion())
    fmt.Printf("Model flags: 0x%02x\n", model.ModelFlags())
    fmt.Printf("Num dirty segments: %d\n\n", model.NDirtySegments())
    
    printMetadata(model)
    printVertices(model)
}
