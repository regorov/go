package main

import (
	"fmt"
	"github.com/kierdavis/go/minecraft"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Error: Not enough arguments\n")
		os.Exit(2)
	}

	addr := os.Args[1]
	description, onlineUsers, maxUsers, err := minecraft.ScanServer(addr)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%s (%d/%d online)\n", description, onlineUsers, maxUsers)
}
