package sound

import (
	"bytes"
	"encoding/binary"
	"io"
)

type wavheader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte
	Subchunk2Size uint32
}

func WriteWAV(w io.Writer, in chan float64) (err error) {
	buffer := new(bytes.Buffer)

	for sample := range in {
		n := int16(sample * 32767.0)
		err := binary.Write(buffer, binary.LittleEndian, n)
		if err != nil {
			return err
		}
	}

	l := uint32(buffer.Len())

	header := &wavheader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     l + 36,
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   1,
		SampleRate:    44100,
		ByteRate:      88200,
		BlockAlign:    2,
		BitsPerSample: 16,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: l,
	}

	err = binary.Write(w, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buffer)
	if err != nil {
		return err
	}

	return nil
}
