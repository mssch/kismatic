package rule

import "fmt"

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

type PackageInstalled struct {
	RuleMeta
	PackageName    string
	PackageVersion string
}

func (p PackageInstalled) Name() string {
	return fmt.Sprintf("%s %s is installed", p.PackageName, p.PackageVersion)
}

func (p PackageInstalled) IsRemoteRule() bool { return false }

type ExecutableInPath struct {
	RuleMeta
	Executable string
}

func (e ExecutableInPath) Name() string {
	return fmt.Sprintf("%s is in the executable path", e.Executable)
}

func (e ExecutableInPath) IsRemoteRule() bool { return false }

type FileContentMatches struct {
	RuleMeta
	File         string
	ContentRegex string
}

func (f FileContentMatches) Name() string {
	return fmt.Sprintf("Contents of %q match the regular expression %s", f.File, f.ContentRegex)
}

func (f FileContentMatches) IsRemoteRule() bool { return false }

type TCPPortAvailable struct {
	RuleMeta
	Port int
}

func (p TCPPortAvailable) Name() string {
	return fmt.Sprintf("Port %d is available", p.Port)
}

func (p TCPPortAvailable) IsRemoteRule() bool { return false }

type TCPPortAccessible struct {
	RuleMeta
	Port    int
	Timeout string
}

func (p TCPPortAccessible) Name() string {
	return fmt.Sprintf("Port %d is accessible via the network", p.Port)
}

func (p TCPPortAccessible) IsRemoteRule() bool { return true }

type RuleResult struct {
	Name        string
	Success     bool
	Error       string
	Remediation string
}
