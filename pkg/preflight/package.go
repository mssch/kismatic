package preflight

import (
	"fmt"
	"os/exec"
)

// PackageInstalledCheck verifies that a given package is installed on the host
type PackageInstalledCheck struct {
	PackageName string
}

// Check returns nil if the package is installed. Otherwise, returns an error message indicating the package was not found.
func (c *PackageInstalledCheck) Check() error {
	cmd, err := getPackageCheckCommand(c.PackageName)
	if err != nil {
		return err
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Install %q, as it was not found on the system.", c.PackageName)
	}
	return nil
}

func (c *PackageInstalledCheck) Name() string {
	return fmt.Sprintf("%s is intalled", c.PackageName)
}

// Returns a distro-specific command for checking installed packages
func getPackageCheckCommand(packageName string) (*exec.Cmd, error) {
	rhelPackageManager := BinaryDependencyCheck{
		BinaryName: "yum",
	}
	if rhelPackageManager.Check() == nil {
		return exec.Command("yum", "list", "installed", packageName), nil
	}
	ubuntuPackageManager := BinaryDependencyCheck{
		BinaryName: "apt-get",
	}
	if ubuntuPackageManager.Check() == nil {
		return exec.Command("apt", "list", "--installed", packageName), nil
	}
	return nil, fmt.Errorf("attempting to check dependency on unsupported distribution.")
}

// PackageAvailableCheck verifies that a given package is available for download
// using the operating system's package manager.
type PackageAvailableCheck struct {
	PackageName string
	Version     string
}

func (c PackageAvailableCheck) Name() string {
	return fmt.Sprintf("%s %s is available", c.PackageName, c.Version)
}

// Check returns nil if the package manager lists the package as available.
// Otherwise returns an error message.
func (c PackageAvailableCheck) Check() error {
	return nil
}
