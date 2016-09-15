package explain

import "io"

// StreamExplainer wraps the Explain method, which reads the incoming stream, and explains to the user what is
// happening
type StreamExplainer interface {
	// Explain the incoming stream
	Explain(in io.Reader) error
}
