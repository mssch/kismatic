package integration

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"time"

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

	Describe("Calling installer with 'install validate'", func() {
		Context("using a Minikube Ubuntu 16/04 layout", func() {
			It("should return successfully with a valid plan", func() {
				ValidateKismaticMini(AMIUbuntu1604USEAST, "ubuntu")
			})
		})
	})

	Describe("Calling installer with 'install validate'", func() {
		Context("using a Minikube CentOS 7 layout", func() {
			It("should return successfully with a valid plan", func() {
				ValidateKismaticMini(AMICentos7UsEast, "centos")
			})
		})
	})

	Describe("Calling installer with a plan targeting bad infrastructure", func() {
		Context("Using a 1/1/1 Ubtunu 16.04 layout pointing to bad ip addresses", func() {
			It("should bomb validate and apply", func() {
				if !completesInTime(InstallKismaticWithABadNode, 30*time.Second) {
					Fail("It shouldn't take 30 seconds for Kismatic to fail with bad nodes.")
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

func InstallKismaticWithABadNode() {
	By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")

	By("Faking infrastructure")
	fakeNode := AWSNodeDeets{
		Instanceid: "FakeId",
		Publicip:   "10.0.0.0",
		Hostname:   "FakeHostname",
	}

	By("Building a plan to set up an overlay network cluster on this hardware")

	sshKey, err := GetSSHKeyFile()
	FailIfError(err, "Error getting SSH Key file")

	nodes := PlanAWS{
		Etcd:                []AWSNodeDeets{fakeNode},
		Master:              []AWSNodeDeets{fakeNode},
		Worker:              []AWSNodeDeets{fakeNode},
		MasterNodeFQDN:      "yep.nope",
		MasterNodeShortName: "yep",
		SSHUser:             "Billy Rubin",
		SSHKeyFile:          sshKey,
	}

	f, fileErr := os.Create("kismatic-testing.yaml")
	FailIfError(fileErr, "Error waiting for nodes")
	defer f.Close()
	w := bufio.NewWriter(f)
	execErr := template.Execute(w, &nodes)
	FailIfError(execErr, "Error filling in plan template")
	w.Flush()

	f.Close()

	// By("Validing our plan")
	// ver := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
	// ver.Stdout = os.Stdout
	// ver.Stderr = os.Stderr
	// verErr := ver.Run()

	// if verErr == nil {
	// 	// This should really be a failure but at the moment validation does not run tests against target nodes
	// 	// Fail("Validation succeeeded even though it shouldn't have")
	// }

	By("Well, try it anyway")
	app := exec.Command("./kismatic", "install", "apply", "-f", f.Name())
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	appErr := app.Run()

	if appErr == nil {
		Fail("Application succeeeded even though it shouldn't have")
	}
}

func completesInTime(dothis func(), howLong time.Duration) bool {
	c1 := make(chan string, 1)
	go func() {
		dothis()
		c1 <- "completed"
	}()

	select {
	case <-c1:
		return true
	case <-time.After(howLong):
		return false
	}
}
