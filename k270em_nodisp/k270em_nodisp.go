// Command k270em_nodisp is a frontend to the k270emlib K270 processor emulator, but without the SDL
// dependencies. This removes the character display element (however, LDV and STV instructions will
// still continue to work).
package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/kierdavis/go/k270emlib"
    "github.com/kierdavis/go/ihex"
    "io"
    "os"
)

var stdinReader = bufio.NewReader(os.Stdin)
var inputBuffer = make([]byte, 0)

var (
    screenDump = flag.String("d", "", "Write a screen dump to the specified file.")
)

// Function die panics with `err` if `err` is not nil.
func die(err error) {
    if err != nil {
        panic(err)
    }
}

// Function getKey is the keyboard handler for the emulator (see k270emlib.Emulator.SetGetKey).
func getKey() (key byte) {
    if len(inputBuffer) == 0 {
        os.Stdout.Write([]byte("Input required: "))
        
        line, _, err := stdinReader.ReadLine()
        
        if err == io.EOF {
            return 0xFF // EOF indicator
        } else {
            die(err)
        }
        
        inputBuffer = line
    }
    
    ch := inputBuffer[0]
    inputBuffer = inputBuffer[1:]
    return ch
}

// Function main is the main entry point in the program.
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
    em.SetGetKey(getKey)
    em.SetMemory(program)
    em.SetTraceFile(os.Stdout)
    
    em.Run()
    fmt.Println("")
    em.DumpState()
    
    if *screenDump != "" {
        sf, err := os.Create(*screenDump); die(err)
        defer sf.Close()
        
        vmem := em.GetVideoMemory()

        for y := 0; y < k270emlib.VMEM_HEIGHT; y++ {
            for x := 0; x < k270emlib.VMEM_WIDTH; x++ {
                ch := vmem[(y * k270emlib.VMEM_WIDTH) + x]
                
                if ch.Char == 0 {
                    _, err := sf.Write([]byte{' '}); die(err)
                } else {
                    _, err := sf.Write([]byte{ch.Char}); die(err)
                }
            }
            
            _, err := sf.Write([]byte{'\n'}); die(err)
        }
        
        fmt.Printf("Wrote '%s'\n", *screenDump)
    }
}
