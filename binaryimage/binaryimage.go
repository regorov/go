package binaryimage

import (
	"encoding/hex"
	"fmt"
	"github.com/kierdavis/goutil"
	"io"
)

func max(a, b uint64) (r uint64) {
	if a > b {
		return a
	}

	return b
}

func min(a, b uint64) (r uint64) {
	if a < b {
		return a
	}

	return b
}

type Image struct {
	data map[uint64]byte
	max  uint64
}

func New() (image *Image) {
	return &Image{make(map[uint64]byte), 0}
}

func (image *Image) Put(addr uint64, data byte) {
	image.data[addr] = data
	image.max = max(image.max, addr)
}

func (image *Image) PutBytes(addr uint64, data []byte) (n uint64) {
	for i, b := range data {
		image.Put(addr+uint64(i), b)
		n++
	}

	return n
}

func (image *Image) Get(addr uint64) (data byte) {
	return image.data[addr]
}

func (image *Image) GetBytes(addr uint64, data []byte) (n uint64) {
	n = uint64(len(data))
	if n > image.max-addr {
		n = image.max - addr
	}

	for i := uint64(0); i < n; i++ {
		data[i] = image.Get(addr + i)
	}

	return n
}

func (image *Image) Max() (max uint64) {
	return image.max
}

func (image *Image) ReadRaw(r io.Reader) (err error) {
	_, err = io.Copy(NewImageWriter(image), r)
	return err
}

func (image *Image) WriteRaw(w io.Writer) (err error) {
	_, err = io.Copy(w, NewImageReader(image))
	return err
}

func (image *Image) ReadIHex(r io.Reader) (err error) {
	lineChan, errChan := util.IterLines(r)
	lineno := 0

	var baseAddress uint64

	for line := range lineChan {
		lineno++

		if len(line) == 0 {
			continue
		}

		if line[0] != ':' {
			return fmt.Errorf("[line %d] Invalid record start byte: expected a colon (:), found %q", lineno, line[0])
		}

		record, err := hex.DecodeString(line[1:])
		if err != nil {
			return err
		}

		last := len(record) - 1

		if ihexChecksum(record[:last]) != record[last] {
			return fmt.Errorf("[line %d] Checksum mismatch", lineno)
		}

		length := record[0]
		address := baseAddress
		address += uint64(record[1]) << 8
		address += uint64(record[2])
		recordType := record[3]
		data := record[4:last]

		if len(data) != int(length) {
			return fmt.Errorf("[line %d] Data length mismatch", lineno)
		}

		switch recordType {
		case 0x00:
			image.PutBytes(address, data)

		case 0x01:
			return nil

		case 0x02:
			if length != 2 {
				return fmt.Errorf("[line %d] Expected data of length 2 for 02 record", lineno)
			}

			baseAddress = uint64(data[0]) << (8 + 4)
			baseAddress += uint64(data[1]) << 4

		case 0x04:
			if length != 2 {
				return fmt.Errorf("[line %d] Expected data of length 2 for 04 record", lineno)
			}

			baseAddress = uint64(data[0]) << (8 + 16)
			baseAddress += uint64(data[1]) << 16
		}
	}

	return <-errChan
}

func (image *Image) WriteIHex(w io.Writer) (err error) {
	var addr, baseAddr uint64

	buffer := make([]byte, 16)

	for addr <= image.max {
		thisBase := addr & 0xFFFF0000
		if thisBase != baseAddr {
			err = emitIHexRecord(w, 0x04, 0, []byte{byte(thisBase >> 24), byte(thisBase >> 16)})
			if err != nil {
				return err
			}

			baseAddr = thisBase
		}

		l := image.GetBytes(addr, buffer)
		err = emitIHexRecord(w, 0x00, uint16(addr&0xFFFF), buffer[:l])
		if err != nil {
			return err
		}

		addr += l
	}

	return emitIHexRecord(w, 0x01, 0, nil)
}

type ImageWriter struct {
	image  *Image
	offset uint64
}

func NewImageWriter(image *Image) (w *ImageWriter) {
	return &ImageWriter{image, 0}
}

func (w *ImageWriter) Offset() (offset uint64) {
	return w.offset
}

func (w *ImageWriter) SetOffset(offset uint64) {
	w.offset = offset
}

func (w *ImageWriter) Write(data []byte) (n int, err error) {
	l := w.image.PutBytes(w.offset, data)
	w.offset += l
	return int(l), nil
}

type ImageReader struct {
	image  *Image
	offset uint64
}

func NewImageReader(image *Image) (r *ImageReader) {
	return &ImageReader{image, 0}
}

func (r *ImageReader) Offset() (offset uint64) {
	return r.offset
}

func (r *ImageReader) SetOffset(offset uint64) {
	r.offset = offset
}

func (r *ImageReader) Seek(offset uint64) {
	r.offset = offset
}

func (r *ImageReader) Read(data []byte) (n int, err error) {
	l := r.image.GetBytes(r.offset, data)
	r.offset += l
	return int(l), nil
}

func ReadRaw(r io.Reader) (image *Image, err error) {
	image = New()
	err = image.ReadRaw(r)
	return image, err
}

func ReadIHex(r io.Reader) (image *Image, err error) {
	image = New()
	err = image.ReadIHex(r)
	return image, err
}

func ihexChecksum(data []byte) (sum byte) {
	for _, b := range data {
		sum += b
	}

	return ^sum
}

func emitIHexRecord(w io.Writer, recordType byte, address uint16, data []byte) (err error) {
	record := make([]byte, len(data)+5)
	record[0] = byte(len(data))
	record[1] = byte(address >> 8)
	record[2] = byte(address)
	record[3] = recordType
	copy(record[4:], data)
	record[len(record)-1] = ihexChecksum(record[:len(record)-1])

	buffer := make([]byte, hex.EncodedLen(len(record))+1)
	buffer[0] = ':'
	hex.Encode(buffer[1:], record)

	_, err = w.Write(buffer)
	return err
}
