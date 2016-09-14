package explain

import (
	"bytes"
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
)

// AnsibleEventStreamExplainer explains the incoming ansible event stream
type AnsibleEventStreamExplainer struct {
	// EventStream is a function that returns a channel of ansible Events
	EventStream func(in io.Reader) <-chan ansible.Event
	// Out is the destination where the explanations are written
	Out io.Writer
	// Verbose is used to control the output level
	Verbose bool
	// ExplainEvent returns a string explanation fo the ansible event.
	// The function returns an empty string if the event should be ignored.
	ExplainEvent func(e ansible.Event, verbose bool) string
}

// Explain the incoming ansible event stream
func (e *AnsibleEventStreamExplainer) Explain(in io.Reader) error {
	events := e.EventStream(in)
	for ev := range events {
		exp := e.ExplainEvent(ev, e.Verbose)
		if exp != "" {
			fmt.Fprint(e.Out, exp)
		}
	}
	return nil
}

// EventExplanationText returns an explanation for the given event
func EventExplanationText(e ansible.Event, verbose bool) string {
	switch event := e.(type) {
	default:
		return fmt.Sprintf("Unhandled event: %T\n", event)
	case *ansible.PlaybookStartEvent:
		return fmt.Sprintf("Running playbook %s\n", event.Name)
	case *ansible.PlayStartEvent:
		return fmt.Sprintf("=> %s\n", event.Name)
	case *ansible.RunnerUnreachableEvent:
		return fmt.Sprintf("[UNREACHABLE] %s\n", event.Host)
	case *ansible.RunnerFailedEvent:
		if event.IgnoreErrors {
			return ""
		}
		buf := bytes.Buffer{}
		buf.WriteString(fmt.Sprintf("Error from %s: %s\n", event.Host, event.Result.Message))
		if event.Result.Stdout != "" {
			buf.WriteString(fmt.Sprintf("---- STDOUT ----\n%s\n", event.Result.Stdout))
		}
		if event.Result.Stderr != "" {
			buf.WriteString(fmt.Sprintf("---- STDERR ----\n%s\n", event.Result.Stderr))
		}
		if event.Result.Stderr != "" || event.Result.Stdout != "" {
			buf.WriteString(fmt.Sprint("---------------\n"))
		}
		return buf.String()
	case *ansible.RunnerItemRetryEvent:
		return ""
	case *ansible.TaskStartEvent:
		if verbose {
			return fmt.Sprintf("- Running task: %s\n", event.Name)
		}
		return ""
	case *ansible.HandlerTaskStartEvent:
		if verbose {
			return fmt.Sprintf("- Running task: %s\n", event.Name)
		}
		return ""
	case *ansible.RunnerItemOKEvent:
		if verbose {
			return fmt.Sprintf("   [OK] %s\n", event.Host)
		}
		return ""
	case *ansible.RunnerSkippedEvent:
		if verbose {
			return fmt.Sprintf("   [SKIPPED] %s\n", event.Host)
		}
		return ""
	case *ansible.RunnerOKEvent:
		if verbose {
			return fmt.Sprintf("   [OK] %s\n", event.Host)
		}
		return ""
	}
}
