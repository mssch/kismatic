package check

import "testing"

func TestExecutableInPathExists(t *testing.T) {
	c := ExecutableInPathCheck{
		Name: "ls",
	}
	ok, err := c.Check()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !ok {
		t.Errorf("Expected ls to be in the path")
	}
}

func TestExecutableMissingFromPath(t *testing.T) {
	c := ExecutableInPathCheck{
		Name: "non-existent-binary",
	}
	ok, err := c.Check()
	if err != nil {
		t.Errorf("Unexpected error when running binary dependency check: %v", err)
	}
	if ok {
		t.Error("check returned OK for a binary that does not exist")
	}
}

func TestExecutableInPathBadName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "echo; exit 0"},
		{name: "1234"},
		{name: "hello$?"},
		{name: "!echo"},
	}
	for _, test := range tests {
		c := ExecutableInPathCheck{
			Name: test.name,
		}
		ok, err := c.Check()
		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
		if ok {
			t.Errorf("check returned OK for an invalid binary name")
		}
	}

}
