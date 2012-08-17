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

type wavsample struct {
	Left  int16
	Right int16
}

func newWavHeader(dataLength uint, numChannels uint, bytesPerSample uint) (header *wavheader) {
	return &wavheader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     uint32(dataLength + 36),
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   uint16(numChannels),
		SampleRate:    uint32(SampleRate),
		ByteRate:      uint32(SampleRate * numChannels * bytesPerSample),
		BlockAlign:    uint16(numChannels * bytesPerSample),
		BitsPerSample: uint16(bytesPerSample * 8),
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: uint32(dataLength),
	}
}

func WriteWAV(w io.Writer, in chan Sample) (err error) {
	buffer := new(bytes.Buffer)

	for sample := range in {
		left := int16(sample.Left * 32767.0)
		right := int16(sample.Right * 32767.0)

		err := binary.Write(buffer, binary.LittleEndian, wavsample{left, right})
		if err != nil {
			return err
		}
	}

	l := uint(buffer.Len())

	err = binary.Write(w, binary.LittleEndian, newWavHeader(l, 2, 2))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buffer)
	if err != nil {
		return err
	}

	return nil
}

func WriteWAVMono(w io.Writer, in chan Sample) (err error) {
	buffer := new(bytes.Buffer)

	for sample := range in {
		x := int16(sample.Mono() * 32767.0)
		err := binary.Write(buffer, binary.LittleEndian, x)
		if err != nil {
			return err
		}
	}

	l := uint(buffer.Len())

	err = binary.Write(w, binary.LittleEndian, newWavHeader(l, 1, 2))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buffer)
	if err != nil {
		return err
	}

	return nil
}
