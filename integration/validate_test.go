package integration

import (
	"time"

	"github.com/apprenda/kismatic-platform/integration/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kismatic install validate tests", func() {
	Context("Targetting AWS infrastructure", func() {
		Describe("Running validation with package installation enabled", func() {
			Context("using a Minikube Ubuntu layout", func() {
				ItOnAWS("should succeed if and only if all packages are installed on the node", func(awsClient *aws.Client) {
					validateMiniPkgInstallEnabled(awsClient, Ubuntu1604LTS, "ubuntu")
				})
			})
			Context("using a Minikube CentOS 7 layout", func() {
				ItOnAWS("should succeed if and only if all packages are installed on the node", func(awsClient *aws.Client) {
					validateMiniPkgInstallEnabled(awsClient, CentOS7, "centos")
				})
			})
		})

		Describe("Running validation with package installation disabled", func() {
			Context("using a Minikube Ubuntu layout", func() {
				ItOnAWS("should succeed if and only if all packages are installed on the node", func(awsClient *aws.Client) {
					validateMiniPkgInstallationDisabled(awsClient, Ubuntu1604LTS, "ubuntu")
				})
			})
			Context("using a Minikube CentOS 7 layout", func() {
				ItOnAWS("should succeed if and only if all packages are installed on the node", func(awsClient *aws.Client) {
					validateMiniPkgInstallationDisabled(awsClient, CentOS7, "centos")
				})
			})
		})
	})
})

func ItOnAWS(description string, f func(*aws.Client)) {
	It(description, func() {
		awsClient, ok := awsClientFromEnvironment()
		if !ok {
			Skip("AWS environment variables were not defined")
		}
		f(awsClient)
	})
}

func validateMiniPkgInstallEnabled(client *aws.Client, distro linuxDistro, sshUser string) {
	By("Provisioning nodes on AWS")
	nodes, err := provisionAWSNodes(client, NodeCount{Worker: 1}, distro)
	defer terminateNodes(client, nodes)
	FailIfError(err, "Failed to provision nodes for test")

	By("Waiting until nodes are SSH-accessible")
	sshKey, err := GetSSHKeyFile()
	FailIfError(err, "Failed to get SSH Key")
	err = waitForSSH(nodes, sshUser, sshKey)
	FailIfError(err, "Error waiting for nodes to become SSH-accessible")
	ValidateKismaticMini(nodes.worker[0], sshUser, sshKey)
}

func validateMiniPkgInstallationDisabled(client *aws.Client, distro linuxDistro, sshUser string) {
	By("Provisioning nodes on AWS")
	nodes, err := provisionAWSNodes(client, NodeCount{Worker: 1}, distro)
	defer terminateNodes(client, nodes)
	FailIfError(err, "Failed to provision nodes for test")

	By("Waiting until nodes are SSH-accessible")
	sshKey, err := GetSSHKeyFile()
	FailIfError(err, "Failed to get SSH Key")
	err = waitForSSH(nodes, sshUser, sshKey)
	FailIfError(err, "Error waiting for nodes to become SSH-accessible")
	theNode := nodes.worker[0]

	if err = ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey); err == nil {
		Fail("Missing dependencies, but still passed")
	}

	By("Prepping nodes for the test")
	prep := getPrepForDistro(distro)
	prepNode := []AWSNodeDeets{theNode}
	RunViaSSH(prep.CommandsToPrepRepo, sshUser, prepNode, 5*time.Minute)
	By("Installing etcd on the node")
	RunViaSSH(prep.CommandsToInstallEtcd, sshUser, prepNode, 5*time.Minute)
	if err = ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey); err == nil {
		Fail("Missing dependencies, but still passed")
	}

	By("Installing Docker")
	RunViaSSH(prep.CommandsToInstallDocker, sshUser, prepNode, 5*time.Minute)
	if err = ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey); err == nil {
		Fail("Missing dependencies, but still passed")
	}

	By("Installing Master")
	RunViaSSH(prep.CommandsToInstallK8sMaster, sshUser, prepNode, 5*time.Minute)
	err = ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey)
	Expect(err).To(BeNil())
}
