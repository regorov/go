package main

import (
    "github.com/kierdavis/go/dcpuem"
    "log"
    "os"
)

var Program = []uint16{
    0x7C01, 0x0030, //          SET A, 0x30
    0x7FC1, 0x0020, 0x1000, //  SET [0x1000], 0x20
    0x7803, 0x1000, //          SUB A, [0x1000]
    0xC413,         //          IFN A, 0x10
    0x7F81, 0x001A, //             SET PC, crash
    0xACC1,         //          SET I, 10
    0x7C01, 0x2000, //          SET A, 0x2000
    0x22C1, 0x2000, // :loop    SET [0x2000+I], [A]
    0x88C3,         //          SUB I, 1
    0x84D3,         //          IFN I, 0
    0x7F81, 0x000D, //             SET PC, loop
    0x9461,         //          SET X, 0x4
    0x7C20, 0x0018, //          JSR testsub
    0x7F81, 0x001A, //          SET PC, crash
    0x946E,         // :testsub SHL X, 4
    0x6381,         //          SET PC, POP
    0x7F81, 0x001A, // :crash   SET PC, crash
}

func main() {
    em := dcpuem.NewEmulator()
    em.LoadProgram(Program, 0)
    em.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

    err := em.Run()
    if err != nil {
        panic(err)
    }

    em.DumpState()
}
