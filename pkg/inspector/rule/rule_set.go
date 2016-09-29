package rule

import (
	"io"
	"strings"
)

// DefaultRuleSet is the list of rules that are built into the inspector
const defaultRuleSet = `---
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

// DefaultRules returns the list of rules that are built into the inspector
func DefaultRules() []Rule {
	rules, err := unmarshalRules([]byte(defaultRuleSet))
	if err != nil {
		// The default rules should not contain errors
		// If they do, panic so that we catch them during tests
		panic(err)
	}
	return rules
}

// DumpDefaultRules writes the default rule set to a file
func DumpDefaultRules(writer io.Writer) error {
	_, err := io.Copy(writer, strings.NewReader(defaultRuleSet))
	if err != nil {
		return err
	}
	return nil
}
