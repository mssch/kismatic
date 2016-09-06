package check

import (
	"fmt"
	"os/exec"
)

// packageInstalledCheck verifies that a given package is installed on the host
type PackageInstalledCheck struct {
	PackageName string
}

// Check returns nil if the package is installed. Otherwise, returns an error message indicating the package was not found.
func (c *PackageInstalledCheck) Check() error {
	cmd := exec.Command("yum", "list", "installed", c.PackageName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Install %q, as it was not found on the system.", c.PackageName)
	}
	return nil
}

func (c *PackageInstalledCheck) Name() string {
	return fmt.Sprintf("%s is intalled", c.PackageName)
}
