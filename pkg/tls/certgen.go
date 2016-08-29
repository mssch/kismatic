package tls

import (
	"fmt"

	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
)

// CA contains information about the Certificate Authority
type CA struct {
	// Key is the CA's private key.
	Key []byte
	// Password is the CA's private key password. Can be empty if not password is set.
	Password string
	// Cert is the CA's public certificate.
	Cert []byte
	// ConfigFile contains a cfssl configuration file for the Certificate Authority
	ConfigFile string
	// Profile to be used when signing with this Certificate Authority
	Profile string
}

// GenerateNewCertificate creates a new certificate/key pair using the CertificateAuthority
// provided.
func GenerateNewCertificate(ca CA, req csr.CertificateRequest) (key, cert []byte, err error) {
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err := g.ProcessRequest(&req)
	if err != nil {
		return nil, nil, fmt.Errorf("error processing CSR: %v", err)
	}

	// Get CA private key
	caPriv, err := helpers.ParsePrivateKeyPEMWithPassword(ca.Key, []byte(ca.Password))
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing privte key: %v", err)
	}

	// Parse CA Cert
	caCert, err := helpers.ParseCertificatePEM(ca.Cert)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing CA cert: %v", err)
	}
	sigAlgo := signer.DefaultSigAlgo(caPriv)

	// Get CA config from file
	caConfig, err := config.LoadFile(ca.ConfigFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading CA Config: %v", err)
	}

	// Create signer using CA
	s, err := local.NewSigner(caPriv, caCert, sigAlgo, caConfig.Signing)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating signer: %v", err)
	}

	// Generate cert using CA signer
	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Profile: ca.Profile,
	}
	cert, err = s.Sign(signReq)
	if err != nil {
		return nil, nil, fmt.Errorf("error signing certificate: %v", err)
	}

	return key, cert, nil
}
