package mgmt

// Connection is a wrapper for the communication between the adapter and the
// engine. It provides a simple interface to send commands to the engine and
// receive its responses.
type Connection struct {
	in  chan string
	out chan string
	Pid int
}

// NewConnection creates a new connection between the adapter and the engine.
// It returns a pointer to the connection.
func NewConnection(pid int, in chan string, out chan string) *Connection {
	return &Connection{in, out, pid}
}

// Expect sends a command to the engine and waits for a confirmation. If the
// confirmation is received, the function returns a slice of strings containing
// the responses from the engine.
func (conn *Connection) Expect(cmd string, cnf string) []string {
	var result []string

	conn.in <- cmd

	for resp := range conn.out {
		if resp == cnf {
			return result
		}

		result = append(result, resp)
	}

	return []string{}
}

// Send sends a command to the engine.
func (conn *Connection) Send(cmd string) {
	conn.in <- cmd
}

// Line sends a command to the engine and waits for a response. The response is
// returned as a string.
func (conn *Connection) Line() string {
	return <-conn.out
}

// Scan sends a command to the engine and waits for a response. The response is
// returned as a string. The response is passed to the filter function. If the
// filter function returns true, the function returns.
func (conn *Connection) Scan(cmd string, filter func(resp string) bool) {
	conn.in <- cmd

	for resp := range conn.out {
		if filter(resp) {
			return
		}
	}
}

// Next returns the next response from the engine.
func (conn *Connection) Next() string {
	return <-conn.out
}
