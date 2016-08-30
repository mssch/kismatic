package tls

import (
	"reflect"
	"testing"
	"time"

	"github.com/cloudflare/cfssl/helpers"
)

func TestNewCACert(t *testing.T) {
	_, cert, err := NewCACert("test/ca-csr.json")
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

	if !reflect.DeepEqual(parsedCert.AuthorityKeyId, parsedCert.SubjectKeyId) {
		t.Errorf("certificate auth key ID is not the subject key ID of the CA")
	}

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
