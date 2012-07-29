package sound

func Clip(threshold float64, in chan float64, out chan float64) {
	defer close(out)

	for sample := range in {
		if sample >= threshold {
			sample = threshold
		} else if sample <= -threshold {
			sample = -threshold
		}

		// Sample will have a max of threshold (and a min of -threshold)
		// We want the max to be 1; therefore we divide by threshold

		out <- sample / threshold
	}
}

func GoClip(threshold float64, in chan float64) (out chan float64) {
	out = make(chan float64)
	go Clip(threshold, in, out)
	return out
}
