package explain

import "github.com/apprenda/kismatic/pkg/ansible"

// StreamExplainer wraps the Explain method, which reads the incoming stream, and explains to the user what is
// happening
type StreamExplainer interface {
	// Explain the incoming stream
	Explain(events <-chan ansible.Event) error
}
