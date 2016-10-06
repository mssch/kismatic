package check

import "fmt"

type PackageQuery struct {
	Name    string
	Version string
}

func (p PackageQuery) String() string {
	return fmt.Sprintf("%s %s", p.Name, p.Version)
}

// PackageAvailableCheck verifies that a given package is available for download
// using the operating system's package manager.
type PackageAvailableCheck struct {
	PackageQuery   PackageQuery
	PackageManager PackageManager
}

// Check returns true if the package is available. Otherwise returns false, or an error
// if the check is unable to determine the condition.
func (c PackageAvailableCheck) Check() (bool, error) {
	ok, err := c.PackageManager.IsAvailable(c.PackageQuery)
	if err != nil {
		return false, fmt.Errorf("failed to determine if %s is available from the operating system's package manager: %v", c.PackageQuery, err)
	}
	return ok, nil
}
