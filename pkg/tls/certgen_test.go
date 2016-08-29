package tls

import (
	"fmt"
	"testing"

	"github.com/cloudflare/cfssl/csr"
)

func TestGenerateNewCertificate(t *testing.T) {
	// Create CA Cert
	key, cert, err := NewCACert("test/ca-csr.json")
	if err != nil {
		t.Fatalf("error creating CA: %v", err)
	}

	ca := &CA{
		Key:        key,
		Cert:       cert,
		ConfigFile: "test/ca-config.json",
		Profile:    "kubernetes",
	}

	req := csr.CertificateRequest{
		CN: "testKube",
		KeyRequest: &csr.BasicKeyRequest{
			A: "rsa",
			S: 2048,
		},
		Hosts: []string{"Alex"},
		Names: []csr.Name{
			{
				C:  "US",
				L:  "Troy",
				O:  "Kubernetes",
				OU: "Cluster",
				ST: "New York",
			},
		},
	}

	if err != nil {
		t.Fatalf("error decoding csr: %v", err)
	}

	cert, key, err = GenerateNewCertificate(ca, req)
	if err != nil {
		t.Errorf("error creating certificate: %v", err)
	}

	fmt.Println(string(cert))
	fmt.Println(string(key))

}
