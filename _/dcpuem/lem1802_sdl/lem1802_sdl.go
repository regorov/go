package lem1802_sdl

import (
    "errors"
    "github.com/banthar/Go-SDL/sdl"
    "github.com/banthar/Go-SDL/ttf"
    "github.com/kierdavis/go/dcpuem"
    "github.com/kierdavis/go/resourcemanager"
)

var (
    ErrSurfaceTooSmall  = errors.New("The specified surface did not meet the size requirements (128 x 96)")
    ErrCouldNotLoadFont = errors.New("Could not load TTF font")
)

type Font [128]*sdl.Surface

type LEM1802 struct {
    // The associated emulator.
    Em *dcpuem.Emulator

    // The surface to draw onto.
    Surf *sdl.Surface

    // The memory address of the start of video RAM, or 0 to disable the display.
    VideoMap uint16

    FontMap uint16

    PaletteMap uint16

    // The precompiled default font (loaded by LoadTTF).
    DefaultFont *Font

    // The precompiled user-defined font.
    UserFont *Font
}

func New(surf *sdl.Surface) (d *LEM1802, err error) {
    if surf.W < 128 || surf.H < 96 {
        return d, ErrSurfaceTooSmall
    }

    d = new(LEM1802)
    d.Em = nil
    d.Surf = surf
    d.VideoMap = 0
    return d, nil
}

func (d *LEM1802) LoadFontFromTTF(filename string) (font *Font, err error) {
    if filename == "" {
        filename = resourcemanager.NewResourceManager("github.com/kierdavis/go/dcpuem/lem1802").GetFilename("FreeMono.ttf")
    }

    f := ttf.OpenFont(filename, 8)
    if f == nil {
        return nil, ErrCouldNotLoadFont
    }

    defer f.Close()

    font = new(Font)

    for ch := uint16(0); ch < 128; ch++ {
        _, _, _, _, advance, e := f.GlyphMetrics(ch)
        if e != 0 {
            return nil, ErrCouldNotLoadFont
        }

        surf := sdl.CreateRGBSurface(sdl.SWSURFACE|sdl.SRCALPHA, advance, 8, 32, 0x000000ff, 0x0000ff00, 0x00ff0000, 0xff000000)

        fontSurf := ttf.RenderText_Blended(f, string(ch), sdl.Color{255, 255, 255, 0})
        fontSurf.Blit(nil, surf, nil)

        font[ch] = surf
    }

    return font, nil
}

func (d *LEM1802) LoadFromFromRAM(start uint16) (font *Font, err error) {
    return nil, nil
}

func (d *LEM1802) AssociateEmulator(em *dcpuem.Emulator) {
    d.Em = em
}

func (d *LEM1802) ID() (id uint32) {
    return 0x7349f615
}

func (d *LEM1802) Version() (ver uint16) {
    return 0x1802
}

func (d *LEM1802) Manufacturer() (manu uint32) {
    return 0x1c6c8b36
}

func (d *LEM1802) Render() (err error) {
    var font *Font

    if d.FontMap == 0 {
        if d.DefaultFont == nil {
            font, err = d.LoadFontFromTTF("")
            if err != nil {
                return err
            }

            d.DefaultFont = font

        } else {
            font = d.DefaultFont
        }

    } else {
        if d.UserFont == nil {
            font, err = d.LoadFromFromRAM(d.FontMap)
            if err != nil {
                return err
            }

            d.UserFont = font

        } else {
            font = d.UserFont
        }
    }

    return nil
}
