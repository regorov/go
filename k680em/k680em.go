package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/kierdavis/go/ihex"
    "github.com/kierdavis/go/k680emlib"
    "os"
)

func main() {
    /*
       defer func() {
           if x := recover(); x != nil {
               fmt.Fprintf(os.Stderr, "Runtime error: %s\n", x)
               os.Exit(1)
           }
       }()
    */

    flag.Parse()

    if flag.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "Not enough arguments\nusage: %s file.hex\n", os.Args[0])
        os.Exit(2)
    }

    f, err := os.Open(flag.Arg(0))
    if err != nil {
        panic(err)
    }

    reader := bufio.NewReader(f)
    ix, err := ihex.ReadIHex(reader)
    if err != nil {
        panic(err)
    }
    program := ix.ExtractDataToEnd(0)

    em := k680emlib.NewEmulator()
    em.TraceFile = os.Stdout
    em.LoadProgram(program, 0)

    err = em.Run()
    if err != nil {
        panic(err)
    }

    em.DumpState()
}
