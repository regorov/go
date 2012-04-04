package ihex

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "strings"
)

type IHex struct {
    areas map[uint][]byte
}

func NewIHex() (ix *IHex) {
    ix = new(IHex)
    ix.areas = make(map[uint][]byte)
    return ix
}

func ReadIHex(reader *bufio.Reader) (ix *IHex, err error) {
    ix = NewIHex()
    
    for {
        line, _, err := reader.ReadLine()
        if err == io.EOF {break}
        if err != nil {return nil, err}
        
        sline := strings.Trim(string(line), " \t")
        if len(sline) == 0 {continue}
        
        t, addr, data, err := parseLine(sline)
        if err != nil {return nil, err}
        
        switch t {
        case 0x00:
            ix.InsertData(addr, data)
        }
    }
    
    return ix, nil
}

func parseLine(rawline string) (t uint, addr uint, data []byte, err error) {
    if rawline[0] != ':' {
        return 0, 0, nil, errors.New(fmt.Sprintf("Invalid line start character: %q", rawline[0]))
    }
    
    line := make([]byte, (len(rawline) - 1) / 2)
    _, err = fmt.Sscanf(rawline, ":%x", &line)
    if err != nil {return 0, 0, nil, err}
    
    length := uint(line[0])
    addr = (uint(line[1]) << 8) | uint(line[2])
    t = uint(line[3])
    
    data = line[4:4 + length]
    cs1 := uint(line[4 + length])
    cs2 := calcChecksum(line[:4 + length])
    
    if cs1 != cs2 {
        return 0, 0, nil, errors.New("Checksums do not match")
    }
    
    return t, addr, data, nil
}

func calcChecksum(data []byte) (checksum uint) {
    total := uint(0)
    
    for _, x := range data {
        total += uint(x)
    }
    
    return (-total) & 0xFF
}

func (ix *IHex) GetSize() (size uint) {
    size = 0
    
    for addr, data := range ix.areas {
        if addr + uint(len(data)) > size {
            size = addr + uint(len(data))
        }
    }
    
    return size
}

func (ix *IHex) ExtractData(start uint, end uint) (result []byte) {
    result = make([]byte, 0, end - start)
    
    for addr, data := range ix.areas {
        if addr >= start && addr < end {
            copy(result[start:end], data[start-addr:end-addr])
        }
    }
    
    return result
}

func (ix *IHex) ExtractDataToEnd(start uint) (result []byte) {
    end := uint(0)
    result = make([]byte, ix.GetSize() - start)
    
    for addr, data := range ix.areas {
        if addr >= start {
            if addr + uint(len(data)) > end {
                end = addr + uint(len(data))
            }
            
            copy(result[start:end], data[start-addr:end-addr])
        }
    }
    
    return result
}

func (ix *IHex) GetArea(addr uint) (area uint, ok bool) {
    for start, data := range ix.areas {
        end := start + uint(len(data))
        
        if addr >= start && addr <= end {
            return start, true
        }
    }
    
    return 0, false
}

func (ix *IHex) InsertData(istart uint, idata []byte) {
    iend := istart + uint(len(idata))
    
    area, ok := ix.GetArea(istart)
    
    if ok {
        data := ix.areas[area]
        newdata := make([]byte, len(data) + len(idata))
        
        copy(newdata, data[:istart-area])
        copy(newdata[istart-area:], idata)
        
        if iend-area < uint(len(data)) {
            copy(newdata[(istart-area) + uint(len(idata)):], data[iend-area:])
        }
        
        ix.areas[area] = newdata
    
    } else {
        ix.areas[istart] = idata
    }
}
