package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/kierdavis/go/k270emlib"
    "github.com/kierdavis/go/ihex"
    "os"
)

func die(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    flag.Parse()
    
    if flag.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "Not enough arguments\nusage: %s file.hex\n", os.Args[0])
        os.Exit(2)
    }
    
    f, err := os.Open(flag.Arg(0)); die(err)
    defer f.Close()
    
    reader := bufio.NewReader(f)
    ix, err := ihex.ReadIHex(reader); die(err)
    program := ix.ExtractDataToEnd(0)
    
    em := k270emlib.NewEmulator()
    em.SetMemory(program)
    em.SetTraceFile(os.Stdout)
    
    em.Run()
}
