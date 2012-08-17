package sound

import (
	"time"
)

func Clip(threshold float64, in chan Sample, out chan Sample) {
	defer close(out)

	for sample := range in {
		if sample.Left >= threshold {
			sample.Left = threshold
		} else if sample.Left <= -threshold {
			sample.Left = -threshold
		}

		if sample.Right >= threshold {
			sample.Right = threshold
		} else if sample.Right <= -threshold {
			sample.Right = -threshold
		}

		// Sample will have a max of threshold (and a min of -threshold)
		// We want the max to be 1; therefore we divide by threshold

		out <- sample.Div(threshold)
	}
}

func GoClip(threshold float64, in chan Sample) (out chan Sample) {
	out = make(chan Sample, ChannelBuffer)
	go Clip(threshold, in, out)
	return out
}

func Flange(freq float64, minPeriod, maxPeriod, windowSize time.Duration, in chan Sample, out chan Sample) {
	defer close(out)

	windowSizeSamples := int(windowSize / SampleTime)
	ctl := sineInput(freq, SampleRate/float64(windowSizeSamples))
	periodRange := maxPeriod - minPeriod

	var lastStream2 chan Sample

	for {
		// Read a window
		window := make([]Sample, 0, windowSizeSamples)
		for i := 0; i < windowSizeSamples; i++ {
			window = append(window, <-in)
		}

		// Create two streams that play the window
		stream1 := make(chan Sample, ChannelBuffer)
		stream2 := make(chan Sample, ChannelBuffer)

		go func() {
			defer close(stream1)
			defer close(stream2)

			for _, sample := range window {
				stream1 <- sample
				stream2 <- sample
			}
		}()

		// Delay the second stream by a set period
		period := (time.Duration(<-ctl) * periodRange) + minPeriod
		stream2 = GoDelay(period, stream2)

		// Mix the last stream2 with the current stream1 and stream2
		// Make sure we cut off the current stream2, so there's some left to mix into the next window.
		var windowOutput chan Sample

		if lastStream2 == nil {
			windowOutput = GoMix(stream1, GoCopyFor(windowSize, stream2))
		} else {
			windowOutput = GoMix(stream1, GoCopyFor(windowSize, stream2), lastStream2)
		}

		go Append(out, windowOutput)
		lastStream2 = stream2
	}
}
