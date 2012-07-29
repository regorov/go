package sound

import (
	"math"
)

const TimeStep = 1.0 / 44100.0

func Sine(freq float64, length float64, out chan float64) {
	defer close(out)

	step := (math.Pi * freq) / 44100.0
	x := 0.0

	if length > 0.0 {
		for t := 0.0; t < length; t += TimeStep {
			out <- x

			x += step
			if x >= math.Pi {
				x -= math.Pi
			}
		}

	} else {
		for {
			out <- x

			x += step
			if x >= math.Pi {
				x -= math.Pi
			}
		}
	}
}

func GoSine(freq float64, length float64) (out chan float64) {
	out = make(chan float64)
	go Sine(freq, length, out)
	return out
}

func Saw(freq float64, length float64, out chan float64) {
	defer close(out)

	step := (2.0 * freq) / 44100.0
	x := -1.0

	if length > 0.0 {
		for t := 0.0; t < length; t += TimeStep {
			out <- x

			x += step
			if x >= 1.0 {
				x -= 2.0
			}
		}

	} else {
		for {
			out <- x

			x += step
			if x >= 1.0 {
				x -= 2.0
			}
		}
	}
}

func GoSaw(freq float64, length float64) (out chan float64) {
	out = make(chan float64)
	go Saw(freq, length, out)
	return out
}

func Triangle(freq float64, length float64, out chan float64) {
	defer close(out)

	step := (4.0 * freq) / 44100.0
	x := -1.0

	if length > 0.0 {
		for t := 0.0; t < length; t += TimeStep {
			out <- x

			x += step
			if x >= 1.0 || x <= -1.0 {
				step = -step
			}
		}

	} else {
		for {
			out <- x

			x += step
			if x >= 1.0 || x <= -1.0 {
				step = -step
			}
		}
	}
}

func GoTriangle(freq float64, length float64) (out chan float64) {
	out = make(chan float64)
	go Triangle(freq, length, out)
	return out
}

func Square(freq float64, length float64, out chan float64) {
	defer close(out)

	in := GoTriangle(freq, length)

	for sample := range in {
		if sample >= 0.0 {
			out <- 1.0
		} else {
			out <- -1.0
		}
	}
}

func GoSquare(freq float64, length float64) (out chan float64) {
	out = make(chan float64)
	go Square(freq, length, out)
	return out
}
