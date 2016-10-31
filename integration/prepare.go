package integration

import (
	"fmt"
	"time"

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

func InstallKismaticRPMs(nodes provisionedNodes, distro linuxDistro, sshUser, sshKey string) {
	prep := getPrepForDistro(distro)
	By("Configuring package repository")
	RunViaSSH(prep.CommandsToPrepRepo, sshUser, append(append(nodes.etcd, nodes.master...), nodes.worker...), 5*time.Minute)

	By("Installing Etcd")
	RunViaSSH(prep.CommandsToInstallEtcd, sshUser, nodes.etcd, 5*time.Minute)

	By("Installing Docker")
	dockerNodes := append(nodes.master, nodes.worker...)
	RunViaSSH(prep.CommandsToInstallDocker, sshUser, dockerNodes, 5*time.Minute)

	By("Installing Master:")
	RunViaSSH(prep.CommandsToInstallK8sMaster, sshUser, nodes.master, 5*time.Minute)

	By("Installing Worker:")
	RunViaSSH(prep.CommandsToInstallK8s, sshUser, nodes.worker, 5*time.Minute)
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
