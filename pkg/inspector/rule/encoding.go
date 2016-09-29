package rule

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// ReadFromFile returns the list of rules contained in the specified file
func ReadFromFile(file string) ([]Rule, error) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, fmt.Errorf("attempted to read rules from non-existent file %q", file)
	}
	rawRules, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading rules from %q: %v", file, err)
	}
	rules, err := UnmarshalRulesYAML(rawRules)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling rules from %q: %v", file, err)
	}
	// TODO: Validate the rules we just read
	return rules, nil
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

// UnmarshalRulesYAML unmarshals the data into a list of rules
func UnmarshalRulesYAML(data []byte) ([]Rule, error) {
	catchAllRules := []catchAllRule{}
	if err := yaml.Unmarshal(data, &catchAllRules); err != nil {
		return nil, err
	}
	return rulesFromCatchAllRules(catchAllRules)
}

// UnmarshalRulesJSON unmarshals the JSON rules into a list of rules
func UnmarshalRulesJSON(data []byte) ([]Rule, error) {
	catchAllRules := []catchAllRule{}
	if err := json.Unmarshal(data, &catchAllRules); err != nil {
		return nil, err
	}
	return rulesFromCatchAllRules(catchAllRules)
}

func rulesFromCatchAllRules(catchAllRules []catchAllRule) ([]Rule, error) {
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
	meta := RuleMeta{
		Kind: kind,
		When: catchAll.When,
	}
	switch kind {
	default:
		return nil, fmt.Errorf("invalid rule kind %q was provided", kind)
	case "packageavailable":
		r := PackageAvailable{
			PackageName:    catchAll.PackageName,
			PackageVersion: catchAll.PackageVersion,
		}
		r.RuleMeta = meta
		return r, nil
	case "packageinstalled":
		r := PackageInstalled{
			PackageName:    catchAll.PackageName,
			PackageVersion: catchAll.PackageVersion,
		}
		r.RuleMeta = meta
		return r, nil
	case "executableinpath":
		r := ExecutableInPath{
			Executable: catchAll.Executable,
		}
		r.RuleMeta = meta
		return r, nil
	case "tcpportavailable":
		r := TCPPortAvailable{
			Port: catchAll.Port,
		}
		r.RuleMeta = meta
		return r, nil
	case "tcpportaccessible":
		r := TCPPortAccessible{
			Port: catchAll.Port,
		}
		r.RuleMeta = meta
		return r, nil
	case "filecontentmatches":
		r := FileContentMatches{
			File:         catchAll.File,
			ContentRegex: catchAll.ContentRegex,
		}
		r.RuleMeta = meta
		return r, nil
	}
}
