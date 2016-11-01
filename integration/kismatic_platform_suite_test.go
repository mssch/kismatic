package integration

import (
	"os"

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
	kisPath, err = ExtractKismaticToTemp()
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
