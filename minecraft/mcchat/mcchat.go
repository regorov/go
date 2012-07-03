package main

import (
	"bufio"
	"fmt"
	"github.com/kierdavis/go/minecraft"
	"io"
	"os"
)

func die(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		panic("foobar")
	}
}

func main() {
	defer func() {
		recover()
	}()

	if len(os.Args) < 2 {
		fmt.Printf("Not enough arguments\n\nusage: %s <server address>\n\nThis program expects the MC_USER and MC_PASSWD environment variables to be set. Otherwise, the user is logged in with an offline account.\n", os.Args[0])
		os.Exit(2)
	}

	fmt.Printf("*** Welcome to mcchat!\n")

	addr := os.Args[1]
	username := os.Getenv("MC_USER")
	password := os.Getenv("MC_PASSWD")

	fmt.Printf("*** Logging in...\n")

	var err error
	var client *minecraft.Client

	if password == "" {
		if username == "" {
			username = "Player"
		}

		client = minecraft.LoginOffline(username, nil)

	} else {
		client, err = minecraft.Login(username, password, nil)
		die(err)
	}

	go func() {
		die(<-client.ErrChan) // Die if there's ever an asynchronous error
	}()

	client.HandleMessage = func(msg string) {
		fmt.Printf("\r%s\n>", minecraft.ANSIEscapes(msg))
	}

	fmt.Printf("*** Connecting to %s...\n", addr)
	die(client.Join(addr))
	fmt.Printf("*** Connected!\n*** Type & press enter to send messages!\n*** Press Ctrl+D to exit\n\n")

	go func() {
		stdinReader := bufio.NewReader(os.Stdin)

		fmt.Printf(">")

		for {
			msg, err := stdinReader.ReadString('\n')
			if err == io.EOF {
				client.Leave()
				client.Logout()
				return
			}

			fmt.Printf("\x1b[1T>")

			die(err)
			client.Chat(msg[:len(msg)-1])
		}
	}()

	client.Run()
}
