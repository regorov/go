package comm

import (
	"io"
)

type ClientConn struct {
	R io.Reader
	W io.WriteCloser
}

func (c *ClientConn) Read(buffer []byte) (n int, err error) {
	return c.R.Read(buffer)
}

func (c *ClientConn) Write(buffer []byte) (n int, err error) {
	return c.W.Write(buffer)
}

func (c *ClientConn) Close() (err error) {
	return c.W.Close()
}

type ServerConn struct {
	R io.Reader
	W io.Writer
}

func (c *ServerConn) Read(buffer []byte) (n int, err error) {
	return c.R.Read(buffer)
}

func (c *ServerConn) Write(buffer []byte) (n int, err error) {
	return c.W.Write(buffer)
}

func (c *ServerConn) Close() (err error) {
	os.Exit(0)
}
