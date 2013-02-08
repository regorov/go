package mcquery

import (
	"fmt"
	"net"
)

const MaxRetries = 3

type Connection struct {
	Conn      net.Conn
	ID        uint32
	Retries   int
	Challenge uint32
}

type Stat struct {
	MOTD       string
	GameType   string
	Map        string
	NumPlayers int
	MaxPlayers int
	HostPort   int
	HostName   string
}

func getUint32(b []byte) (n uint32) {
	n = uint32(b[0])
	n |= uint32(b[1]) << 8
	n |= uint32(b[2]) << 16
	n |= uint32(b[3]) << 24
	return n
}

func putUint32(b []byte, n uint32) {
	b[0] = byte(n)
	b[1] = byte(n >> 8)
	b[2] = byte(n >> 16)
	b[3] = byte(n >> 24)
}

func NewSession(addr string) (c *Connection, err error) {
	conn, err := net.Conn("udp", addr)
	if err != nil {
		return nil, err
	}

	c = &Connection{
		Conn:      conn,
		ID:        0,
		Retries:   0,
		Challenge: 0,
	}

	err = c.Handshake()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Connection) WritePacket(t byte, payload []byte) (err error) {
	message := make([]byte, len(payload)+7)
	message[0] = 0xFE
	message[1] = 0xFD
	message[2] = t

	putUint32(message[3:7], c.ID)
	copy(message[7:], payload)

	_, err = c.Conn.Write(message)
	return err
}

func (c *Connection) ReadPacket() (t byte, id uint32, payload []byte, err error) {
	buffer := make([]byte, 2048)
	n, err := c.Conn.Read(buffer)
	if err != nil {
		return 0, 0, nil, err
	}

	buffer = buffer[:n]
	t = buffer[0]
	id = getUint32(buffer[1:5])
	payload = buffer[5:]

	return t, id, payload, nil
}

func (c *Connection) Handshake() (err error) {
	c.ID += 1
	err = c.WritePacket(9, nil)
	if err != nil {
		return err
	}

	t, id, payload, err := c.ReadPacket()
	if err != nil {
		e, ok := err.(net.Error)
		if ok && e.Timeout() {
			c.Retries++

			if c.Retries == MaxRetries {
				return fmt.Errorf("Retry limit reached - server down?")
			}

			return c.Handshake()
		}

		return err
	}

	c.Retries = 0
	c.Challenge, err = strconv.ParseUint(string(payload[:len(payload)-1]), 10, 32)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) BasicStat() (r *Stat, err error) {
	payload := make([]byte, 4)
	putUint32(payload, c.Challenge)

	err = c.WritePacket(0, payload)
	if err != nil {
		return nil, err
	}

	t, id, payload, err := c.ReadPacket()
	if err != nil {
		err = c.Handshake()
		if err != nil {
			return nil, err
		}

		return c.BasicStat()
	}

	r = new(Stat)
	parts := bytes.SplitN(payload, []byte{0}, 6)

	r.MOTD = string(parts[0])
	r.GameType = string(parts[1])
	r.Map = string(parts[2])

	r.NumPlayers, err = strconv.ParseInt(string(parts[3]), 10, 0)
	if err != nil {
		return nil, err
	}

	r.MaxPlayers, err = strconv.ParseInt(string(parts[4]), 10, 0)
	if err != nil {
		return nil, err
	}

	payload = parts[5]
	r.HostPort = int(uint16(payload[0]) | (uint16(payload[1]) << 8))
	r.HostName = string(payload[2 : len(payload)-1])

	return r, nil
}

/*
func (c *Connection) FullStat() (r *Stat, err error) {
	payload := make([]byte, 8)
	putUint32(payload, c.Challenge)

	err = c.WritePacket(0, payload)
	if err != nil {
		return nil, err
	}

	t, id, payload, err := c.ReadPacket()
	if err != nil {
		err = c.Handshake()
		if err != nil {
			return nil, err
		}

		return c.BasicStat()
	}

	payload = payload[11:]

	r = new(Stat)

	return r, nil
}
*/
