Command: dcpuem
===============

Command dcpuem is a frontend to the dcpuemlib DCPU-16 emulator. It takes one required argument,
the filename of the program to load. The format of the input file can be either Intel Hex or
raw binary and is determined by the file's extension (defaulting to raw binary if it is
unrecognised). There is also one command-line option: -b. If this is present, the input file
is parsed into words in a big-endian fashion rather than the default (little-endian).

Package Dependencies
--------------------

* [github.com/kierdavis/go/dcpuemlib](https://github.com/kierdavis/go/tree/master/dcpuemlib) ([doc](http://gopkgdoc.appspot.com/pkg/github.com/kierdavis/go/dcpuemlib))
* [github.com/kierdavis/go/ihex](https://github.com/kierdavis/go/tree/master/ihex) ([doc](http://gopkgdoc.appspot.com/pkg/github.com/kierdavis/go/ihex))

(documentation provided by [GoPkgDoc](http://gopkgdoc.appspot.com/index))

