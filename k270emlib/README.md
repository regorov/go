Package: github.com/kierdavis/go/k270emlib
==========================================

[doc](http://gopkgdoc.appspot.com/pkg/github.com/kierdavis/go/k270emlib)

Package k270emlib emulates a K270 processor. Typical usage:

em := k270emlib.NewEmulator()    // Create an emulator
em.LoadProgram(myprogram)        // Load a program
em.SetTraceFile(os.Stdout)       // (optional) log instructions to stdout
em.Run()                         // Go!


Package Dependencies
--------------------

None

