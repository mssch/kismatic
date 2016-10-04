package rule

import "testing"

func TestPacakgeAvailableRuleValidation(t *testing.T) {
	p := PackageAvailable{}
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

func TestPackageInstalledRuleValidation(t *testing.T) {
	p := PackageInstalled{}
	if errs := p.Validate(); len(errs) != 2 {
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

func TestExecutableInPathRuleValidation(t *testing.T) {
	e := ExecutableInPath{}
	if errs := e.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	e.Executable = "foo"
	if errs := e.Validate(); len(errs) != 0 {
		t.Errorf("Expected 0 errors, but got %d", len(errs))
	}
}

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

func TestTCPPortAvailableRuleValidation(t *testing.T) {
	p := TCPPortAvailable{}
	if errs := p.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	p.Port = -1
	if errs := p.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	p.Port = 901283
	if errs := p.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	p.Port = 1024
	if errs := p.Validate(); len(errs) != 0 {
		t.Errorf("expected 0 error, but got %d", len(errs))
	}
}

func TestTCPPortAccessibleRuleValidation(t *testing.T) {
	p := TCPPortAccessible{}
	if errs := p.Validate(); len(errs) != 2 {
		t.Errorf("expected 2 error, but got %d", len(errs))
	}
	p.Port = -1
	if errs := p.Validate(); len(errs) != 2 {
		t.Errorf("expected 2 error, but got %d", len(errs))
	}
	p.Port = 901283
	if errs := p.Validate(); len(errs) != 2 {
		t.Errorf("expected 2 error, but got %d", len(errs))
	}
	p.Port = 1024
	if errs := p.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	p.Timeout = "nonDuration"
	if errs := p.Validate(); len(errs) != 1 {
		t.Errorf("expected 1 error, but got %d", len(errs))
	}
	p.Timeout = "3s"
	if errs := p.Validate(); len(errs) != 0 {
		t.Errorf("expected 0 error, but got %d", len(errs))
	}
}
