package main

import (
	"flag"
	"fmt"
	"github.com/kierdavis/ansi"
	"github.com/kierdavis/go/minecraft"
	"io"
	"os"
	"regexp"
	"time"
)

var WhisperRegexp = regexp.MustCompile("^\xC2\xA77([a-zA-Z0-9_]+) whispers (.+)")

var (
	usernameP = flag.String("username", "Woodcutter", "The username the bot will log in with.")
	passwordP = flag.String("password", "", "The password the bot will log in with. If not specified, no authentication occurs and the server is expected to be in offline mode.")
	debugP    = flag.Bool("debug", false, "Whether to show debug messages.")
)

func die(err error) {
	if err != nil {
		ansi.Fprintf(os.Stderr, ansi.RedBold, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

type DebugWriter struct{}

func (w DebugWriter) Write(s []byte) (n int, err error) {
	return ansi.Print(ansi.White, string(s))
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		ansi.Fprintf(os.Stderr, ansi.RedBold, "usage: %s [options] <server address>\n", os.Args[0])
		os.Exit(2)
	}

	addr := flag.Arg(0)
	username := *usernameP
	password := *passwordP
	debug := *debugP

	var debugWriter io.Writer
	var client *minecraft.Client
	var err error

	if debug {
		debugWriter = DebugWriter{}
	}

	if password == "" {
		if username == "" {
			username = "Player"
		}

		ansi.Printf(ansi.Green, "Authenticating offline as %s...\n", username)
		client = minecraft.LoginOffline(username, debugWriter)

	} else {
		ansi.Printf(ansi.Green, "Authenticating as %s...\n", username)
		client, err = minecraft.Login(username, password, debugWriter)
		die(err)
	}

	client.StoreWorld = true
	client.HandleMessage = func(msg string) {
		matches := WhisperRegexp.FindStringSubmatch(msg)
		if matches != nil {
			ansi.Printf(ansi.YellowBold, "Message from %s: %s\n", matches[1], matches[2])
			client.Chat(fmt.Sprintf("/tell %s %s", matches[1], matches[2]))
		}

		//fmt.Printf("# %s\n", minecraft.ANSIEscapes(msg))
		//fmt.Printf("# %s\n", msg)
	}

	go func() {
		/*
			for err := range client.ErrChan {
				ansi.Printf(ansi.RedBold, "Error: %s\n", err.Error())
			}
		*/

		die(<-client.ErrChan)
	}()

	ansi.Printf(ansi.Green, "Connecting to %s...\n", addr)
	die(client.Join(addr))

	ansi.Printf(ansi.Green, "Connected!\n")

	go bot(client)

	kickMessage := client.Run()
	ansi.Printf(ansi.Green, "Disconnected: %s\n", kickMessage)
}

func bot(client *minecraft.Client) {
	for {
		time.Sleep(time.Second * 5)

		p, ok := findNearestTree(client)
		if !ok {
			return
		}

		ansi.Printf(ansi.Green, "Found tree: %d, %d, %d\n", p.x, p.y, p.z)

		chop(client, p)
	}
}

type xyz struct {
	x, y, z int
}

func findNearestTree(client *minecraft.Client) (p xyz, ok bool) {
	p = xyz{int(client.PlayerX), int(client.PlayerY), int(client.PlayerZ)}
	radius := 1
	x := -radius
	z := -radius
	p.x -= radius
	p.z -= radius
	dir := 1

mainloop:
	for {
		//fmt.Scanln()

		block, _, _, _, _, ok := client.GetBlock(p.x, p.y, p.z)
		if !ok {
			ansi.Printf(ansi.RedBold, "No log found!\n")
			return p, false
		}

		//println(p.x, p.y, p.z, block)

		if block == 0 {
			under, _, _, _, _, _ := client.GetBlock(p.x, p.y-1, p.z)

			if under == 0 {
				p.y--

			} else {
				if (dir > 0 && x == radius) || (dir < 0 && x == -radius) {
					if (dir > 0 && z == radius) || (dir < 0 && z == -radius) {
						if dir < 0 {
							radius++
							x = -radius
							z = -radius

							p.x -= 1
							p.z -= 1
						}

						dir = -dir

					} else {
						z += dir
						p.z += dir
					}

				} else {
					x += dir
					p.x += dir
				}
			}

		} else if block == 17 {
			// We found log!

			n := p

			for {
				n.y++

				block, _, _, _, _, ok = client.GetBlock(n.x, n.y, n.z)

				if block == 17 { // Log
					continue

				} else if block == 18 { // Leaves, always found above a tree
					break

				} else { // Not a log
					continue mainloop
				}
			}

			for {
				block, _, _, _, _, ok = client.GetBlock(p.x, p.y-1, p.z)

				if block == 17 {
					p.y--
					continue

				} else { // Bottom block of trunk
					return p, true
				}
			}

		} else {
			p.y++
		}
	}

	return p, false
}

func moveTo(client *minecraft.Client, p xyz) {
	p.x++ // Dont stand in the tree!

	ticker := time.NewTicker(minecraft.Tick)

	ansi.Printf(ansi.Green, "Moving from (%d, %d, %d) to (%d, %d, %d)\n", int(client.PlayerX), int(client.PlayerY), int(client.PlayerZ), p.x, p.y, p.z)

	for (int(client.PlayerX) != p.x) || (int(client.PlayerY) != p.y) || (int(client.PlayerZ) != p.z) {

		//println(client.PlayerX, client.PlayerY, client.PlayerZ)

		if int(client.PlayerX) < p.x {
			client.PlayerX += 0.2
		} else if int(client.PlayerX) > p.x {
			client.PlayerX -= 0.2
		}

		if int(client.PlayerY) < p.y {
			client.PlayerY += 0.2
		} else if int(client.PlayerY) > p.y {
			client.PlayerY -= 0.2
		}

		if int(client.PlayerZ) < p.z {
			client.PlayerZ += 0.2
		} else if int(client.PlayerZ) > p.z {
			client.PlayerZ -= 0.2
		}

		client.PlayerStance = client.PlayerY + 1.0
		die(client.SendPacket(0x0D, client.PlayerX, client.PlayerY, client.PlayerStance, client.PlayerZ, client.PlayerYaw, client.PlayerPitch, client.PlayerOnGround))

		<-ticker.C
	}

	ansi.Printf(ansi.Green, "Done moving\n")
}

func chop(client *minecraft.Client, p xyz) {
	for {
		block, _, _, _, _, _ := client.GetBlock(p.x, p.y, p.z)

		if block == 17 {
			moveTo(client, p)
			ansi.Printf(ansi.Green, "Breaking block at (%d, %d, %d)\n", p.x, p.y, p.z)
			die(client.SendPacket(0x0E, int8(0), int32(p.x), int8(p.y), int32(p.z), int8(5)))
			die(client.SendPacket(0x0E, int8(2), int32(p.x), int8(p.y), int32(p.z), int8(5)))
			time.Sleep(time.Second * 3)
			p.y++

		} else {
			break
		}
	}
}
