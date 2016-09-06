package tls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/log"
)

func init() {
	log.Level = log.LevelError
}

// NewCACert creates a new Certificate Authority and returns it's private key and public certificate.
func NewCACert(csrFile string) (key, cert []byte, err error) {
	// Open CSR file
	f, err := os.Open(csrFile)
	if os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("%q does not exist", csrFile)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("error opening %q", csrFile)
	}
	// Create CSR struct
	csr := &csr.CertificateRequest{
		KeyRequest: csr.NewBasicKeyRequest(),
	}
	err = json.NewDecoder(f).Decode(csr)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding CSR: %v", err)
	}
	// Generate CA Cert according to CSR
	cert, _, key, err = initca.New(csr)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating CA cert: %v", err)
	}

	return key, cert, nil
}

// ReadCACert read CA file
func ReadCACert(name, dir string) (key, cert []byte, err error) {
	keyName := fmt.Sprintf("%s-key.pem", name)
	dest := filepath.Join(dir, keyName)
	key, errKey := ioutil.ReadFile(dest)
	if errKey != nil {
		return nil, nil, fmt.Errorf("error reading private key: %v", errKey)
	}

	certName := fmt.Sprintf("%s.pem", name)
	dest = filepath.Join(dir, certName)
	cert, errCert := ioutil.ReadFile(dest)
	if errCert != nil {
		return nil, nil, fmt.Errorf("error reading certificate: %v", errKey)
	}

	return key, cert, nil
}
