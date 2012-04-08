// Command k270em is a frontend to the k270emlib K270 processor emulator, including a SDL-based
// character display. It takes one command-ling argument, the program to load (which should be in
// Intel Hex format).
package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
    "github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
    "github.com/kierdavis/go/k270emlib"
    "github.com/kierdavis/go/ihex"
    "github.com/kierdavis/go/resourcemanager"
    "os"
    "time"
)

// Function die panics with `err` if `err` is not nil.
func die(err error) {
    if err != nil {
        panic(err)
    }
}

// Function main is the main entry point in the program.
func main() {
    flag.Parse()
    
    if flag.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "Not enough arguments\nusage: %s file.hex\n", os.Args[0])
        os.Exit(2)
    }
    
    if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
        panic(sdl.GetError())
    }
    
    defer sdl.Quit()
    
    if ttf.Init() != 0 {
        panic(sdl.GetError())
    }
    
    defer ttf.Quit()
    
    rm := resourcemanager.NewResourceManager("github.com/kierdavis/go/k270em")
    
    font := ttf.OpenFont(rm.GetFilename("FreeMono.ttf"), 12)
    if font == nil {panic(sdl.GetError())}
    
    defer font.Close()
    
    _, _, _, _, fontWidth, errID := font.GlyphMetrics('m')
    if errID != 0 {panic(sdl.GetError())}
    
    fontHeight := font.Height()
    
    fmt.Printf("h: %d, w: %d\n", fontHeight, fontWidth)
    
    screenHeight := fontHeight * 48
    screenWidth := fontWidth * 128
    
    screen := sdl.SetVideoMode(screenWidth, screenHeight, 32, sdl.RESIZABLE)
    if screen == nil {panic(sdl.GetError())}
    
    sdl.WM_SetCaption("K270 Emulator Display", "")
    
    f, err := os.Open(flag.Arg(0)); die(err)
    defer f.Close()
    
    reader := bufio.NewReader(f)
    ix, err := ihex.ReadIHex(reader); die(err)
    program := ix.ExtractDataToEnd(0)
    
    em := k270emlib.NewEmulator()
    em.SetTraceFile(os.Stdout)
    em.SetMemory(program)
    
    running := true
    stopRequest := make(chan bool)
    vmem := em.GetVideoMemory()
    
    go em.Run()
    go func() {
        scanTicker := time.NewTicker(time.Second / 24.0) // 24 hz
        color := sdl.Color{255, 255, 255, 0}
        
        for {
            if !running {break}
            
            select {
            case <-scanTicker.C:
                address := uint16(0)
                for y := 0; y < 48; y++ {
                    chars := make([]byte, 128)
                    
                    for x := 0; x < 128; x++ {
                        c := vmem[address].Char
                        if c == 0 {c = ' '}
                        chars[x] = c
                        address++
                    }
                    
                    surf := ttf.RenderText_Solid(font, string(chars), color)
                    screen.Blit(&sdl.Rect{0, int16(y * fontHeight), 0, 0}, surf, nil)
                }
            
            case <-stopRequest:
                return
            }
        }
    }()
    
    ticker := time.NewTicker(time.Second / 24.0) // 24 fps
    
    for running {
        select {
        case <-ticker.C:
            screen.Flip()
        
        case ev := <-sdl.Events:
            switch e := ev.(type) {
            case sdl.QuitEvent:
                running = false
            
            case sdl.KeyboardEvent:
                if e.Keysym.Sym == sdl.K_ESCAPE {
                    running = false
                }
            
            case sdl.ResizeEvent:
                screen = sdl.SetVideoMode(int(e.W), int(e.H), 32, sdl.RESIZABLE)
                if screen == nil {panic(sdl.GetError())}
            }
        }
    }
    
    stopRequest <- true
    
    //fmt.Println("Locking...")
    em.Mutex.Lock()
    //fmt.Println("Locked!")
    em.SetRunning(false)
    //fmt.Println("Unlocking...")
    em.Mutex.Unlock()
    //fmt.Println("Unlocked!")
}
