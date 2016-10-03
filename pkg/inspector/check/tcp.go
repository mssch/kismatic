package check

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
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

	testMsg := "ECHO\n"
	fmt.Fprint(conn, testMsg)
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if resp != testMsg {
		return false, fmt.Errorf("Port %d on host %q did not send the expected response. Response was %q", c.PortNumber, c.IPAddress, resp)
	}
	return true, nil
}

// TCPPortServerCheck ensures that the given port is free, and stands up a TCP server that can be used to
// check TCP connectivity to the host using TCPPortClientCheck
type TCPPortServerCheck struct {
	PortNumber     int
	closeListener  func() error
	listenerClosed chan interface{}
}

// Check returns true if the port is available for the server. Otherwise returns false
// and an error message
func (c *TCPPortServerCheck) Check() (bool, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", c.PortNumber))
	if err != nil {
		// TODO: We could check if the port is being used here..
		return false, fmt.Errorf("Failed to bind port %d. This could mean the port is in use by another process. Error was: %v", c.PortNumber, err)
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
	return true, nil
}

// Close the TCP server
func (c *TCPPortServerCheck) Close() error {
	close(c.listenerClosed)
	return c.closeListener()
}
