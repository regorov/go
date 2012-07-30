Package: github.com/kierdavis/go/dcpuem
=======================================

[doc](http://gopkgdoc.appspot.com/pkg/github.com/kierdavis/go/dcpuem)

Package dcpuem is a DCPU-16 emulator.


Package dcpuem implements a DCPU-16 emulator. It currently conforms to revision 1.7 of the
specification.

Example usage:

    // Create the emulator
    em := dcpuem.NewEmulator()

    // Set up logging
    em.Logger = log.NewLogger(os.Stdout, "", log.LstdFlags | log.Lshortfile)

    // Load the program (for this example, using binaryloader)
    program, err := binaryloader.Load(os.Args[1])
    if err != nil {panic(err)}
    em.LoadProgramBytesBE(program)

    // Set the clock frequency to 120 Hz (optional)
    em.ClockTicker = time.NewTicker(time.Second / 120.0)

    // Create some devices & attach them to the emulator (optional)
    clock_device := clock.New()
    disk_drive := hmd2043.New()
    em.AttachDevice(clock_device)
    em.AttachDevice(disk_drive)

    // Start the devices
    em.StartDevices()

    // Run the emulator
    err = em.Run()
    if err != nil {panic(err)}

    // Stop the devices
    em.StopDevices()




Install
-------

    $ go get github.com/kierdavis/dcpuem

Package Dependencies
--------------------

None

