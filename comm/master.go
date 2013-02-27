package comm

import (
	"io"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
)

type Slave struct {
	proc   *os.Process
	stdin  io.WriteCloser
	stdout io.ReadCloser
	client *rpc.Client
}

func SpawnSlave(filename string) (slave *Slave, err error) {
	if strings.HasSuffix(filename, ".go") {
		binFile, err := Compile(filename)
		if err != nil {
			return nil, err
		}

		filename = binFile
	}

	cmd := exec.Command(filename)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	c := ClientConn{stdout, stdin}
	client := rpc.NewClient(c)

	slave = &Slave{
		proc:   cmd.Process,
		stdin:  stdin,
		stdout: stdout,
		client: client,
	}

	return slave, nil
}

func (slave *Slave) Stop() (err error) {
	err = slave.stdin.Close()
	if err != nil {
		return err
	}

	_, err = slave.proc.Wait()
	return err
}

func (slave *Slave) Kill() (err error) {
	return slave.proc.Kill()
}
