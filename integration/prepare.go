package integration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/cloudflare/cfssl/csr"
	. "github.com/onsi/ginkgo"
)

const (
	copyKismaticYumRepo        = `sudo curl https://s3.amazonaws.com/kismatic-rpm/kismatic.repo -o /etc/yum.repos.d/kismatic.repo`
	installEtcdYum             = `sudo yum -y install kismatic-etcd`
	installDockerEngineYum     = `sudo yum -y install kismatic-docker-engine`
	installKubernetesMasterYum = `sudo yum -y install kismatic-kubernetes-master`
	installKubernetesYum       = `sudo yum -y install kismatic-kubernetes-node`

	copyKismaticKeyDeb         = `wget -qO - https://kismatic-deb.s3.amazonaws.com/public.key | sudo apt-key add - `
	copyKismaticRepoDeb        = `sudo add-apt-repository "deb https://kismatic-deb.s3.amazonaws.com/ xenial main"`
	updateAptGet               = `sudo apt-get update`
	installEtcdApt             = `sudo apt-get -y install kismatic-etcd`
	installDockerApt           = `sudo apt-get -y install kismatic-docker-engine`
	installKubernetesMasterApt = `sudo apt-get -y install kismatic-kubernetes-master`
	installKubernetesApt       = `sudo apt-get -y install kismatic-kubernetes-node`
)

type nodePrep struct {
	CommandsToPrepRepo         []string
	CommandsToInstallEtcd      []string
	CommandsToInstallDocker    []string
	CommandsToInstallK8sMaster []string
	CommandsToInstallK8s       []string
}

var ubuntu1604Prep = nodePrep{
	CommandsToPrepRepo:         []string{copyKismaticKeyDeb, copyKismaticRepoDeb, updateAptGet},
	CommandsToInstallEtcd:      []string{installEtcdApt},
	CommandsToInstallDocker:    []string{installDockerApt},
	CommandsToInstallK8sMaster: []string{installKubernetesMasterApt},
	CommandsToInstallK8s:       []string{installKubernetesApt},
}

var centos7Prep = nodePrep{
	CommandsToPrepRepo:         []string{copyKismaticYumRepo},
	CommandsToInstallEtcd:      []string{installEtcdYum},
	CommandsToInstallDocker:    []string{installDockerEngineYum},
	CommandsToInstallK8sMaster: []string{installKubernetesMasterYum},
	CommandsToInstallK8s:       []string{installKubernetesYum},
}

func InstallKismaticRPMs(nodes provisionedNodes, distro linuxDistro, sshKey string) {
	prep := getPrepForDistro(distro)
	sshUser := nodes.master[0].SSHUser
	By("Configuring package repository")
	runViaSSH(prep.CommandsToPrepRepo, sshUser, append(append(nodes.etcd, nodes.master...), nodes.worker...), 5*time.Minute)

	By("Installing Etcd")
	runViaSSH(prep.CommandsToInstallEtcd, sshUser, nodes.etcd, 5*time.Minute)

	By("Installing Docker")
	dockerNodes := append(nodes.master, nodes.worker...)
	runViaSSH(prep.CommandsToInstallDocker, sshUser, dockerNodes, 5*time.Minute)

	By("Installing Master:")
	runViaSSH(prep.CommandsToInstallK8sMaster, sshUser, nodes.master, 5*time.Minute)

	By("Installing Worker:")
	runViaSSH(prep.CommandsToInstallK8s, sshUser, nodes.worker, 5*time.Minute)
}

func getPrepForDistro(distro linuxDistro) nodePrep {
	switch distro {
	case Ubuntu1604LTS:
		return ubuntu1604Prep
	case CentOS7:
		return centos7Prep
	default:
		panic(fmt.Sprintf("Unsupported distro %s", distro))
	}
}

func deployDockerRegistry(node NodeDeets, listeningPort int, sshKey string) (string, error) {
	// Install Docker on the node
	installDockerCmds := []string{
		"sudo curl -sSL https://get.docker.com/ | sh",
		"sudo systemctl start docker",
		"mkdir ~/certs",
	}
	ok := runViaSSH(installDockerCmds, node.SSHUser, []NodeDeets{node}, 10*time.Minute)
	if !ok {
		return "", errors.New("Failed to install Docker on the node")
	}
	// Generate CA
	subject := tls.Subject{
		Organization:       "someOrg",
		OrganizationalUnit: "someOrgUnit",
	}
	key, caCert, err := tls.NewCACert("test-tls/ca-csr.json", "someCommonName", subject)
	if err != nil {
		return "", fmt.Errorf("error generating CA cert for Docker: %v", err)
	}
	err = ioutil.WriteFile("docker-ca.pem", caCert, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing CA cert to file")
	}
	// Generate Certificate
	ca := &tls.CA{
		Key:        key,
		Cert:       caCert,
		ConfigFile: "test-tls/ca-config.json",
		Profile:    "kubernetes",
	}
	certHosts := []string{node.Hostname, node.PrivateIP, node.PublicIP}
	req := csr.CertificateRequest{
		CN: node.Hostname,
		KeyRequest: &csr.BasicKeyRequest{
			A: "rsa",
			S: 2048,
		},
		Hosts: certHosts,
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
	key, cert, err := tls.NewCert(ca, req)
	if err != nil {
		return "", fmt.Errorf("error generating certificate for Docker registry: %v", err)
	}
	if err = ioutil.WriteFile("docker.pem", cert, 0644); err != nil {
		return "", fmt.Errorf("error writing certificate to file: %v", err)
	}
	if err = ioutil.WriteFile("docker-key.pem", key, 0644); err != nil {
		return "", fmt.Errorf("error writing private key to file: %v", err)
	}

	copyFileToRemote("docker.pem", "~/certs/docker.pem", node.SSHUser, []NodeDeets{node}, 1*time.Minute)
	copyFileToRemote("docker-key.pem", "~/certs/docker-key.pem", node.SSHUser, []NodeDeets{node}, 1*time.Minute)

	startDockerRegistryCmd := []string{fmt.Sprintf("sudo docker run -d -p %d:5000 --restart=always ", listeningPort) +
		"--name registry -v ~/certs:/certs -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/docker.pem " +
		"-e REGISTRY_HTTP_TLS_KEY=/certs/docker-key.pem registry"}
	if ok := runViaSSH(startDockerRegistryCmd, node.SSHUser, []NodeDeets{node}, 1*time.Minute); !ok {
		return "", fmt.Errorf("failed to start docker registry")
	}

	// Need the full path, otherwise ansible looks for it in the wrong place
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "docker-ca.pem"), nil
}
