package minecraft

import (
	"encoding/binary"
	"fmt"
	"unicode/utf8"
)

// Sends a packet on the channel. The types of the fields are determined by runtime reflection.
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

// Reads a packet ID and returns it.
func (client *Client) RecvAnyPacket() (id byte, err error) {
	err = binary.Read(client.conn, binary.BigEndian, &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Reads a packet ID and ensures that it matches one of the ones passed as arguments.
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

// Reads data from the connection and stores in the arguments (which should be pointers to
// appropriately typed values)
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

// Reads an entity metadata block, returning it as a map of uint8s to appropriately typed values.
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
