package rule

import "testing"

func TestFileContentMatchesRuleValidation(t *testing.T) {
	f := FileContentMatches{}
	if errs := f.Validate(); len(errs) != 2 {
		t.Errorf("expected 2 error, but got %d", len(errs))
	}
	f.File = "foo"
	if errs := f.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	f.ContentRegex = "\\i"
	if errs := f.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	f.ContentRegex = "foo"
	if errs := f.Validate(); len(errs) != 0 {
		t.Errorf("expected 0 error, but got %d", len(errs))
	}
}
