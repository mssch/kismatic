package check

import (
	"fmt"
	"os/exec"
	"regexp"
)

// BinaryDependencyCheck checks whether the binary is on the executable path
type BinaryDependencyCheck struct {
	BinaryName string
}

// Check returns true if the binary dependency was found
func (c BinaryDependencyCheck) Check() (bool, error) {
	// Need to explicitly call bash when running against Ubuntu
	if err := c.validateBinaryName(); err != nil {
		return false, err
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf("command -v %s", c.BinaryName))
	if err := cmd.Run(); err != nil {
		return false, nil
	}
	return true, nil
}

func (c BinaryDependencyCheck) validateBinaryName() error {
	// We need to ensure that we are not executing arbitrary code here...
	// Only allow single word binary names. Allow dashes.
	r := regexp.MustCompile("^[a-zA-Z-]+$")
	if !r.MatchString(c.BinaryName) {
		return fmt.Errorf("invalid binary name used in check: %s. Names must adhere to the following regexp: %s", c.BinaryName, r.String())
	}
	return nil
}
