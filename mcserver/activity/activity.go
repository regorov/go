package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/kierdavis/go/mcserver/logscanner"
    "io"
    "os"
    "time"
)

func error(msg string) {
    fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
    os.Exit(1)
}

func main() {
    flag.Parse()

    if flag.NArg() < 1 {
        error(fmt.Sprintf("Not enough arguments (usage: %s <server.log>)", os.Args[0]))
    }

    filename := flag.Arg(0)

    f, err := os.Open(filename)
    if err != nil {
        error(err.Error())
    }
    defer f.Close()

    reader := bufio.NewReader(f)
    scanner := logscanner.NewLogScanner(reader, time.Local)

    for {
        event, err := scanner.ReadEvent()

        if err == io.EOF {
            break
        } else if err != nil {
            error(err.Error())
        }

        fmt.Printf("%d %+v\n", event.Type(), event)
    }
}
