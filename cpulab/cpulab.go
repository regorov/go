package main

import (
    "fmt"
    "github.com/kierdavis/go/termdialog"
    "github.com/nsf/termbox-go"
    "os"
)

var MainMenuDialog = termdialog.NewSelectionDialog("Main Menu", []termdialog.Option{termdialog.Option{"foo", func(option termdialog.Option) { termbox.Close(); println("foo") }}, termdialog.Option{"bar", nil}})

func die(err error) {
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
    }
}

func main() {
    err := termbox.Init()
    die(err)
    defer termbox.Close()

    MainMenuDialog.Run()

    termbox.Close()

    fmt.Fprintf(os.Stderr, "Result: %d\n", (MainMenuDialog.GetSelectedIndex()))
}
