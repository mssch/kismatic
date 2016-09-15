package explain

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
	"github.com/apprenda/kismatic-platform/pkg/preflight"
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
		results := []preflight.CheckResult{}
		if err := json.Unmarshal([]byte(event.Result.Stdout), &results); err != nil {
			// Something actually went wrong running the play... use the default explainer
			return explainer.DefaultExplainer.ExplainEvent(event, verbose)
		}

		// We got the results from the inspector... explain results:
		if !explainer.DefaultExplainer.FirstErrorPrinted {
			util.PrintError(buf, "[ERROR]\n")
			explainer.DefaultExplainer.FirstErrorPrinted = true
		}
		util.PrettyPrintErrf(buf, "=> Pre-Flight Checks failed on %q:", event.Host)
		for _, r := range results {

			if !r.Success {
				buf.WriteString(fmt.Sprintf("   - %s\n", r.Error))
			}
		}
		if verbose {
			buf.WriteString("\n=> Successful pre-flight checks:\n")
			for _, r := range results {
				if r.Success {
					buf.WriteString(fmt.Sprintf("   - %q\n", r.Name))
				}
			}
		}
		buf.WriteString("\n")
		return buf.String()
	}
}
