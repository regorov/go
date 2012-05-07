package binaryloader

import (
    "bufio"
    "github.com/kierdavis/go/ihex"
    "os"
    "path/filepath"
)

func Load(filename string) (data []byte, err error) {
    ext := filepath.Ext(filename)

    switch ext {
    case ".hex":
        return LoadHex(filename)

    default:
        return LoadBinary(filename)
    }

    return nil, nil
}

func LoadHex(filename string) (data []byte, err error) {
    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    reader := bufio.NewReader(f)
    ix, err := ihex.ReadIHex(reader)
    if err != nil {
        return nil, err
    }

    data = ix.ExtractDataToEnd(0)
    return data, nil
}

func LoadBinary(filename string) (data []byte, err error) {
    fi, err := os.Stat(filename)
    if err != nil {
        return nil, err
    }

    data = make([]byte, fi.Size())

    f, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    _, err = f.Read(data)
    if err != nil {
        return nil, err
    }

    return data, nil
}
