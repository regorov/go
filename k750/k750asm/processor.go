package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os"
    "sync"
)

type AsmError struct {
    Coord   Coord
    Message string
}

func (e *AsmError) Error() (msg string) {
    return fmt.Sprintf("%s %s", e.Coord.String(), e.Message)
}

var parserOutput = make(chan Item)

var errChan = make(chan *AsmError)
var errControl = make(chan bool)

var waitGroup sync.WaitGroup

func errorMonitor() {
    wasErrors := false

    for {
        select {
        case asmerr := <-errChan:
            log.Println(asmerr.Error())
            wasErrors = true

        case <-errControl:
            errControl <- wasErrors
            return
        }
    }
}

func stage1() (items []Item, ok bool) {
    // Run the first stage - verify, reduce aliases and compute lengths

    go errorMonitor()

    items = make([]Item, 0)

    for item := range parserOutput {
        waitGroup.Add(1)
        go item.VerifyAndReduce()

        items = append(items, item)
    }

    waitGroup.Wait()

    errControl <- true
    if <-errControl { // Reply indicates whether there were errors
        return nil, false
    }

    return items, true
}

func stage2(items []Item) (labelMap map[string]uint32, ok bool) {
    // Run the second stage - label mapping
    // Requires that lengths have been computed (in stage 1)
    // Ensure that the items are assigned offsets in the corrent order.

    //go errorMonitor()

    offset := uint32(0)
    labelMap = make(map[string]uint32)

    for _, item := range items {
        label, ok := item.Label()
        if ok {
            labelMap[label] = offset
        }

        item.SetOffset(offset)
        offset += item.Length()
    }

    return labelMap, true
}

func stage3(items []Item, labelMap map[string]uint32) (ok bool) {
    // Run the third stage - encoding
    // Requires that label offsets have been computed

    go errorMonitor()

    for _, item := range items {
        waitGroup.Add(1)
        go item.Encode(labelMap)
    }

    waitGroup.Wait()

    errControl <- true
    if <-errControl {
        return false
    }

    return true
}

func RunAssembler(reader io.Reader, filename string) {
    pushCoord(Coord{filename, 1})

    lexer := newLexer(bufio.NewReader(reader))
    go yyParse(lexer)

    // To-do:
    //  * convert items to a channel
    //  * start stage2 in a goroutine, creating and passing a new labelMap as an argument
    //  * allows stage2 to process items as soon as they are finished with by stage1
    //  * collect the items list in stage2 rather than stage1

    items, ok := stage1()
    if !ok {
        return
    }

    labelMap, ok := stage2(items)
    if !ok {
        return
    }

    ok = stage3(items, labelMap)
    if !ok {
        return
    }

    for _, item := range items {
        fmt.Printf("%s %s %v\n", item.GetCoord().String(), item.String(), item.Encoded())
    }
}

func main() {
    fname := os.Args[1]

    f, err := os.Open(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    RunAssembler(f, fname)
}
