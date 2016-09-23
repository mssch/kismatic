package preflight

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

// Name of the check
func (c FileContentCheck) Name() string {
	return fmt.Sprintf("Contents of %q match %q", c.File, c.SearchString)
}

// Check returns nil if the search string was found in the file.
func (c FileContentCheck) Check() error {
	if _, err := os.Stat(c.File); os.IsNotExist(err) {
		return fmt.Errorf("Attempted to validate file %q, but it doesn't exist.", c.File)
	}
	r, err := regexp.Compile(c.SearchString)
	if err != nil {
		return fmt.Errorf("Invalid search string provided %q: %v", c.SearchString, err)
	}
	b, err := ioutil.ReadFile(c.File)
	if err != nil {
		return fmt.Errorf("Attempted to validate %q, but got an error: %v", c.File, err)
	}
	if !r.Match(b) {
		return fmt.Errorf("Searched %q with the expression %q, but no matches were found.", c.File, c.SearchString)
	}
	return nil
}
