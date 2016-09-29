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
