package sound

import (
	"math"
	"time"
)

func sineInput(freq float64, outFreq float64) (out chan float64) {
	out = make(chan float64)

	go func() {
		step := (math.Pi * freq * 2.0) / outFreq
		x := 0.0

		for {
			out <- math.Sin(x)

			x += step
			if x >= math.Pi {
				x -= math.Pi * 2
			}
		}
	}()

	return out
}

func Sine(freq float64, length time.Duration, out chan Sample) {
	defer close(out)

	ctl := sineInput(freq, SampleRate)

	var t time.Duration

	for t = 0; t < length; t += SampleTime {
		x := <-ctl
		out <- Sample{x, x}
	}
}

func GoSine(freq float64, length time.Duration) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Sine(freq, length, out)
	return out
}

func Saw(freq float64, length time.Duration, out chan Sample) {
	defer close(out)

	step := (2.0 * freq) / 44100.0
	x := 0.0

	var t time.Duration

	for t = 0; t < length; t += SampleTime {
		out <- Sample{x, x}

		x += step
		if x >= 1.0 {
			x -= 2.0
		}
	}
}

func GoSaw(freq float64, length time.Duration) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Saw(freq, length, out)
	return out
}

func Triangle(freq float64, length time.Duration, out chan Sample) {
	defer close(out)

	step := (4.0 * freq) / 44100.0
	x := 0.0

	var t time.Duration

	for t = 0; t < length; t += SampleTime {
		out <- Sample{x, x}

		x += step
		if x >= 1.0 || x <= -1.0 {
			step = -step

			// Reverse the change we made before the if statement, and then go another step
			// in the new direction.
			x += step * 2
		}
	}
}

func GoTriangle(freq float64, length time.Duration) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Triangle(freq, length, out)
	return out
}

func Square(freq float64, length time.Duration, out chan Sample) {
	Clip(1e-6, GoTriangle(freq, length), out)
}

func GoSquare(freq float64, length time.Duration) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Square(freq, length, out)
	return out
}

func Silence(length time.Duration, out chan Sample) {
	defer close(out)

	var t time.Duration

	for t = 0; t < length; t += SampleTime {
		out <- Sample{0.0, 0.0}
	}
}

func GoSilence(length time.Duration) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Silence(length, out)
	return out
}
