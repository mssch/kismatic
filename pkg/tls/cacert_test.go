package tls

import (
	"fmt"
	"testing"
)

func TestNewCACert(t *testing.T) {
	key, cert, err := NewCACert("test/ca-csr.json")
	if err != nil {
		t.Fatalf("error creating CA cert: %v", err)
	}
	fmt.Println(string(key))
	fmt.Println(string(cert))
}
