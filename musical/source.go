package musical

import (
    "math"
)

type Chunk []float64

func Silence(length float64, rate float64) (chunk Chunk) {
    return make(Chunk, int(length*rate))
}

func generateWaveInput(freq float64, length float64, rate float64, phase float64) (chunk Chunk) {
    ilength := int(length * rate)
    factor := (freq * math.Pi * 2) / rate
    phase *= rate / 2.0

    chunk = make(Chunk, ilength)
    for i := 0; i < ilength; i++ {
        chunk[i] = (float64(i) + phase) * factor
    }

    return chunk
}

func Sine(freq float64, length float64, rate float64, phase float64) (chunk Chunk) {
    chunk = generateWaveInput(freq, length, rate, phase)

    for i := 0; i < len(chunk); i++ {
        chunk[i] = math.Sin(chunk[i])
    }

    return chunk
}
