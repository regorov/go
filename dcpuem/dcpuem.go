// Command dcpuem is a frontend to the dcpuemlib DCPU-16 emulator.
package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/kierdavis/go/dcpuemlib"
    "github.com/kierdavis/go/ihex"
    "os"
    "strings"
)

// -b (big-endian) flag - Reads source file as big-endian rather than little-endian
var bigEndian *bool = flag.Bool("b", false, "Reads source file as big-endian rather than little-endian")

// Function die panics with `err` if `err` is not nil.
func die(err error) {
    if err != nil {
        panic(err)
    }
}

// Function loadHex loads all the bytes from the file named `filename` in Intel Hex format and returns it.
func loadHex(filename string) (program []byte) {
    f, err := os.Open(filename); die(err)
    defer f.Close()
    
    reader := bufio.NewReader(f)
    ix, err := ihex.ReadIHex(reader); die(err)
    program = ix.ExtractDataToEnd(0)
    
    return program
}

// Function loadBin loads all the bytes from the file named `filename` in raw binary format and returns it.
func loadBin(filename string) (program []byte) {
    st, err := os.Stat(filename); die(err)
    program = make([]byte, st.Size())
    
    f, err := os.Open(filename); die(err)
    defer f.Close()
    
    _, err = f.Read(program); die(err)
    
    return program
}

// Function main is the main entry point in the program.
func main() {
    flag.Parse()
    
    if flag.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "Not enough arguments\nusage: %s file.hex\n", os.Args[0])
        os.Exit(2)
    }
    
    filename := flag.Arg(0)
    parts := strings.Split(filename, ".")
    ext := parts[len(parts) - 1]
    
    var program []byte
    
    if ext == ".hex" {
        program = loadHex(filename)
    } else if ext == ".bin" {
        program = loadBin(filename)
    } else {
        // Default to raw binary
        program = loadBin(filename)
    }
    
    em := dcpuemlib.NewEmulator()
    
    if *bigEndian {
        em.LoadProgramBytesBE(program)
    } else {
        em.LoadProgramBytesLE(program)
    }
    
    em.TraceFile = os.Stdout
    
    em.Run()
    em.DumpState()
}
