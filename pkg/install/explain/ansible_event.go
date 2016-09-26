package explain

import (
	"bytes"
	"fmt"
	"io"
	"strings"

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
	printPlayMessage bool
	printPlayStatus  bool
	lastPlay         string
	currentTask      string
	playCount        int
	currentPlayCount int
}

func (explainer *DefaultEventExplainer) getCount() string {
	return rightPadToLen(fmt.Sprintf("%d/%d", explainer.currentPlayCount, explainer.playCount), ".", 7)
}

func rightPadToLen(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

// ExplainEvent returns an explanation for the given event
func (explainer *DefaultEventExplainer) ExplainEvent(e ansible.Event, verbose bool) string {
	buf := &bytes.Buffer{}
	switch event := e.(type) {
	case *ansible.PlayStartEvent:
		// On a play start the previos play ends
		// Print a success status, but only when there were no errors
		if verbose {
			// In verbose mode the status is printed as a whole line after all the tasks
			// Dont print message before first play
			if explainer.printPlayMessage {
				// No tasks were printed, add a new line: something is wrong
				if explainer.printPlayStatus {
					fmt.Fprintln(buf)
					util.PrintColor(buf, util.Red, "%s  %s Finished with no tasks, are hosts reachable?\n", explainer.getCount(), explainer.lastPlay)
				} else {
					util.PrintColor(buf, util.Green, "%s  %s Finished\n", explainer.getCount(), explainer.lastPlay)
				}
				explainer.currentPlayCount = explainer.currentPlayCount + 1
			}
			// Print the play name in color
			util.PrintColor(buf, util.White, "%s  %s", explainer.getCount(), event.Name)
		} else {
			// Do not print status on the first start event or when there is an ERROR
			if explainer.printPlayStatus {
				// In regular mode print the status
				util.PrintOkln(buf)
			}
			// Print the play name
			util.PrettyPrint(buf, "%s  %s", explainer.getCount(), event.Name)
			explainer.currentPlayCount = explainer.currentPlayCount + 1
		}
		// Set default state for the play
		explainer.lastPlay = event.Name
		explainer.printPlayStatus = true
		explainer.printPlayMessage = true
	case *ansible.RunnerFailedEvent:
		// An error
		// Print newline before first task status
		if explainer.printPlayStatus {
			fmt.Fprintln(buf)
			// Dont print play success status on error
			explainer.printPlayStatus = false
		}
		// Tasks only print at verbose level, on ERROR also print task name
		if !verbose {
			fmt.Fprintf(buf, "- Running task: %s\n", explainer.currentTask)
		}
		if event.IgnoreErrors {
			util.PrettyPrintErrorIgnored(buf, "  %s", event.Host)
		} else {
			util.PrettyPrintErr(buf, "  %s %s", event.Host, event.Result.Message)
		}
		if event.Result.Stdout != "" {
			util.PrintColor(buf, util.Red, "---- STDOUT ----\n%s\n", event.Result.Stdout)
		}
		if event.Result.Stderr != "" {
			util.PrintColor(buf, util.Red, "---- STDERR ----\n%s\n", event.Result.Stderr)
		}
		if event.Result.Stderr != "" || event.Result.Stdout != "" {
			util.PrintColor(buf, util.Red, "---------------\n")
		}
	case *ansible.RunnerUnreachableEvent:
		// Host is unreachable
		// Print newline before first task
		if explainer.printPlayStatus {
			fmt.Fprintln(buf)
			// Dont print play success status on error
			explainer.printPlayStatus = false
		}
		util.PrettyPrintUnreachable(buf, "  %s", event.Host)
	case *ansible.TaskStartEvent:
		if verbose {
			// Print newline before first task status
			if explainer.printPlayStatus {
				fmt.Fprintln(buf)
				// Dont print play success status on error
				explainer.printPlayStatus = false
			}
			fmt.Fprintf(buf, "- Running task: %s\n", event.Name)
		}
		// Set current task name
		explainer.currentTask = event.Name
	case *ansible.HandlerTaskStartEvent:
		if verbose {
			// Print newline before first task
			if explainer.printPlayStatus {
				fmt.Fprintln(buf)
				// Dont print play success status on error
				explainer.printPlayStatus = false
			}
			fmt.Fprintf(buf, "- Running task: %s\n", event.Name)
		}
		// Set current task name
		explainer.currentTask = event.Name
	case *ansible.PlaybookEndEvent:
		// Playbook ends, print the last play status
		if explainer.printPlayStatus {
			if verbose {
				util.PrintColor(buf, util.Green, "%s Finished\n", explainer.lastPlay)
			} else {
				util.PrintOkln(buf)
			}
		}
	case *ansible.RunnerSkippedEvent:
		if verbose {
			util.PrettyPrintSkipped(buf, "  %s", event.Host)
		}
	case *ansible.RunnerOKEvent:
		if verbose {
			util.PrettyPrintOk(buf, "  %s", event.Host)
		}
	case *ansible.RunnerItemOKEvent:
		if verbose {
			util.PrettyPrintOk(buf, "  %s", event.Host)
		}
	case *ansible.RunnerItemRetryEvent:
		return ""
	case *ansible.PlaybookStartEvent:
		explainer.playCount = event.Count
		explainer.currentPlayCount = 1
		return ""
	default:
		if verbose {
			util.PrintColor(buf, util.Orange, "Unhandled event: %T\n", event)
		}
	}
	return buf.String()
}
