package minecraft

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
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

	debug     bool
	username  string
	sessionId string
	serverId  string

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

func Login(username string, password string, debug bool) (client *Client, err error) {
	params := url.Values{
		"user":     {username},
		"password": {password},
		"version":  {"13"},
	}

	if debug {
		fmt.Printf("Logging in to minecraft.net as '%s'\n", username)
		fmt.Printf("POST https://login.minecraft.net username=%s&password=...&version=13\n", username)
	}

	resp, err := http.PostForm("https://login.minecraft.net", params)
	if err != nil {
		return nil, err
	}

	respdata, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	respparts := strings.Split(string(respdata), ":")

	client = newClient(respparts[2], respparts[3], debug)

	if debug {
		fmt.Printf("Session ID: %s\n\n", client.sessionId)
	}

	return client, nil
}

func LoginOffline(debug bool) (client *Client) {
	return newClient("Player", "", debug)
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

func (client *Client) HTTPKeepAlive() {
	params := url.Values{
		"name":    {client.username},
		"session": {client.sessionId},
	}

	sessionURL := "https://login.minecraft.net/session?" + params.Encode()
	ticker := time.NewTicker(Tick * 6000)

	for {
		select {
		case <-ticker.C:
			if client.debug {
				fmt.Printf("(HTTP keep-alive) GET %s\n", sessionURL)
			}

			resp, err := http.Get(sessionURL)
			if err != nil {
				client.ErrChan <- err
			}

			resp.Body.Close()

		case <-client.stopHTTPKeepAlive:
			return
		}
	}
}

func (client *Client) Join(addr string) (err error) {
	if client.conn != nil {
		client.Leave()
	}

	if client.debug {
		fmt.Printf("Joining server %s\n", addr)
	}

	// **************** Connect to server

	if strings.Index(addr, ":") < 0 {
		addr += ":25565"
	}

	if client.debug {
		fmt.Printf("Connecting to %s via TCP\n", addr)
	}

	client.conn, err = net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	// **************** 0x02 Handshake

	if client.debug {
		fmt.Printf("Sending handshake packet\n")
	}

	err = client.SendPacket(0x02, client.username+";"+addr)
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

	// **************** Register with minecraft.net

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

	// **************** 0x01 Login Request

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

	if client.debug {
		fmt.Printf("Joined!\n\n")
	}

	go client.Receiver()

	return nil
}

func (client *Client) Chat(msg string) (err error) {
	return client.SendPacket(0x03, msg)
}

func (client *Client) PositionSender() {
	ticker := time.NewTicker(time.Millisecond * 50)

	for {
		select {
		case <-client.stopPositionSender:
			client.stopPositionSender <- struct{}{}
			return

		case <-ticker.C:
			err := client.SendPacket(0x0D, client.playerX, client.playerY, client.playerStance, client.playerZ, client.playerYaw, client.playerPitch, client.playerOnGround)
			if err != nil {
				client.ErrChan <- err
				continue
			}
		}
	}
}

func (client *Client) Receiver() {
	defer func() {
		if client.conn != nil {
			client.Leave()
		}
	}()

	for {
		id, err := client.RecvAnyPacket()
		if err != nil {
			client.ErrChan <- err
			continue
		}

		err = nil
		switch id {
		case 0x00:
			err = client.handleKeepAlivePacket()
		case 0x03:
			err = client.handleChatMessagePacket()
		case 0x04:
			err = client.handleTimeUpdatePacket()
		case 0x05:
			err = client.handleEntityEquipmentPacket()
		case 0x06:
			err = client.handleSpawnPositionPacket()
		case 0x08:
			err = client.handleUpdateHealthPacket()
		case 0x09:
			err = client.handleRespawnPacket()
		case 0x0D:
			err = client.handlePlayerPositionLookPacket()
		case 0x11:
			err = client.handleUseBedPacket()
		case 0x12:
			err = client.handleAnimationPacket()
		case 0x14:
			err = client.handleSpawnNamedEntityPacket()
		case 0x15:
			err = client.handleSpawnDroppedItemPacket()
		case 0x16:
			err = client.handleCollectItemPacket()
		case 0x17:
			err = client.handleSpawnObjectPacket()
		case 0x18:
			err = client.handleSpawnMobPacket()
		case 0x19:
			err = client.handleSpawnPaintingPacket()
		case 0x1A:
			err = client.handleSpawnExperienceOrbPacket()
		case 0x1C:
			err = client.handleEntityVelocityPacket()
		case 0x1D:
			err = client.handleDestroyEntityPacket()
		case 0x1E:
			err = client.handleEntityPacket()
		case 0x1F:
			err = client.handleEntityRelativeMovePacket()
		case 0x20:
			err = client.handleEntityLookPacket()
		case 0x21:
			err = client.handleEntityLookRelativeMovePacket()
		case 0x22:
			err = client.handleEntityTeleportPacket()
		case 0x23:
			err = client.handleEntityHeadLookPacket()
		case 0x26:
			err = client.handleEntityStatusPacket()
		case 0x27:
			err = client.handleAttachEntityPacket()
		case 0x28:
			err = client.handleEntityMetadataPacket()
		case 0x29:
			err = client.handleEntityEffectPacket()
		case 0x2A:
			err = client.handleRemoveEntityEffectPacket()
		case 0x2B:
			err = client.handleSetExperiencePacket()
		case 0x32:
			err = client.handleMapColumnAllocationPacket()
		case 0x33:
			err = client.handleMapChunksPacket()
		case 0x34:
			err = client.handleMultiBlockChangePacket()
		case 0x35:
			err = client.handleBlockChangePacket()
		case 0x36:
			err = client.handleBlockActionPacket()
		case 0x3C:
			err = client.handleExplosionPacket()
		case 0x3D:
			err = client.handleSoundParticleEffectPacket()
		case 0x46:
			err = client.handleChangeGameStatePacket()
		case 0x47:
			err = client.handleThunderboltPacket()
		case 0x64:
			err = client.handleOpenWindowPacket()
		case 0x65:
			err = client.handleCloseWindowPacket()
		case 0x67:
			err = client.handleSetSlotPacket()
		case 0x68:
			err = client.handleSetWindowItemsPacket()
		case 0x69:
			err = client.handleUpdateWindowPropertyPacket()
		case 0x6A:
			err = client.handleConfirmTransactionPacket()
		case 0x6B:
			err = client.handleCreativeInventoryActionPacket()
		case 0x82:
			err = client.handleUpdateSignPacket()
		case 0x83:
			err = client.handleItemDataPacket()
		case 0x84:
			err = client.handleUpdateTileEntityPacket()
		case 0xC8:
			err = client.handleIncrementStatisticPacket()
		case 0xC9:
			err = client.handlePlayerListItemPacket()
		case 0xCA:
			err = client.handlePlayerAbilitiesPacket()
		case 0xFA:
			err = client.handlePluginMessagePacket()
		case 0xFF:
			err = client.handleKickPacket()
		default:
			fmt.Fprintf(os.Stderr, "Ignoring unhandled packet with id 0x%02X", id)
		}

		if err == Stop {
			return
		}

		if err != nil {
			client.ErrChan <- err
		}
	}
}

func (client *Client) Leave() (err error) {
	err = client.SendPacket(0xFF, "github.com/kierdavis/minecraft woz 'ere")
	if err != nil {
		return err
	}

	return client.LeaveNoKick()
}

func (client *Client) LeaveNoKick() (err error) {
	client.stopPositionSender <- struct{}{}

	<-client.stopPositionSender

	client.conn.Close()
	client.conn = nil

	return nil
}

func (client *Client) SendPacket(id byte, fields ...interface{}) (err error) {
	err = binary.Write(client.conn, binary.BigEndian, id)
	if err != nil {
		return err
	}

	for _, ifield := range fields {
		switch field := ifield.(type) {
		case uint8:
			err = binary.Write(client.conn, binary.BigEndian, field)
		case uint16:
			err = binary.Write(client.conn, binary.BigEndian, field)
		case uint32:
			err = binary.Write(client.conn, binary.BigEndian, field)
		case uint64:
			err = binary.Write(client.conn, binary.BigEndian, field)

		case int8:
			err = binary.Write(client.conn, binary.BigEndian, field)
		case int16:
			err = binary.Write(client.conn, binary.BigEndian, field)
		case int32:
			err = binary.Write(client.conn, binary.BigEndian, field)
		case int64:
			err = binary.Write(client.conn, binary.BigEndian, field)

		case float32:
			err = binary.Write(client.conn, binary.BigEndian, field)
		case float64:
			err = binary.Write(client.conn, binary.BigEndian, field)

		case string:
			err = binary.Write(client.conn, binary.BigEndian, uint16(len(field)))

			i := 0
			for i < len(field) {
				if err != nil {
					return err
				}

				r, n := utf8.DecodeRuneInString(field[i:])
				i += n
				err = binary.Write(client.conn, binary.BigEndian, uint16(r))
			}

		case bool:
			u := uint8(0)
			if field {
				u = 1
			}

			err = binary.Write(client.conn, binary.BigEndian, u)

		default:
			err = fmt.Errorf("Invalid type for SendPacket: %T", ifield)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) RecvAnyPacket() (id byte, err error) {
	err = binary.Read(client.conn, binary.BigEndian, &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (client *Client) RecvPacket(acceptIds ...byte) (id byte, err error) {
	id, err = client.RecvAnyPacket()
	if err != nil {
		return 0, err
	}

	accepted := false
	for _, acceptId := range acceptIds {
		if acceptId == id {
			accepted = true
		}
	}

	if !accepted {
		return 0, fmt.Errorf("Unexpected 0x%02X packet", id)
	}

	return id, nil
}

func (client *Client) RecvPacketData(fields ...interface{}) (err error) {
	for _, ifield := range fields {
		switch field := ifield.(type) {
		case *uint8:
			err = binary.Read(client.conn, binary.BigEndian, field)
		case *uint16:
			err = binary.Read(client.conn, binary.BigEndian, field)
		case *uint32:
			err = binary.Read(client.conn, binary.BigEndian, field)
		case *uint64:
			err = binary.Read(client.conn, binary.BigEndian, field)

		case *int8:
			err = binary.Read(client.conn, binary.BigEndian, field)
		case *int16:
			err = binary.Read(client.conn, binary.BigEndian, field)
		case *int32:
			err = binary.Read(client.conn, binary.BigEndian, field)
		case *int64:
			err = binary.Read(client.conn, binary.BigEndian, field)

		case *float32:
			err = binary.Read(client.conn, binary.BigEndian, field)
		case *float64:
			err = binary.Read(client.conn, binary.BigEndian, field)

		case *string:
			var l uint16
			err = binary.Read(client.conn, binary.BigEndian, &l)

			runes := make([]rune, l)
			reallen := 0

			for i := uint16(0); i < l; i++ {
				if err != nil {
					return err
				}

				var r16 uint16
				err = binary.Read(client.conn, binary.BigEndian, &r16)

				r := rune(r16)
				runes[i] = r
				reallen += utf8.RuneLen(r)
			}

			b := make([]byte, reallen)
			pos := 0

			for _, r := range runes {
				pos += utf8.EncodeRune(b[pos:], r)
			}

			*field = string(b)

		case *bool:
			var b uint8
			err = binary.Read(client.conn, binary.BigEndian, &b)
			*field = b == 1

		case *Slot:
			err = binary.Read(client.conn, binary.BigEndian, &field.ID)
			if err != nil {
				return err
			}

			if field.ID != -1 {
				err = binary.Read(client.conn, binary.BigEndian, &field.Count)
				if err != nil {
					return err
				}

				err = binary.Read(client.conn, binary.BigEndian, &field.Damage)
				if err != nil {
					return err
				}

				if (256 <= field.ID && field.ID <= 259) || (267 <= field.ID && field.ID <= 279) || (283 <= field.ID && field.ID <= 286) || (290 <= field.ID && field.ID <= 294) || (298 <= field.ID && field.ID <= 317) || field.ID == 261 || field.ID == 359 || field.ID == 346 {
					var l int16
					err = binary.Read(client.conn, binary.BigEndian, &l)
					if err != nil {
						return err
					}

					if l == -1 {
						field.Data = make([]byte, 0)

					} else {
						field.Data = make([]byte, l)
						_, err = client.conn.Read(field.Data)
						if err != nil {
							return err
						}
					}
				}
			}

		default:
			err = fmt.Errorf("Invalid type for RecvPacketData: %T", ifield)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) RecvEntityMetadata() (metadata Metadata, err error) {
	metadata = make(Metadata)

	for {
		var b uint8
		err = client.RecvPacketData(&b)
		if err != nil {
			return nil, err
		}

		if b == 127 {
			break
		}

		key := b & 0x1f

		switch b >> 5 {
		case 0:
			var value int8
			err = client.RecvPacketData(&value)
			metadata[key] = value

		case 1:
			var value int16
			err = client.RecvPacketData(&value)
			metadata[key] = value

		case 2:
			var value int32
			err = client.RecvPacketData(&value)
			metadata[key] = value

		case 3:
			var value float32
			err = client.RecvPacketData(&value)
			metadata[key] = value

		case 4:
			var value string
			err = client.RecvPacketData(&value)
			metadata[key] = value

		case 5:
			var id, damage int16
			var count int8

			err = client.RecvPacketData(&id, &count, &damage)
			metadata[key] = Slot{id, count, damage, nil}

		case 6:
			var x, y, z int32
			err = client.RecvPacketData(&x, &y, &z)
			metadata[key] = Position{x, y, z}
		}

		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}
