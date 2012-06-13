package argo

import (
	"io"
)

type Parser interface {
	Parse(io.Reader) (chan *Triple, chan error)
}

type Serializer interface {
	Serialize(io.Writer, chan *Triple) error
}

type ParserSerializer interface {
	Parser
	Serializer
}
