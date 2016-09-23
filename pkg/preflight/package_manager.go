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

type yumManager struct {
	run func(string, ...string) ([]byte, error)
}

func (m yumManager) isInstalled(p pkg) (bool, error) {
	out, err := m.run("yum", "list", "installed", "-q", p.name)
	if err != nil && strings.Contains(string(out), "No matching Packages to list") {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("unable to determine if %s %s is installed: %v", p.name, p.version, err)
	}
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		line := s.Text()
		f := strings.Fields(line)
		if len(f) != 3 {
			// Ignore lines that don't match the expected format
			continue
		}
		if p.name == f[0] && p.version == f[1] {
			return true, nil
		}
	}
	return false, nil
}

type aptManager struct {
	run func(string, ...string) ([]byte, error)
}

func (m aptManager) isInstalled(p pkg) (bool, error) {
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
		fmt.Println(f)
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
