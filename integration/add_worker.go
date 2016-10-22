package integration

import (
	"bufio"
	"html/template"
	"os"
	"os/exec"

	g "github.com/onsi/ginkgo"
)

func AddWorkerToKismaticMini(awsos AWSOSDetails) {
	g.By("Making infrastructure")
	node, err := MakeWorkerNode(awsos.AWSAMI)
	FailIfError(err, "Error making node")
	defer TerminateInstances(node.Instanceid)
	// Create worker that will be added
	newWorker, err := MakeWorkerNode(awsos.AWSAMI)
	FailIfError(err, "Error making node")
	defer TerminateInstances(newWorker.Instanceid)
	// Get SSH key for accessing nodes
	sshKey, err := GetSSHKeyFile()
	FailIfError(err, "Error getting SSH Key file")
	// Wait for both nodes to be up
	err = WaitForInstanceToStart(awsos.AWSUser, sshKey, &node, &newWorker)
	FailIfError(err, "Error waiting for nodes")

	g.By("Building a plan to set up an overlay network cluster on this hardware")
	nodes := PlanAWS{
		Etcd:                     []AWSNodeDeets{node},
		Master:                   []AWSNodeDeets{node},
		Worker:                   []AWSNodeDeets{node},
		MasterNodeFQDN:           node.Publicip,
		MasterNodeShortName:      node.Publicip,
		SSHKeyFile:               sshKey,
		SSHUser:                  awsos.AWSUser,
		AllowPackageInstallation: true,
	}

	g.By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")
	f, err := os.Create("kismatic-testing.yaml")
	FailIfError(err, "Error waiting for nodes")
	defer f.Close()
	w := bufio.NewWriter(f)
	err = template.Execute(w, &nodes)
	FailIfError(err, "Error filling in plan template")
	w.Flush()

	g.By("Punch it Chewie!")
	app := exec.Command("./kismatic", "install", "apply", "-f", f.Name())
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	err = app.Run()
	FailIfError(err, "Error applying plan")

	g.By("Adding new worker")
	app = exec.Command("./kismatic", "install", "add-worker", "-f", f.Name(), newWorker.Hostname, newWorker.Publicip, newWorker.Privateip)
	app.Stdout = os.Stdout
	app.Stderr = os.Stderr
	err = app.Run()
	FailIfError(err, "Error adding worker node")
}
