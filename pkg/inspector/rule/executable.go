package rule

import (
	"errors"
	"fmt"
)

// ExecutableInPath is a rule that ensures the given executable is in
// the system's path
type ExecutableInPath struct {
	Meta
	Executable string
}

// Name is the name of the rule
func (e ExecutableInPath) Name() string {
	return fmt.Sprintf("Executable In Path: %s", e.Executable)
}

// IsRemoteRule returns true if the rule is to be run from outside of the node
func (e ExecutableInPath) IsRemoteRule() bool { return false }

// Validate the rule
func (e ExecutableInPath) Validate() []error {
	if e.Executable == "" {
		return []error{errors.New("Executable cannot be empty")}
	}
	return nil
}
