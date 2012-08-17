package sound

import (
	"io"
	"time"
)

const SampleRate = 44100.0
const SampleTime = time.Second / SampleRate

var ChannelBuffer = 10

type Sample struct {
	Left  float64
	Right float64
}

func (s Sample) Mono() (x float64) {
	return (s.Left + s.Right) / 2.0
}

func (s Sample) Add(other Sample) (result Sample) {
	return Sample{s.Left + other.Left, s.Right + other.Right}
}

func (s Sample) Sub(other Sample) (result Sample) {
	return Sample{s.Left - other.Left, s.Right - other.Right}
}

func (s Sample) Mul(x float64) (result Sample) {
	return Sample{s.Left * x, s.Right * x}
}

func (s Sample) Div(x float64) (result Sample) {
	return Sample{s.Left / x, s.Right / x}
}

type Encoder func(io.Writer, chan Sample) error
type Decoder func(io.Reader, chan Sample) error
