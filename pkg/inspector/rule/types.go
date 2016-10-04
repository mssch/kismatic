package rule

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type RuleMeta struct {
	Kind string
	When []string
}

func (rm RuleMeta) GetRuleMeta() RuleMeta {
	return rm
}

type Rule interface {
	Name() string
	GetRuleMeta() RuleMeta
	IsRemoteRule() bool
	Validate() []error
}

type PackageAvailable struct {
	RuleMeta
	PackageName    string
	PackageVersion string
}

func (p PackageAvailable) Name() string {
	return fmt.Sprintf("%s %s is available", p.PackageName, p.PackageVersion)
}

func (p PackageAvailable) IsRemoteRule() bool { return false }

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

type PackageInstalled struct {
	RuleMeta
	PackageName    string
	PackageVersion string
}

func (p PackageInstalled) Name() string {
	return fmt.Sprintf("%s %s is installed", p.PackageName, p.PackageVersion)
}

func (p PackageInstalled) IsRemoteRule() bool { return false }

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

type ExecutableInPath struct {
	RuleMeta
	Executable string
}

func (e ExecutableInPath) Name() string {
	return fmt.Sprintf("%s is in the executable path", e.Executable)
}

func (e ExecutableInPath) IsRemoteRule() bool { return false }

func (e ExecutableInPath) Validate() []error {
	if e.Executable == "" {
		return []error{errors.New("Executable cannot be empty")}
	}
	return nil
}

type FileContentMatches struct {
	RuleMeta
	File         string
	ContentRegex string
}

func (f FileContentMatches) Name() string {
	return fmt.Sprintf("Contents of %q match the regular expression %s", f.File, f.ContentRegex)
}

func (f FileContentMatches) IsRemoteRule() bool { return false }

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

type TCPPortAvailable struct {
	RuleMeta
	Port int
}

func (p TCPPortAvailable) Name() string {
	return fmt.Sprintf("Port %d is available", p.Port)
}

func (p TCPPortAvailable) IsRemoteRule() bool { return false }

func (p TCPPortAvailable) Validate() []error {
	if p.Port < 1 || p.Port > 65535 {
		return []error{fmt.Errorf("Invalid port number %d specified", p.Port)}
	}
	return nil
}

type TCPPortAccessible struct {
	RuleMeta
	Port    int
	Timeout string
}

func (p TCPPortAccessible) Name() string {
	return fmt.Sprintf("Port %d is accessible via the network", p.Port)
}

func (p TCPPortAccessible) IsRemoteRule() bool { return true }

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

type RuleResult struct {
	Name        string
	Success     bool
	Error       string
	Remediation string
}
