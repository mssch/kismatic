package integration

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"

	"os/exec"

	"github.com/apprenda/kismatic-platform/integration/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Happy Path Installation Tests", func() {
	Describe("Calling installer with no input", func() {
		It("should output help text", func() {
			c := exec.Command("./kismatic")
			helpbytes, helperr := c.Output()
			Expect(helperr).To(BeNil())
			helpText := string(helpbytes)
			Expect(helpText).To(ContainSubstring("Usage"))
		})
	})

	Describe("Calling installer with 'install plan'", func() {
		Context("and just hitting enter", func() {
			It("should result in the output of a well formed default plan file", func() {
				By("Outputing a file")
				c := exec.Command("./kismatic", "install", "plan")
				helpbytes, helperr := c.Output()
				Expect(helperr).To(BeNil())
				helpText := string(helpbytes)
				Expect(helpText).To(ContainSubstring("Generating installation plan file with 3 etcd nodes, 2 master nodes and 3 worker nodes"))
				Expect(FileExists("kismatic-cluster.yaml")).To(Equal(true))

				By("Outputing a file with valid YAML")
				yamlBytes, err := ioutil.ReadFile("kismatic-cluster.yaml")
				if err != nil {
					Fail("Could not read cluster file")
				}
				yamlBlob := string(yamlBytes)

				planFromYaml := ClusterPlan{}

				unmarshallErr := yaml.Unmarshal([]byte(yamlBlob), &planFromYaml)
				if unmarshallErr != nil {
					Fail("Could not unmarshall cluster yaml: %v")
				}
			})
		})
	})
	Describe("Calling installer with a plan targeting bad infrastructure", func() {
		Context("Using a 1/1/1 Ubtunu 16.04 layout pointing to bad ip addresses", func() {
			It("should bomb validate and apply", func() {
				if !completesInTime(installKismaticWithABadNode, 30*time.Second) {
					Fail("It shouldn't take 30 seconds for Kismatic to fail with bad nodes.")
				}
			})
		})
	})

	Describe("Installing with package installation enabled", func() {
		Context("Targetting AWS infrastructure", func() {
			Context("Using a 1/1/1 Ubuntu 16.04 layout", func() {
				ItOnAWS("should result in a working cluster", func(awsClient *aws.Client) {
					installKismaticPkgInstallEnabled(awsClient, NodeCount{1, 1, 1}, Ubuntu1604LTS, "ubuntu")
				})
			})

			Context("Using a 1/1/1 CentOS 7 layout", func() {
				ItOnAWS("should result in a working cluster", func(awsClient *aws.Client) {
					installKismaticPkgInstallEnabled(awsClient, NodeCount{1, 1, 1}, CentOS7, "centos")
				})
			})

			Context("Using a 3/2/3 CentOS 7 layout", func() {
				ItOnAWS("should result in a working cluster", func(awsClient *aws.Client) {
					installKismaticPkgInstallEnabled(awsClient, NodeCount{3, 2, 3}, CentOS7, "centos")
				})
			})
		})
	})

	Describe("Installing against a minikube layout", func() {
		Context("Using CentOS 7", func() {
			ItOnAWS("should result in a working cluster", func(awsClient *aws.Client) {
				By("Provisioning nodes")
				nodes, err := provisionAWSNodes(awsClient, NodeCount{Worker: 1}, CentOS7)
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
			})
		})
		Context("Using Ubuntu 16.04", func() {
			ItOnAWS("should result in a working cluster", func(awsClient *aws.Client) {
				By("Provisioning nodes")
				nodes, err := provisionAWSNodes(awsClient, NodeCount{Worker: 1}, Ubuntu1604LTS)
				defer terminateNodes(awsClient, nodes)
				Expect(err).ToNot(HaveOccurred())

				By("Waiting until nodes are SSH-accessible")
				sshUser := "ubuntu"
				sshKey, err := GetSSHKeyFile()
				Expect(err).ToNot(HaveOccurred())
				err = waitForSSH(nodes, sshUser, sshKey)
				Expect(err).ToNot(HaveOccurred())

				theNode := nodes.worker[0]
				err = installKismaticMini(theNode, sshUser, sshKey)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Installing with package installation disabled", func() {
		Context("Targetting AWS infrastructure", func() {
			Context("Using a 1/1/1 Ubuntu 16.04 layout", func() {
				ItOnAWS("Should result in a working cluster", func(awsClient *aws.Client) {
					By("Provisioning nodes")
					distro := Ubuntu1604LTS
					nodes, err := provisionAWSNodes(awsClient, NodeCount{1, 1, 1}, distro)
					defer terminateNodes(awsClient, nodes)
					Expect(err).ToNot(HaveOccurred())

					By("Waiting until nodes are SSH-accessible")
					sshUser := "ubuntu"
					sshKey, err := GetSSHKeyFile()
					Expect(err).ToNot(HaveOccurred())
					err = waitForSSH(nodes, sshUser, sshKey)
					Expect(err).ToNot(HaveOccurred())

					By("Installing the Kismatic RPMs")
					InstallKismaticRPMs(nodes, distro, sshUser, sshKey)

					installOpts := installOptions{
						allowPackageInstallation: false,
					}
					err = installKismatic(nodes, installOpts, sshUser, sshKey)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("Using a 1/1/1 CentOS 7 layout", func() {
				ItOnAWS("Should result in a working cluster", func(awsClient *aws.Client) {
					By("Provisioning nodes")
					distro := CentOS7
					nodes, err := provisionAWSNodes(awsClient, NodeCount{1, 1, 1}, distro)
					defer terminateNodes(awsClient, nodes)
					Expect(err).ToNot(HaveOccurred())

					By("Waiting until nodes are SSH-accessible")
					sshUser := "centos"
					sshKey, err := GetSSHKeyFile()
					Expect(err).ToNot(HaveOccurred())
					err = waitForSSH(nodes, sshUser, sshKey)
					Expect(err).ToNot(HaveOccurred())

					By("Installing the Kismatic RPMs")
					InstallKismaticRPMs(nodes, distro, sshUser, sshKey)

					installOpts := installOptions{
						allowPackageInstallation: false,
					}
					err = installKismatic(nodes, installOpts, sshUser, sshKey)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})

func installKismaticPkgInstallEnabled(awsClient *aws.Client, nodeCount NodeCount, distro linuxDistro, sshUser string) {
	By("Provisioning nodes")
	nodes, err := provisionAWSNodes(awsClient, nodeCount, distro)
	defer terminateNodes(awsClient, nodes)
	Expect(err).ToNot(HaveOccurred())

	By("Waiting until nodes are SSH-accessible")
	sshKey, err := GetSSHKeyFile()
	Expect(err).ToNot(HaveOccurred())
	err = waitForSSH(nodes, sshUser, sshKey)
	Expect(err).ToNot(HaveOccurred())

	By("Running installation against infrastructure")
	installOpts := installOptions{
		allowPackageInstallation: true,
	}
	err = installKismatic(nodes, installOpts, sshUser, sshKey)
	Expect(err).ToNot(HaveOccurred())
}
