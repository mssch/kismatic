package check

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// FileContentCheck runs a search against the contents of the specified file.
// The SearchString is a regular expression that is in accordance with the
// RE2 syntax defined by the Go regexp package.
type FileContentCheck struct {
	File         string
	SearchString string
}

// Check returns true if file contents match the regular expression. Otherwise,
// returns false. If an error occurrs, returns false and the error.
func (c FileContentCheck) Check() (bool, error) {
	if _, err := os.Stat(c.File); os.IsNotExist(err) {
		return false, fmt.Errorf("Attempted to validate file %q, but it doesn't exist.", c.File)
	}
	r, err := regexp.Compile(c.SearchString)
	if err != nil {
		return false, fmt.Errorf("Invalid search string provided %q: %v", c.SearchString, err)
	}
	b, err := ioutil.ReadFile(c.File)
	if err != nil {
		return false, fmt.Errorf("Error reading file %q: %v", c.File, err)
	}
	return r.Match(b), nil
}
