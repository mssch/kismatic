package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKismaticPlatform(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KismaticPlatform Suite")
}

var kisPath string
var _ = BeforeSuite(func() {
	var err error
	kisPath, err = extractKismaticToTemp()
	if err != nil {
		Fail("Failed to extract kismatic")
	}
	os.Chdir(kisPath)
})

var _ = AfterSuite(func() {
	if !leaveIt() {
		os.RemoveAll(kisPath)
	}
})

func extractKismaticToTemp() (string, error) {
	tmpDir, err := ioutil.TempDir("", "kisint")
	if err != nil {
		log.Fatal("Error making temp dir: ", err)
	}
	By(fmt.Sprintf("Extracting Kismatic to temp directory %q", tmpDir))
	cmd := exec.Command("tar", "-zxf", "../out/kismatic.tar.gz", "-C", tmpDir)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error extracting kismatic to temp dir: %v", err)
	}
	return tmpDir, nil
}
