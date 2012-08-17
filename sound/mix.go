package sound

import (
	"time"
)

func Mix(out chan Sample, streams ...chan Sample) {
	defer close(out)

	intermediate := make(chan Sample, ChannelBuffer)
	go Clip(1.0, intermediate, out)

	for {
		var sample Sample

		numOpenStreams := 0

		for _, stream := range streams {
			x, streamOpen := <-stream
			sample = sample.Add(x)

			if streamOpen {
				numOpenStreams++
			}
		}

		if numOpenStreams == 0 {
			break
		}

		intermediate <- sample.Div(float64(numOpenStreams))
	}
}

func GoMix(streams ...chan Sample) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Mix(out, streams...)
	return out
}

func Concatenate(out chan Sample, streams ...chan Sample) {
	defer close(out)

	for _, stream := range streams {
		for sample := range stream {
			out <- sample
		}
	}
}

func GoConcatenate(streams ...chan Sample) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Concatenate(out, streams...)
	return out
}

func Append(out chan Sample, in chan Sample) {
	for sample := range in {
		out <- sample
	}
}

func Delay(d time.Duration, in chan Sample, out chan Sample) {
	Concatenate(out, GoSilence(d), in)
}

func GoDelay(d time.Duration, in chan Sample) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Delay(d, in, out)
	return out
}

func CopyFor(d time.Duration, in chan Sample, out chan Sample) {
	defer close(out)

	var t time.Duration

	for t = 0; t < d; t++ {
		out <- <-in
	}

	/*
		for _ = range in {
			// do nothing
		}
	*/
}

func GoCopyFor(d time.Duration, in chan Sample) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go CopyFor(d, in, out)
	return out
}
