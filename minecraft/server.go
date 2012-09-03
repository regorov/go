package minecraft

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type encryptedStream struct {
	conn   io.ReadWriter
	cipher cipher.Block
}

func min(a, b int) (min int) {
	if a < b {
		return a
	}

	return b
}

func (s *encryptedStream) Read(plain []byte) (n int, err error) {
	encrypted := make([]byte, len(plain))
	n, err = s.conn.Read(encrypted)
	if err != nil {
		return 0, err
	}

	fmt.Println("r encrypted ", encrypted)

	blockSize := s.cipher.BlockSize()

	for start := 0; start < n; start += blockSize {
		end := min(n, start+blockSize)
		encryptedBlock := encrypted[start:end]
		blockLength := len(encryptedBlock)
		plainBlock := make([]byte, blockSize)

		for len(encryptedBlock) < blockSize {
			encryptedBlock = append(encryptedBlock, 0)
		}

		s.cipher.Decrypt(plainBlock, encryptedBlock)

		copy(plain[start:], plainBlock[:blockLength])
	}

	fmt.Println("r plain ", plain)

	return n, nil
}

func (s *encryptedStream) Write(plain []byte) (n int, err error) {
	fmt.Println("w plain ", plain)

	encrypted := make([]byte, len(plain))
	blockSize := s.cipher.BlockSize()
	n = len(plain)

	for start := 0; start < n; start += blockSize {
		end := min(n, start+blockSize)
		plainBlock := plain[start:end]
		blockLength := len(plainBlock)
		encryptedBlock := make([]byte, blockSize)

		for len(plainBlock) < blockSize {
			plainBlock = append(plainBlock, 0)
		}

		s.cipher.Encrypt(encryptedBlock, plainBlock)

		copy(encrypted[start:], encryptedBlock[:blockLength])
	}

	fmt.Println("w encrypted ", encrypted)

	n, err = s.conn.Write(encrypted)
	if err != nil {
		return 0, err
	}

	return n, nil
}

type publicKeyInfo struct {
	Algorithm struct {
		Algorithm  asn1.ObjectIdentifier
		Parameters interface{} `asn1:"optional"`
	}

	SubjectPublicKey asn1.BitString
}

// Starts a connection to the specified address.
func (client *Client) connect() (err error) {
	if strings.Index(client.serverAddr, ":") < 0 {
		client.serverAddr += ":25565"
	}

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Connecting to %s via TCP\n", client.serverAddr)
	}

	client.netConn, err = net.Dial("tcp", client.serverAddr)
	if err != nil {
		return err
	}

	client.conn = client.netConn

	return nil
}

// Performs the 0x02 handshake transfer.
func (client *Client) handshake() (err error) {
	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Sending handshake packet\n")
	}

	host, portStr, err := net.SplitHostPort(client.serverAddr)
	if err != nil {
		return err
	}

	port, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return err
	}

	err = client.SendPacket(0x02, uint8(39), client.username, host, int32(port))
	if err != nil {
		return err
	}

	_, err = client.RecvPacket(0xFD)
	if err != nil {
		return err
	}

	err = client.RecvPacketData(&client.serverId, &client.serverKeyMessage, &client.serverVerifyToken)
	if err != nil {
		return err
	}

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Received encryption request\n")
	}

	return nil
}

// Generates a symmetric key and encrypts the verification token
func (client *Client) genKey() (err error) {
	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Decoding public key\n")
	}

	var pki publicKeyInfo
	_, err = asn1.Unmarshal(client.serverKeyMessage, &pki)
	if err != nil {
		return err
	}

	client.serverKey = new(rsa.PublicKey)
	_, err = asn1.Unmarshal(pki.SubjectPublicKey.Bytes, client.serverKey)
	if err != nil {
		return err
	}

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Generating encryption key\n")
	}

	client.sharedSecret = make([]byte, 16)

	_, err = rand.Reader.Read(client.sharedSecret)
	if err != nil {
		return err
	}

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Encrypting shared secret\n")
	}

	client.encryptedSharedSecret, err = rsa.EncryptPKCS1v15(rand.Reader, client.serverKey, client.sharedSecret)
	if err != nil {
		return err
	}

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Encrypting verification token\n")
	}

	client.encryptedVerifyToken, err = rsa.EncryptPKCS1v15(rand.Reader, client.serverKey, client.serverVerifyToken)
	if err != nil {
		return err
	}

	return nil
}

// Registers the server join with session.minecraft.net
func (client *Client) registerJoin() (err error) {
	if client.serverId != "-" {
		h := crypto.SHA1.New()
		fmt.Fprint(h, client.serverId)
		fmt.Fprint(h, client.sharedSecret)
		fmt.Fprint(h, client.serverKeyMessage)
		sum := h.Sum(nil)

		negative := false

		if sum[0] >= 0x80 {
			negative = true

			for i := 0; i < h.Size(); i++ {
				sum[i] = 255 - sum[i]
			}

			for i := h.Size() - 1; i >= 0; i-- {
				sum[i] += 1

				if sum[i] != 0 { // no overflow
					break
				}
			}
		}

		hexSum := hex.EncodeToString(sum)
		hexSum = strings.TrimLeft(hexSum, " ")
		if negative {
			hexSum = "-" + hexSum
		}

		params := url.Values{
			"user":      {client.username},
			"sessionId": {client.sessionId},
			"serverId":  {hexSum},
		}

		if client.DebugWriter != nil {
			fmt.Fprintf(client.DebugWriter, "Registering join with minecraft.net\n")
			fmt.Fprintf(client.DebugWriter, "GET http://session.minecraft.net/game/joinserver.jsp?%s\n", params.Encode())
		}

		resp, err := http.Get("http://session.minecraft.net/game/joinserver.jsp?" + params.Encode())
		if err != nil {
			return err
		}

		resp.Body.Close()
	}

	return nil
}

// Performs the 0xFC encryption response.
func (client *Client) encryptionResponse() (err error) {
	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Sending encryption response\n")
	}

	err = client.SendPacket(0xFC, client.encryptedSharedSecret, client.encryptedVerifyToken)
	if err != nil {
		return err
	}

	_, err = client.RecvPacket(0xFC)
	if err != nil {
		return err
	}

	var x, y []byte

	err = client.RecvPacketData(&x, &y)
	if err != nil {
		return err
	}

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Received encryption acknowledgement\nSwitching to encrypted transfer\n")
	}

	c, err := aes.NewCipher(client.sharedSecret)
	if err != nil {
		return err
	}

	client.conn = &encryptedStream{client.conn, c}

	return nil
}

// Performs the 0x01 login request.
func (client *Client) login() (err error) {
	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Sending client status packet\n")
	}

	//err = client.SendPacket(0x01, int32(29), client.username, "", int32(0), int32(0), int8(0), uint8(0), uint8(0))
	err = client.SendPacket(0xCD, uint8(0))
	if err != nil {
		return err
	}

	id, err := client.RecvPacket(0x01, 0xFF)
	if err != nil {
		return err
	}

	switch id {
	case 0xFF:
		if client.DebugWriter != nil {
			fmt.Fprintf(client.DebugWriter, "Received kick: login was rejected\n")
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

		if client.DebugWriter != nil {
			fmt.Fprintf(client.DebugWriter, "Received login packet\n")
		}
	}

	return nil
}

// Connects to a server.
func (client *Client) Join(addr string) (err error) {
	if client.conn != nil {
		client.Leave()
	}

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Joining server %s\n", addr)
	}

	client.serverAddr = addr
	client.PacketLogging = client.DebugWriter != nil

	err = client.connect()
	if err != nil {
		return err
	}

	err = client.handshake()
	if err != nil {
		return err
	}

	err = client.genKey()
	if err != nil {
		return err
	}

	err = client.registerJoin()
	if err != nil {
		return err
	}

	err = client.encryptionResponse()
	if err != nil {
		return err
	}

	err = client.login()
	if err != nil {
		return err
	}

	client.PacketLogging = false

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Joined!\n\nStarting receiver...\nStarting position sender...\n\n")
	}

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
				if !client.PlayerOnGround && client.serverMode == 0 {
					client.PlayerY -= 0.2
				}
			*/

			//fmt.Printf("sending...\n")

			err := client.SendPacket(0x0D, client.PlayerX, client.PlayerY, client.PlayerStance, client.PlayerZ, client.PlayerYaw, client.PlayerPitch, client.PlayerOnGround)
			if err != nil {
				client.ErrChan <- err
				continue
			}
		}
	}
}

// Sends a kick packet to the server before calling LeaveNoKick
func (client *Client) Leave() (err error) {
	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Disconnecting...\n")
	}

	err = client.SendPacket(0xFF, "github.com/kierdavis/go/minecraft woz 'ere")
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 100)

	return client.LeaveNoKick()
}

// Shuts down background processes before closing the connection.
func (client *Client) LeaveNoKick() (err error) {
	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Stopping position sender...\n")
	}

	// Tell PositionSender to stop
	client.stopPositionSender <- struct{}{}

	// Wait for a reply
	<-client.stopPositionSender

	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Closing connection...\n")
	}

	client.netConn.Close()
	client.netConn = nil
	client.conn = nil
	if client.DebugWriter != nil {
		fmt.Fprintf(client.DebugWriter, "Done!\n\n")
	}

	return nil
}
