package preflight

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

// tcpPortClientCheck is a client that can be used to verify the TCPPortServerCheck
type tcpPortClientCheck struct {
	portNumber int
	ipAddress  string
}

// Check returns nil if the TCP connection is established, and the TCPPortServerCheck is running on the other side.
// Otherwise, returns an error message indicating the problem.
func (c *tcpPortClientCheck) Check() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.ipAddress, c.portNumber))
	if err != nil {
		return fmt.Errorf("Port %d on host %q is unreachable. This might mean there is a firewall blocking access to the port,"+
			"or there is nothing listening on the other end. Error was: %v", c.portNumber, c.ipAddress, err)
	}

	testMsg := "ECHO\n"
	fmt.Fprint(conn, testMsg)
	resp, err := bufio.NewReader(conn).ReadString('\n')
	if resp != testMsg {
		return fmt.Errorf("Port %d on host %q did not send the expected response. Response was %q", c.portNumber, c.ipAddress, resp)
	}
	return nil
}

func (c *tcpPortClientCheck) Name() string {
	return fmt.Sprintf("TCP Port %d accessible", c.portNumber)
}

// tcpPortServerCheck ensures that the given port is free, and stands up a TCP server that can be used to
// check TCP connectivity to the host using TCPPortClientCheck
type tcpPortServerCheck struct {
	PortNumber   int
	serverCloser func() error
}

// Check returns nil if the port is available, and the TCP listener is up and running.
// Otherwise, returns an error message indicating the failure.
func (c *tcpPortServerCheck) Check() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", c.PortNumber))
	if err != nil {
		// TODO: We could check if the port is being used here..
		return fmt.Errorf("Attempted to bind port %d but failed. This might mean the port is in use by another process. Error was: %v", c.PortNumber, err)
	}
	c.serverCloser = ln.Close
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				io.Copy(c, c)
				c.Close()
			}(conn)
		}
	}()
	return nil
}

// Close the TCP server
func (c *tcpPortServerCheck) Close() error {
	return c.serverCloser()
}

// Name of the check
func (c *tcpPortServerCheck) Name() string {
	return fmt.Sprintf("TCP Port %d bindable", c.PortNumber)
}
