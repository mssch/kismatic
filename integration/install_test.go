package integration

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

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

	Describe("Calling installer with install plan", func() {
		Context("and just hitting enter", func() {
			It("should result in the output of a well formed default plan file", func() {
				dir, _ := exec.Command("pwd").Output()
				log.Println(string(dir))
				c := exec.Command("./kismatic", "install", "plan")
				helpbytes, helperr := c.Output()
				Expect(helperr).To(BeNil())
				helpText := string(helpbytes)
				Expect(helpText).To(ContainSubstring("Generating installation plan file with 3 etcd nodes, 2 master nodes and 3 worker nodes"))

				Expect(FileExists("kismatic-cluster.yaml")).To(Equal(true))
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
	kisPath := tmpDir + "kisint/" + randomness
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

func new_scanner(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		fmt.Printf("nn\n")
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '?'); i >= 0 {
		// We have a full ?-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
