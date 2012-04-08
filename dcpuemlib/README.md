Package: github.com/kierdavis/go/dcpuemlib
==========================================

[doc](http://gopkgdoc.appspot.com/pkg/github.com/kierdavis/go/dcpuemlib)

Package dcpuemlib emulates a DCPU-16 processor. Typical usage is as follows:

em := dcpuemlib.NewEmulator()      // Create an emulator
em.LoadProgramBytesLE(program)     // Load a program (in little-endian format)
em.TraceFile = os.Stdout           // Set the debug trace file to standard output
em.Run()                           // Run until a halt-like instruction (or an error) is
encountered.
em.DumpState()                     // Dump the registers to standard output

A word about RAM:

Because the architecture requires a relatively large amount of RAM (128Kb to be exact),
allocating this amount whenever the emulator is reset will result in a large memory footprint.
To counter this, a dynamically expanding array is used, that starts at 1 Kb and expands if
more is needed.

Package Dependencies
--------------------

None

