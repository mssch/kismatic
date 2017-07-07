package util

import (
	"bufio"
	"io"
)

// LineReader provides a way to read lines of text with arbitrary length.
// The line is buffered in memory.
type LineReader struct {
	reader *bufio.Reader
}

// NewLineReader returns a LineReader to read from r. It wraps the reader
// with a buffered reader that has a buffer of size bufSizeBytes.
func NewLineReader(r io.Reader, bufSizeBytes int) LineReader {
	return LineReader{bufio.NewReaderSize(r, bufSizeBytes)}
}

// Read an entire line from the reader.
func (lr LineReader) Read() ([]byte, error) {
	var (
		isPrefix       bool  = true
		err            error = nil
		line, fullLine []byte
	)
	for isPrefix && err == nil { // read until we have the whole line, or we get an err
		line, isPrefix, err = lr.reader.ReadLine()
		fullLine = append(fullLine, line...)
	}
	return fullLine, err
}
