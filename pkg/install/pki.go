package install

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/cloudflare/cfssl/csr"
)

// The PKI provides a way for generating certificates for the cluster described by the Plan
type PKI interface {
	GenerateClusterCerts(p *Plan) error
}

// LocalPKI is a file-based PKI
type LocalPKI struct {
	CACsr            string
	CAConfigFile     string
	CASigningProfile string
	DestinationDir   string
	Log              io.Writer
}

// GenerateClusterCerts creates a Certificate Authority and Certificates
// for all nodes on the cluster.
func (lp *LocalPKI) GenerateClusterCerts(p *Plan) error {
	if lp.Log == nil {
		lp.Log = ioutil.Discard
	}
	// First, generate a CA
	key, cert, err := tls.NewCACert(lp.CACsr)
	if err != nil {
		return fmt.Errorf("failed to create CA Cert: %v", err)
	}

	err = lp.writeFiles(key, cert, "ca")
	if err != nil {
		return fmt.Errorf("error writing CA files: %v", err)
	}

	ca := &tls.CA{
		Key:        key,
		Cert:       cert,
		ConfigFile: lp.CAConfigFile,
		Profile:    lp.CASigningProfile,
	}

	// Add kubernetes service IP to certificates
	kubeServiceIP, err := getKubernetesServiceIP(p)
	if err != nil {
		return fmt.Errorf("Error getting kubernetes service IP: %v", err)
	}

	defaultCertHosts := []string{
		p.Cluster.Name,
		p.Cluster.Name + ".default",
		p.Cluster.Name + ".default.svc",
		p.Cluster.Name + ".default.svc.cluster.local",
		"127.0.0.1",
		kubeServiceIP,
	}

	// Then, create certs for all nodes
	nodes := []Node{}
	nodes = append(nodes, p.Etcd.Nodes...)
	nodes = append(nodes, p.Master.Nodes...)
	nodes = append(nodes, p.Worker.Nodes...)

	seenNodes := []string{}
	for _, n := range nodes {
		// Only generate certs once for each node, nodes can be in more than one group
		if contains(seenNodes, n.Host) {
			continue
		}
		seenNodes = append(seenNodes, n.Host)
		fmt.Fprintf(lp.Log, "Generating certificates for %q\n", n.Host)
		key, cert, err := generateNodeCert(p, &n, ca, defaultCertHosts)
		if err != nil {
			return fmt.Errorf("error during cluster cert generation: %v", err)
		}
		err = lp.writeFiles(key, cert, n.Host)
		if err != nil {
			return fmt.Errorf("error writing cert files for host %q: %v", n.Host, err)
		}
	}
	// Finally, create cert for user `admin`
	adminUser := "admin"
	fmt.Fprintf(lp.Log, "Generating certificates for user %q\n", adminUser)
	adminKey, adminCert, err := generateClientCert(p, adminUser, ca)
	if err != nil {
		return fmt.Errorf("error during admin cert generation: %v", err)
	}
	err = lp.writeFiles(adminKey, adminCert, adminUser)
	if err != nil {
		return fmt.Errorf("error writing cert files for user %q: %v", adminUser, err)
	}

	return nil
}

func (lp *LocalPKI) writeFiles(key, cert []byte, name string) error {
	// Create destination dir if it doesn't exist
	if _, err := os.Stat(lp.DestinationDir); os.IsNotExist(err) {
		err := os.Mkdir(lp.DestinationDir, 0744)
		if err != nil {
			return fmt.Errorf("error creating destination dir: %v", err)
		}
	}

	// Write private key with read-only for user
	keyName := fmt.Sprintf("%s-key.pem", name)
	dest := filepath.Join(lp.DestinationDir, keyName)
	err := ioutil.WriteFile(dest, key, 0600)
	if err != nil {
		return fmt.Errorf("error writing private key: %v", err)
	}

	// Write cert
	certName := fmt.Sprintf("%s.pem", name)
	dest = filepath.Join(lp.DestinationDir, certName)
	err = ioutil.WriteFile(dest, cert, 0644)
	if err != nil {
		return fmt.Errorf("error writing certificate: %v", err)
	}
	return nil
}

func generateNodeCert(p *Plan, n *Node, ca *tls.CA, initialHostList []string) (key, cert []byte, err error) {
	hosts := append(initialHostList, n.Host, n.InternalIP, n.IP)
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

func generateClientCert(p *Plan, user string, ca *tls.CA) (key, cert []byte, err error) {
	req := csr.CertificateRequest{
		CN: user,
		KeyRequest: &csr.BasicKeyRequest{
			A: "rsa",
			S: 2048,
		},
		Hosts: []string{},
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
		return nil, nil, fmt.Errorf("error generating certs for user %q: %v", user, err)
	}

	return key, cert, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
