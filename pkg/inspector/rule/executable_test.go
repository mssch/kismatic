package rule

import "testing"

func TestExecutableInPathRuleValidation(t *testing.T) {
	e := ExecutableInPath{}
	if errs := e.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	e.Executable = "123"
	if errs := e.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	e.Executable = "foo"
	if errs := e.Validate(); len(errs) != 0 {
		t.Errorf("Expected 0 errors, but got %d", len(errs))
	}
}
