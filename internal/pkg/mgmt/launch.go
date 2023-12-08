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
func (e *EngineInstance) Launch(scb func(string), rcb func(string)) (*Connection, error) {
	path := e.Path()
	return LaunchEngine(path, scb, rcb)
}

// LaunchEngine launches an engine executable at the given path and returns a connection to it.
// The connection is used to communicate with the engine using channels.
// If the engine instance could not be launched, an error is returned and the connection is nil.
func LaunchEngine(path string, scb func(string), rcb func(string)) (*Connection, error) {
	proc := exec.Command(path)

	inPipe, _ := proc.StdinPipe()
	outPipe, _ := proc.StdoutPipe()

	in := make(chan string)
	out := make(chan string)

	os.Chmod(path, 0700)

	if e := proc.Start(); e != nil {
		return nil, e
	}

	return bind(proc.Process.Pid, in, out, inPipe, outPipe, scb, rcb), nil
}

func bind(pid int, in chan string, out chan string, wr io.Writer, rd io.Reader, scb func(string), rcb func(string)) *Connection {
	go distribute(in, wr, scb)
	go listen(rd, out, rcb)

	return &Connection{in, out, pid}
}

func distribute(in chan string, wr io.Writer, cb func(string)) {
	for cmd := range in {
		wr.Write([]byte(cmd + "\n"))

		if cb != nil {
			cb(cmd)
		}
	}
}

func listen(rd io.Reader, out chan string, cb func(string)) {
	scanner := bufio.NewScanner(rd)

	for scanner.Scan() {
		text := scanner.Text()
		out <- text

		if cb != nil {
			cb(text)
		}
	}
}
