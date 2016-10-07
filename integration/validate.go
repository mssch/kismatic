package integration

import (
	"bufio"
	"html/template"
	"log"
	"os"
	"os/exec"

	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
)

func ValidateKismaticMini(nodeType string, user string) PlanAWS {
	By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")

	By("Making infrastructure")
	etcdNode, etcErr := MakeWorkerNode(nodeType)
	FailIfError(etcErr, "Error making etcd node")

	defer TerminateInstances(etcdNode.Instanceid)
	descErr := WaitForInstanceToStart(&etcdNode)
	masterNode := etcdNode
	workerNode := etcdNode
	FailIfError(descErr, "Error waiting for nodes")
	log.Printf("Created etcd nodes: %v (%v), master nodes %v (%v), workerNodes %v (%v)",
		etcdNode.Instanceid, etcdNode.Publicip,
		masterNode.Instanceid, masterNode.Publicip,
		workerNode.Instanceid, workerNode.Publicip)

	By("Building a plan to set up an overlay network cluster on this hardware")
	nodes := PlanAWS{
		Etcd:                []AWSNodeDeets{etcdNode},
		Master:              []AWSNodeDeets{masterNode},
		Worker:              []AWSNodeDeets{workerNode},
		MasterNodeFQDN:      masterNode.Hostname,
		MasterNodeShortName: masterNode.Hostname,
		User:                user,
	}
	var hdErr error
	nodes.HomeDirectory, hdErr = homedir.Dir()
	FailIfError(hdErr, "Error getting home directory")

	f, fileErr := os.Create("kismatic-testing.yaml")
	FailIfError(fileErr, "Error waiting for nodes")
	defer f.Close()
	w := bufio.NewWriter(f)
	execErr := template.Execute(w, &nodes)
	FailIfError(execErr, "Error filling in plan template")
	w.Flush()

	By("Validate our plan")
	ver := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
	ver.Stdout = os.Stdout
	ver.Stderr = os.Stderr
	err = ver.Run()

	FailIfError(err, "Error validating plan")

	if bailBeforeAnsible() == true {
		return nodes
	}
	return nodes
}
