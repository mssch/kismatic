package check

import (
	"fmt"
	"os/exec"
	"regexp"
)

// ExecutableInPathCheck checks whether the binary is on the executable path
type ExecutableInPathCheck struct {
	Name string
}

// Check returns true if the executable is in the path
func (c ExecutableInPathCheck) Check() (bool, error) {
	// Need to explicitly call bash when running against Ubuntu
	if err := c.validateExecutableName(); err != nil {
		return false, err
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf("command -v %s", c.Name))
	if err := cmd.Run(); err != nil {
		return false, nil
	}
	return true, nil
}

func (c ExecutableInPathCheck) validateExecutableName() error {
	// We need to ensure that we are not executing arbitrary code here...
	// Only allow single word binary names. Allow dashes.
	r := regexp.MustCompile("^[a-zA-Z-]+$")
	if !r.MatchString(c.Name) {
		return fmt.Errorf("invalid binary name used in check: %s. Names must adhere to the following regexp: %s", c.Name, r.String())
	}
	return nil
}
