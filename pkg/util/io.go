package util

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// PromptForInt read command line input
func PromptForInt(in io.Reader, out io.Writer, prompt string, defaultValue int) (int, error) {
	fmt.Fprintf(out, "=> %s [%d]: ", prompt, defaultValue)
	s := bufio.NewScanner(in)
	// Scan the first token
	s.Scan()
	if s.Err() != nil {
		return defaultValue, fmt.Errorf("error reading number: %v", s.Err())
	}
	ans := s.Text()
	if ans == "" {
		return defaultValue, nil
	}
	// Convert input into integer
	i, err := strconv.Atoi(ans)
	if err != nil {
		return defaultValue, fmt.Errorf("%q is not a number", ans)
	}
	return i, nil
}
