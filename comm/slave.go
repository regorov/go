package comm

import (
	"net/rpc"
	"os"
)

type Comm struct{}

func Serve() {
	rpc.Register(&Comm{})
	rpc.ServeConn(ServerConn{os.Stdin, os.Stdout})
}
