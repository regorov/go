package main

import (
	"encoding/json"
	"fmt"
	"github.com/kierdavis/go/reiminimap"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Not enough arguments\nusage: %s <input.points> <output.json>\n", os.Args[0])
		os.Exit(2)
	}

	waypointsFile := os.Args[1]

	f, err := os.Open(waypointsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	waypoints, err := reiminimap.ReadWaypoints(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	data, err := json.Marshal(waypoints)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	out, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer out.Close()

	_, err = out.Write(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
