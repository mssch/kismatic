package check

import (
	"fmt"
	"os/exec"
)

// BinaryDependencyCheck checks whether the binary is on the executable path
type BinaryDependencyCheck struct {
	BinaryName string
}

// Check returns true if the binary dependency was found
func (c BinaryDependencyCheck) Check() (bool, error) {
	// Need to explicitly call bash when running against Ubuntu
	cmd := exec.Command("bash", "-c", fmt.Sprintf("command -v %s", c.BinaryName))
	if err := cmd.Run(); err != nil {
		return false, nil
	}
	return true, nil
}
