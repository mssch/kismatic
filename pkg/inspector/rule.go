package inspector

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
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

type catchAllRule struct {
	RuleMeta       `yaml:",inline"`
	PackageName    string `yaml:"packageName"`
	PackageVersion string `yaml:"packageVersion"`
	Executable     string `yaml:"executable"`
	Port           int    `yaml:"port"`
	File           string `yaml:"file"`
	ContentRegex   string `yaml:"contentRegex"`
}

func DefaultRules() []Rule {
	y := `---
- kind: PackageAvailable
  when: ["centos"]
  packageName: somePackage
  packageVersion: 1.0

- kind: PackageAvailable
  when: ["ubuntu"]
  packageName: otherPackage
  packageVersion: 1.2

- kind: PackageInstalled
  when: []
  packageName: docker
  packageVersion: 1.11
`
	rules, err := unmarshalRules([]byte(y))
	if err != nil {
		// The default rules should not contain errors
		// If they do, panic so that we catch them during tests
		panic(err)
	}
	return rules
}

func unmarshalRules(data []byte) ([]Rule, error) {
	catchAllRules := []catchAllRule{}
	if err := yaml.Unmarshal(data, &catchAllRules); err != nil {
		return nil, err
	}
	rules := []Rule{}
	for _, catchAllRule := range catchAllRules {
		r, err := buildRule(catchAllRule)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func buildRule(catchAll catchAllRule) (Rule, error) {
	kind := strings.ToLower(strings.TrimSpace(catchAll.Kind))
	switch kind {
	default:
		return nil, fmt.Errorf("invalid rule kind %q was provided", kind)
	case "packageavailable":
		r := PackageAvailable{
			PackageName:    catchAll.PackageName,
			PackageVersion: catchAll.PackageVersion,
		}
		r.When = catchAll.When
		return r, nil
	case "packageinstalled":
		r := PackageInstalled{
			PackageName:    catchAll.PackageName,
			PackageVersion: catchAll.PackageVersion,
		}
		r.When = catchAll.When
		return r, nil
	case "executableinpath":
		r := ExecutableInPath{
			Executable: catchAll.Executable,
		}
		r.When = catchAll.When
		return r, nil
	case "tcpportavailable":
		r := TCPPortAvailable{
			Port: catchAll.Port,
		}
		r.When = catchAll.When
		return r, nil
	case "tcpportaccessible":
		r := TCPPortAccessible{
			Port: catchAll.Port,
		}
		r.When = catchAll.When
		return r, nil
	case "filecontentmatches":
		r := FileContentMatches{
			File:         catchAll.File,
			ContentRegex: catchAll.ContentRegex,
		}
		r.When = catchAll.When
		return r, nil
	}
}
