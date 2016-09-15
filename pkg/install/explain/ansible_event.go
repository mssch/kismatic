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
	// EventStream is a function that returns a channel of ansible Events
	EventStream func(in io.Reader) <-chan ansible.Event
	// Out is the destination where the explanations are written
	Out io.Writer
	// Verbose is used to control the output level
	Verbose bool
	// ExplainEvent returns a string explanation fo the ansible event.
	// The function returns an empty string if the event should be ignored.
	// ExplainEvent   func(e ansible.Event, verbose bool) string
	EventExplainer AnsibleEventExplainer
}

// Explain the incoming ansible event stream
func (e *AnsibleEventStreamExplainer) Explain(in io.Reader) error {
	events := e.EventStream(in)
	for ev := range events {
		exp := e.EventExplainer.ExplainEvent(ev, e.Verbose)
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
	lastTask          string
	firstTaskPrinted  bool
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
			if verbose {
				util.PrintOkf(buf, "%s Finished", explainer.lastPlay)
			} else {
				util.PrintOk(buf, "[OK]")
			}
		}
		util.PrettyPrintf(buf, "%s\t", event.Name)

		explainer.lastPlay = event.Name
		explainer.firstTaskPrinted = false
	case *ansible.RunnerUnreachableEvent:
		if !verbose {
			util.PrintWarn(buf, "[WARNING]")
			explainer.lastPlay = ""
		}
		util.PrettyPrintUnreachablef(buf, "  %s", event.Host)
	case *ansible.RunnerFailedEvent:
		if event.IgnoreErrors {
			if !verbose {
				util.PrintWarn(buf, "[WARNING]")
				explainer.lastPlay = ""
			}
			util.PrettyPrintf(buf, "- Running task: %s\n", explainer.lastTask)
			util.PrettyPrintErrorIgnoredf(buf, "  %s", event.Host)
		} else {
			util.PrettyPrintErrf(buf, "  %s %s", event.Host, event.Result.Message)
		}
		if event.Result.Stdout != "" {
			util.PrintErrorf(buf, "---- STDOUT ----\n%s", event.Result.Stdout)
		}
		if event.Result.Stderr != "" {
			util.PrintErrorf(buf, "---- STDERR ----\n%s", event.Result.Stderr)
		}
		if event.Result.Stderr != "" || event.Result.Stdout != "" {
			util.PrintErrorf(buf, "---------------")
		}
	// Do nothing with the following events
	case *ansible.RunnerItemRetryEvent:
		return ""
	case *ansible.PlaybookStartEvent:
		if verbose {
			util.PrettyPrintf(buf, "Running playbook %s\n", event.Name)
		}
	case *ansible.TaskStartEvent:
		explainer.lastTask = event.Name
		if verbose {
			if !explainer.firstTaskPrinted {
				util.PrettyPrint(buf, "\n")
				explainer.firstTaskPrinted = true
			}
			util.PrettyPrintf(buf, "- Running task: %s\n", event.Name)
		}
	case *ansible.HandlerTaskStartEvent:
		explainer.lastTask = event.Name
		if verbose {
			if !explainer.firstTaskPrinted {
				util.PrettyPrint(buf, "\n")
				explainer.firstTaskPrinted = true
			}
			util.PrettyPrintf(buf, "- Running task: %s\n", event.Name)
		}
	case *ansible.RunnerItemOKEvent:
		if verbose {
			util.PrettyPrintOkf(buf, "  %s", event.Host)
		}
	case *ansible.RunnerSkippedEvent:
		if verbose {
			util.PrettyPrintSkippedf(buf, "  %s", event.Host)
		}
	case *ansible.RunnerOKEvent:
		if verbose {
			util.PrettyPrintOkf(buf, "  %s", event.Host)
		}
	}
	return buf.String()
}
