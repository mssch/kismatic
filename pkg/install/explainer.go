package install

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
	"github.com/apprenda/kismatic-platform/pkg/preflight"
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

// EventExplainer explains a single Ansible event
type EventExplainer interface {
	// Explain returns a string explanation of the given event. If the event is ignored by the explainer,
	// it returns empty string.
	Explain(e ansible.Event, verbose bool) string
}

// AnsibleEventExplainer explains the incoming ansible event stream
type AnsibleEventExplainer struct {
	EventStream    func(in io.Reader) <-chan ansible.Event
	Out            io.Writer
	Verbose        bool
	EventExplainer EventExplainer
}

// Explain the incoming ansible event stream
func (e *AnsibleEventExplainer) Explain(in io.Reader) error {
	events := e.EventStream(in)
	for ev := range events {
		exp := e.EventExplainer.Explain(ev, e.Verbose)
		if exp != "" {
			fmt.Fprint(e.Out, exp)
		}
	}
	return nil
}

// EventExplainerFunc is an adapter to allow use of ordinary functions
// as explainers
type EventExplainerFunc func(e ansible.Event, verbose bool) string

// Explain runs the explainer functions
func (f EventExplainerFunc) Explain(e ansible.Event, verbose bool) string {
	return f(e, verbose)
}

// CLIEventExplanation returns an explanation for the given event suitable
// for printing to a CLI
func CLIEventExplanation(e ansible.Event, verbose bool) string {
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
		fmt.Sprintf("Error from %s: %s\n", event.Host, event.Result.Message)
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

// PreFlightCLIExplanation explains the pre-flight check error events,
// while delegating other event types to the regular CLI event explainer
func PreFlightCLIExplanation(e ansible.Event, verbose bool) string {
	switch event := e.(type) {
	default:
		return CLIEventExplanation(event, verbose)
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
