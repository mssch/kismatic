package rule

import (
	"fmt"
	"io/ioutil"
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
func DumpDefaultRules(file string) error {
	err := ioutil.WriteFile(file, []byte(defaultRuleSet), 0644)
	if err != nil {
		return fmt.Errorf("error writing default rule set to %q: %v", file, err)
	}
	return nil
}
