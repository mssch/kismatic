package preflight

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

type pkg struct {
	name    string
	version string
}

type packageManager interface {
	isInstalled(pkg) (bool, error)
	isAvailable(pkg) (bool, error)
}

type rpmManager struct {
	run func(string, ...string) ([]byte, error)
}

func (m rpmManager) isInstalled(p pkg) (bool, error) {
	out, err := m.run("yum", "list", "installed", "-q", p.name)
	if err != nil && strings.Contains(string(out), "No matching Packages to list") {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to determine if %s %s is installed: %v", p.name, p.version, err)
	}
	return m.isPackageListed(p, out), nil
}

func (m rpmManager) isAvailable(p pkg) (bool, error) {
	out, err := m.run("yum", "list", "-q", p.name)
	if err != nil && strings.Contains(string(out), "No matching Packages to list") {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to determine if %s %s is available for download: %v", p.name, p.version, err)
	}
	return m.isPackageListed(p, out), nil
}

func (m rpmManager) isPackageListed(p pkg, list []byte) bool {
	s := bufio.NewScanner(bytes.NewReader(list))
	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) != 3 {
			// Ignore lines that don't match the expected format
			continue
		}
		if p.name == f[0] && p.version == f[1] {
			return true
		}
	}
	return false
}

type debManager struct {
	run func(string, ...string) ([]byte, error)
}

func (m debManager) isInstalled(p pkg) (bool, error) {
	out, err := m.run("dpkg", "-l", p.name)
	if err != nil && strings.Contains(string(out), "no packages found matching") {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to determine if %s %s is installed: %v", p.name, p.version, err)
	}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) < 5 {
			// Ignore lines with unexpected format
			continue
		}
		if p.name == f[1] && p.version == f[2] {
			return true, nil
		}
	}
	return false, nil
}

func (m debManager) isAvailable(p pkg) (bool, error) {
	// Attempt to install using --dry-run. If exit status is zero, we
	// know the package is available for download
	_, err := m.run("apt-get", "install", "-q", "--dry-run", fmt.Sprintf("%s=%s", p.name, p.version))
	return err == nil, err
}
