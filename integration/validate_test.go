package integration

import (
	"time"

	"github.com/apprenda/kismatic-platform/integration/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kismatic install validate tests", func() {
	Describe("Running validation with package installation enabled", func() {
		Context("Targetting AWS infrastructure", func() {
			Context("using a minikube layout with Ubuntu 16.04 LTS", func() {
				ItOnAWS("should succeed", func(awsClient *aws.Client) {
					WithMiniInfrastructure(Ubuntu1604LTS, awsClient, "ubuntu", func(node AWSNodeDeets, sshUser, sshKey string) {
						ValidateKismaticMini(node, sshUser, sshKey)
					})
				})
			})
			Context("using a minikube layout with CentOS 7", func() {
				ItOnAWS("should succeed", func(awsClient *aws.Client) {
					WithMiniInfrastructure(CentOS7, awsClient, "centos", func(node AWSNodeDeets, sshUser, sshKey string) {
						ValidateKismaticMini(node, sshUser, sshKey)
					})
				})
			})
		})

		Describe("Running validation with package installation disabled", func() {
			Context("using a minikube layout with Ubuntu 16.04 LTS", func() {
				ItOnAWS("should succeed if and only if all packages are installed on the node", func(awsClient *aws.Client) {
					WithMiniInfrastructure(Ubuntu1604LTS, awsClient, "ubuntu", func(node AWSNodeDeets, sshUser, sshKey string) {
						validateMiniPkgInstallationDisabled(node, Ubuntu1604LTS, sshUser, sshKey)
					})
				})
			})
			Context("Using a minikube layout with CentOS 7", func() {
				ItOnAWS("should succeed if and only if all packages are installed on the node", func(awsClient *aws.Client) {
					WithMiniInfrastructure(CentOS7, awsClient, "centos", func(node AWSNodeDeets, sshUser, sshKey string) {
						validateMiniPkgInstallationDisabled(node, CentOS7, sshUser, sshKey)
					})
				})
			})
		})
	})
})

func validateMiniPkgInstallationDisabled(theNode AWSNodeDeets, distro linuxDistro, sshUser, sshKey string) {
	if err := ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey); err == nil {
		Fail("Missing dependencies, but still passed")
	}

	By("Prepping nodes for the test")
	prep := getPrepForDistro(distro)
	prepNode := []AWSNodeDeets{theNode}
	RunViaSSH(prep.CommandsToPrepRepo, sshUser, prepNode, 5*time.Minute)

	By("Installing etcd on the node")
	RunViaSSH(prep.CommandsToInstallEtcd, sshUser, prepNode, 5*time.Minute)
	if err := ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey); err == nil {
		Fail("Missing dependencies, but still passed")
	}

	By("Installing Docker")
	RunViaSSH(prep.CommandsToInstallDocker, sshUser, prepNode, 5*time.Minute)
	if err := ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey); err == nil {
		Fail("Missing dependencies, but still passed")
	}

	By("Installing Master")
	RunViaSSH(prep.CommandsToInstallK8sMaster, sshUser, prepNode, 5*time.Minute)
	err := ValidateKismaticMiniDenyPkgInstallation(theNode, sshUser, sshKey)
	Expect(err).To(BeNil())
}
