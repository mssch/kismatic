package integration

import (
	"os"
	"os/exec"

	g "github.com/onsi/ginkgo"
)

func addWorkerToKismaticMini(newWorker NodeDeets) error {
	g.By("Adding new worker")
	app := exec.Command("./kismatic", "install", "add-worker", "-f", "kismatic-testing.yaml", newWorker.Hostname, newWorker.PublicIP, newWorker.PrivateIP)
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	return app.Run()
}
