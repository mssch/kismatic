package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	yaml "gopkg.in/yaml.v2"

	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jmcvetta/guid"
)

var guidMaker = guid.SimpleGenerator()

var _ = Describe("Happy Path Installation Tests", func() {
	kisPath := CopyKismaticToTemp()

	BeforeSuite(func() {
		fmt.Println("Unpacking kismatic to ", kisPath)
		c := exec.Command("tar", "-zxf", "../out/kismatic.tar.gz", "-C", kisPath)
		tarOut, tarErr := c.CombinedOutput()
		if tarErr != nil {
			log.Fatal("Error unpacking installer", string(tarOut), tarErr)
		}
		os.Chdir(kisPath)
	})

	AfterSuite(func() {
		os.RemoveAll(kisPath)
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
				dir, _ := exec.Command("pwd").Output()
				log.Println(string(dir))
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
})

func CopyKismaticToTemp() string {
	tmpDir := os.TempDir()
	randomness, randomErr := GenerateGUIDString()
	if randomErr != nil {
		log.Fatal("Error making a GUID: ", randomErr)
	}
	kisPath := tmpDir + "/kisint/" + randomness
	err := os.MkdirAll(kisPath, 0777)
	if err != nil {
		log.Fatal("Error making temp dir: ", err)
	}

	return kisPath
}

func GenerateGUIDString() (string, error) {
	randomness, randomErr := guidMaker.NextId()

	if randomErr != nil {
		return "", randomErr
	}

	return strconv.FormatInt(randomness, 16), nil
}

func AssertKismaticDirectory(kisPath string) {
	if FileExists(kisPath + "/kismatic") {
		log.Fatal("Installer unpacked but kismatic wasn't there")
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
