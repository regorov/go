package minecraft

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const Tick = time.Second / 20

var Stop = errors.New("STOP")

type Signal chan struct{}
type Metadata map[uint8]interface{}

type Slot struct {
	ID     int16
	Count  int8
	Damage int16
	Data   []byte
}

type Position struct {
	X int32
	Y int32
	Z int32
}

type ExplosionRecord struct {
	X int8
	Y int8
	Z int8
}

type Client struct {
	ErrChan chan error
	conn    net.Conn

	stopHTTPKeepAlive  Signal
	stopPositionSender Signal

	debug      bool
	username   string
	sessionId  string
	serverId   string
	serverAddr string

	entityID   int32
	levelType  string
	serverMode int32
	dimension  int32
	difficulty int8
	maxPlayers uint8

	playerX        float64
	playerY        float64
	playerZ        float64
	playerStance   float64
	playerYaw      float32
	playerPitch    float32
	playerOnGround bool
}

func newClient(username string, sessionId string, debug bool) (client *Client) {
	client = &Client{
		ErrChan:            make(chan error),
		stopHTTPKeepAlive:  make(Signal),
		stopPositionSender: make(Signal),
		debug:              debug,
		username:           username,
		sessionId:          sessionId,
	}

	go client.HTTPKeepAlive()

	return client
}

func ANSIEscapes(input string) (output string) {
	start := 0

	for {
		end := strings.Index(input[start:], "\xC2\xA7")
		if end < 0 {
			break
		}

		output += input[start:end]

		switch input[end+2] {
		case '0':
			output += "\x1b[21m\x1b[30m"
		case '1':
			output += "\x1b[21m\x1b[34m"
		case '2':
			output += "\x1b[21m\x1b[32m"
		case '3':
			output += "\x1b[21m\x1b[36m"
		case '4':
			output += "\x1b[21m\x1b[31m"
		case '5':
			output += "\x1b[21m\x1b[35m"
		case '6':
			output += "\x1b[21m\x1b[33m"
		case '7':
			output += "\x1b[21m\x1b[37m"
		case '8':
			output += "\x1b[1m\x1b[30m"
		case '9':
			output += "\x1b[1m\x1b[34m"
		case 'a', 'A':
			output += "\x1b[1m\x1b[32m"
		case 'b', 'B':
			output += "\x1b[1m\x1b[36m"
		case 'c', 'C':
			output += "\x1b[1m\x1b[31m"
		case 'd', 'D':
			output += "\x1b[1m\x1b[35m"
		case 'e', 'E':
			output += "\x1b[1m\x1b[33m"
		case 'f', 'F':
			output += "\x1b[1m\x1b[37m"
		}

		start = end + 3
	}

	output += input[start:] + "\x1b[21m\x1b[39m"
	return output
}

func ScanServer(addr string) (description string, onlineUsers int, maxUsers int, err error) {
	if strings.Index(addr, ":") < 0 {
		addr += ":25565"
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", 0, 0, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte{0xFE})
	if err != nil {
		return "", 0, 0, err
	}

	data := make([]byte, 256)
	n, err := conn.Read(data)
	if err != nil {
		return "", 0, 0, err
	}

	if data[0] != 0xFF {
		return "", 0, 0, fmt.Errorf("Expected kick packet (0xFF)")
	}

	realLen := 0
	runes := make([]rune, 0, (n-3)/2)

	for i := 3; i < n; i += 2 {
		r := (rune(data[i]) << 8) | rune(data[i+1])
		runes = append(runes, r)
		realLen += utf8.RuneLen(r)
	}

	b := make([]byte, realLen)
	pos := 0

	for _, r := range runes {
		pos += utf8.EncodeRune(b[pos:], r)
	}

	parts := strings.Split(string(b), "\xC2\xA7")

	onlineUsers64, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return "", 0, 0, err
	}

	maxUsers64, err := strconv.ParseInt(parts[2], 10, 0)
	if err != nil {
		return "", 0, 0, err
	}

	return parts[0], int(onlineUsers64), int(maxUsers64), nil
}

// Sends a chat message
func (client *Client) Chat(msg string) (err error) {
	return client.SendPacket(0x03, msg)
}
