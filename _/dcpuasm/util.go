package main

import (
    "fmt"
    "os"
)

var wasErrors = false

func error(msg string) {
    fmt.Printf("*** %s\n", msg)
    wasErrors = true
}

func exitIfErrors() {
    if wasErrors {
        fmt.Printf("*** Exiting due to errors\n")
        os.Exit(1)
    }
}