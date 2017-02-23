package rule

import "testing"

func TestPackageDependencyRuleValidation(t *testing.T) {
	p := PackageDependency{}
	errs := p.Validate()
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, but got %d", len(errs))
	}
	p.PackageName = "foo"
	if errs := p.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	p.PackageVersion = "1.0"
	if errs := p.Validate(); len(errs) != 0 {
		t.Errorf("expected 0 error, but got %d", len(errs))
	}
}
