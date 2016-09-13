package explain

import "io"

// RawExplainer dumps out the stream to the user, without performing any interpretation
type RawExplainer struct {
	Out io.Writer
}

// Explain copies the incoming stream to the out stream
func (e *RawExplainer) Explain(in io.Reader) error {
	_, err := io.Copy(e.Out, in)
	return err
}
