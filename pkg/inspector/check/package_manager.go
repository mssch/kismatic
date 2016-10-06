package check

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// PackageManager runs queries against the underlying operating system's
// package manager
type PackageManager interface {
	IsAvailable(PackageQuery) (bool, error)
}

// NewPackageManager returns a package manager for the given distribution
func NewPackageManager(distro Distro) (PackageManager, error) {
	run := func(name string, arg ...string) ([]byte, error) {
		return exec.Command(name, arg...).CombinedOutput()
	}
	switch distro {
	case RHEL, CentOS:
		return &rpmManager{
			run: run,
		}, nil
	case Ubuntu:
		return &debManager{
			run: run,
		}, nil
	default:
		return nil, fmt.Errorf("%s is not supported", distro)
	}
}

// package manager for EL-based distributions
type rpmManager struct {
	run func(string, ...string) ([]byte, error)
}

func (m rpmManager) IsAvailable(p PackageQuery) (bool, error) {
	out, err := m.run("yum", "list", "-q", p.Name)
	if err != nil && strings.Contains(string(out), "No matching Packages to list") {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to determine if %s %s is available: %v", p.Name, p.Version, err)
	}
	return m.isPackageListed(p, out), nil
}

func (m rpmManager) isPackageListed(p PackageQuery, list []byte) bool {
	s := bufio.NewScanner(bytes.NewReader(list))
	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) != 3 {
			// Ignore lines that don't match the expected format
			continue
		}
		if p.Name == f[0] && p.Version == f[1] {
			return true
		}
	}
	return false
}

// package manager for debian-based distributions
type debManager struct {
	run func(string, ...string) ([]byte, error)
}

func (m debManager) IsAvailable(p PackageQuery) (bool, error) {
	// First check if the package is installed
	installed, err := m.isInstalled(p)
	if err != nil {
		return false, err
	}
	if installed {
		return true, nil
	}
	// If it's not installed, ensure that it is available via the
	// package manager. We attempt to install using --dry-run. If exit status is zero, we
	// know the package is available for download
	_, err = m.run("apt-get", "install", "-q", "--dry-run", fmt.Sprintf("%s=%s", p.Name, p.Version))
	return err == nil, err
}

func (m debManager) isInstalled(p PackageQuery) (bool, error) {
	out, err := m.run("dpkg", "-l", p.Name)
	if err != nil && strings.Contains(string(out), "no packages found matching") {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to determine if %s %s is installed: %v", p.Name, p.Version, err)
	}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) < 5 {
			// Ignore lines with unexpected format
			continue
		}
		if p.Name == f[1] && p.Version == f[2] {
			return true, nil
		}
	}
	return false, nil
}
