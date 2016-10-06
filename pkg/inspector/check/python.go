package check

import (
	"errors"
	"os/exec"
	"strings"
)

// Python2Check returns true if python 2 is installed on the node
// and the version prefix-matches one of the supported versions
type Python2Check struct {
	SupportedVersions []string
}

func (c Python2Check) Check() (bool, error) {
	cmd := exec.Command("python", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New("Python 2 doesn't seem to be installed")
	}
	for _, version := range c.SupportedVersions {
		if strings.HasPrefix(string(out), version) {
			return true, nil
		}
	}
	return false, nil
}
