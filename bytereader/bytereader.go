package bytereader

import (
    "io"
)

const (
    BIG_ENDIAN = 0
    LITTLE_ENDIAN = 1
)

type ByteReader struct {
    reader io.ReadSeeker
    endian int8
    pos int64
    posStack []int64
}

func NewByteReader(reader io.ReadSeeker, endian int8) (byteReader *ByteReader) {
    return &ByteReader{reader: reader, endian: endian, pos: 0, posStack: make([]int64, 0, 4)}
}

func (byteReader *ByteReader) ReadBytes(num int) (value []byte, err error) {
    value = make([]byte, num)
    _, err = byteReader.reader.Read(value); if err != nil {return}
    
    byteReader.pos += int64(num)
    return value, nil
}

func (byteReader *ByteReader) ReadUnsigned8() (value uint8, err error) {
    buffer, err := byteReader.ReadBytes(1); if err != nil {return}
    return buffer[0], nil
}

func (byteReader *ByteReader) ReadUnsigned16() (value uint16, err error) {
    buffer, err := byteReader.ReadBytes(2); if err != nil {return}
    
    if byteReader.endian == BIG_ENDIAN {
        return  (uint16(buffer[0]) << 8) |
                 uint16(buffer[1]), nil
    
    } else if byteReader.endian == LITTLE_ENDIAN {
        return  (uint16(buffer[1]) << 8) |
                 uint16(buffer[0]), nil
    }
    
    return 0, NewError(E_INVALID_ENDIANNESS, "ByteReader.ReadUnsigned16: Invalid value for endianness parameter specified at initialisation (should be either BIG_ENDIAN or LITTLE_ENDIAN)")
}

func (byteReader *ByteReader) ReadUnsigned32() (value uint32, err error) {
    buffer, err := byteReader.ReadBytes(4); if err != nil {return}
    
    if byteReader.endian == BIG_ENDIAN {
        return  (uint32(buffer[0]) << 24) |
                (uint32(buffer[1]) << 16) |
                (uint32(buffer[2]) << 8) |
                 uint32(buffer[3]), nil
    
    } else if byteReader.endian == LITTLE_ENDIAN {
        return  (uint32(buffer[3]) << 24) |
                (uint32(buffer[2]) << 16) |
                (uint32(buffer[1]) << 8) |
                 uint32(buffer[0]), nil
    }
    
    return 0, NewError(E_INVALID_ENDIANNESS, "ByteReader.ReadUnsigned32: Invalid value for endianness parameter specified at initialisation (should be either BIG_ENDIAN or LITTLE_ENDIAN)")
}

func (byteReader *ByteReader) ReadUnsigned64() (value uint64, err error) {
    buffer, err := byteReader.ReadBytes(8); if err != nil {return}
    
    if byteReader.endian == BIG_ENDIAN {
        return  (uint64(buffer[0]) << 56) |
                (uint64(buffer[1]) << 48) |
                (uint64(buffer[2]) << 40) |
                (uint64(buffer[3]) << 32) |
                (uint64(buffer[4]) << 24) |
                (uint64(buffer[5]) << 16) |
                (uint64(buffer[6]) << 8) |
                 uint64(buffer[7]), nil
    
    } else if byteReader.endian == LITTLE_ENDIAN {
        return  (uint64(buffer[7]) << 56) |
                (uint64(buffer[6]) << 48) |
                (uint64(buffer[5]) << 40) |
                (uint64(buffer[4]) << 32) |
                (uint64(buffer[3]) << 16) |
                (uint64(buffer[2]) << 16) |
                (uint64(buffer[1]) << 8) |
                 uint64(buffer[0]), nil
    }
    
    return 0, NewError(E_INVALID_ENDIANNESS, "ByteReader.ReadUnsigned64: Invalid value for endianness parameter specified at initialisation (should be either BIG_ENDIAN or LITTLE_ENDIAN)")
}

func (byteReader *ByteReader) ReadSigned8() (value int8, err error) {
    x, err := byteReader.ReadUnsigned8(); if err != nil {return}
    return int8(x), nil
}

func (byteReader *ByteReader) ReadSigned16() (value int16, err error) {
    x, err := byteReader.ReadUnsigned16(); if err != nil {return}
    return int16(x), nil
}

func (byteReader *ByteReader) ReadSigned32() (value int32, err error) {
    x, err := byteReader.ReadUnsigned32(); if err != nil {return}
    return int32(x), nil
}

func (byteReader *ByteReader) ReadSigned64() (value int64, err error) {
    x, err := byteReader.ReadUnsigned64(); if err != nil {return}
    return int64(x), nil
}

func (byteReader *ByteReader) Pos() (pos int64) {
    return byteReader.pos
}

func (byteReader *ByteReader) JumpTo(offset int64) (err error) {
    _, err = byteReader.reader.Seek(offset, 0); if err != nil {return}
    byteReader.pos = offset
    return nil
}

func (byteReader *ByteReader) PushPos() {
    byteReader.posStack = append(byteReader.posStack, byteReader.pos)
}

func (byteReader *ByteReader) PopPos() (err error) {
    idx := len(byteReader.posStack) - 1
    byteReader.JumpTo(byteReader.posStack[idx]); if err != nil {return}
    byteReader.posStack = byteReader.posStack[:idx]
    return nil
}
