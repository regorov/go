package main

import (
    "math/rand"
)

const MAX_GENOME_LENGTH = 128

type Genome []float64

func RandomGenome(seed int64) (genome Genome) {
    r := rand.New(rand.NewSource(seed))
    size := r.Int31n(MAX_GENOME_LENGTH)

    genome = make(Genome, size)

    for i := 0; i < size; i++ {
        genome[i] = r.Float64()
    }

    return genome
}

func (genome Genome) Pump(c chan float64, stop chan bool) {
    i := 0

    for {
        select {
        case <-stop:
            return

        case c <- genome[i]:
            i = (i + 1) % len(genome)
        }
    }
}
