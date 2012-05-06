package main

import (
    "flag"
    "fmt"
    "github.com/kierdavis/go/dcpuem"
    "github.com/kierdavis/go/dcpuem/keyboard_tb"
    "github.com/kierdavis/go/dcpuem/lem1802_tb"
    "github.com/nsf/termbox-go"
    "log"
    "os"
    "time"
)

var termboxInitialised bool = false

func die(err error) {
    if err != nil {
        if termboxInitialised {
            termbox.Close()
        }

        fmt.Fprintf(os.Stderr, "Runtime error: %s\n", err.Error())
        os.Exit(1)
    }
}

func main() {
    flag.Parse()

    if flag.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "usage: %s <file>\n", os.Args[0])
        os.Exit(2)
    }

    fname := flag.Arg(0)
    fi, err := os.Stat(fname)
    die(err)

    f, err := os.Open(fname)
    die(err)
    defer f.Close()

    buffer := make([]byte, fi.Size())
    _, err = f.Read(buffer)
    die(err)

    logfile, err := os.Create(fname + ".trace")
    die(err)
    defer logfile.Close()

    em := dcpuem.NewEmulator()
    em.LoadProgramBytesBE(buffer, 0)
    em.Logger = log.New(logfile, "", log.LstdFlags|log.Lshortfile)
    em.ClockTicker = time.NewTicker(time.Second / 240.0) // 60Hz

    disp := lem1802_tb.New()
    em.AttachDevice(disp)

    kbd := keyboard_tb.New()
    em.AttachDevice(kbd)

    err = termbox.Init()
    die(err)
    termboxInitialised = true
    defer termbox.Close()

    em.StartDevices()
    defer em.StopDevices()

    go em.Run()

mainloop:
    for {
        ev := termbox.PollEvent()

        if kbd.HandleEvent(ev) {
            continue
        }

        switch ev.Type {
        case termbox.EventKey:
            switch ev.Key {
            case termbox.KeyEsc:
                break mainloop
            }
        }
    }

    //    em.DumpState()
}
