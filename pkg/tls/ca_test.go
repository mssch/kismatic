package tls

import (
	"reflect"
	"testing"
	"time"

	"github.com/cloudflare/cfssl/helpers"
)

func TestNewCACert(t *testing.T) {
	subject := Subject{
		Country:  "someCountry",
		State:    "someState",
		Locality: "someLocality",
	}
	_, cert, err := NewCACert("test/ca-csr.json", subject)
	if err != nil {
		t.Fatalf("error creating CA cert: %v", err)
	}

	parsedCert, err := helpers.ParseCertificatePEM(cert)
	if err != nil {
		t.Fatalf("error parsing certificate: %v", err)
	}

	if !parsedCert.IsCA {
		t.Errorf("Genereated CA cert is not CA")
	}

	if !reflect.DeepEqual(parsedCert.Issuer, parsedCert.Subject) {
		t.Errorf("cert issuer is not equal to the CA's subject")
	}

	// You might be tempted to test for this, but it seems like the AuthKeyID doesn't have to be set
	// for self-signed certificates. https://go.googlesource.com/go/+/b623b71509b2d24df915d5bc68602e1c6edf38ca
	// if !bytes.Equal(parsedCert.AuthorityKeyId, parsedCert.SubjectKeyId) {
	// 	t.Errorf("certificate auth key ID %q is not the subject key ID of the CA %q", string(parsedCert.AuthorityKeyId), string(parsedCert.SubjectKeyId))
	// }

	// Verify expiration
	now := time.Now().UTC()
	d, err := time.ParseDuration("8760h")
	if err != nil {
		t.Fatalf("error parsing duration: %v", err)
	}
	expectedExpiration := now.Add(d)
	if expectedExpiration.Year() != parsedCert.NotAfter.Year() || expectedExpiration.YearDay() != parsedCert.NotAfter.YearDay() {
		t.Errorf("expected expiration date %q, got %q", expectedExpiration, parsedCert.NotAfter)
	}
}
