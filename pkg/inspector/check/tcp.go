package check

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

// TCPPortClientCheck verifies that a given port on a remote node
// is accessible through the network
type TCPPortClientCheck struct {
	// IPAddress is the IP of the remote node
	IPAddress string
	// PortNumber is the target service port
	PortNumber int
	// Timeout is the maximum amount of time the check will
	// wait when connecting to the server before bailing out
	Timeout time.Duration
}

// Check returns true if the TCP connection is established and the server
// returns the expected response. Otherwise, returns false and an error message
func (c *TCPPortClientCheck) Check() (bool, error) {
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.IPAddress, c.PortNumber), timeout)
	if err != nil {
		return false, fmt.Errorf("Port %d on host %q is unreachable. Error was: %v", c.PortNumber, c.IPAddress, err)
	}
	conn.Close()
	return true, nil
}

// TCPPortServerCheck ensures that the given port is free, or bound to the right
// process. In the case that it is free, it stands up a TCP server that can be
// used to check TCP connectivity to the host using TCPPortClientCheck
type TCPPortServerCheck struct {
	PortNumber     int
	ProcName       string
	started        bool
	closeListener  func() error
	listenerClosed chan interface{}
}

// Check returns true if the port is free, or taken by the expected process.
func (c *TCPPortServerCheck) Check() (bool, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", c.PortNumber))
	if err != nil && strings.Contains(err.Error(), "address already in use") {
		return portTakenByProc(c.PortNumber, c.ProcName)
	}
	if err != nil {
		return false, fmt.Errorf("error listening on port %d: %v", c.PortNumber, err)
	}
	c.closeListener = ln.Close
	// Setup go routine for accepting connections
	c.listenerClosed = make(chan interface{})
	go func(closed <-chan interface{}) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-closed:
					// don't log the error, as we have closed the server and the error
					// is related to that.
					return
				default:
					log.Println(fmt.Sprintf("error occurred accepting request: %v", err))
					continue
				}
			}
			// Setup go routine that behaves as an echo server
			go func(c net.Conn) {
				io.Copy(c, c)
				c.Close()
			}(conn)
		}
	}(c.listenerClosed)
	c.started = true
	return true, nil
}

// Close the TCP server if it was started. Otherwise this is a noop.
func (c *TCPPortServerCheck) Close() error {
	if c.started {
		close(c.listenerClosed)
		return c.closeListener()
	}
	return nil
}

// Returns true if the port is taken by a process with the given name.
func portTakenByProc(port int, procName string) (bool, error) {
	// Use lsof to find the process that is bound to the tcp port in listen
	// mode.
	// ~# lsof -i TCP:2379 -s TCP:LISTEN -Pn +c 0
	// COMMAND       PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
	// docker-proxy 7294 root    4u  IPv6  43407      0t0  TCP *:2379 (LISTEN)
	portArg := fmt.Sprintf("TCP:%d", port)
	cmd := exec.Command("lsof", "-i", portArg, "-s", "TCP:LISTEN", "-Pn", "+c", "0")
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error running lsof: %v", err)
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return false, fmt.Errorf("expected lsof to return at least 2 lines, but returned %d", len(lines))
	}
	// There are cases where lsof will return multiple lines for the same port. For example:
	// ~# lsof -i TCP:$port -s TCP:LISTEN -Pn +c 0
	// COMMAND   PID   USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
	// nginx   18611   root   10u  IPv4  82841      0t0  TCP *:80 (LISTEN)
	// nginx   18628 nobody   10u  IPv4  82841      0t0  TCP *:80 (LISTEN)
	// nginx   18629 nobody   10u  IPv4  82841      0t0  TCP *:80 (LISTEN)
	//
	// Use the first line after the header for verifying the proc name.
	lsofFields := strings.Fields(lines[1])
	return lsofFields[0] == procName, nil
}
