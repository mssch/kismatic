package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kismatic install add-worker tests", func() {
	Context("Targetting AWS infrastructure", func() {
		Context("Using Ubuntu 16.04 LTS", func() {
			ItOnAWS("should successfully add a new worker", func(provisioner infrastructureProvisioner) {
				WithInfrastructure(NodeCount{Worker: 2}, Ubuntu1604LTS, provisioner, "ubuntu", func(nodes provisionedNodes, sshUser, sshKey string) {
					theNode := nodes.worker[0]
					err := installKismaticMini(theNode, sshUser, sshKey)
					Expect(err).ToNot(HaveOccurred())

					newWorker := nodes.worker[1]
					err = addWorkerToKismaticMini(newWorker)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})
