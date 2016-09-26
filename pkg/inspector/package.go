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
	packageManager packageManager
}

// Check returns nil if the package is installed. Otherwise, returns an error message indicating the package was not found.
func (c PackageInstalledCheck) Check() error {
	ok, err := c.packageManager.isInstalled(c.pkgQuery)
	if !ok || err != nil {
		return fmt.Errorf("Install %q, as it was not found on the system.", c.pkgQuery)
	}
	return nil
}

func (c PackageInstalledCheck) Name() string {
	return fmt.Sprintf("%s is intalled", c.pkgQuery)
}

// PackageAvailableCheck verifies that a given package is available for download
// using the operating system's package manager.
type PackageAvailableCheck struct {
	pkgQuery       packageQuery
	packageManager packageManager
}

func (c PackageAvailableCheck) Name() string {
	return fmt.Sprintf("%s is available for download.", c.pkgQuery)
}

// Check returns nil if the package manager lists the package as available.
// Otherwise returns an error message.
func (c PackageAvailableCheck) Check() error {
	ok, err := c.packageManager.isAvailable(c.pkgQuery)
	if !ok || err != nil {
		return fmt.Errorf("%s is not available from the operating system's package manager", c.pkgQuery)
	}
	return nil
}
