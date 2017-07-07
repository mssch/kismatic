package util

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestLineReaderReadEmptyString(t *testing.T) {
	lr := NewLineReader(bytes.NewBufferString(""), 64)
	line, err := lr.Read()
	if err != io.EOF {
		t.Errorf("unexpected error: %v", err)
	}
	if string(line) != "" {
		t.Errorf("unexpected line: %v", err)
	}
}

func TestLineReaderReadLine(t *testing.T) {
	expectedLine := "foo"
	r := bytes.NewBufferString(expectedLine)
	lr := NewLineReader(r, 64)
	line, err := lr.Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(line) != expectedLine {
		t.Errorf("got %s, but wanted %s", line, expectedLine)
	}
}

func TestLineReaderReadLongLine(t *testing.T) {
	expectedLine := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum sed tellus sit amet ligula cursus pretium feugiat ut lorem. In eu vestibulum turpis. Sed et mauris ut massa finibus pharetra. Suspendisse blandit eu est vitae congue. Praesent elementum blandit sapien et convallis. Sed quam nibh, vestibulum eu varius quis, consectetur in purus. Vestibulum ac justo in leo tempus pellentesque et in metus. Sed in est sed magna consequat hendrerit non et lectus. Sed sodales dictum tortor vel gravida. Phasellus egestas metus nec massa imperdiet, non gravida ligula cursus. Ut non diam sit amet ante porttitor tristique auctor vel leo. Morbi lectus tortor, scelerisque a laoreet id, aliquet sit amet nibh. Sed ullamcorper lectus dictum tellus ultrices iaculis. Integer rutrum, est eu mattis laoreet, eros tellus interdum tortor, a pharetra eros metus et lectus.`
	r := bytes.NewBufferString(expectedLine)
	lr := NewLineReader(r, 64)
	line, err := lr.Read()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(line) != expectedLine {
		t.Errorf("got %s, but wanted %s", line, expectedLine)
	}
}

func TestLineReaderReadLineEOFWhenDone(t *testing.T) {
	lr := NewLineReader(bytes.NewBufferString("foo"), 64)
	l, err := lr.Read()
	if string(l) != "foo" {
		t.Errorf("unexpected line: %v", l)
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	_, err = lr.Read()
	if err != io.EOF {
		t.Errorf("got %v but wanted %v", err, io.EOF)
	}
}

func TestLineReaderReadLineMultipleLines(t *testing.T) {
	lines := `foo
bar
baz`

	lr := NewLineReader(bytes.NewBufferString(lines), 64)
	var (
		gotLines []string
		line     []byte
		err      error
	)
	for {
		line, err = lr.Read()
		if err != nil {
			break
		}
		gotLines = append(gotLines, string(line))
	}
	if err != io.EOF {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedLines := []string{"foo", "bar", "baz"}
	if !reflect.DeepEqual(gotLines, expectedLines) {
		t.Log(len(gotLines))
		t.Errorf("got %v, but expected %v", gotLines, expectedLines)
	}
}
