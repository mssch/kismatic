package install

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/apprenda/kismatic-platform/pkg/util"
	"github.com/cloudflare/cfssl/csr"
)

const (
	certOrganization = "Apprenda"
	certOrgUnit      = "Kismatic"
	certCountry      = "US"
	certState        = "NY"
	certLocality     = "Troy"
)

// The PKI provides a way for generating certificates for the cluster described by the Plan
type PKI interface {
	GenerateClusterCA(p *Plan) (*tls.CA, error)
	GenerateClusterCertificates(p *Plan, ca *tls.CA, users []string) error
}

// LocalPKI is a file-based PKI
type LocalPKI struct {
	CACsr                   string
	CAConfigFile            string
	CASigningProfile        string
	GeneratedCertsDirectory string
	Log                     io.Writer
}

// GenerateClusterCA creates a Certificate Authority for the cluster
func (lp *LocalPKI) GenerateClusterCA(p *Plan) (*tls.CA, error) {
	ca := &tls.CA{
		ConfigFile: lp.CAConfigFile,
		Profile:    lp.CASigningProfile,
	}
	exists, err := tls.CertKeyPairExists("ca", lp.GeneratedCertsDirectory)
	if err != nil {
		return nil, fmt.Errorf("error verifying CA certificate/key: %v", err)
	}
	if exists {
		key, cert, err := tls.ReadCACert("ca", lp.GeneratedCertsDirectory)
		if err != nil {
			return nil, fmt.Errorf("error reading CA certificate/key: %v", err)
		}
		util.PrettyPrintOk(lp.Log, "Found a cluster Certificate Authority")
		ca.Cert = cert
		ca.Key = key
		return ca, nil
	}

	util.PrettyPrintOk(lp.Log, "Generating cluster Certificate Authority")
	// It doesn't exist, generate one
	caSubject := tls.Subject{
		Organization:       certOrganization,
		OrganizationalUnit: certOrgUnit,
		Country:            certCountry,
		State:              certState,
		Locality:           certLocality,
	}
	key, cert, err := tls.NewCACert(lp.CACsr, p.Cluster.Name, caSubject)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA Cert: %v", err)
	}
	if err = tls.WriteCert(key, cert, "ca", lp.GeneratedCertsDirectory); err != nil {
		return nil, fmt.Errorf("error writing CA files: %v", err)
	}
	ca.Cert = cert
	ca.Key = key
	return ca, nil
}

// GenerateClusterCertificates creates a Certificates for all nodes on the cluster
func (lp *LocalPKI) GenerateClusterCertificates(p *Plan, ca *tls.CA, users []string) error {
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

	// Create certs for master nodes.. they include the load balanced names
	seenNodes := map[string]bool{}
	for _, n := range p.Master.Nodes {
		if _, ok := seenNodes[n.Host]; ok {
			continue
		}
		seenNodes[n.Host] = true
		names := []string{}
		names = append(names, defaultCertHosts...)
		if p.Master.LoadBalancedFQDN != "" {
			names = append(names, p.Master.LoadBalancedFQDN)
		}
		if p.Master.LoadBalancedShortName != "" {
			names = append(names, p.Master.LoadBalancedShortName)
		}
		if err := lp.generateNodeCert(p, n, ca, names); err != nil {
			return err
		}
	}

	// Then, create certs for rest of nodes
	nodes := []Node{}
	nodes = append(nodes, p.Etcd.Nodes...)
	nodes = append(nodes, p.Worker.Nodes...)

	for _, n := range nodes {
		// Only generate certs once for each node, nodes can be in more than one group
		if _, ok := seenNodes[n.Host]; ok {
			continue
		}
		seenNodes[n.Host] = true
		if err := lp.generateNodeCert(p, n, ca, defaultCertHosts); err != nil {
			return err
		}
	}
	// Create certs for docker registry if it's missing
	if p.DockerRegistry.UseInternal {
		if err := lp.generateDockerRegistryCert(p, ca); err != nil {
			return err
		}
	}
	// Finally, create certs for user if they are missing
	for _, user := range users {
		if err := lp.generateUserCert(p, user, ca); err != nil {
			return err
		}
	}
	return nil
}

func (lp *LocalPKI) generateNodeCert(plan *Plan, node Node, ca *tls.CA, defaultCertHosts []string) error {
	// Don't generate if the key pair is already there
	exist, err := tls.CertKeyPairExists(node.Host, lp.GeneratedCertsDirectory)
	if err != nil {
		return fmt.Errorf("error verifying certificates for node %q: %v", node.Host, err)
	}
	if exist {
		util.PrettyPrintOk(lp.Log, "Found key and certificate for node %q", node.Host)
		return nil
	}

	util.PrettyPrintOk(lp.Log, "Generating certificates for host %q", node.Host)
	nodeList := append(defaultCertHosts, node.Host, node.IP, node.InternalIP)
	key, cert, err := generateCert(node.Host, plan, nodeList, ca)
	if err != nil {
		return fmt.Errorf("error during cluster cert generation: %v", err)
	}
	err = tls.WriteCert(key, cert, node.Host, lp.GeneratedCertsDirectory)
	if err != nil {
		return fmt.Errorf("error writing cert files for host %q: %v", node.Host, err)
	}
	return nil
}

func (lp *LocalPKI) generateDockerRegistryCert(p *Plan, ca *tls.CA) error {
	// Skip generation if already exist
	exist, err := tls.CertKeyPairExists("docker", lp.GeneratedCertsDirectory)
	if err != nil {
		return fmt.Errorf("error verifying certificates for docker registry: %v", err)
	}
	if exist {
		util.PrettyPrintOk(lp.Log, "Found certificate for docker registry")
	}
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
	return nil
}

func (lp *LocalPKI) generateUserCert(p *Plan, user string, ca *tls.CA) error {
	exists, err := tls.CertKeyPairExists(user, lp.GeneratedCertsDirectory)
	if err != nil {
		return fmt.Errorf("error verifying user certificates: %v", err)
	}
	if exists {
		util.PrettyPrintOk(lp.Log, "Found key and certificate for user %q", user)
		return nil
	}
	util.PrettyPrintOk(lp.Log, "Generating certificates for user %q", user)
	adminKey, adminCert, err := generateCert(user, p, []string{user}, ca)
	if err != nil {
		return fmt.Errorf("error during user cert generation: %v", err)
	}
	err = tls.WriteCert(adminKey, adminCert, user, lp.GeneratedCertsDirectory)
	if err != nil {
		return fmt.Errorf("error writing cert files for user %q: %v", user, err)
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
				O:  certOrganization,
				OU: certOrgUnit,
				C:  certCountry,
				ST: certState,
				L:  certLocality,
			},
		},
	}
	key, cert, err = tls.NewCert(ca, req)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating certs for %q: %v", cnName, err)
	}
	return key, cert, err
}
