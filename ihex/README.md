Package: github.com/kierdavis/go/ihex
=====================================

[doc](http://gopkgdoc.appspot.com/pkg/github.com/kierdavis/go/ihex)

Package ihex reads and writes Intel Hex files. The Intel Hex format is used to store binary data
in an ASCII (hexadecimal) format, and supports storing data at arbitrary addresses rather than
in a linear form.

Example: reading a file

f, err := os.Open("input.hex")
if err != nil {
panic(err)
}
defer f.Close()

reader := bufio.NewReader(f)
ix, err := ihex.ReadIHex(reader)
if err != nil {
panic(err)
}

f.Close()
data := ix.ExtractDataToEnd(0) // data from 0 -> end

Example: writing a file

ix := ihex.NewIHex()
ix.InsertData(0, data)

f, err := os.Create("output.hex")
if err != nil {
panic(err)
}

writer := bufio.NewWriter(f)
err = ix.Write(writer)
if err != nil {
panic(err)
}

writer.Flush()
f.Close()

Package Dependencies
--------------------

None

