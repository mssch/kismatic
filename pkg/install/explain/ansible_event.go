package explain

import (
	"bytes"
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
	"github.com/apprenda/kismatic-platform/pkg/util"
)

// AnsibleEventStreamExplainer explains the incoming ansible event stream
type AnsibleEventStreamExplainer struct {
	// Out is the destination where the explanations are written
	Out io.Writer
	// Verbose is used to control the output level
	Verbose bool
	// EventExplainer for processing ansible events
	EventExplainer AnsibleEventExplainer
}

// Explain the incoming ansible event stream
func (e *AnsibleEventStreamExplainer) Explain(events <-chan ansible.Event) error {
	for event := range events {
		exp := e.EventExplainer.ExplainEvent(event, e.Verbose)
		if exp != "" {
			fmt.Fprint(e.Out, exp)
		}
	}
	return nil
}

// AnsibleEventExplainer explains a single event
type AnsibleEventExplainer interface {
	ExplainEvent(e ansible.Event, verbose bool) string
}

// DefaultEventExplainer returns the default string explanation of a given event
type DefaultEventExplainer struct {
	// Keeping this state is necessary for supporting the current way of
	// printing output to the console... I am not a fan of this, but it'll
	// do for now...
	lastPlay          string
	FirstErrorPrinted bool
}

// ExplainEvent returns an explanation for the given event
func (explainer *DefaultEventExplainer) ExplainEvent(e ansible.Event, verbose bool) string {
	buf := &bytes.Buffer{}
	switch event := e.(type) {
	default:
		if verbose {
			util.PrettyPrintWarnf(buf, "Unhandled event: %T", event)
		}
	case *ansible.PlayStartEvent:
		if explainer.lastPlay != "" {
			util.PrintOk(buf, "[OK]")
		}
		util.PrettyPrintf(buf, "%s\t", event.Name)
		explainer.lastPlay = event.Name
	case *ansible.RunnerUnreachableEvent:
		util.PrintErrorf(buf, "[UNREACHABLE] %s", event.Host)
	case *ansible.RunnerFailedEvent:
		if event.IgnoreErrors {
			return ""
		}
		if !explainer.FirstErrorPrinted {
			util.PrintError(buf, "[ERROR]\n")
		}
		util.PrintErrorf(buf, "Error from %s: %s", event.Host, event.Result.Message)
		if event.Result.Stdout != "" {
			util.PrettyPrintf(buf, "---- STDOUT ----\n%s\n", event.Result.Stdout)
		}
		if event.Result.Stderr != "" {
			util.PrettyPrintf(buf, "---- STDERR ----\n%s\n", event.Result.Stderr)
		}
		if event.Result.Stderr != "" || event.Result.Stdout != "" {
			util.PrettyPrintf(buf, "---------------\n")
		}

	// Do nothing with the following events
	case *ansible.RunnerItemRetryEvent:
		return ""
	case *ansible.PlaybookStartEvent:
		if verbose {
			util.PrettyPrintf(buf, "Running playbook %s\n", event.Name)
		}
	case *ansible.TaskStartEvent:
		if verbose {
			util.PrettyPrintf(buf, "- Running task: %s\n", event.Name)
		}
	case *ansible.HandlerTaskStartEvent:
		if verbose {
			util.PrettyPrintf(buf, "- Running task: %s\n", event.Name)
		}
	case *ansible.RunnerItemOKEvent:
		if verbose {
			util.PrettyPrintf(buf, "   [OK] %s\n", event.Host)
		}
	case *ansible.RunnerSkippedEvent:
		if verbose {
			util.PrettyPrintf(buf, "   [SKIPPED] %s\n", event.Host)
		}
	case *ansible.RunnerOKEvent:
		if verbose {
			util.PrettyPrintf(buf, "   [OK] %s\n", event.Host)
		}
	}
	return buf.String()
}
