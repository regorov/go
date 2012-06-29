package minecraft

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Starts a connection to the specified address.
func (client *Client) connect() (err error) {
	if strings.Index(client.serverAddr, ":") < 0 {
		client.serverAddr += ":25565"
	}

	if client.debug {
		fmt.Printf("Connecting to %s via TCP\n", client.serverAddr)
	}

	client.conn, err = net.Dial("tcp", client.serverAddr)
	if err != nil {
		return err
	}

	return nil
}

// Performs the 0x02 handshake transfer.
func (client *Client) handshake() (err error) {
	if client.debug {
		fmt.Printf("Sending handshake packet\n")
	}

	err = client.SendPacket(0x02, client.username+";"+client.serverAddr)
	if err != nil {
		return err
	}

	_, err = client.RecvPacket(0x02)
	if err != nil {
		return err
	}

	err = client.RecvPacketData(&client.serverId)
	if err != nil {
		return err
	}

	if client.debug {
		fmt.Printf("Received handshake packet\n")
	}

	return nil
}

// Registers the server join with session.minecraft.net
func (client *Client) registerJoin() (err error) {
	params := url.Values{
		"user":      {client.username},
		"sessionId": {client.sessionId},
		"serverId":  {client.serverId},
	}

	if client.debug {
		fmt.Printf("Registering join with minecraft.net\n")
		fmt.Printf("GET http://session.minecraft.net/game/joinserver.jsp?%s\n", params.Encode())
	}

	resp, err := http.Get("http://session.minecraft.net/game/joinserver.jsp?" + params.Encode())
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

// Performs the 0x01 login request.
func (client *Client) login() (err error) {
	if client.debug {
		fmt.Printf("Sending login packet\n")
	}

	err = client.SendPacket(0x01, int32(29), client.username, "", int32(0), int32(0), int8(0), uint8(0), uint8(0))
	if err != nil {
		return err
	}

	id, err := client.RecvPacket(0x01, 0xFF)
	if err != nil {
		return err
	}

	switch id {
	case 0xFF:
		if client.debug {
			fmt.Printf("Received kick: login was rejected\n")
		}

		var msg string
		err = client.RecvPacketData(&msg)
		if err != nil {
			return err
		}

		return fmt.Errorf("Login rejected: %s\n", msg)

	case 0x01:
		var unusedStr string
		var unusedByte uint8
		err = client.RecvPacketData(&client.entityID, &unusedStr, &client.levelType, &client.serverMode, &client.dimension, &client.difficulty, &unusedByte, &client.maxPlayers)
		if err != nil {
			return err
		}

		if client.debug {
			fmt.Printf("Received login packet\n")
		}
	}

	return nil
}

// Connects to a server.
func (client *Client) Join(addr string) (err error) {
	if client.conn != nil {
		client.Leave()
	}

	if client.debug {
		fmt.Printf("Joining server %s\n", addr)
	}

	client.serverAddr = addr

	err = client.connect()
	if err != nil {
		return err
	}

	err = client.handshake()
	if err != nil {
		return err
	}

	err = client.registerJoin()
	if err != nil {
		return err
	}

	err = client.login()
	if err != nil {
		return err
	}

	if client.debug {
		fmt.Printf("Joined!\n\nStarting receiver...\nStarting position sender...\n\n")
	}

	// Start the receiver background process.
	go client.Receiver()

	// Start the position sender background process.
	go client.PositionSender()

	return nil
}

// Runs in the background, sending an 0x0D packet every 50 ms
func (client *Client) PositionSender() {
	ticker := time.NewTicker(time.Millisecond * 50)

	for {
		select {
		case <-client.stopPositionSender:
			client.stopPositionSender <- struct{}{}
			return

		case <-ticker.C:
			/*
				if !client.playerOnGround && client.serverMode == 0 {
					if client.debug {
						fmt.Printf("Falling\n")
					}

					client.playerY -= 0.2
				}
			*/

			//fmt.Printf("sending...\n")

			err := client.SendPacket(0x0D, client.playerX, client.playerY, client.playerStance, client.playerZ, client.playerYaw, client.playerPitch, client.playerOnGround)
			if err != nil {
				client.ErrChan <- err
				continue
			}
		}
	}
}

// Sends a kick packet to the server before calling LeaveNoKick
func (client *Client) Leave() (err error) {
	if client.debug {
		fmt.Printf("Disconnecting...\n")
	}

	err = client.SendPacket(0xFF, "github.com/kierdavis/go/minecraft woz 'ere")
	if err != nil {
		return err
	}

	return client.LeaveNoKick()
}

// Shuts down background processes before closing the connection.
func (client *Client) LeaveNoKick() (err error) {
	if client.debug {
		fmt.Printf("Stopping position sender...\n")
	}

	// Tell PositionSender to stop
	client.stopPositionSender <- struct{}{}

	// Wait for a reply
	<-client.stopPositionSender
	if client.debug {
		fmt.Printf("Closing connection...\n")
	}

	client.conn.Close()
	client.conn = nil
	if client.debug {
		fmt.Printf("Done!\n\n")
	}

	return nil
}
