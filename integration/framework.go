package integration

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// ItOnAWS runs a spec if the AWS details have been provided
func ItOnAWS(description string, f func(infrastructureProvisioner)) {
	Context("when using AWS infrastructure [aws]", func() {
		It(description, func() {
			awsClient, ok := AWSClientFromEnvironment()
			if !ok {
				Skip("AWS environment variables were not defined")
			}
			f(awsClient)
		})
	})
}

// ItOnPacket runs a spec if the Packet.Net details have been provided
func ItOnPacket(description string, f func(infrastructureProvisioner)) {
	Context("when using Packet.net infrastructure [packet]", func() {
		It(description, func() {
			packetClient, ok := packetClientFromEnv()
			if !ok {
				Skip("Packet environment variables were not defined")
			}
			f(packetClient)
		})
	})
}

type infraDependentTest func(nodes provisionedNodes, sshKey string)

// WithInfrastructure runs the spec with the requested infrastructure
func WithInfrastructure(nodeCount NodeCount, distro linuxDistro, provisioner infrastructureProvisioner, f infraDependentTest) {
	By(fmt.Sprintf("Provisioning nodes: %+v", nodeCount))
	start := time.Now()
	nodes, err := provisioner.ProvisionNodes(nodeCount, distro)
	if !leaveIt() {
		defer provisioner.TerminateNodes(nodes)
	}
	Expect(err).ToNot(HaveOccurred())
	fmt.Println("Provisioning infrastructure took", time.Since(start))

	By("Waiting until nodes are SSH-accessible")
	start = time.Now()
	sshKey := provisioner.SSHKey()
	err = waitForSSH(nodes, sshKey)
	Expect(err).ToNot(HaveOccurred())
	fmt.Println("Waiting for SSH took", time.Since(start))

	f(nodes, sshKey)
}

// WithInfrastructureAndDNS runs the spec with the requested infrastructure and DNS
func WithInfrastructureAndDNS(nodeCount NodeCount, distro linuxDistro, provisioner infrastructureProvisioner, f infraDependentTest) {
	By(fmt.Sprintf("Provisioning nodes and DNS: %+v", nodeCount))
	start := time.Now()
	nodes, err := provisioner.ProvisionNodes(nodeCount, distro)
	if !leaveIt() {
		defer provisioner.TerminateNodes(nodes)
	}
	Expect(err).ToNot(HaveOccurred())
	fmt.Println("Provisioning infrastructure took", time.Since(start))

	By("Waiting until nodes are SSH-accessible")
	start = time.Now()
	sshKey := provisioner.SSHKey()
	err = waitForSSH(nodes, sshKey)
	Expect(err).ToNot(HaveOccurred())
	fmt.Println("Waiting for SSH took", time.Since(start))

	By("Configuring DNS entries")
	start = time.Now()
	var masterIPs []string
	for _, node := range nodes.master {
		masterIPs = append(masterIPs, node.PrivateIP)
	}
	dnsRecord, err := provisioner.ConfigureDNS(masterIPs)
	nodes.dnsRecord = dnsRecord
	Expect(err).ToNot(HaveOccurred())
	if !leaveIt() {
		By("Removing DNS entries")
		defer provisioner.RemoveDNS(dnsRecord)
	}
	fmt.Println("Configuring DNS entries took", time.Since(start))

	f(nodes, sshKey)
}

type miniInfraDependentTest func(node NodeDeets, sshKey string)

// WithMiniInfrastructure runs the spec with a Minikube-like infrastructure setup.
func WithMiniInfrastructure(distro linuxDistro, provisioner infrastructureProvisioner, f miniInfraDependentTest) {
	By("Provisioning minikube node")
	start := time.Now()
	nodes, err := provisioner.ProvisionNodes(NodeCount{Worker: 1}, distro)
	if !leaveIt() {
		defer provisioner.TerminateNodes(nodes)
	}
	Expect(err).ToNot(HaveOccurred())
	fmt.Println("Provisioning node took", time.Since(start))

	By("Waiting until nodes are SSH-accessible")
	start = time.Now()
	sshKey := provisioner.SSHKey()
	err = waitForSSH(nodes, sshKey)
	Expect(err).ToNot(HaveOccurred())
	fmt.Println("Waiting for SSH took", time.Since(start))

	f(nodes.worker[0], sshKey)
}

// SubDescribe allows you to define specifications inside another spec.
// We have found the need for this because Gingko does not support
// serializing a subset of tests when running in parallel. This means
// that we must define multiple specs inside a parent It() block.
// Use this when it is truly needed.
//
// Example:
// Describe("the foo service", func() {
//	It("should be deployed successfully", func() {
//		// some assertions here...
//		sub := SubDescribe("should return 200", func() error {
//			// call service and return error if not 200
//		})
//	})
// })
func SubDescribe(name string) *subTest {
	return &subTest{name: name}
}

type subTest struct {
	name   string
	specs  []string
	errors []error
}

func (sub *subTest) It(name string, f func() error) {
	By(fmt.Sprintf("Running spec: %s - %s", sub.name, name))
	sub.specs = append(sub.specs, name)
	sub.errors = append(sub.errors, f())
}

func (sub *subTest) Check() {
	var failed bool
	for _, e := range sub.errors {
		if e != nil {
			failed = true
			break
		}
	}
	if failed {
		// Print report and fail test
		sub.printReport()
		Fail("Failed: " + sub.name)
	}
}

func (sub *subTest) printReport() {
	fmt.Println(sub.name)
	for i, spec := range sub.specs {
		if sub.errors[i] != nil {
			fmt.Printf("%s: FAILED: %v\n", spec, sub.errors[i])
		} else {
			fmt.Printf("%s: PASSED\n", spec)
		}
	}
}
