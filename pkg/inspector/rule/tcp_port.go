package rule

import (
	"errors"
	"fmt"
	"time"
)

// TCPPortAvailable is a rule that ensures that a given port is available
// on the node. Available means that the port is not being used by another
// process.
type TCPPortAvailable struct {
	Meta
	Port int
}

// Name is the name of the rule
func (p TCPPortAvailable) Name() string {
	return fmt.Sprintf("Port Available: %d", p.Port)
}

// IsRemoteRule returns true if the rule is to be run from outside the node
func (p TCPPortAvailable) IsRemoteRule() bool { return false }

// Validate the rule
func (p TCPPortAvailable) Validate() []error {
	if p.Port < 1 || p.Port > 65535 {
		return []error{fmt.Errorf("Invalid port number %d specified", p.Port)}
	}
	return nil
}

// TCPPortAccessible is a rule that ensures the given port on a remote node
// is accessible from the network
type TCPPortAccessible struct {
	Meta
	Port    int
	Timeout string
}

// Name returns the name of the rule
func (p TCPPortAccessible) Name() string {
	return fmt.Sprintf("Port Accessible: %d", p.Port)
}

// IsRemoteRule returns true if the rule is to be run from a remote node
func (p TCPPortAccessible) IsRemoteRule() bool { return true }

// Validate the rule
func (p TCPPortAccessible) Validate() []error {
	errs := []error{}
	if p.Port < 1 || p.Port > 65535 {
		errs = append(errs, fmt.Errorf("Invalid port number %d specified", p.Port))
	}
	if p.Timeout == "" {
		errs = append(errs, errors.New("Timeout cannot be empty"))
	}
	if p.Timeout != "" {
		if _, err := time.ParseDuration(p.Timeout); err != nil {
			errs = append(errs, fmt.Errorf("Invalid duration provided %q", p.Timeout))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
