package integration

import (
	"github.com/apprenda/kismatic-platform/integration/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kismatic install add-worker tests", func() {
	Context("Targetting AWS infrastructure", func() {
		Context("Using a Minikube Ubuntu layout", func() {
			ItOnAWS("should successfully add a new worker", func(awsClient *aws.Client) {
				By("Provisioning nodes")
				nodes, err := provisionAWSNodes(awsClient, NodeCount{Worker: 2}, CentOS7)
				defer terminateNodes(awsClient, nodes)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting until nodes are SSH-accessible")
				sshUser := "centos"
				sshKey, err := GetSSHKeyFile()
				Expect(err).ToNot(HaveOccurred())
				err = waitForSSH(nodes, sshUser, sshKey)
				Expect(err).ToNot(HaveOccurred())

				theNode := nodes.worker[0]
				err = installKismaticMini(theNode, sshUser, sshKey)
				Expect(err).ToNot(HaveOccurred())

				newWorker := nodes.worker[1]
				err = addWorkerToKismaticMini(newWorker)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
