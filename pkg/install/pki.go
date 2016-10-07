package install

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/apprenda/kismatic-platform/pkg/util"
	"github.com/cloudflare/cfssl/csr"
)

// The PKI provides a way for generating certificates for the cluster described by the Plan
type PKI interface {
	ReadClusterCA(p *Plan) (*tls.CA, error)
	GenerateClusterCA(p *Plan) (*tls.CA, error)
	GenerateClusterCerts(p *Plan, ca *tls.CA, users []string) error
}

// LocalPKI is a file-based PKI
type LocalPKI struct {
	CACsr                   string
	CAConfigFile            string
	CASigningProfile        string
	GeneratedCertsDirectory string
	Log                     io.Writer
}

// ReadClusterCA read a Certificate Authority from a file
func (lp *LocalPKI) ReadClusterCA(p *Plan) (*tls.CA, error) {
	key, cert, err := tls.ReadCACert("ca", lp.GeneratedCertsDirectory)
	if err != nil {
		return nil, err
	}
	ca := &tls.CA{
		Key:        key,
		Cert:       cert,
		ConfigFile: lp.CAConfigFile,
		Profile:    lp.CASigningProfile,
	}

	return ca, nil
}

// GenerateClusterCA creates a Certificate Authority for the cluster
func (lp *LocalPKI) GenerateClusterCA(p *Plan) (*tls.CA, error) {
	// First, generate a CA
	key, cert, err := tls.NewCACert(lp.CACsr)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA Cert: %v", err)
	}

	err = tls.WriteCert(key, cert, "ca", lp.GeneratedCertsDirectory)
	if err != nil {
		return nil, fmt.Errorf("error writing CA files: %v", err)
	}

	ca := &tls.CA{
		Key:        key,
		Cert:       cert,
		ConfigFile: lp.CAConfigFile,
		Profile:    lp.CASigningProfile,
	}

	return ca, nil
}

// GenerateClusterCerts creates a Certificates for all nodes on the cluster
func (lp *LocalPKI) GenerateClusterCerts(p *Plan, ca *tls.CA, users []string) error {
	if lp.Log == nil {
		lp.Log = ioutil.Discard
	}

	// Add kubernetes service IP to certificates
	kubeServiceIP, err := getKubernetesServiceIP(p)
	if err != nil {
		return fmt.Errorf("Error getting kubernetes service IP: %v", err)
	}

	defaultCertHosts := []string{
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		"kubernetes.default.svc.cluster.local",
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
		if util.ContainsString(seenNodes, n.Host) {
			continue
		}
		seenNodes = append(seenNodes, n.Host)
		util.PrettyPrintOk(lp.Log, "Generating certificates for host %q", n.Host)
		nodeList := append(defaultCertHosts, n.Host, n.IP, n.InternalIP)
		key, cert, err := generateCert(p.Cluster.Name, p, nodeList, ca)
		if err != nil {
			return fmt.Errorf("error during cluster cert generation: %v", err)
		}
		err = tls.WriteCert(key, cert, n.Host, lp.GeneratedCertsDirectory)
		if err != nil {
			return fmt.Errorf("error writing cert files for host %q: %v", n.Host, err)
		}
	}

	// Create certs for docker registry
	if p.DockerRegistry.UseInternal {
		util.PrettyPrintOk(lp.Log, "Generating certificates for docker registry")
		// Default registry will be deployed on the first master
		n := p.Master.Nodes[0]
		dockerKey, dockerCert, err := generateCert(n.Host, p, []string{n.Host, n.IP, n.InternalIP}, ca)
		if err != nil {
			return fmt.Errorf("error during user cert generation: %v", err)
		}
		err = tls.WriteCert(dockerKey, dockerCert, "docker", lp.GeneratedCertsDirectory)
		if err != nil {
			return fmt.Errorf("error writing cert files for docker registry")
		}
	}

	// Finally, create certs for user
	for _, user := range users {
		util.PrettyPrintOk(lp.Log, "Generating certificates for user %q", user)
		adminKey, adminCert, err := generateCert(user, p, []string{user}, ca)
		if err != nil {
			return fmt.Errorf("error during user cert generation: %v", err)
		}
		err = tls.WriteCert(adminKey, adminCert, user, lp.GeneratedCertsDirectory)
		if err != nil {
			return fmt.Errorf("error writing cert files for user %q: %v", user, err)
		}
	}

	return nil
}

func generateCert(cnName string, p *Plan, hostList []string, ca *tls.CA) (key, cert []byte, err error) {
	req := csr.CertificateRequest{
		CN: cnName,
		KeyRequest: &csr.BasicKeyRequest{
			A: "rsa",
			S: 2048,
		},
		Hosts: hostList,
		Names: []csr.Name{
			{
				C:  p.Cluster.Certificates.LocationCountry,
				ST: p.Cluster.Certificates.LocationState,
				L:  p.Cluster.Certificates.LocationCity,
			},
		},
	}

	key, cert, err = tls.NewCert(ca, req)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating certs for %q: %v", cnName, err)
	}

	return key, cert, err
}
