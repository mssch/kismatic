package install

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
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
}

// Explain the incoming ansible event stream
func (e *AnsibleEventExplainer) Explain(in io.Reader) error {
	events := e.EventStream(in)
	for ev := range events {
		switch event := ev.(type) {
		default:
			fmt.Fprintf(e.Out, "Unhandled event: %T\n", event)
		case *ansible.PlaybookStartEvent:
			fmt.Fprintf(e.Out, "Running playbook %s\n", event.Name)
		case *ansible.PlayStartEvent:
			fmt.Fprintf(e.Out, "- %s\n", event.Name)
		case *ansible.RunnerUnreachableEvent:
			fmt.Fprintf(e.Out, "[UNREACHABLE] %s\n", event.Host)
		case *ansible.RunnerFailedEvent:
			fmt.Fprintf(e.Out, "Error from %s: %s\n", event.Host, event.Result.Message)
			if event.Result.Stdout != "" {
				fmt.Fprintf(e.Out, "---- STDOUT ----\n%s\n", event.Result.Stdout)
			}
			if event.Result.Stderr != "" {
				fmt.Fprintf(e.Out, "---- STDERR ----\n%s\n", event.Result.Stderr)
			}
			if event.Result.Stderr != "" || event.Result.Stdout != "" {
				fmt.Fprint(e.Out, "---------------\n")
			}

		// Do nothing with the following events
		case *ansible.RunnerItemRetryEvent:
			continue
		case *ansible.TaskStartEvent:
			continue
		case *ansible.HandlerTaskStartEvent:
			continue
		case *ansible.RunnerItemOKEvent:
			continue
		case *ansible.RunnerSkippedEvent:
			continue
		case *ansible.RunnerOKEvent:
			continue
		}
	}
	return nil
}
