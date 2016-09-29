package inspector

import "fmt"

type packageQuery struct {
	name    string
	version string
}

func (p packageQuery) String() string {
	return fmt.Sprintf("%s %s", p.name, p.version)
}

// PackageInstalledCheck verifies that a given package is installed on the host
type PackageInstalledCheck struct {
	pkgQuery       packageQuery
	packageManager PackageManager
}

// Check returns true if the package is installed. If an error occurrs,
// returns false and the error.
func (c PackageInstalledCheck) Check() (bool, error) {
	ok, err := c.packageManager.IsInstalled(c.pkgQuery)
	if err != nil {
		return false, fmt.Errorf("Failed to determine if %q is installed on the system: %v", c.pkgQuery, err)
	}
	return ok, nil
}

// PackageAvailableCheck verifies that a given package is available for download
// using the operating system's package manager.
type PackageAvailableCheck struct {
	pkgQuery       packageQuery
	packageManager PackageManager
}

// Check returns true if the package is available. Otherwise returns false, or an error
// if the check is unable to determine the condition.
func (c PackageAvailableCheck) Check() (bool, error) {
	ok, err := c.packageManager.IsAvailable(c.pkgQuery)
	if err != nil {
		return false, fmt.Errorf("failed to determine if %s is available from the operating system's package manager: %v", c.pkgQuery, err)
	}
	return ok, nil
}
