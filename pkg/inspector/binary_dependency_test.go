package inspector

import "testing"

func TestBinaryDependencyExists(t *testing.T) {
	c := BinaryDependencyCheck{
		BinaryName: "ls",
	}
	if err := c.Check(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestBinaryDependency(t *testing.T) {
	c := BinaryDependencyCheck{
		BinaryName: "non-existent-binary",
	}
	if err := c.Check(); err == nil {
		t.Errorf("Expected an error, but didn't get one")
	}
}
