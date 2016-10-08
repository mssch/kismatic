package integration

import (
	"bufio"
	"html/template"
	"log"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
)

// ValidateKismaticMini runs validation against a mini kismatic cluster
func ValidateKismaticMini(nodeType string, user string) PlanAWS {
	By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")

	By("Making infrastructure")
	node, err := MakeWorkerNode(nodeType)
	FailIfError(err, "Error making etcd node")
	defer TerminateInstances(node.Instanceid)

	sshKey, err := GetSSHKeyFile()
	FailIfError(err, "Error getting SSH Key")

	descErr := WaitForInstanceToStart(user, sshKey, &node)
	FailIfError(descErr, "Error waiting for nodes")

	log.Printf("Created single node for Kismatic Mini: %s (%s)", node.Instanceid, node.Publicip)
	By("Building a plan to set up an overlay network cluster on this hardware")
	nodes := PlanAWS{
		Etcd:                []AWSNodeDeets{node},
		Master:              []AWSNodeDeets{node},
		Worker:              []AWSNodeDeets{node},
		MasterNodeFQDN:      node.Hostname,
		MasterNodeShortName: node.Hostname,
		SSHUser:             user,
		SSHKeyFile:          sshKey,
	}

	// Create template file
	f, fileErr := os.Create("kismatic-testing.yaml")
	FailIfError(fileErr, "Error waiting for nodes")
	defer f.Close()
	w := bufio.NewWriter(f)
	execErr := template.Execute(w, &nodes)
	FailIfError(execErr, "Error filling in plan template")
	w.Flush()

	// Run validation
	By("Validate our plan")
	ver := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
	ver.Stdout = os.Stdout
	ver.Stderr = os.Stderr
	err = ver.Run()
	FailIfError(err, "Error validating plan")
	return nodes
}
