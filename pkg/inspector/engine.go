package inspector

import "sync"

type Engine struct {
	PackageManager PackageManager
	mu             sync.Mutex
	closableChecks []ClosableCheck
}

func (e *Engine) ExecuteRules(rules []Rule, facts []string) []RuleResult {
	results := []RuleResult{}
	closable := []ClosableCheck{}
	for _, rule := range rules {
		if !shouldExecuteRule(rule, facts) {
			continue
		}

		// Map the rule to a check
		var c Check
		switch r := rule.(type) {
		case PackageInstalled:
			pkgQuery := packageQuery{name: r.PackageName, version: r.PackageVersion}
			c = &PackageInstalledCheck{pkgQuery, e.PackageManager}
		case PackageAvailable:
			pkgQuery := packageQuery{name: r.PackageName, version: r.PackageVersion}
			c = &PackageAvailableCheck{pkgQuery, e.PackageManager}
		case ExecutableInPath:
			c = &BinaryDependencyCheck{r.Executable}
		case TCPPortAvailable:
			check := TCPPortServerCheck{PortNumber: r.Port}
			closable = append(closable, &check)
			c = &check
		}

		// Run the check and report result
		err := c.Check()
		res := RuleResult{
			Name:        rule.Name(),
			Success:     err == nil,
			Error:       err,
			Remediation: "",
		}
		results = append(results, res)
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.closableChecks = closable

	return results
}

func (e *Engine) CloseChecks() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, c := range e.closableChecks {
		if err := c.Close(); err != nil {
			// TODO: Figure out what to do with the error here
		}
	}
	return nil
}

func shouldExecuteRule(rule Rule, facts []string) bool {
	if len(rule.GetRuleMeta().When) == 0 {
		// No conditions on the rule => always run
		return true
	}
	// Run if and only if the all the conditions on the rule are
	// satisfied by the facts
	for _, whenCondition := range rule.GetRuleMeta().When {
		found := false
		for _, l := range facts {
			if whenCondition == l {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
