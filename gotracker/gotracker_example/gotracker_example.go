package main

import (
	"github.com/kierdavis/go/gotracker"
	"github.com/kierdavis/go/musical"
)

func main() {
	inNote := musical.ParseNote("c4")
	outNote := musical.ParseNote("e4")

	input := make(musical.Stream)
	musical.Sine(inNote.Frequency(), 44100, 0.0, input)

	output := make(musical.Stream)
	gotracker.ConvertPitch(input, output, inNote, outNote, 44100)
}
