package integration

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Happy Path Installation Tests", func() {
	kisPath := CopyKismaticToTemp()

	BeforeSuite(func() {
		fmt.Println("Unpacking kismatic to", kisPath)
		c := exec.Command("tar", "-zxf", "../out/kismatic.tar.gz", "-C", kisPath)
		tarOut, tarErr := c.CombinedOutput()
		if tarErr != nil {
			log.Fatal("Error unpacking installer", string(tarOut), tarErr)
		}
		os.Chdir(kisPath)
	})

	AfterSuite(func() {
		if !leaveIt() {
			os.RemoveAll(kisPath)
		}
	})

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

	Describe("Calling installer with a plan targetting AWS", func() {
		Context("Using a 1/1/1 Ubtunu 16.04 layout", func() {
			It("should result in a working cluster", func() {
				InstallKismatic(AMIUbuntu1604USEAST, "ubuntu")
			})
		})
		Context("Using a 1/1/1 CentOS 7 layout", func() {
			It("should result in a working cluster", func() {
				InstallKismatic(AMICentos7UsEast, "centos")
			})
		})
		Context("Using a Minikube CentOS 7 layout", func() {
			It("should result in a working cluster", func() {
				InstallKismaticMini(AMICentos7UsEast, "centos")
			})
		})
		Context("Using a 3/2/3 CentOS 7 layout", func() {
			It("should result in a working cluster", func() {
				InstallBigKismatic(
					NodeCount{
						Etcd:   3,
						Master: 2,
						Worker: 3,
					},
					AMICentos7UsEast, "centos")
			})
		})
	})
})
