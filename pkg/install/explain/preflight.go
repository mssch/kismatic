package explain

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
	"github.com/apprenda/kismatic-platform/pkg/preflight"
)

// PreFlightEventExplanationText explains the pre-flight check error events,
// while delegating other event types to the regular, text-based, event explainer
func PreFlightEventExplanationText(e ansible.Event, verbose bool) string {
	switch event := e.(type) {
	default:
		return EventExplanationText(event, verbose)
	case *ansible.RunnerFailedEvent:
		if event.IgnoreErrors {
			return ""
		}
		results := []preflight.CheckResult{}
		err := json.Unmarshal([]byte(event.Result.Stdout), &results)
		if err != nil {
			return fmt.Sprintf("error explaining pre-flight check result: %v", err)
		}
		buf := bytes.Buffer{}
		buf.WriteString(fmt.Sprintf("\nPre-flight Checks failed on %q\n", event.Host))
		for _, r := range results {
			if r.Success && verbose {
				buf.WriteString(fmt.Sprintf("[OK] %q\n", r.Name))
			}
			if !r.Success {
				buf.WriteString(fmt.Sprintf("[ERROR] %s\n", r.Error))
			}
		}
		buf.WriteString("\n")
		return buf.String()
	}
}
