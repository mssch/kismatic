package rule

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Meta contains the rule's metadata
type Meta struct {
	Kind string
	When []string
}

// GetRuleMeta returns the rule's metadata
func (rm Meta) GetRuleMeta() Meta {
	return rm
}

// Rule is an inspector rule
type Rule interface {
	Name() string
	GetRuleMeta() Meta
	IsRemoteRule() bool
	Validate() []error
}

// PackageAvailable is a rule that ensures that the given package
// is available for download using the operating system's package manager
type PackageAvailable struct {
	Meta
	PackageName    string
	PackageVersion string
}

// Name returns the name of the rule
func (p PackageAvailable) Name() string {
	return fmt.Sprintf("Package Available: %s %s", p.PackageName, p.PackageVersion)
}

// IsRemoteRule returns true if the rule is to be run from outside of the node
func (p PackageAvailable) IsRemoteRule() bool { return false }

// Validate the rule
func (p PackageAvailable) Validate() []error {
	err := []error{}
	if p.PackageName == "" {
		err = append(err, errors.New("PackageName cannot be empty"))
	}
	if p.PackageVersion == "" {
		err = append(err, errors.New("PackageVersion cannot be empty"))
	}
	if len(err) > 0 {
		return err
	}
	return nil
}

// PackageInstalled is a rule that ensures the given package is installed
type PackageInstalled struct {
	Meta
	PackageName    string
	PackageVersion string
}

// Name is the name of the rule
func (p PackageInstalled) Name() string {
	return fmt.Sprintf("Package Installed: %s %s", p.PackageName, p.PackageVersion)
}

// IsRemoteRule returns true if the rule is to be run from outside the node
func (p PackageInstalled) IsRemoteRule() bool { return false }

// Validate the rule
func (p PackageInstalled) Validate() []error {
	err := []error{}
	if p.PackageName == "" {
		err = append(err, errors.New("PackageName cannot be empty"))
	}
	if p.PackageVersion == "" {
		err = append(err, errors.New("PackageVersion cannot be empty"))
	}
	if len(err) > 0 {
		return err
	}
	return nil
}

// ExecutableInPath is a rule that ensures the given executable is in
// the system's path
type ExecutableInPath struct {
	Meta
	Executable string
}

// Name is the name of the rule
func (e ExecutableInPath) Name() string {
	return fmt.Sprintf("Executable In Path: %s", e.Executable)
}

// IsRemoteRule returns true if the rule is to be run from outside of the node
func (e ExecutableInPath) IsRemoteRule() bool { return false }

// Validate the rule
func (e ExecutableInPath) Validate() []error {
	if e.Executable == "" {
		return []error{errors.New("Executable cannot be empty")}
	}
	return nil
}

// FileContentMatches is a rule that verifies that the contents of a file
// match the regular expression provided
type FileContentMatches struct {
	Meta
	File         string
	ContentRegex string
}

// Name is the name of the rule
func (f FileContentMatches) Name() string {
	return fmt.Sprintf("Contents of %q match: %s", f.File, f.ContentRegex)
}

// IsRemoteRule returns true if the rule is to be run from outside of the node
func (f FileContentMatches) IsRemoteRule() bool { return false }

// Validate the rule
func (f FileContentMatches) Validate() []error {
	errs := []error{}
	if f.File == "" {
		errs = append(errs, errors.New("File cannot be empty"))
	}
	if f.ContentRegex == "" {
		errs = append(errs, errors.New("ContentRegex cannot be empty"))
	}
	if f.ContentRegex != "" {
		if _, err := regexp.Compile(f.ContentRegex); err != nil {
			errs = append(errs, fmt.Errorf("ContentRegex contains an invalid regular expression: %v", err))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// TCPPortAvailable is a rule that ensures that a given port is available
// on the node. Available means that the port is not being used by another
// process.
type TCPPortAvailable struct {
	Meta
	Port int
}

// Name is the name of the rule
func (p TCPPortAvailable) Name() string {
	return fmt.Sprintf("Port Available: %d", p.Port)
}

// IsRemoteRule returns true if the rule is to be run from outside the node
func (p TCPPortAvailable) IsRemoteRule() bool { return false }

// Validate the rule
func (p TCPPortAvailable) Validate() []error {
	if p.Port < 1 || p.Port > 65535 {
		return []error{fmt.Errorf("Invalid port number %d specified", p.Port)}
	}
	return nil
}

// TCPPortAccessible is a rule that ensures the given port on a remote node
// is accessible from the network
type TCPPortAccessible struct {
	Meta
	Port    int
	Timeout string
}

// Name returns the name of the rule
func (p TCPPortAccessible) Name() string {
	return fmt.Sprintf("Port Accessible: %d", p.Port)
}

// IsRemoteRule returns true if the rule is to be run from a remote node
func (p TCPPortAccessible) IsRemoteRule() bool { return true }

// Validate the rule
func (p TCPPortAccessible) Validate() []error {
	errs := []error{}
	if p.Port < 1 || p.Port > 65535 {
		errs = append(errs, fmt.Errorf("Invalid port number %d specified", p.Port))
	}
	if p.Timeout == "" {
		errs = append(errs, errors.New("Timeout cannot be empty"))
	}
	if p.Timeout != "" {
		if _, err := time.ParseDuration(p.Timeout); err != nil {
			errs = append(errs, fmt.Errorf("Invalid duration provided %q", p.Timeout))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Result contains the results from executing the rule
type Result struct {
	// Name is the rule's name
	Name string
	// Success is true when the rule was asserted
	Success bool
	// Error message if there was an error executing the rule
	Error string
	// Remediation contains potential remediation steps for the rule
	Remediation string
}
