package integration

import (
	"github.com/apprenda/kismatic-platform/integration/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// ItOnAWS runs a spec if the AWS details have been provided
func ItOnAWS(description string, f func(*aws.Client)) {
	It(description, func() {
		awsClient, ok := awsClientFromEnvironment()
		if !ok {
			Skip("AWS environment variables were not defined")
		}
		f(awsClient)
	})
}

type infraDependentTest func(nodes provisionedNodes, sshUser, sshKey string)

// WithInfrastructure runs the spec with the requested infrastructure
func WithInfrastructure(nodeCount NodeCount, distro linuxDistro, client *aws.Client, sshUser string, f infraDependentTest) {
	By("Provisioning nodes")
	nodes, err := provisionAWSNodes(client, nodeCount, distro)
	defer terminateNodes(client, nodes)
	Expect(err).ToNot(HaveOccurred())

	By("Waiting until nodes are SSH-accessible")
	sshKey, err := GetSSHKeyFile()
	Expect(err).ToNot(HaveOccurred())
	err = waitForSSH(nodes, sshUser, sshKey)
	Expect(err).ToNot(HaveOccurred())

	f(nodes, sshUser, sshKey)
}

type miniInfraDependentTest func(node AWSNodeDeets, sshUser, sshKey string)

// WithMiniInfrastructure runs the spec with a Minikube-like infrastructure setup.
func WithMiniInfrastructure(distro linuxDistro, client *aws.Client, sshUser string, f miniInfraDependentTest) {
	By("Provisioning nodes")
	nodes, err := provisionAWSNodes(client, NodeCount{Worker: 1}, distro)
	defer terminateNodes(client, nodes)
	Expect(err).ToNot(HaveOccurred())

	By("Waiting until nodes are SSH-accessible")
	sshKey, err := GetSSHKeyFile()
	Expect(err).ToNot(HaveOccurred())
	err = waitForSSH(nodes, sshUser, sshKey)
	Expect(err).ToNot(HaveOccurred())

	f(nodes.worker[0], sshUser, sshKey)
}
