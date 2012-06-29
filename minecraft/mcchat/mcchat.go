package main

import (
	"fmt"
	"github.com/kierdavis/go/minecraft"
	"os"
	"time"
)

func die(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	var err error

	//client, err := minecraft.Login(os.Getenv("MC_USER"), os.Getenv("MC_PASSWD"), true)
	//die(err)
	client := minecraft.LoginOffline(true)

	go func() {
		die(<-client.ErrChan)
	}()

	err = client.Join("localhost")
	die(err)
	defer client.Leave()

	time.Sleep(time.Second * 5)

	//client.Chat("hello")

	for {
	}
}
