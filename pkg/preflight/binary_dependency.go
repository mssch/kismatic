package preflight

import (
	"fmt"
	"os/exec"
)

// binaryDependencyCheck checks whether the binary is on the executable path.
type binaryDependencyCheck struct {
	BinaryName string
}

// Check returns true if the binary dependency was found. Otherwise, returns false with an error message that includes remediation information.
func (c *binaryDependencyCheck) Check() error {
	cmd := exec.Command("command", "-v", c.BinaryName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Install %q, as it was not found in the system", c.BinaryName)
	}
	return nil
}

func (c *binaryDependencyCheck) Name() string {
	return fmt.Sprintf("%s exists", c.BinaryName)
}
