package rule

import (
	"errors"
	"fmt"
	"regexp"
)

// FileContentMatches is a rule that verifies that the contents of a file
// match the regular expression provided
type FileContentMatches struct {
	Meta
	File         string
	ContentRegex string
}

// Name is the name of the rule
func (f FileContentMatches) Name() string {
	return fmt.Sprintf("Contents of %q match: %s", f.File, f.ContentRegex)
}

// IsRemoteRule returns true if the rule is to be run from outside of the node
func (f FileContentMatches) IsRemoteRule() bool { return false }

// Validate the rule
func (f FileContentMatches) Validate() []error {
	errs := []error{}
	if f.File == "" {
		errs = append(errs, errors.New("File cannot be empty"))
	}
	if f.ContentRegex == "" {
		errs = append(errs, errors.New("ContentRegex cannot be empty"))
	}
	if f.ContentRegex != "" {
		if _, err := regexp.Compile(f.ContentRegex); err != nil {
			errs = append(errs, fmt.Errorf("ContentRegex contains an invalid regular expression: %v", err))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
