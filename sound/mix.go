package sound

func Mix(out chan float64, streams ...chan float64) {
	defer close(out)

	for {
		sample := 0.0
		ok := true

		for _, stream := range streams {
			x, y := <-stream
			sample += x
			ok = ok || y
		}

		if !ok {
			break
		}

		if sample < -1.0 {
			sample = -1.0
		} else if sample > 1.0 {
			sample = 1.0
		}

		out <- sample
	}
}

func GoMix(streams ...chan float64) (out chan float64) {
	out = make(chan float64)
	go Mix(out, streams...)
	return out
}

func Concatenate(out chan float64, streams ...chan float64) {
	defer close(out)

	for _, stream := range streams {
		for sample := range stream {
			out <- sample
		}
	}
}

func GoConcatenate(streams ...chan float64) (out chan float64) {
	out = make(chan float64)
	go Concatenate(out, streams...)
	return out
}
