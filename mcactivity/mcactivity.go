package main

import (
    "bufio"
    "encoding/csv"
    "flag"
    "fmt"
    "github.com/kierdavis/go/mclogscanner"
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

    mostRecentDates := make(map[string]int64)

    for {
        event, err := scanner.ReadEvent()

        if err == io.EOF {
            break
        } else if err != nil {
            error(err.Error())
        }

        if event.Type() == logscanner.PlayerConnectEventType {
            player := event.(*logscanner.PlayerConnectEvent).Player()
            secs := event.Date().Unix()
            if secs > mostRecentDates[player] {
                mostRecentDates[player] = secs
            }
        }
    }

    aWeekAgo := time.Now().Add(-time.Hour * 24 * 7).Unix()

    out := csv.NewWriter(os.Stdout)

    for player, secs := range mostRecentDates {
        if secs < aWeekAgo {
            delete(mostRecentDates, player)
        } else {
            out.Write([]string{player, time.Unix(secs, 0).String()})
        }
    }

    out.Flush()
}
