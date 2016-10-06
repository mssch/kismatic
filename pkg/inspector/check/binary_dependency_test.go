package check

import "testing"

func TestBinaryDependencyExists(t *testing.T) {
	c := BinaryDependencyCheck{
		BinaryName: "ls",
	}
	ok, err := c.Check()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !ok {
		t.Errorf("Expected ls to be in the path")
	}
}

func TestBinaryDependency(t *testing.T) {
	c := BinaryDependencyCheck{
		BinaryName: "non-existent-binary",
	}
	ok, err := c.Check()
	if err != nil {
		t.Errorf("Unexpected error when running binary dependency check: %v", err)
	}
	if ok {
		t.Error("check returned OK for a binary that does not exist")
	}
}

func TestBinaryDependencyBadBinary(t *testing.T) {
	tests := []struct {
		binaryName string
	}{
		{binaryName: "echo; exit 0"},
		{binaryName: "1234"},
		{binaryName: "hello$?"},
		{binaryName: "!echo"},
	}
	for _, test := range tests {
		c := BinaryDependencyCheck{
			BinaryName: test.binaryName,
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
