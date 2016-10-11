package rule

import (
	"errors"
	"fmt"
)

// PythonVersion rule for checking the host's python version
type Python2Version struct {
	Meta
	SupportedVersions []string
}

func (p Python2Version) Name() string {
	return fmt.Sprintf("Python 2 version in %v", p.SupportedVersions)
}

func (p Python2Version) IsRemoteRule() bool { return false }

func (p Python2Version) Validate() []error {
	if len(p.SupportedVersions) == 0 {
		return []error{errors.New("List of supported versions is empty")}
	}
	return nil
}
