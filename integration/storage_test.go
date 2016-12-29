package integration

import (
	"os"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("Storage", func() {
	BeforeEach(func() {
		os.Chdir(kisPath)
	})

	Context("Specifying multiple storage nodes in the plan file", func() {
		Context("targetting CentOS", func() {
			ItOnAWS("should result in a working storage cluster", func(aws infrastructureProvisioner) {
				testGlusterCluster(aws, CentOS7)
			})
		})
		Context("targetting Ubuntu", func() {
			ItOnAWS("should result in a working storage cluster", func(aws infrastructureProvisioner) {
				testGlusterCluster(aws, Ubuntu1604LTS)
			})
		})
	})
})
