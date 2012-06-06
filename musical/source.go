package musical

import (
	"math"
)

type Stream chan float64

func Silence(length float64, rate float64, output Stream) {
	go func() {
		output <- 0.0
	}()
}

func generateWaveInput(freq float64, rate uint, phase float64, output Stream) {
	factor := (freq * math.Pi * 2) / float64(rate)
	phase *= float64(rate) / 2.0

	go func() {
		i := 0

		for {
			output <- (float64(i) + phase) * factor
			i++
		}
	}()
}

func Sine(freq float64, rate uint, phase float64, output Stream) {
	input := make(Stream)
	generateWaveInput(freq, rate, phase, input)

	go func() {
		for {
			output <- math.Sin(<-input)
		}
	}()
}

func Sawtooth(freq float64, rate uint, phase float64, output Stream) {
	input := make(Stream)
	generateWaveInput(freq, rate, phase, input)

	go func() {
		for {
			tmod := math.Mod(<-input, 2*math.Pi)
			output <- (tmod / math.Pi) - 1
		}
	}()
}

func Square(freq float64, rate uint, phase float64, output Stream) {
	input := make(Stream)
	generateWaveInput(freq, rate, phase, input)

	go func() {
		for {
			tmod := math.Mod(<-input, 2*math.Pi)

			if tmod < math.Pi {
				output <- 1.0
			} else {
				output <- -1.0
			}
		}
	}()
}
