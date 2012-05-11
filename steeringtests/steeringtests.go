package main

import (
    "github.com/banthar/Go-SDL/sdl"
    "github.com/kierdavis/go/resourcemanager"
    "image"
)

var rm = resourcemanager.NewResourceManager("github.com/kierdavis/go/steeringtests")

const WIDTH = 400
const HEIGHT = 300

const (
    X = iota
    Y
    Z
)

func truncate(x, max float32) float32 {
    if x > max {
        return max
    }
    return x
}

type Wolf struct {
    Mass     float32
    X        float32
    Y        float32
    VX       float32
    VY       float32
    MaxForce float32
    MaxSpeed float32
    Heading  float32
}

func (wolf *Wolf) Move(fx, fy float32) {
    wolf.VX += fx / wolf.Mass
    wolf.VY += fy / wolf.Mass
    wolf.X += wolf.VX
    wolf.Y += wolf.VY
}

func main() {
    if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
        panic(sdl.GetError())
    }
    defer sdl.Quit()

    screen := sdl.SetVideoMode(WIDTH, HEIGHT, 32, sdl.RESIZABLE)
    if screen == nil {
        panic(sdl.GetError())
    }

    sdl.WM_SetCaption("Steering tests", "")

    f, err := rm.GetFile("wolf.png")
    if err != nil {
        panic(err)
    }

    wolf_image, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    wolf := sdl.CreateSurfaceFromImage(wolf_image)

}
