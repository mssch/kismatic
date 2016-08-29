package install

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/cloudflare/cfssl/csr"
)

var defaultCertHosts = []string{
	"kubernetes",
	"kubernetes.default",
	"kubernetes.default.svc",
	"kubernetes.default.svc.cluster.local",
	"10.3.0.1",
	"10.3.0.5",
	"10.3.0.10",
	"127.0.0.1",
}

// The PKI provides a way for generating certificates for the cluster described by the Plan
type PKI interface {
	GenerateClusterCerts(p *Plan) error
}

// LocalPKI is a file-based PKI
type LocalPKI struct {
	CACsr            string
	CAConfigFile     string
	CASigningProfile string
}

// GenerateClusterCerts creates a Certificate Authority and Certificates
// for all nodes on the cluster.
func (lp *LocalPKI) GenerateClusterCerts(p *Plan) error {
	// First, generate a CA
	key, cert, err := tls.NewCACert(lp.CACsr)
	if err != nil {
		return fmt.Errorf("failed to create CA Cert: %v", err)
	}

	err = writeFiles(key, cert, "ca")
	if err != nil {
		return fmt.Errorf("error writing CA files: %v", err)
	}

	ca := &tls.CA{
		Key:        key,
		Cert:       cert,
		ConfigFile: lp.CAConfigFile,
		Profile:    lp.CASigningProfile,
	}

	// Then, create certs for all nodes
	nodes := []Node{}
	nodes = append(nodes, p.Etcd.Nodes...)
	nodes = append(nodes, p.Master.Nodes...)
	nodes = append(nodes, p.Worker.Nodes...)

	for _, n := range nodes {
		key, cert, err := generateNodeCert(p, &n, ca)
		if err != nil {
			return fmt.Errorf("error during cluster cert generation: %v", err)
		}
		err = writeFiles(key, cert, n.Host)
		if err != nil {
			return fmt.Errorf("error writing cert files for host %q: %v", n.Host, err)
		}
	}
	return nil
}

func writeFiles(key, cert []byte, name string) error {
	// Write into ansible directory for now...
	destDir := filepath.Join("ansible", "playbooks", "tls")
	keyName := fmt.Sprintf("%s-key.pem", name)
	dest := filepath.Join(destDir, keyName)
	err := ioutil.WriteFile(dest, key, 0600)
	if err != nil {
		return fmt.Errorf("error writing private key: %v", err)
	}
	certName := fmt.Sprintf("%s.pem", name)
	dest = filepath.Join(destDir, certName)
	err = ioutil.WriteFile(dest, cert, 0644)
	if err != nil {
		return fmt.Errorf("error writing certificate: %v", err)
	}
	return nil
}

func generateNodeCert(p *Plan, n *Node, ca *tls.CA) (key, cert []byte, err error) {
	hosts := append(defaultCertHosts, n.Host, n.InternalIP, n.IP)
	req := csr.CertificateRequest{
		CN: p.Cluster.Name,
		KeyRequest: &csr.BasicKeyRequest{
			A: "rsa",
			S: 2048,
		},
		Hosts: hosts,
		Names: []csr.Name{
			{
				C:  p.Cluster.Certificates.LocationCountry,
				ST: p.Cluster.Certificates.LocationState,
				L:  p.Cluster.Certificates.LocationCity,
			},
		},
	}

	key, cert, err = tls.GenerateNewCertificate(ca, req)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating certs for node %q: %v", n.Host, err)
	}

	return key, cert, err
}
