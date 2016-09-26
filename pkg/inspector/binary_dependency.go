package inspector

import (
	"fmt"
	"os/exec"
)

// BinaryDependencyCheck checks whether the binary is on the executable path.
type BinaryDependencyCheck struct {
	BinaryName string
}

// Check returns true if the binary dependency was found. Otherwise, returns false with an error message that includes remediation information.
func (c BinaryDependencyCheck) Check() error {
	// Need to explicitly use bash when running against Ubuntu
	cmd := exec.Command("bash", "-c", fmt.Sprintf("command -v %s", c.BinaryName))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Install %q, as it was not found in the system", c.BinaryName)
	}
	return nil
}

func (c BinaryDependencyCheck) Name() string {
	return fmt.Sprintf("%s exists", c.BinaryName)
}
