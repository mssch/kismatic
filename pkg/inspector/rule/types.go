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
}

type PackageAvailable struct {
	RuleMeta
	PackageName    string
	PackageVersion string
}

func (p PackageAvailable) Name() string {
	return fmt.Sprintf("%s %s is available", p.PackageName, p.PackageVersion)
}

type PackageInstalled struct {
	RuleMeta
	PackageName    string
	PackageVersion string
}

func (p PackageInstalled) Name() string {
	return fmt.Sprintf("%s %s is installed", p.PackageName, p.PackageVersion)
}

type ExecutableInPath struct {
	RuleMeta
	Executable string
}

func (e ExecutableInPath) Name() string {
	return fmt.Sprintf("%s is in the executable path", e.Executable)
}

type FileContentMatches struct {
	RuleMeta
	File         string
	ContentRegex string
}

func (f FileContentMatches) Name() string {
	return fmt.Sprintf("Contents of %q match the regular expression %s", f.File, f.ContentRegex)
}

type TCPPortAvailable struct {
	RuleMeta
	Port int
}

func (p TCPPortAvailable) Name() string {
	return fmt.Sprintf("Port %d is available", p.Port)
}

type TCPPortAccessible struct {
	RuleMeta
	Port int
}

func (p TCPPortAccessible) Name() string {
	return fmt.Sprintf("Port %d is accessible via the network", p.Port)
}

type RuleResult struct {
	Name        string
	Success     bool
	Error       error
	Remediation string
}
