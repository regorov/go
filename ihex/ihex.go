// Package ihex reads and writes Intel Hex files. The Intel Hex format is used to store binary data
// in an ASCII (hexadecimal) format, and supports storing data at arbitrary addresses rather than
// in a linear form.
// 
// Example: reading a file
//
//     f, err := os.Open("input.hex")
//     if err != nil {
//         panic(err)
//     }
//     defer f.Close()
//     
//     reader := bufio.NewReader(f)
//     ix, err := ihex.ReadIHex(reader)
//     if err != nil {
//         panic(err)
//     }
//     
//     f.Close()
//     data := ix.ExtractDataToEnd(0) // data from 0 -> end
// 
// Example: writing a file
// 
//     ix := ihex.NewIHex()
//     ix.InsertData(0, data)
//     
//     f, err := os.Create("output.hex")
//     if err != nil {
//         panic(err)
//     }
//     
//     writer := bufio.NewWriter(f)
//     err = ix.Write(writer)
//     if err != nil {
//         panic(err)
//     }
//     
//     writer.Flush()
//     f.Close()
package ihex

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "strings"
)

// Type IHex represents an Intel Hex file.
type IHex struct {
    areas map[uint][]byte   // A map of area addresses to byte slices.
}

// Function makeLine constructs an Intel Hex record, including the final newline, from the type `t`,
// address `addr` and content `data` of the record.
func makeLine(t uint, addr uint, data []byte) (line string) {
    lineBytes := make([]byte, 5 + len(data))
    lineBytes[0] = uint8(len(data))
    lineBytes[1] = uint8(addr >> 8)
    lineBytes[2] = uint8(addr)
    lineBytes[3] = uint8(t)
    
    dataend := 4 + len(data)
    copy(lineBytes[4:dataend], data)
    
    lineBytes[dataend] = CalcChecksum(lineBytes[:dataend])
    return fmt.Sprintf(":%x\n", lineBytes)
}

// Function parseLine parses an Intel Hex record into its type, address and content.
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
    cs1 := line[4 + length]
    cs2 := CalcChecksum(line[:4 + length])
    
    if cs1 != cs2 {
        return 0, 0, nil, errors.New("Checksums do not match")
    }
    
    return t, addr, data, nil
}

// Function CalcChecksum calculates an Intel Hex checksum from the data `data`.
func CalcChecksum(data []byte) (checksum uint8) {
    total := uint(0)
    
    for _, x := range data {
        total += uint(x)
    }
    
    return uint8(-total)
}

// Function NewIHex creates and returns a new IHex object.
func NewIHex() (ix *IHex) {
    ix = new(IHex)
    ix.areas = make(map[uint][]byte)
    return ix
}

// Function ReadIHex reads an Intel Hex file from `reader` and returns it.
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

// Function IHex.getArea finds the address of the area that would contain the address `addr`. If a
// suitable area already exists, it returns its address and true. If not, it returns 0 and false.
func (ix *IHex) getArea(addr uint) (area uint, ok bool) {
    for start, data := range ix.areas {
        end := start + uint(len(data))
        
        if addr >= start && addr <= end {
            return start, true
        }
    }
    
    return 0, false
}

// Function IHex.GetSize returns the maximum address of the data, plus one.
func (ix *IHex) GetSize() (size uint) {
    size = 0
    
    for addr, data := range ix.areas {
        if addr + uint(len(data)) > size {
            size = addr + uint(len(data))
        }
    }
    
    return size
}

// Function IHex.ExtractData copies data out of the IHex file, starting at the address `start` and
// ending at the address before `end`. It returns the copied data.
func (ix *IHex) ExtractData(start uint, end uint) (result []byte) {
    result = make([]byte, 0, end - start)
    
    for addr, data := range ix.areas {
        if addr >= start && addr < end {
            copy(result[addr:], data)
        }
    }
    
    return result
}

// Function IHex.ExtractDataToEnd copies data out of the IHex file, starting at the address `start`
// and ending at the last address in the data.
func (ix *IHex) ExtractDataToEnd(start uint) (result []byte) {
    end := uint(0)
    result = make([]byte, ix.GetSize() - start)
    
    for addr, data := range ix.areas {
        if addr >= start {
            if addr + uint(len(data)) > end {
                end = addr + uint(len(data))
            }
            
            copy(result[addr:], data)
        }
    }
    
    return result
}

// Function IHex.InsertData inserts the data contained in `idata` into the IHex object, starting at
// address `istart`.
func (ix *IHex) InsertData(istart uint, idata []byte) {
    iend := istart + uint(len(idata))
    
    area, ok := ix.getArea(istart)
    
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

// Function IHex.Write writes the representation of the file to the writer `wr`.
func (ix *IHex) Write(wr *bufio.Writer) (err error) {
    for start, data := range ix.areas {
        i := 0
        
        for i < len(data) {
            var chunk []byte
            
            if i + 16 <= len(data) {
                chunk = data[i:i+16]
            } else {
                chunk = data[i:]
            }
            
            _, err := wr.Write([]byte(makeLine(0x00, start, chunk)))
            if err != nil {return err}
            
            start += 16
        }
    }
    
    _, err = wr.Write([]byte(makeLine(0x01, 0, []byte{})))
    if err != nil {return err}
    
    return nil
}
