package ansible

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// EventStream reads JSON lines from the incoming stream, and convert them
// into a stream of events.
func EventStream(in io.Reader) <-chan Event {
	scanner := bufio.NewScanner(in)
	out := make(chan Event)
	go func() {
		for scanner.Scan() {
			jl := scanner.Bytes()
			event, err := eventFromJSONLine(jl)
			if err != nil {
				// handle this error? Maybe have an outErr channel
				continue
			}
			out <- event
		}
		// Close the channel, as the stream is done
		close(out)
	}()
	return out
}

// eventEnvelope contains event data for a specific event type
type eventEnvelope struct {
	Type string      `json:"eventType"`
	Data interface{} `json:"eventData"`
}

func eventFromJSONLine(line []byte) (Event, error) {
	// Unmarshal the event type, but defer unmarshaling the data
	var data json.RawMessage
	env := &eventEnvelope{
		Data: &data,
	}
	if err := json.Unmarshal(line, env); err != nil {
		return nil, fmt.Errorf("error parsing event: %v\nline was:\n%s\n", err, string(line))
	}

	// Unmarshal the data according to the event type
	switch env.Type {
	case "PLAYBOOK_START":
		e := &PlaybookStartEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "PLAYBOOK_END":
		e := &PlaybookEndEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "PLAY_START":
		e := &PlayStartEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "TASK_START":
		e := &TaskStartEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "HANDLER_TASK_START":
		e := &HandlerTaskStartEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "RUNNER_OK":
		e := &RunnerOKEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "RUNNER_ITEM_OK":
		e := &RunnerItemOKEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "RUNNER_ITEM_RETRY":
		e := &RunnerItemRetryEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "RUNNER_FAILED":
		e := &RunnerFailedEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "RUNNER_SKIPPED":
		e := &RunnerSkippedEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	case "RUNNER_UNREACHABLE":
		e := &RunnerUnreachableEvent{}
		if err := json.Unmarshal(data, e); err != nil {
			return nil, fmt.Errorf("error reading event data: %v\nline was:\n%s\n", err, string(line))
		}
		return e, nil
	default:
		return nil, fmt.Errorf("unhandled ansible event type %q", env.Type)
	}
}
