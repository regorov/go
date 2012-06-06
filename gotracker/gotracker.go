package gotracker

import (
	"github.com/kierdavis/go/musical"

//	"math"
)

type BeatType int

const (
	None	BeatType	= iota
	Hit
	Rest
)

type Song struct {
	SongData	*SongData
	Patterns	[]*Pattern
}

func (song *Song) Render(stream musical.Stream) {
	for _, pattern := range song.Patterns {
		pattern.Render(stream)
	}
}

type Pattern struct {
	SongData	*SongData
	Channels	[]*Channel
}

func (pattern *Pattern) Render(stream musical.Stream) {
	maxLength := 0
	chanStreams := make([]musical.Stream, len(pattern.Channels))

	for i, channel := range pattern.Channels {
		length := len(channel.Beats)
		if length > maxLength {
			maxLength = length
		}

		chanStreams[i] = make(musical.Stream)
		go channel.Render(chanStreams[i])
	}

	numSamples := maxLength * int(pattern.SongData.SamplesPerBeat)

	for i := 0; i < numSamples; i++ {
		var sample float64

		for _, chanStream := range chanStreams {
			// If the channel is closed, 0.0 is received
			sample += <-chanStream
		}

		stream <- sample
	}
}

type Channel struct {
	SongData	*SongData
	Instrument	*Instrument
	Beats		[]Beat
}

func (channel *Channel) Render(stream musical.Stream) {
	var notelength uint
	var note musical.Note

	for _, beat := range channel.Beats {
		if beat.Type == None {
			notelength++
		} else {
			if notelength > 0 {
				channel.Instrument.Render(note, channel.SongData.SampleRate, channel.SongData.SamplesPerBeat*notelength, stream)

				notelength = 0
			}

			if beat.Type == Hit {
				notelength = 1
				note = beat.Note

			} else {
				for i := 0; i < int(channel.SongData.SamplesPerBeat); i++ {
					stream <- 0.0
				}
			}
		}
	}

	if notelength > 0 {
		channel.Instrument.Render(note, channel.SongData.SampleRate, channel.SongData.SamplesPerBeat*notelength, stream)
	}

	close(stream)
}

type Beat struct {
	Type	BeatType
	Note	musical.Note
}

type Instrument struct {
	Samples []NoteSamplePair
}

func (inst *Instrument) Render(note musical.Note, sampleRate uint, numSamples uint, stream musical.Stream) {
	// Find the closest sample

	var closestNote musical.Note
	var closestDiff musical.Note = -1
	var closestSample Sample

	for _, pair := range inst.Samples {
		diff := pair.Note - note
		if diff < 0 {
			diff = -diff
		}

		if closestDiff < 0 || diff < closestDiff {
			closestDiff = diff
			closestNote = pair.Note
			closestSample = pair.Sample
		}
	}

	if closestNote == note {
		closestSample.Play(sampleRate, numSamples, stream)

	} else {
		raw := make(musical.Stream)
		closestSample.Play(sampleRate, numSamples, raw)

		ConvertPitch(raw, stream, closestNote, note, sampleRate)
	}
}

type NoteSamplePair struct {
	musical.Note
	Sample
}

type Sample struct {
	Data []float64
}

func (sample *Sample) Play(sampleRate uint, numSamples uint, stream musical.Stream) {
	data := sample.Data

	if numSamples < uint(len(data)) {
		data = data[:numSamples]
	}

	go func() {
		for _, v := range data {
			stream <- v
		}

		close(stream)
	}()
}

type SongData struct {
	SampleRate	uint
	Tempo		uint
	BeatLength	float32
	SamplesPerBeat	uint
}

func (sd *SongData) Calc() {
	sd.BeatLength = 1 / (float32(sd.Tempo) / 60)
	sd.SamplesPerBeat = uint(float32(sd.SampleRate) * sd.BeatLength)
}

func ConvertPitch(in, out musical.Stream, inNote, outNote musical.Note, sampleRate uint) {
	newRate := uint((outNote.Frequency() * float64(sampleRate)) / inNote.Frequency())

	// Playing the audio at sampleRate makes it sound like inNote
	// Playing the audio at newRate makes it sound like outNote
	// So we need to convert it back to sampleRate

	ConvertRate(in, out, newRate, sampleRate)
}

func ConvertRate(in, out musical.Stream, inRate, outRate uint) {
	gcd := GCD(inRate, outRate)
	inRate /= gcd
	outRate /= gcd

	partial := make(musical.Stream)

	go func() {
		a := <-in

		for {
			b := <-in
			step := (b - a) / float64(outRate)

			partial <- a

			for i := uint(0); i < outRate-1; i++ {
				a += step
				partial <- a
			}

			a = b
		}
	}()

	go func() {
		for {
			out <- <-partial

			for i := uint(0); i < inRate-1; i++ {
				<-partial
			}
		}
	}()
}

func GCD(a, b uint) (gcd uint) {
	for b > 0 {
		a, b = b, a%b
	}

	return a
}

func LCM(a, b uint) (lcm uint) {
	return (a / GCD(a, b)) * b
}
