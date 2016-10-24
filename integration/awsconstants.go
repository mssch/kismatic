package integration

const TARGET_REGION = "us-east-1"

const SUBNETID = "subnet-edab6fd1"

//"subnet-edab6fd1"

const KEYNAME = "kismatic-integration-testing"
const SECURITYGROUPID = "sg-d1dc4dab"

const AMIUbuntu1604USEAST = "ami-29f96d3e"
const AMICentos7UsEast = "ami-6d1c2007"

var UbuntuEast = AWSOSDetails{
	AWSAMI:  AMIUbuntu1604USEAST,
	AWSUser: "ubuntu",

	CommandsToPrepRepo:         []string{CopyKismaticKeyDeb, CopyKismaticRepoDeb, UpdateAptGet},
	CommandsToInstallEtcd:      []string{InstallEtcdApt},
	CommandsToInstallDocker:    []string{InstallDockerApt},
	CommandsToInstallK8sMaster: []string{InstallKubernetesMasterApt},
	CommandsToInstallK8s:       []string{InstallKubernetesApt},
}

var CentosEast = AWSOSDetails{
	AWSAMI:  AMICentos7UsEast,
	AWSUser: "centos",

	CommandsToPrepRepo:         []string{CopyKismaticYumRepo},
	CommandsToInstallEtcd:      []string{InstallEtcdYum},
	CommandsToInstallDocker:    []string{InstallDockerEngineYum},
	CommandsToInstallK8sMaster: []string{InstallKubernetesMasterYum},
	CommandsToInstallK8s:       []string{InstallKubernetesYum},
}

type AWSOSDetails struct {
	AWSAMI  string
	AWSUser string

	CommandsToPrepRepo         []string
	CommandsToInstallEtcd      []string
	CommandsToInstallDocker    []string
	CommandsToInstallK8sMaster []string
	CommandsToInstallK8s       []string
}
