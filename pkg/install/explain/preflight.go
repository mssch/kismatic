package explain

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
	"github.com/apprenda/kismatic-platform/pkg/inspector"
	"github.com/apprenda/kismatic-platform/pkg/util"
)

// PreflightEventExplainer explains the Ansible events that run
// when doing the preflight checks
type PreflightEventExplainer struct {
	DefaultExplainer *DefaultEventExplainer
}

// ExplainEvent explains the pre-flight check error events,
// while delegating other event types to the regular, text-based, event explainer
func (explainer *PreflightEventExplainer) ExplainEvent(e ansible.Event, verbose bool) string {
	switch event := e.(type) {
	default:
		return explainer.DefaultExplainer.ExplainEvent(event, verbose)
	case *ansible.RunnerFailedEvent:
		if event.IgnoreErrors {
			return ""
		}
		buf := &bytes.Buffer{}
		results := []inspector.CheckResult{}
		if err := json.Unmarshal([]byte(event.Result.Stdout), &results); err != nil {
			// Something actually went wrong running the play... use the default explainer
			return explainer.DefaultExplainer.ExplainEvent(event, verbose)
		}
		util.PrintColor(buf, util.Red, "\n=> Pre-Flight Checks failed on %q:\n", event.Host)
		for _, r := range results {
			if !r.Success {
				util.PrintColor(buf, util.Red, "   - %s\n", r.Error)
			}
		}
		if verbose {
			util.PrintColor(buf, util.Green, "=> Successful pre-flight checks:\n")
			for _, r := range results {
				if r.Success {
					util.PrintColor(buf, util.Green, "   - %s\n", r.Name)
				}
			}
		}
		fmt.Fprintln(buf)
		explainer.DefaultExplainer.printPlayStatus = false
		return buf.String()
	}
}
