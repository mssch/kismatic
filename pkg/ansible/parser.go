package ansible

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// OutputParser reads JSON lines from ansible's stdout, and converts each line into an Event.
// Ansible must be configured to use the JSON lines stdout callback plugin.
type OutputParser struct {
	// Out is the write destination of the parsed events
	Out io.Writer
}

// eventEnvelope contains event data for a specific event type
type eventEnvelope struct {
	Type string      `json:"eventType"`
	Data interface{} `json:"eventData"`
}

type namedEvent struct {
	Name string
}

type runnerResult struct {
	Command []string `json:"cmd"`
	Stdout  string
	Stderr  string
}

type runnerResultEvent struct {
	Host   string
	Result runnerResult
}

// Transform parses the incoming stream, and writes the parsed events into it's
// output destination.
func (p *OutputParser) Transform(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// Read the event type, but defer unmarshaling of data
		var data json.RawMessage
		e := &eventEnvelope{
			Data: &data,
		}

		jsonLine := scanner.Bytes()
		if err := json.Unmarshal(jsonLine, e); err != nil {
			fmt.Fprintf(p.Out, "Error reading output from ansible: %v\n line was: \n%s\n", err, string(jsonLine))
			continue
		}

		// Unmarshal data according to type
		switch e.Type {
		case "PLAYBOOK_START":
			ne := &namedEvent{}
			if err := json.Unmarshal(data, ne); err != nil {
				fmt.Fprintf(p.Out, "Error reading playbook start event %v", err)
				continue
			}
			fmt.Fprintf(p.Out, "Starting playbook %s\n", ne.Name)
		case "PLAY_START":
			ne := &namedEvent{}
			if err := json.Unmarshal(data, ne); err != nil {
				fmt.Fprintf(p.Out, "error reading play start event: %v", err)
				continue
			}
			fmt.Fprintf(p.Out, "- %s\n", ne.Name)
		case "TASK_START", "HANDLER_TASK_START":
			ne := &namedEvent{}
			if err := json.Unmarshal(data, ne); err != nil {
				fmt.Fprintf(p.Out, "Error reading task start event: %v", err)
				continue
			}
			// fmt.Fprintf(p.Out, "    * %s\n", ne.Name)
		case "RUNNER_OK", "RUNNER_ITEM_OK":
			re := &runnerResultEvent{}
			if err := json.Unmarshal(data, re); err != nil {
				fmt.Fprintf(p.Out, "Error reading task start event: %v", err)
				continue
			}
			// fmt.Fprintf(p.Out, "    [OK] %s\n", re.Host)
		case "RUNNER_ITEM_RETRY":
			re := &runnerResultEvent{}
			if err := json.Unmarshal(data, re); err != nil {
				fmt.Fprintf(p.Out, "Error reading task start event: %v", err)
				continue
			}
			fmt.Fprintf(p.Out, "    [RETRYING] %s\n", re.Host)
		case "RUNNER_FAILED":
			re := &runnerResultEvent{}
			if err := json.Unmarshal(data, re); err != nil {
				fmt.Fprintf(p.Out, "Error reading task start event: %v", err)
				continue
			}
			fmt.Fprintf(p.Out, "    [FAILED] %s\n", re.Host)
			fmt.Fprintf(p.Out, "    |- stdout: %s\n", re.Result.Stdout)
			fmt.Fprintf(p.Out, "    |- stderr: %s\n", re.Result.Stderr)
		case "RUNNER_SKIPPED":
			re := &runnerResultEvent{}
			if err := json.Unmarshal(data, re); err != nil {
				fmt.Fprintf(p.Out, "Error reading task start event: %v", err)
				continue
			}
			// fmt.Fprintf(p.Out, "    [SKIPPED] %s\n", e.Host)
		case "RUNNER_UNREACHABLE":
			re := &runnerResultEvent{}
			if err := json.Unmarshal(data, e); err != nil {
				fmt.Fprintf(p.Out, "Error reading task start event: %v", err)
				continue
			}
			fmt.Fprintf(p.Out, "    [ERROR] Node %q is unreachable\n", re.Host)
		default:
			fmt.Fprintf(p.Out, "Unhandled ansible event type %q\n", e.Type)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from pipe\n")
	}
}
