package install

import (
	"io"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
	"github.com/apprenda/kismatic-platform/pkg/util"
)

// Explainer reads the incoming stream, and explains to the user what is
// happening
type Explainer interface {
	// Explain the incoming stream
	Explain(in io.Reader) error
}

// RawExplainer dumps out the stream to the user, without performing any interpretation
type RawExplainer struct {
	Out io.Writer
}

// Explain copies the incoming stream to the out stream
func (e *RawExplainer) Explain(in io.Reader) error {
	_, err := io.Copy(e.Out, in)
	return err
}

// AnsibleEventExplainer explains the incoming ansible event stream
type AnsibleEventExplainer struct {
	EventStream func(in io.Reader) <-chan ansible.Event
	Out         io.Writer
	Verbose     bool
}

// Explain the incoming ansible event stream
func (e *AnsibleEventExplainer) Explain(in io.Reader) error {
	eventName := ""
	events := e.EventStream(in)
	for ev := range events {
		switch event := ev.(type) {
		default:
			if e.Verbose {
				util.PrettyPrintWarnf(e.Out, "Unhandled event: %T", event)
			}
		case *ansible.PlayStartEvent:
			if eventName != "" {
				util.PrintOk(e.Out, "[OK]")
			}
			util.PrettyPrintf(e.Out, "%s\t", event.Name)
			eventName = event.Name
		case *ansible.RunnerUnreachableEvent:
			util.PrintErrorf(e.Out, "[UNREACHABLE] %s", event.Host)
		case *ansible.RunnerFailedEvent:
			util.PrintError(e.Out, "[ERROR]")
			util.PrintErrorf(e.Out, "Error from %s: %s", event.Host, event.Result.Message)
			if event.Result.Stdout != "" {
				util.PrettyPrintf(e.Out, "---- STDOUT ----\n%s\n", event.Result.Stdout)
			}
			if event.Result.Stderr != "" {
				util.PrettyPrintf(e.Out, "---- STDERR ----\n%s\n", event.Result.Stderr)
			}
			if event.Result.Stderr != "" || event.Result.Stdout != "" {
				util.PrettyPrintf(e.Out, "---------------\n")
			}

		// Do nothing with the following events
		case *ansible.RunnerItemRetryEvent:
			continue
		case *ansible.PlaybookStartEvent:
			if e.Verbose {
				util.PrettyPrintf(e.Out, "Running playbook %s\n", event.Name)
			}
		case *ansible.TaskStartEvent:
			if e.Verbose {
				util.PrettyPrintf(e.Out, "- Running task: %s\n", event.Name)
			}
		case *ansible.HandlerTaskStartEvent:
			if e.Verbose {
				util.PrettyPrintf(e.Out, "- Running task: %s\n", event.Name)
			}
		case *ansible.RunnerItemOKEvent:
			if e.Verbose {
				util.PrettyPrintf(e.Out, "   [OK] %s\n", event.Host)
			}
		case *ansible.RunnerSkippedEvent:
			if e.Verbose {
				util.PrettyPrintf(e.Out, "   [SKIPPED] %s\n", event.Host)
			}
		case *ansible.RunnerOKEvent:
			if e.Verbose {
				util.PrettyPrintf(e.Out, "   [OK] %s\n", event.Host)
			}
		}
	}
	return nil
}
