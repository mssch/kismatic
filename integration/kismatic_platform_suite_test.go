package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKismaticPlatform(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KismaticPlatform Suite")
}
