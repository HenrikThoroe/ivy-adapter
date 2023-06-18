package mgmt

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

// Launch launches the engine instance and returns a connection to it.
// The connection is used to communicate with the engine using channels.
// If the engine instance could not be launched, an error is returned and the connection is nil.
func (e *EngineInstance) Launch() (*Connection, error) {
	path := e.Path()
	proc := exec.Command(path)

	inPipe, _ := proc.StdinPipe()
	outPipe, _ := proc.StdoutPipe()

	in := make(chan string)
	out := make(chan string)

	os.Chmod(path, 0700)

	if e := proc.Start(); e != nil {
		return nil, e
	}

	return bind(proc.Process.Pid, in, out, inPipe, outPipe), nil
}

func bind(pid int, in chan string, out chan string, wr io.Writer, rd io.Reader) *Connection {
	go distribute(in, wr)
	go listen(rd, out)

	return &Connection{in, out, pid}
}

func distribute(in chan string, wr io.Writer) {
	for cmd := range in {
		wr.Write([]byte(cmd + "\n"))
	}
}

func listen(rd io.Reader, out chan string) {
	scanner := bufio.NewScanner(rd)

	for scanner.Scan() {
		out <- scanner.Text()
	}
}
