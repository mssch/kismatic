package explain

import "github.com/apprenda/kismatic/pkg/ansible"

// AnsibleEventStreamExplainer explains the incoming ansible event stream
type AnsibleEventStreamExplainer struct {
	// EventExplainer for processing ansible events
	EventExplainer AnsibleEventExplainer
}

// Explain the incoming ansible event stream
func (e *AnsibleEventStreamExplainer) Explain(events <-chan ansible.Event) error {
	for event := range events {
		e.EventExplainer.ExplainEvent(event)
	}
	return nil
}

// AnsibleEventExplainer explains a single event
type AnsibleEventExplainer interface {
	ExplainEvent(e ansible.Event)
}
